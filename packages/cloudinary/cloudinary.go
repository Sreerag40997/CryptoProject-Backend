package cloudinary

import (
	"context"
	"log"
	"os"

	"github.com/cloudinary/cloudinary-go/v2"
)

var CLD *cloudinary.Cloudinary
var Ctx = context.Background()

func InitCloudinary() {
	var err error

	CLD, err = cloudinary.NewFromParams(
		os.Getenv("CLOUDINARY_CLOUD_NAME"),
		os.Getenv("CLOUDINARY_API_KEY"),
		os.Getenv("CLOUDINARY_API_SECRET"),
	)

	if err != nil {
		log.Fatal("Cloudinary init failed:", err)
	}

	log.Println("Cloudinary connected ")
}