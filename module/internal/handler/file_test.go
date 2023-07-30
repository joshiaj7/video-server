package handler_test

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/julienschmidt/httprouter"
	"github.com/stretchr/testify/assert"

	"video-server/internal/testutil"
	"video-server/internal/util"
	"video-server/module/entity"
	"video-server/module/fixture"
	"video-server/module/response"
)

func TestFileHandler_CreateFile(t *testing.T) {
	type Request struct {
		req *http.Request
	}

	type Response struct {
		statusCode int
		err        error
	}

	testcases := map[string]struct {
		request  Request
		response Response
		mockFn   func(*fixture.MockFileHandler, Request)
	}{
		"success": {
			response: Response{
				statusCode: 201,
				err:        nil,
			},
			mockFn: func(m *fixture.MockFileHandler, req Request) {
				reqFile, reqFileHeader, _ := req.req.FormFile("data")
				defer reqFile.Close()

				m.FileUsecase.EXPECT().CreateFile(req.req.Context(), util.NewFileReader(reqFile, reqFileHeader)).
					Return(&entity.File{ID: 1}, nil)
			},
		},
		"File exists": {
			response: Response{
				statusCode: 409,
				err:        entity.ErrorFileExists,
			},
			mockFn: func(m *fixture.MockFileHandler, req Request) {
				reqFile, reqFileHeader, _ := req.req.FormFile("data")
				defer reqFile.Close()

				m.FileUsecase.EXPECT().CreateFile(req.req.Context(), util.NewFileReader(reqFile, reqFileHeader)).
					Return(nil, entity.ErrorFileExists)
			},
		},
		"Unsupported File Type": {
			response: Response{
				statusCode: 415,
				err:        entity.ErrorFileUnsupported,
			},
			mockFn: func(m *fixture.MockFileHandler, req Request) {
				reqFile, reqFileHeader, _ := req.req.FormFile("data")
				defer reqFile.Close()

				m.FileUsecase.EXPECT().CreateFile(req.req.Context(), util.NewFileReader(reqFile, reqFileHeader)).
					Return(nil, entity.ErrorFileUnsupported)
			},
		},
	}

	for name, tc := range testcases {
		t.Run(name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			handler, mocks := fixture.NewFileHandler(ctrl)

			req := testutil.RequestPayloadCreateFile("./../../../test/post_1/sample.mp4")
			tc.mockFn(mocks, Request{req: req})

			responseWriter := httptest.NewRecorder()
			handler.CreateFile(responseWriter, req, nil)
			assert.Equal(t, responseWriter.Code, tc.response.statusCode)
		})
	}
}

func TestFileHandler_ListFiles(t *testing.T) {
	type Request struct {
		req *http.Request
		ctx context.Context
	}

	type Response struct {
		body interface{}
	}

	createdAt := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
	file := &entity.File{
		ID:        1,
		Name:      "Some Name",
		MimeType:  "video/mp4",
		Size:      100,
		CreatedAt: createdAt,
	}

	resp := &response.File{
		ID:        "1",
		Name:      "Some Name",
		Size:      100,
		CreatedAt: createdAt,
	}

	testcases := map[string]struct {
		request  Request
		response Response
		mockFn   func(*fixture.MockFileHandler, Request)
	}{
		"success": {
			request: Request{
				ctx: context.Background(),
			},
			response: Response{
				body: []*response.File{resp},
			},
			mockFn: func(m *fixture.MockFileHandler, r Request) {
				m.FileUsecase.EXPECT().ListFiles(r.req.Context()).
					Return([]*entity.File{file}, nil)
			},
		},
		"ListFiles error": {
			request: Request{
				ctx: context.Background(),
			},
			response: Response{
				body: map[string]interface{}{"message": "DB Error"},
			},
			mockFn: func(m *fixture.MockFileHandler, r Request) {
				m.FileUsecase.EXPECT().ListFiles(r.req.Context()).
					Return(nil, testutil.ErrDB)
			},
		},
	}
	for name, tc := range testcases {
		t.Run(name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			req, _ := http.NewRequest(http.MethodGet, "http://example.com/", nil)
			tc.request.req = req
			handler, mocks := fixture.NewFileHandler(ctrl)
			tc.mockFn(mocks, tc.request)

			responseWriter := httptest.NewRecorder()
			handler.ListFiles(responseWriter, req, nil)
			resultBody, _ := io.ReadAll(responseWriter.Body)
			body, _ := json.Marshal(tc.response.body)
			assert.Equal(t, string(body)+"\n", string(resultBody))
		})
	}
}

