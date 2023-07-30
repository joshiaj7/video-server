package fixture

import (
	sqlmock "github.com/DATA-DOG/go-sqlmock"

	"video-server/internal/testutil"
	repository "video-server/module/internal/repository"
)

type MockFileRepository struct {
	SQLMock sqlmock.Sqlmock
}

func NewFileRepository() (repository.FileRepository, *MockFileRepository) {
	db, sqlMock := testutil.NewDatabase()
	mocks := &MockFileRepository{SQLMock: sqlMock}
	repo := repository.NewFileRepository(db)
	return repo, mocks
}
