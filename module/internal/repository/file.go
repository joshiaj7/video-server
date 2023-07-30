package repository

//go:generate mockgen -source file.go -destination mock/file.go

import (
	"context"
	"errors"
	"time"

	"gorm.io/gorm"

	"video-server/module/entity"
	"video-server/module/param"

	"github.com/go-sql-driver/mysql"
)

var (
	FileColumnsInsert = []string{
		"name",
		"size",
		"mime_type",
		"created_at",
	}
	FileColumns = append([]string{"id"}, FileColumnsInsert...)
)

type FileRepository interface {
	CreateFile(ctx context.Context, params *param.CreateFile) (*entity.File, error)
	ListFiles(ctx context.Context) ([]*entity.File, error)
	GetFile(ctx context.Context, id int) (*entity.File, error)
	DeleteFile(ctx context.Context, id int) error
}

type fileRepository struct {
	database *gorm.DB
}

func NewFileRepository(database *gorm.DB) *fileRepository {
	return &fileRepository{
		database: database,
	}
}

func (r *fileRepository) CreateFile(ctx context.Context, params *param.CreateFile) (*entity.File, error) {
	timeNow := time.Now()
	file := &entity.File{
		Name:      params.Name,
		Size:      params.Size,
		MimeType:  params.MimeType,
		CreatedAt: timeNow,
	}
	err := r.database.Select(FileColumnsInsert).Create(file).Error
	if err != nil {
		var mysqlErr *mysql.MySQLError
		if errors.As(err, &mysqlErr) && mysqlErr.Number == 1062 {
			return nil, entity.ErrorFileExists
		}
		return nil, err
	}

	return file, err
}

func (r *fileRepository) ListFiles(ctx context.Context) ([]*entity.File, error) {
	files := []*entity.File{}
	err := r.database.Select(FileColumns).Find(&files).Error

	return files, err
}

func (r *fileRepository) GetFile(ctx context.Context, id int) (*entity.File, error) {
	file := &entity.File{ID: id}
	err := r.database.Select(FileColumns).First(&file).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, entity.ErrorFileNotFound
		}
		return nil, err
	}

	return file, err
}

func (r *fileRepository) DeleteFile(ctx context.Context, id int) error {
	file := &entity.File{ID: id}

	return r.database.Delete(&file).Error
}
