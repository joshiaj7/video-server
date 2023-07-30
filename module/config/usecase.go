package config

import "video-server/module/internal/usecase"

type Usecase struct {
	FileUsecase usecase.FileUsecase
}

func RegisterUsecase(repository *Repository) *Usecase {
	fileUcs := usecase.NewFileUsecase(repository.FileRepository)

	return &Usecase{
		FileUsecase: fileUcs,
	}
}
