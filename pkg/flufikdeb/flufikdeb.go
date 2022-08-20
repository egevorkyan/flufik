package flufikdeb

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"crypto/md5"
	"fmt"
	"github.com/blakesmith/ar"
	"github.com/egevorkyan/flufik/crypto/pgp"
	"github.com/egevorkyan/flufik/pkg/logging"
	"io"
	"strings"
	"time"
)

type tarGzWriter struct {
	buffer *bytes.Buffer
	gz     *gzip.Writer
	tar    *tar.Writer
}

type md5Writer struct {
	buffer *bytes.Buffer
	tmp    []byte
}

type FlufikDeb struct {
	FlufikDebMetaData
	Signature FlufikDebSignature

	files []FlufikDebFile

	preIn  string
	postIn string
	preUn  string
	postUn string

	configFiles *bytes.Buffer
	logger      *logging.Logger
	debugging   string
}

type FlufikDebSignature struct {
	FlufikPackageSignature
	Type string
}

type FlufikPackageSignature struct {
	PgpName string
}

func (flufikTgz *tarGzWriter) WriteHeader(header *tar.Header) error {
	return flufikTgz.tar.WriteHeader(header)
}

func (flufikTgz *tarGzWriter) Write(header *tar.Header, b []byte) (int, error) {
	if err := flufikTgz.WriteHeader(header); err != nil {
		return -1, fmt.Errorf("write header failed: %w", err)
	}

	if size, err := flufikTgz.tar.Write(b); err != nil {
		return -1, fmt.Errorf("write body failed: %w", err)
	} else {
		return size, nil
	}
}

func (flufikTgz *tarGzWriter) Close() error {
	if err := flufikTgz.tar.Close(); err != nil {
		return fmt.Errorf("close tar failed: %w", err)
	}

	if err := flufikTgz.gz.Close(); err != nil {
		return fmt.Errorf("close gz failed: %w", err)
	}

	return nil
}

func (flufikTgz *tarGzWriter) Bytes() []byte {
	return flufikTgz.buffer.Bytes()
}

func newFlufikTgz() *tarGzWriter {
	buffer := new(bytes.Buffer)
	gzWriter := gzip.NewWriter(buffer)
	tarWriter := tar.NewWriter(gzWriter)

	return &tarGzWriter{
		buffer: buffer,
		gz:     gzWriter,
		tar:    tarWriter,
	}
}

func (flufikMd5 *md5Writer) Record(b []byte, name string) error {
	digest := md5.New()

	if _, err := digest.Write(b); err != nil {
		return fmt.Errorf("writing digest failed: %w", err)
	}

	if _, err := fmt.Fprintf(flufikMd5.buffer, "%x %s\n", digest.Sum(flufikMd5.tmp), name); err != nil {
		return fmt.Errorf("storing md5 digest failed: %w", err)
	}
	return nil
}

func (flufikMd5 *md5Writer) MD5Sums() []byte {
	return flufikMd5.buffer.Bytes()
}

func newFlufikMd5() *md5Writer {
	return &md5Writer{
		buffer: new(bytes.Buffer),
		tmp:    make([]byte, 0, md5.Size),
	}
}

func (d *FlufikDeb) internalFilePath(flufikFile *FlufikDebFile) (string, error) {
	if !strings.HasPrefix(flufikFile.Name, "/") {
		return "", fmt.Errorf("input file path is not an absolute path: %s", flufikFile.Name)
	}
	return "." + flufikFile.Name, nil
}

func (d *FlufikDeb) compressFile(flufikFile *FlufikDebFile, data *tarGzWriter, md5sum *md5Writer) error {
	installPath, err := d.internalFilePath(flufikFile)
	if err != nil {
		return err
	}

	head := tar.Header{
		Name:     installPath,
		Size:     int64(len(flufikFile.Body)),
		Mode:     int64(flufikFile.Mode),
		ModTime:  flufikFile.MTime,
		Typeflag: tar.TypeReg,
	}

	if _, err = data.Write(&head, flufikFile.Body); err != nil {
		return fmt.Errorf("compress file failed: %w", err)
	}

	if err = md5sum.Record(flufikFile.Body, installPath[2:]); err != nil {
		return fmt.Errorf("generate md5 information for %s failed: %w", installPath, err)
	}

	if flufikFile.isConfig() {
		if _, err = fmt.Fprintf(d.configFiles, flufikFile.Name); err != nil {
			return fmt.Errorf("generate config file information for %s failed: %w", installPath, err)
		}
	}
	return nil
}

func (d *FlufikDeb) compressDir(flufikFile *FlufikDebFile, data *tarGzWriter) error {
	installPath, err := d.internalFilePath(flufikFile)
	if err != nil {
		return err
	}

	if !strings.HasSuffix(installPath, "/") {
		installPath += "/"
	}

	head := tar.Header{
		Name:     installPath,
		Mode:     int64(flufikFile.Mode),
		ModTime:  flufikFile.MTime,
		Typeflag: tar.TypeDir,
	}

	if err = data.WriteHeader(&head); err != nil {
		return fmt.Errorf("compress dir failed: %w", err)
	}
	return nil
}

func (d *FlufikDeb) compressMetaData(filename string, body []byte, meta *tarGzWriter) error {
	head := tar.Header{
		Name:     filename,
		Size:     int64(len(body)),
		Mode:     0644,
		ModTime:  time.Now(),
		Typeflag: tar.TypeReg,
	}

	if _, err := meta.Write(&head, body); err != nil {
		return fmt.Errorf("compress %s file failed: %w", filename, err)
	}
	return nil
}

