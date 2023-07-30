package usecase_test

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"

	"video-server/internal/testutil"
	"video-server/internal/util"
	"video-server/module/entity"
	"video-server/module/fixture"
	"video-server/module/param"
)

func TestFileUsecase_CreateFile(t *testing.T) {
	type Request struct {
		ctx      context.Context
		filePath string
	}

	type Response struct {
		result interface{}
		err    error
	}

	testcases := map[string]struct {
		request  Request
		response Response
		mockFn   func(*fixture.MockFileUsecase, context.Context, util.FileReader)
	}{
		"success": {
			request: Request{
				ctx:      context.Background(),
				filePath: "./../../../test/post_1/sample.mp4",
			},
			response: Response{
				result: map[string]interface{}{"ID": 1},
				err:    nil,
			},
			mockFn: func(m *fixture.MockFileUsecase, ctx context.Context, fileReader util.FileReader) {
				m.FileRepository.EXPECT().CreateFile(ctx, &param.CreateFile{
					Name:     fileReader.GetName(),
					Size:     fileReader.GetSize(),
					MimeType: "video/mp4",
				}).Return(&entity.File{ID: 1}, nil)
			},
		},
		"Unsupported type error": {
			request: Request{
				ctx:      context.Background(),
				filePath: "./../../../test/post_4/test.txt",
			},
			response: Response{
				result: nil,
				err:    entity.ErrorFileUnsupported,
			},
			mockFn: func(m *fixture.MockFileUsecase, ctx context.Context, fileReader util.FileReader) {},
		},
		"CreateFile error": {
			request: Request{
				ctx:      context.Background(),
				filePath: "./../../../test/post_1/sample.mp4",
			},
			response: Response{
				result: nil,
				err:    testutil.ErrDB,
			},
			mockFn: func(m *fixture.MockFileUsecase, ctx context.Context, fileReader util.FileReader) {
				m.FileRepository.EXPECT().CreateFile(ctx, &param.CreateFile{
					Name:     fileReader.GetName(),
					Size:     fileReader.GetSize(),
					MimeType: "video/mp4",
				}).Return(nil, testutil.ErrDB)
			},
		},
	}
	for name, tc := range testcases {
		t.Run(name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			ucs, mocks := fixture.NewFileUsecase(ctrl)

			httpRequest := testutil.RequestPayloadCreateFile(tc.request.filePath)
			reqFile, reqFileHeader, _ := httpRequest.FormFile("data")
			fileReader := util.NewFileReader(reqFile, reqFileHeader)
			defer reqFile.Close()
			tc.mockFn(mocks, tc.request.ctx, fileReader)

			result, err := ucs.CreateFile(tc.request.ctx, fileReader)
			testutil.AssertErrorExAc(t, tc.response.err, err)
			testutil.AssertStructExAc(t, tc.response.result, result)
		})
	}
}

func TestFileUsecase_ListFiles(t *testing.T) {
	type Request struct {
		ctx context.Context
	}

	type Response struct {
		result []*entity.File
		err    error
	}

	file := &entity.File{
		ID:       1,
		Name:     "Some Name",
		Size:     100,
		MimeType: "video/mp4",
	}

	testcases := map[string]struct {
		request  Request
		response Response
		mockFn   func(*fixture.MockFileUsecase, Request)
	}{
		"success": {
			request: Request{
				ctx: context.Background(),
			},
			response: Response{
				result: []*entity.File{file},
				err:    nil,
			},
			mockFn: func(m *fixture.MockFileUsecase, req Request) {
				m.FileRepository.EXPECT().ListFiles(req.ctx).
					Return([]*entity.File{file}, nil)
			},
		},
		"ListFiles error": {
			request: Request{
				ctx: context.Background(),
			},
			response: Response{
				result: nil,
				err:    testutil.ErrDB,
			},
			mockFn: func(m *fixture.MockFileUsecase, req Request) {
				m.FileRepository.EXPECT().ListFiles(req.ctx).
					Return(nil, testutil.ErrDB)
			},
		},
	}
	for name, tc := range testcases {
		t.Run(name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			ucs, mocks := fixture.NewFileUsecase(ctrl)
			tc.mockFn(mocks, tc.request)

			result, err := ucs.ListFiles(tc.request.ctx)
			testutil.AssertErrorExAc(t, tc.response.err, err)
			if len(tc.response.result) > 0 {
				testutil.AssertStructExAc(t, tc.response.result[0], result[0])
			}
		})
	}
}

func TestFileUsecase_GetFile(t *testing.T) {
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
		mockFn   func(*fixture.MockFileUsecase, Request)
	}{
		"success": {
			request: Request{
				ctx: context.Background(),
				id:  1,
			},
			response: Response{
				result: map[string]interface{}{"ID": 1},
				err:    nil,
			},
			mockFn: func(m *fixture.MockFileUsecase, req Request) {
				m.FileRepository.EXPECT().GetFile(req.ctx, req.id).
					Return(&entity.File{ID: 1}, nil)
			},
		},
	}
	for name, tc := range testcases {
		t.Run(name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			ucs, mocks := fixture.NewFileUsecase(ctrl)
			tc.mockFn(mocks, tc.request)

			result, err := ucs.GetFile(tc.request.ctx, tc.request.id)
			testutil.AssertErrorExAc(t, tc.response.err, err)
			testutil.AssertStructExAc(t, tc.response.result, result)
		})
	}
}

func TestFileUsecase_UpdateFile(t *testing.T) {
	type Request struct {
		ctx context.Context
		id  int
	}

	type Response struct {
		err error
	}

	file := &entity.File{
		ID:       1,
		Name:     "The Name",
		Size:     100,
		MimeType: "video/mp4",
	}

	testcases := map[string]struct {
		request  Request
		response Response
		mockFn   func(*fixture.MockFileUsecase, Request)
	}{
		"success": {
			request: Request{
				ctx: context.Background(),
				id:  1,
			},
			response: Response{
				err: nil,
			},
			mockFn: func(m *fixture.MockFileUsecase, req Request) {
				m.FileRepository.EXPECT().GetFile(req.ctx, req.id).
					Return(file, nil)
				m.FileRepository.EXPECT().DeleteFile(req.ctx, req.id).
					Return(nil)
			},
		},
		"GetFile error": {
			request: Request{
				ctx: context.Background(),
				id:  1,
			},
			response: Response{
				err: testutil.ErrDB,
			},
			mockFn: func(m *fixture.MockFileUsecase, req Request) {
				m.FileRepository.EXPECT().GetFile(req.ctx, req.id).
					Return(nil, testutil.ErrDB)
			},
		},
	}
	for name, tc := range testcases {
		t.Run(name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			ucs, mocks := fixture.NewFileUsecase(ctrl)
			tc.mockFn(mocks, tc.request)

			err := ucs.DeleteFile(tc.request.ctx, tc.request.id)
			testutil.AssertErrorExAc(t, tc.response.err, err)
		})
	}
}
