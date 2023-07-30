package usecase

//go:generate mockgen -source file.go -destination mock/file.go

import (
	"context"
	"fmt"
	"os"
	"strings"

	"video-server/internal/util"
	"video-server/module/entity"
	repository "video-server/module/internal/repository"
	"video-server/module/param"
)

type FileUsecase interface {
	CreateFile(ctx context.Context, filereader util.FileReader) (*entity.File, error)
	ListFiles(ctx context.Context) ([]*entity.File, error)
	GetFile(ctx context.Context, id int) (*entity.File, error)
	DeleteFile(ctx context.Context, id int) error
}

type fileUsecaseRepository struct {
	file repository.FileRepository
}

type fileUsecase struct {
	repository fileUsecaseRepository
}

func NewFileUsecase(
	fileRepository repository.FileRepository,
) *fileUsecase {
	return &fileUsecase{
		repository: fileUsecaseRepository{
			file: fileRepository,
		},
	}
}

func (u *fileUsecase) CreateFile(ctx context.Context, fileReader util.FileReader) (*entity.File, error) {
	fileMimeType, err := fileReader.GetFileMimeType()
	if err != nil {
		return nil, err
	}

	if strings.Split(fileMimeType, "/")[0] != "video" {
		return nil, entity.ErrorFileUnsupported
	}

	file, err := u.repository.file.CreateFile(ctx, &param.CreateFile{
		Name:     fileReader.GetName(),
		Size:     fileReader.GetSize(),
		MimeType: fileMimeType,
	})
	if err != nil {
		return nil, err
	}

	err = fileReader.Store()
	return file, err
}

func (u *fileUsecase) ListFiles(ctx context.Context) ([]*entity.File, error) {
	return u.repository.file.ListFiles(ctx)
}

func (u *fileUsecase) GetFile(ctx context.Context, id int) (*entity.File, error) {
	return u.repository.file.GetFile(ctx, id)

}

func (u *fileUsecase) DeleteFile(ctx context.Context, id int) error {
	file, err := u.repository.file.GetFile(ctx, id)
	if err != nil {
		return err
	}

	os.Remove(fmt.Sprintf("%s/%s", util.StoragePath, file.Name))
	return u.repository.file.DeleteFile(ctx, id)
}
