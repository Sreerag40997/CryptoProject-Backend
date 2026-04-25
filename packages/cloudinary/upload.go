package cloudinary

import (
	"mime/multipart"

	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

func UploadFile(file multipart.File, filename string) (string, error) {

	resp, err := CLD.Upload.Upload(Ctx, file, uploader.UploadParams{
		PublicID: "kyc/" + filename,
	})

	if err != nil {
		return "", err
	}

	return resp.SecureURL, nil
}

