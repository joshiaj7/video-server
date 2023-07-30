package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"

	"video-server/internal/util"
	"video-server/module/entity"
	"video-server/module/internal/usecase"
	"video-server/module/response"
)

type FileRepository interface {
	CreateFile(w http.ResponseWriter, r *http.Request, params httprouter.Params)
	ListFiles(w http.ResponseWriter, r *http.Request, params httprouter.Params)
	GetFile(w http.ResponseWriter, r *http.Request, params httprouter.Params)
	DeleteFile(w http.ResponseWriter, r *http.Request, params httprouter.Params)
}

type FileHandler struct {
	usecase usecase.FileUsecase
}

func NewFileHandler(uc usecase.FileUsecase) *FileHandler {
	return &FileHandler{
		usecase: uc,
	}
}

func (h *FileHandler) Register(router *httprouter.Router) {
	router.POST("/v1/files", h.CreateFile)
	router.GET("/v1/files", h.ListFiles)
	router.GET("/v1/files/:fileid", h.GetFile)
	router.DELETE("/v1/files/:fileid", h.DeleteFile)
}

func (h *FileHandler) CreateFile(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	reqFile, reqFileHeader, err := r.FormFile("data")
	if err != nil {
		BuildErrorResponse(w, entity.ErrorBadRequest)
		return
	}
	defer reqFile.Close()

	result, err := h.usecase.CreateFile(r.Context(), util.NewFileReader(reqFile, reqFileHeader))
	if err != nil {
		BuildErrorResponse(w, err)
		return
	}

	w.Header().Set("Location", fmt.Sprintf("%s/%d", util.StoragePath, result.ID))
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusCreated)
}

func (h *FileHandler) ListFiles(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	files, err := h.usecase.ListFiles(r.Context())
	if err != nil {
		BuildErrorResponse(w, err)
		return
	}

	result := []*response.File{}
	for _, obj := range files {
		result = append(result, fileEntityToResponse(obj))
	}

	WriteHTTPResponse(w, result, http.StatusOK)
}

func (h *FileHandler) GetFile(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	var id int
	var err error

	id, err = strconv.Atoi(params.ByName("fileid"))
	if err != nil {
		BuildErrorResponse(w, entity.ErrorFileNotFound)
		return
	}

	result, err := h.usecase.GetFile(r.Context(), id)
	if err != nil {
		BuildErrorResponse(w, err)
		return
	}

	w.Header().Set("Content-Disposition", "attachment; filename="+result.Name)
	w.Header().Set("Content-Type", result.MimeType)
	http.ServeFile(w, r, fmt.Sprintf("%s/%s", util.StoragePath, result.Name))
}

func (h *FileHandler) DeleteFile(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	var id int
	var err error

	id, err = strconv.Atoi(params.ByName("fileid"))
	if err != nil {
		BuildErrorResponse(w, entity.ErrorFileNotFound)
		return
	}

	err = h.usecase.DeleteFile(r.Context(), id)
	if err != nil {
		BuildErrorResponse(w, err)
		return
	}

	WriteHTTPResponse(w, nil, http.StatusNoContent)
}

func fileEntityToResponse(eObj *entity.File) *response.File {
	return &response.File{
		ID:        fmt.Sprint(eObj.ID),
		Name:      eObj.Name,
		Size:      eObj.Size,
		CreatedAt: eObj.CreatedAt,
	}
}
