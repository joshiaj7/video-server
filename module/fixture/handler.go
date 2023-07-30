package fixture

import (
	"github.com/golang/mock/gomock"

	"video-server/module/internal/handler"
	mock_usecase "video-server/module/internal/usecase/mock"
)

type MockFileHandler struct {
	// Usecase
	FileUsecase *mock_usecase.MockFileUsecase
}

func NewFileHandler(
	ctrl *gomock.Controller,
) (*handler.FileHandler, *MockFileHandler) {
	mocks := &MockFileHandler{
		FileUsecase: mock_usecase.NewMockFileUsecase(ctrl),
	}

	svc := handler.NewFileHandler(
		mocks.FileUsecase,
	)

	return svc, mocks
}
