package storage

import (
	"os"
	"path/filepath"
)

type LocalStorageService struct {
	Path string
}

func NewLocalStorageService(path string) StorageService {
	return LocalStorageService{
		Path: path,
	}
}

func (l LocalStorageService) UploadFoto(fileName string, fileBytes []byte) (string, error) {
	storagePath := filepath.Join(l.Path, fileName)

	dst, err := os.Create(storagePath)
	if err != nil {
		return "", err
	}

	if _, err := dst.Write(fileBytes); err != nil {
		return "", err
	}

	return storagePath, nil
}

func (l LocalStorageService) GetFoto(filePath string) ([]byte, error) {
	fileBytes, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	return fileBytes, nil
}
