package config

import (
	"gorm.io/gorm"

	repository "video-server/module/internal/repository"
)

type Repository struct {
	FileRepository repository.FileRepository
}

func RegisterRepository(db *gorm.DB) *Repository {
	fileRepo := repository.NewFileRepository(db)

	return &Repository{
		FileRepository: fileRepo,
	}
}
