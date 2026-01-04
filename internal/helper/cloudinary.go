package helper

import (
	"context"
	"errors"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

func DeleteImageFromCloud(publicID string) error {
	if publicID == "" {
		return errors.New("public_id is empty")
	}

	cld, err := cloudinary.New()
	if err != nil {
		return err
	}

	_, err = cld.Upload.Destroy(
		context.Background(),
		uploader.DestroyParams{
			PublicID:     publicID,
			ResourceType: "image",
		},
	)

	return err
}
