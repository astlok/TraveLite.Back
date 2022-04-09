package files

import (
	files "TraveLite/internal/app/files/usecase"
	"TraveLite/internal/models"
	"github.com/labstack/echo/v4"
	"io/ioutil"
	"net/http"
)

type FilesHandler struct {
	filesUseCase *files.FilesUseCase
}

func NewFileHandler(fileUseCase *files.FilesUseCase) *FilesHandler {
	return &FilesHandler{
		filesUseCase: fileUseCase,
	}
}

// Create File godoc
// @Summary Create file metadata info
// @Description Create file metadata info like uuid
// @Tags Files
// @Accept json
// @Produce json
// @Param file body models.FileMeta true "add file info, id return"
// @Success 201 {object} models.FileMeta
// @Failure 500 {object} echo.HTTPError
// @Router /files [post]
func (f FilesHandler) Create(c echo.Context) error {
	var file models.FileMeta

	err := c.Bind(&file)
	if err != nil {
		return err
	}

	file, err = f.filesUseCase.Create(file)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, file)
}

// Upload File godoc
// @Summary Upload file to storage
// @Description Upload file to remote storage
// @Tags Files
// @Accept application/octet-stream
// @Produce json
// @Param file body string true "attach file in application/octet-stream body"
// @Param id path string true "file id, type uuid"
// @Success 200 {object} models.FileInfo
// @Failure 500 {object} echo.HTTPError
// @Router /files/{id} [put]
func (f FilesHandler) Upload(c echo.Context) error {
	fileBody, err := ioutil.ReadAll(c.Request().Body)
	if err != nil {
		return err
	}

	fileID := c.Param("id")

	fileInfo, err := f.filesUseCase.Upload(fileBody, fileID)

	return c.JSON(http.StatusOK, fileInfo)
}
