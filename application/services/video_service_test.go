package services_test

import (
	"encoder/application/repositories"
	"encoder/application/services"
	"encoder/domain"
	"encoder/framework/database"
	"log"
	"testing"
	"time"

	"github.com/joho/godotenv"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/require"
)

func init() {
	err := godotenv.Load("../../.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func prepare() (*domain.Video, *repositories.VideoRepositoryDb) {
	db := database.NewDbTest()
	// defer db.Close()

	video := domain.NewVideo()
	video.ID = uuid.NewV4().String()
	video.FilePath = "pudim-amassado.mp4"
	video.CreatedAt = time.Now()

	repo := repositories.NewVideoRepository(db)

	return video, repo
}

func TestVideoServiceDownload(t *testing.T) {
	video, repo := prepare()

	service := services.NewVideoService()
	service.Video = video
	service.VideoRepository = repo

	err := service.Download("encoder-full-cicly")
	require.Nil(t, err)

	err = service.Fragment()
	require.Nil(t, err)

	err = service.Encode()
	require.Nil(t, err)

	err = service.Finish()
	require.Nil(t, err)
}
