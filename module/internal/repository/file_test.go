package repository_test

import (
	"context"
	"database/sql/driver"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"

	"video-server/internal/testutil"
	"video-server/module/entity"
	"video-server/module/fixture"
	"video-server/module/param"
)

func TestFileRepository_CreateFile(t *testing.T) {
	query := "INSERT INTO `files` (`name`,`size`,`mime_type`,`created_at`) VALUES (?,?,?,?)"

	type Request struct {
		ctx    context.Context
		params *param.CreateFile
	}

	type Response struct {
		result interface{}
		err    error
	}

	testcases := map[string]struct {
		request  Request
		response Response
		mockFn   func(*fixture.MockFileRepository, Request, Response)
	}{
		"success": {
			request: Request{
				ctx: context.Background(),
				params: &param.CreateFile{
					MimeType: "video/mp4",
					Name:     "Some Name",
					Size:     100,
				},
			},
			response: Response{
				result: map[string]interface{}{"ID": 1},
				err:    nil,
			},
			mockFn: func(m *fixture.MockFileRepository, req Request, res Response) {
				m.SQLMock.ExpectBegin()
				m.SQLMock.ExpectExec(regexp.QuoteMeta(query)).
					WithArgs("Some Name", 100, "video/mp4", testutil.AnyTime{}).
					WillReturnResult(sqlmock.NewResult(1, 1))
				m.SQLMock.ExpectCommit()
			},
		},
		"db error file exists": {
			request: Request{
				ctx: context.Background(),
				params: &param.CreateFile{
					MimeType: "video/mp4",
					Name:     "Some Name",
					Size:     100,
				},
			},
			response: Response{
				result: nil,
				err:    entity.ErrorFileExists,
			},
			mockFn: func(m *fixture.MockFileRepository, req Request, res Response) {
				m.SQLMock.ExpectBegin()
				m.SQLMock.ExpectExec(regexp.QuoteMeta(query)).
					WithArgs("Some Name", 100, "video/mp4", testutil.AnyTime{}).
					WillReturnResult(sqlmock.NewResult(0, 0)).
					WillReturnError(&mysql.MySQLError{Number: 1062})
				m.SQLMock.ExpectRollback()
			},
		},
		"db error": {
			request: Request{
				ctx: context.Background(),
				params: &param.CreateFile{
					MimeType: "video/mp4",
					Name:     "Some Name",
					Size:     100,
				},
			},
			response: Response{
				result: nil,
				err:    testutil.ErrDB,
			},
			mockFn: func(m *fixture.MockFileRepository, req Request, res Response) {
				m.SQLMock.ExpectBegin()
				m.SQLMock.ExpectExec(regexp.QuoteMeta(query)).
					WithArgs("Some Name", 100, "video/mp4", testutil.AnyTime{}).
					WillReturnResult(sqlmock.NewResult(0, 0)).
					WillReturnError(testutil.ErrDB)
				m.SQLMock.ExpectRollback()
			},
		},
	}
	for name, tc := range testcases {
		t.Run(name, func(t *testing.T) {
			repo, mocks := fixture.NewFileRepository()
			tc.mockFn(mocks, tc.request, tc.response)
			result, err := repo.CreateFile(context.Background(), tc.request.params)
			testutil.AssertErrorExAc(t, tc.response.err, err)
			if tc.response.result != nil {
				assert.NotNil(t, result)
			}
		})
	}
}

func TestFileRepository_ListFiles(t *testing.T) {
	rowColumns := []string{"id", "name", "size", "mime_type", "created_at"}
	rowValues := []driver.Value{1, "Some Name", 100, "video/mp4", testutil.CreatedAt}
	query := "SELECT `id`,`name`,`size`,`mime_type`,`created_at` FROM `files`"

	type Request struct {
		ctx context.Context
	}

	type Response struct {
		result []*entity.File
		err    error
	}

	testcases := map[string]struct {
		request  Request
		response Response
		mockFn   func(*fixture.MockFileRepository, Request, Response)
	}{
		"success": {
			request: Request{
				ctx: context.Background(),
			},
			response: Response{
				result: []*entity.File{
					{
						ID:        1,
						Name:      "Some Name",
						Size:      100,
						MimeType:  "video/mp4",
						CreatedAt: testutil.CreatedAt,
					},
				},
				err: nil,
			},
			mockFn: func(m *fixture.MockFileRepository, req Request, res Response) {

				rows := m.SQLMock.NewRows(rowColumns)
				rows.AddRow(rowValues...)
				m.SQLMock.ExpectQuery(regexp.QuoteMeta(query)).WillReturnRows(rows)

			},
		},
		"db error": {
			request: Request{
				ctx: context.Background(),
			},
			response: Response{
				result: nil,
				err:    testutil.ErrDB,
			},
			mockFn: func(m *fixture.MockFileRepository, req Request, res Response) {
				m.SQLMock.ExpectQuery(regexp.QuoteMeta(query)).
					WillReturnError(testutil.ErrDB)
			},
		},
	}
	for name, tc := range testcases {
		t.Run(name, func(t *testing.T) {
			repo, mocks := fixture.NewFileRepository()
			tc.mockFn(mocks, tc.request, tc.response)
			result, err := repo.ListFiles(context.Background())
			testutil.AssertErrorExAc(t, tc.response.err, err)
			if len(tc.response.result) > 0 {
				testutil.AssertStructExAc(t, tc.response.result[0], result[0])
			}
		})
	}
}

