package storage

import "time"

type UploadInput struct {
	Directory   string
	ContentType string
	Filename    string
	Body        []byte
}

type UploadOutput struct {
	Key string
	URL string
}

type PresignedUploadInput struct {
	Key         string
	ContentType string
	Expired     time.Duration
}

type PresignedUploadOutput struct {
	URL    string
	Method string
	Header map[string]string
}
