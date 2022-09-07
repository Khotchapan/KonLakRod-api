package test

import (
	"mime/multipart"
)

type UploadForm struct {
	Path string                `form:"path"`
	Mime string                `form:"mime"`
	File *multipart.FileHeader `form:"file"`
}

type GetOneGoogleCloudBooksForm struct {
	ID   *string `param:"id"`
	Name *string `param:"name"`
}