func TestFileRepository_GetFile(t *testing.T) {
	rowColumns := []string{"id", "name", "size", "mime_type", "created_at"}
	rowValues := []driver.Value{123, "Some Name", 100, "video/mp4", testutil.CreatedAt}
	query := "SELECT `id`,`name`,`size`,`mime_type`,`created_at` FROM `files`"

	type Request struct {
		ctx context.Context
		id  int
	}

	type Response struct {
		result interface{}
		err    error
	}

	testcases := map[string]struct {
		request  Request
		response Response
		mockFn   func(*fixture.MockFileRepository, Request, Response)
	}{
		"success": {
			request: Request{
				ctx: context.Background(),
				id:  123,
			},
			response: Response{
				result: map[string]interface{}{"ID": 123},
				err:    nil,
			},
			mockFn: func(m *fixture.MockFileRepository, req Request, res Response) {
				rows := m.SQLMock.NewRows(rowColumns)
				rows.AddRow(rowValues...)
				m.SQLMock.ExpectQuery(regexp.QuoteMeta(query)).WillReturnRows(rows)
			},
		},
		"db error not found": {
			request: Request{
				ctx: context.Background(),
				id:  123,
			},
			response: Response{
				result: nil,
				err:    entity.ErrorFileNotFound,
			},
			mockFn: func(m *fixture.MockFileRepository, req Request, res Response) {
				rows := m.SQLMock.NewRows(rowColumns)
				rows.AddRow(rowValues...)
				m.SQLMock.ExpectQuery(regexp.QuoteMeta(query)).WillReturnError(gorm.ErrRecordNotFound)
			},
		},
		"db error others": {
			request: Request{
				ctx: context.Background(),
				id:  123,
			},
			response: Response{
				result: nil,
				err:    testutil.ErrDB,
			},
			mockFn: func(m *fixture.MockFileRepository, req Request, res Response) {
				rows := m.SQLMock.NewRows(rowColumns)
				rows.AddRow(rowValues...)
				m.SQLMock.ExpectQuery(regexp.QuoteMeta(query)).WillReturnError(testutil.ErrDB)
			},
		},
	}
	for name, tc := range testcases {
		t.Run(name, func(t *testing.T) {
			repo, mocks := fixture.NewFileRepository()
			tc.mockFn(mocks, tc.request, tc.response)
			result, err := repo.GetFile(context.Background(), tc.request.id)
			testutil.AssertErrorExAc(t, tc.response.err, err)
			testutil.AssertStructExAc(t, tc.response.result, result)
		})
	}
}

func TestFileRepository_UpdateFile(t *testing.T) {
	query := "DELETE FROM `files` WHERE `files`.`id` = ?"

	type Request struct {
		ctx context.Context
		id  int
	}

	type Response struct {
		err error
	}

	testcases := map[string]struct {
		request  Request
		response Response
		mockFn   func(*fixture.MockFileRepository, Request, Response)
	}{
		"success": {
			request: Request{
				ctx: context.Background(),
				id:  123,
			},
			response: Response{
				err: nil,
			},
			mockFn: func(m *fixture.MockFileRepository, req Request, res Response) {
				m.SQLMock.ExpectBegin()
				m.SQLMock.ExpectExec(regexp.QuoteMeta(query)).
					WithArgs(123).
					WillReturnResult(sqlmock.NewResult(1, 1))
				m.SQLMock.ExpectCommit()
			},
		},
	}
	for name, tc := range testcases {
		t.Run(name, func(t *testing.T) {
			repo, mocks := fixture.NewFileRepository()
			tc.mockFn(mocks, tc.request, tc.response)
			err := repo.DeleteFile(context.Background(), tc.request.id)
			testutil.AssertErrorExAc(t, tc.response.err, err)
		})
	}
}
