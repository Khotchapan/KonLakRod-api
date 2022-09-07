package test

import (
	"mime/multipart"
)

type Book struct {
	Name string `json:"name"`
}

type UploadForm struct {
	Path string                `form:"path"`
	Mime string                `form:"mime"`
	File *multipart.FileHeader `form:"file"`
}

type GetOneGoogleCloudForm struct {
	ID *string `param:"id"`
}
