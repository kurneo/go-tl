package filesystem

import "strings"

type PublicDisk struct {
	DiskContract
	url string
}

func (p PublicDisk) Url(path string) string {
	if p.url == "" {
		return ""
	}
	return strings.TrimRight(p.url, "\\/") + "/" + path
}

func NewPublicDisk(driver DriverContract, url string) PublicDiskContract {
	return PublicDisk{
		disk{
			driver: driver,
		},
		url,
	}
}
