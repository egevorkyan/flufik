package core

import (
	"fmt"
	"time"
)

const majorVersion = "0.1"

func VersionBuilder() string {
	t := time.Now()
	y, m, d := t.Date()
	h := t.Hour()
	min := t.Minute()
	version := fmt.Sprintf("%s-%v%v%v%v%v", majorVersion, d, int(m), y, h, min)
	return version
}