func (d *FlufikDeb) compressControl(meta *tarGzWriter) error {
	return d.compressMetaData("control", d.MakeControl(), meta)
}

func (d *FlufikDeb) compressConfigurationFiles(meta *tarGzWriter) error {
	return d.compressMetaData("configuration_files", d.configFiles.Bytes(), meta)
}

func (d *FlufikDeb) compressMD5(meta *tarGzWriter, md5sum *md5Writer) error {
	return d.compressMetaData("md5sums", md5sum.MD5Sums(), meta)
}

func (d *FlufikDeb) compressScripts(meta *tarGzWriter) error {
	if err := d.compressMetaData("preinst", []byte(d.preIn), meta); err != nil {
		return err
	}

	if err := d.compressMetaData("postinst", []byte(d.postIn), meta); err != nil {
		return err
	}

	if err := d.compressMetaData("prerm", []byte(d.preUn), meta); err != nil {
		return err
	}

	if err := d.compressMetaData("postrm", []byte(d.postUn), meta); err != nil {
		return err
	}
	return nil
}

func (d *FlufikDeb) arCompress(w *ar.Writer, filename string, body []byte) error {
	head := ar.Header{
		Name:    filename,
		Size:    int64(len(body)),
		Mode:    0644,
		ModTime: time.Now(),
	}

	if err := w.WriteHeader(&head); err != nil {
		return fmt.Errorf("can not write file header: %w", err)
	}

	_, err := w.Write(body)

	return err
}

func (d *FlufikDeb) Write(w io.Writer) error {

	flufikMeta := newFlufikTgz()
	flufikData := newFlufikTgz()
	flufikMd5Sum := newFlufikMd5()

	for _, f := range d.files {
		if f.isDir() {
			if err := d.compressDir(&f, flufikData); err != nil {
				return err
			}
		} else {
			if err := d.compressFile(&f, flufikData, flufikMd5Sum); err != nil {
				return nil
			}
		}
	}

	if err := d.compressControl(flufikMeta); err != nil {
		return err
	}

	if err := d.compressMD5(flufikMeta, flufikMd5Sum); err != nil {
		return err
	}

	if err := d.compressConfigurationFiles(flufikMeta); err != nil {
		return err
	}

	if err := d.compressScripts(flufikMeta); err != nil {
		return err
	}

	if err := flufikMeta.Close(); err != nil {
		return fmt.Errorf("can't close meta: %w", err)
	}
	if err := flufikData.Close(); err != nil {
		return fmt.Errorf("can't close data: %w", err)
	}

	writer := ar.NewWriter(w)

	if err := writer.WriteGlobalHeader(); err != nil {
		return fmt.Errorf("can not write ar header to deb file: %w", err)
	}

	debianBinary := []byte("2.0\n")

	if err := d.arCompress(writer, "debian-binary", debianBinary); err != nil {
		return fmt.Errorf("can not write ar header to deb file: %w", err)
	}

	if err := d.arCompress(writer, "control.tar.gz", flufikMeta.Bytes()); err != nil {
		return fmt.Errorf("can not add control.tar.gz to deb: %w", err)
	}

	if err := d.arCompress(writer, "data.tar.gz", flufikData.Bytes()); err != nil {
		return fmt.Errorf("can not add data.tar.gz to deb: %w", err)
	}

	var pgpKeyName string
	if d.Signature.PgpName == "" {
		pgpKeyName = "flufik"
	} else {
		pgpKeyName = d.Signature.PgpName
	}

	data := io.MultiReader(bytes.NewReader(debianBinary), bytes.NewReader(flufikMeta.Bytes()),
		bytes.NewReader(flufikData.Bytes()))
	signer := pgp.NewSigner(d.logger, d.debugging)
	sig, err := signer.FlufikDebSigner(data, pgpKeyName)
	if err != nil {
		return fmt.Errorf("signing failure: %w", err)
	}
	sigType := "origin"
	if d.Signature.Type != "" {
		sigType = d.Signature.Type
	}

	if sigType != "origin" && sigType != "maint" && sigType != "archive" {
		return fmt.Errorf("invalid signature type")
	}

	if err = d.arCompress(writer, "_gpg"+sigType, sig); err != nil {
		return fmt.Errorf("something went wrong with writing signed file: %w", err)
	}
	return nil
}

func (d *FlufikDeb) AddFile(flufikFile FlufikDebFile) { d.files = append(d.files, flufikFile) }
func (d *FlufikDeb) AddPreIn(s string)                { d.preIn = s }
func (d *FlufikDeb) AddPostIn(s string)               { d.postIn = s }
func (d *FlufikDeb) AddPreUn(s string)                { d.preUn = s }
func (d *FlufikDeb) AddPostUn(s string)               { d.postUn = s }

//Signature related functions

func (d *FlufikDeb) AddSignatureKey(k string)  { d.Signature.PgpName = k }
func (d *FlufikDeb) AddSignatureType(t string) { d.Signature.Type = t }

//func (d *FlufikDeb) AddSignaturePassPhrase(p string) { d.Signature.PassPhrase = p }

func NewDeb(flufikMeta FlufikDebMetaData, logger *logging.Logger, debugging string) (*FlufikDeb, error) {
	return &FlufikDeb{
		FlufikDebMetaData: flufikMeta,
		configFiles:       bytes.NewBufferString(""),
		logger:            logger,
		debugging:         debugging,
	}, nil
}
