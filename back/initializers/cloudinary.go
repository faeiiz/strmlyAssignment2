package initializers

import (
	"log"
	"os"

	"github.com/cloudinary/cloudinary-go/v2"
)

var Cloud *cloudinary.Cloudinary

func ConnectCloudinary() {
	cld, err := cloudinary.NewFromURL(os.Getenv("CLOUDINARY_URL"))
	if err != nil {
		log.Fatal("Failed to connect to Cloudinary: ", err)
	}
	Cloud = cld
	log.Println("Connected to Cloudinary!")
}
