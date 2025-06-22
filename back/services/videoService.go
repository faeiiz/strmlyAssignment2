package services

import (
	"back/initializers"
	"back/repository"
	"context"
	"mime/multipart"
	"time"

	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

type VideoService interface {
	UploadVideo(title, description, uploaderID string, file multipart.File) error
	GetAllVideos() ([]repository.Video, error)
	GetVideosPaginated(page, limit int) ([]repository.Video, error)
}

type videoService struct {
	repo repository.VideoRepository
}

func NewVideoService(r repository.VideoRepository) VideoService {
	return &videoService{repo: r}
}

func (v *videoService) UploadVideo(title, description, uploaderID string, file multipart.File) error {
	ctx := context.Background()

	uploadResult, err := initializers.Cloud.Upload.Upload(ctx, file, uploader.UploadParams{
		Folder: "strmly_videos",
	})

	if err != nil {
		return err
	}

	video := repository.Video{
		Title:       title,
		Description: description,
		URL:         uploadResult.SecureURL,
		UploaderID:  uploaderID,
		UploadDate:  time.Now(),
	}
	return v.repo.Create(video)
}

func (v *videoService) GetAllVideos() ([]repository.Video, error) {
	return v.repo.GetAll()
}

func (v *videoService) GetVideosPaginated(page, limit int) ([]repository.Video, error) {
	return v.repo.GetPaginated(page, limit)
}
