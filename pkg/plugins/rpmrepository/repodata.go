package rpmrepository

import "github.com/jfrog/go-rpm"

type RepoData map[string][]byte

type PackageInfos map[string]PackageInfo

type PackageInfo struct {
	Path string
	rpm.PackageFile
}

func (r *RpmRepo) CreateRepoData(packages PackageInfos) (RepoData, error) {
	rd := RepoData{}
	req := make(map[string]RepomdRequirements)

	pdata := r.GetPrimary(packages)
	pstr, pstrsize, pstrsum, err := r.GetXML(pdata)
	if err != nil {
		return rd, err
	}
	pzip, pzipsize, pzipsum := r.GetZip(pstr)
	rd["primary.xml"] = pstr
	rd["primary.xml.gz"] = pzip
	req["primary"] = RepomdRequirements{
		OpenSum:  pstrsum,
		OpenSize: pstrsize,
		Sum:      pzipsum,
		Size:     pzipsize,
	}

	fdata := r.GetFilelists(packages)
	fstr, fstrsize, fstrsum, err := r.GetXML(fdata)
	if err != nil {
		return rd, err
	}
	fzip, fzipsize, fzipsum := r.GetZip(fstr)
	rd["filelists.xml"] = fstr
	rd["filelists.xml.gz"] = fzip
	req["filelists"] = RepomdRequirements{
		OpenSum:  fstrsum,
		OpenSize: fstrsize,
		Sum:      fzipsum,
		Size:     fzipsize,
	}

	odata := r.GetOther(packages)
	ostr, ostrsize, ostrsum, err := r.GetXML(odata)
	if err != nil {
		return rd, err
	}
	ozip, ozipsize, ozipsum := r.GetZip(ostr)
	rd["other.xml"] = ostr
	rd["other.xml.gz"] = ozip
	req["other"] = RepomdRequirements{
		OpenSum:  ostrsum,
		OpenSize: ostrsize,
		Sum:      ozipsum,
		Size:     ozipsize,
	}

	rdata := r.GetRepomd(req)
	rstr, _, _, err := r.GetXML(rdata)
	if err != nil {
		return rd, err
	}
	rzip, _, _ := r.GetZip(rstr)
	rd["repomd.xml"] = rstr
	rd["repomd.xml.gz"] = rzip

	return rd, nil
}