func TestFileHandler_GetFile(t *testing.T) {
	type Request struct {
		req    *http.Request
		params httprouter.Params
	}

	createdAt := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
	file := &entity.File{
		ID:        1,
		Name:      "test/sample.mp4",
		MimeType:  "video/mp4",
		Size:      100,
		CreatedAt: createdAt,
	}

	testcases := map[string]struct {
		request Request
		mockFn  func(*fixture.MockFileHandler, Request)
	}{
		"success": {
			request: Request{
				params: httprouter.Params{httprouter.Param{
					Key:   "fileid",
					Value: "123",
				}},
			},
			mockFn: func(m *fixture.MockFileHandler, req Request) {
				m.FileUsecase.EXPECT().GetFile(req.req.Context(), 123).
					Return(file, nil)
			},
		},
		"no request param": {
			request: Request{
				params: httprouter.Params{},
			},
			mockFn: func(m *fixture.MockFileHandler, req Request) {},
		},
		"GetFile error": {
			request: Request{
				params: httprouter.Params{httprouter.Param{
					Key:   "fileid",
					Value: "123",
				}},
			},
			mockFn: func(m *fixture.MockFileHandler, req Request) {
				m.FileUsecase.EXPECT().GetFile(req.req.Context(), 123).
					Return(nil, entity.ErrorFileNotFound)
			},
		},
	}
	for name, tc := range testcases {
		t.Run(name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			req, _ := http.NewRequest(http.MethodGet, "http://example.com/", nil)
			handler, mocks := fixture.NewFileHandler(ctrl)
			tc.request.req = req
			tc.mockFn(mocks, tc.request)

			responseWriter := httptest.NewRecorder()
			handler.GetFile(responseWriter, req, tc.request.params)
		})
	}
}

func TestFileHandler_DeleteFile(t *testing.T) {
	type Request struct {
		req    *http.Request
		params httprouter.Params
	}

	type Response struct {
		statusCode int
	}

	testcases := map[string]struct {
		request  Request
		response Response
		mockFn   func(*fixture.MockFileHandler, Request)
	}{
		"success": {
			request: Request{
				params: httprouter.Params{httprouter.Param{
					Key:   "fileid",
					Value: "123",
				}},
			},
			response: Response{
				statusCode: 204,
			},
			mockFn: func(m *fixture.MockFileHandler, req Request) {
				m.FileUsecase.EXPECT().DeleteFile(req.req.Context(), 123).Return(nil)
			},
		},
		"no param": {
			request: Request{
				params: httprouter.Params{},
			},
			response: Response{
				statusCode: 404,
			},
			mockFn: func(m *fixture.MockFileHandler, req Request) {},
		},
		"DeleteFile error": {
			request: Request{
				params: httprouter.Params{httprouter.Param{
					Key:   "fileid",
					Value: "123",
				}},
			},
			response: Response{
				statusCode: 404,
			},
			mockFn: func(m *fixture.MockFileHandler, req Request) {
				m.FileUsecase.EXPECT().DeleteFile(req.req.Context(), 123).Return(entity.ErrorFileNotFound)
			},
		},
	}
	for name, tc := range testcases {
		t.Run(name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			handler, mocks := fixture.NewFileHandler(ctrl)
			req, _ := http.NewRequest(http.MethodPost, "http://example.com/", nil)
			tc.request.req = req
			tc.mockFn(mocks, tc.request)

			responseWriter := httptest.NewRecorder()
			handler.DeleteFile(responseWriter, req, tc.request.params)
			assert.Equal(t, responseWriter.Code, tc.response.statusCode)
		})
	}
}
