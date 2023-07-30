// Code generated by MockGen. DO NOT EDIT.
// Source: file.go

// Package mock_repository is a generated GoMock package.
package mock_repository

import (
	context "context"
	reflect "reflect"

	entity "video-server/module/entity"
	param "video-server/module/param"
	gomock "github.com/golang/mock/gomock"
)

// MockFileRepository is a mock of FileRepository interface.
type MockFileRepository struct {
	ctrl     *gomock.Controller
	recorder *MockFileRepositoryMockRecorder
}

// MockFileRepositoryMockRecorder is the mock recorder for MockFileRepository.
type MockFileRepositoryMockRecorder struct {
	mock *MockFileRepository
}

// NewMockFileRepository creates a new mock instance.
func NewMockFileRepository(ctrl *gomock.Controller) *MockFileRepository {
	mock := &MockFileRepository{ctrl: ctrl}
	mock.recorder = &MockFileRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockFileRepository) EXPECT() *MockFileRepositoryMockRecorder {
	return m.recorder
}

// CreateFile mocks base method.
func (m *MockFileRepository) CreateFile(ctx context.Context, params *param.CreateFile) (*entity.File, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateFile", ctx, params)
	ret0, _ := ret[0].(*entity.File)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateFile indicates an expected call of CreateFile.
func (mr *MockFileRepositoryMockRecorder) CreateFile(ctx, params interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateFile", reflect.TypeOf((*MockFileRepository)(nil).CreateFile), ctx, params)
}

// DeleteFile mocks base method.
func (m *MockFileRepository) DeleteFile(ctx context.Context, id int) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteFile", ctx, id)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteFile indicates an expected call of DeleteFile.
func (mr *MockFileRepositoryMockRecorder) DeleteFile(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteFile", reflect.TypeOf((*MockFileRepository)(nil).DeleteFile), ctx, id)
}

// GetFile mocks base method.
func (m *MockFileRepository) GetFile(ctx context.Context, id int) (*entity.File, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetFile", ctx, id)
	ret0, _ := ret[0].(*entity.File)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetFile indicates an expected call of GetFile.
func (mr *MockFileRepositoryMockRecorder) GetFile(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetFile", reflect.TypeOf((*MockFileRepository)(nil).GetFile), ctx, id)
}

// ListFiles mocks base method.
func (m *MockFileRepository) ListFiles(ctx context.Context) ([]*entity.File, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListFiles", ctx)
	ret0, _ := ret[0].([]*entity.File)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListFiles indicates an expected call of ListFiles.
func (mr *MockFileRepositoryMockRecorder) ListFiles(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListFiles", reflect.TypeOf((*MockFileRepository)(nil).ListFiles), ctx)
}
