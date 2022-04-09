package files

import (
	files "TraveLite/internal/app/files/repository"
	"TraveLite/internal/models"
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

type FilesUseCase struct {
	filesRepo *files.FilesRepo
	fileURL   string
}

func NewFilesUseCase(filesRepo *files.FilesRepo, fileUrl string) *FilesUseCase {
	return &FilesUseCase{
		filesRepo: filesRepo,
		fileURL:   fileUrl,
	}
}

func (f FilesUseCase) Create(file models.FileMeta) (models.FileMeta, error) {
	var err error

	file.ID, err = uuid.NewUUID()
	if err != nil {
		return models.FileMeta{}, errors.Wrap(err, "can't create uuid")
	}

	err = f.filesRepo.Create(file)
	if err != nil {
		return models.FileMeta{}, errors.Wrap(err, "can't insert file to db")
	}

	return file, nil
}

func (f FilesUseCase) Upload(file []byte, fileID string) (models.FileInfo, error) {
	var fileInfo models.FileInfo

	err := f.filesRepo.Upload(file, fileID)
	if err != nil {
		return models.FileInfo{}, err
	}

	fileInfo.Link = f.fileURL + "/" + fileID

	return fileInfo, nil
}
