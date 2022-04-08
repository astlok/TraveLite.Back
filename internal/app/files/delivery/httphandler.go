package files

import (
	files "TraveLite/internal/app/files/usecase"
	"github.com/labstack/echo/v4"
)

type FilesHandler struct {
	filesUseCase *files.FilesUseCase
}

func NewFileHandler(fileUseCase *files.FilesUseCase) *FilesHandler {
	return &FilesHandler{
		filesUseCase: fileUseCase,
	}
}

func (f FilesHandler) Create(c echo.Context) error {
	return nil
}
