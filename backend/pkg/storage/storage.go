package storage

type StorageService interface {
	UploadFoto(fileName string, fileBytes []byte) (string, error)
	GetFoto(filePath string) ([]byte, error)
}
