package fixture

import (
	"github.com/golang/mock/gomock"

	mock_repository "video-server/module/internal/repository/mock"
	"video-server/module/internal/usecase"
)

type MockFileUsecase struct {
	// Repository
	FileRepository *mock_repository.MockFileRepository
}

func NewFileUsecase(ctrl *gomock.Controller) (usecase.FileUsecase, *MockFileUsecase) {
	mocks := &MockFileUsecase{
		FileRepository: mock_repository.NewMockFileRepository(ctrl),
	}
	ucs := usecase.NewFileUsecase(mocks.FileRepository)
	return ucs, mocks
}
