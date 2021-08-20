package flufikinfo

import "time"

type FlufikPackageMeta struct {
	Name        string `yaml:"name,omitempty"`
	Version     string `yaml:"version,omitempty"`
	Release     string `yaml:"release,omitempty"`
	Arch        string `yaml:"arch,omitempty"`
	Summary     string `yaml:"summary"`
	Description string `yaml:"description"`
	OS          string `yaml:"os,omitempty"`
	Vendor      string `yaml:"vendor,omitempty"`
	URL         string `yaml:"url,omitempty"`
	License     string `yaml:"license,omitempty"`
	Maintainer  string `yaml:"maintainer,omitempty"`
	buildTime   time.Time
}

func (flufiPkgkMeta *FlufikPackageMeta) BuildTime() time.Time {
	return flufiPkgkMeta.buildTime.UTC()
}

func (flufiPkgkMeta *FlufikPackageMeta) UpdateBuildTime(buildTime time.Time) {
	flufiPkgkMeta.buildTime = buildTime
}
