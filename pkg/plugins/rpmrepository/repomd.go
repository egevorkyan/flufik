package rpmrepository

import (
	"encoding/xml"
	"time"
)

const repomdXmlns = "http://linux.duke.edu/metadata/repo"
const repomdXmlnsRpm = "http://linux.duke.edu/metadata/rpm"

type Repomd struct {
	XMLName  xml.Name     `xml:"repomd"`
	Revision int64        `xml:"revision"`
	Data     []RepomdData `xml:"data"`
	Xmlns    string       `xml:"xmlns,attr"`
	XmlnsRpm string       `xml:"xmlns:rpm,attr"`
}

type RepomdData struct {
	Checksum     Checksum `xml:"checksum"`
	OpenChecksum Checksum `xml:"open-checksum"`
	Location     Location `xml:"location"`
	Type         string   `xml:"type,attr"`
	Timestamp    int64    `xml:"timestamp"`
	Size         int      `xml:"size"`
	OpenSize     int      `xml:"open-size"`
}

type RepomdRequirements struct {
	Size     int
	OpenSize int
	Sum      string
	OpenSum  string
}

func (r *RpmRepo) GetRepomd(in map[string]RepomdRequirements) Repomd {
	timestamp := time.Now().Unix()

	repomd := Repomd{
		Revision: timestamp,
		Xmlns:    repomdXmlns,
		XmlnsRpm: repomdXmlnsRpm,
		Data:     []RepomdData{},
	}

	for t, d := range in {
		data := RepomdData{
			Type: t,
			Checksum: Checksum{
				Type:  "sha256",
				Value: d.Sum,
			},
			OpenChecksum: Checksum{
				Type:  "sha256",
				Value: d.OpenSum,
			},
			Location: Location{
				Href: "repodata/" + t + ".xml.gz",
			},
			Timestamp: timestamp,
			Size:      d.Size,
			OpenSize:  d.OpenSize,
		}
		repomd.Data = append(repomd.Data, data)
	}
	return repomd
}
