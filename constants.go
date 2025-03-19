package dtorm

import "fmt"

const (
	majorVersion   = 0
	minorVersion   = 1
	releaseVersion = 0
)

func Version() string {
	return fmt.Sprintf("DTORM v%d.%d.%d ©2023 I Have a Hat", majorVersion, minorVersion, releaseVersion)
}
