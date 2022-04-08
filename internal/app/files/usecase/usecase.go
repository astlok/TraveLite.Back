package files

import (
	files "TraveLite/internal/app/files/repository"
)

type FilesUseCase struct {
	filesRepo *files.FilesRepo
}

func NewFilesUseCase(filesRepo *files.FilesRepo) *FilesUseCase {
	return &FilesUseCase{
		filesRepo: filesRepo,
	}
}
