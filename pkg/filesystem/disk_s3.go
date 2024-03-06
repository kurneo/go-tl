package filesystem

import "strings"

type S3Disk struct {
	DiskContract
	url string
}

func (s S3Disk) Url(path string) string {
	if s.url == "" {
		return ""
	}
	return strings.TrimRight(s.url, "\\/") + "/" + path
}

func NewS3Disk(driver DriverContract, url string) PublicDiskContract {
	return PublicDisk{
		disk{
			driver: driver,
		},
		url,
	}
}
