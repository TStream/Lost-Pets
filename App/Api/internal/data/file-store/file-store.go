package filestore

import (
	"fmt"
	"io"
	"lostpets/internal"
	"mime/multipart"
	"os"
	"path"
)

type (
	Config struct {
		Location string `json:"location"`
	}

	FileStore struct {
		FilePath string
	}
)

func NewFileStore(config Config) (*FileStore, error) {
	if _, err := os.Stat(config.Location); os.IsNotExist(err) {
		err := os.MkdirAll(config.Location, os.ModePerm)
		if err != nil {
			return nil, err
		}
	}

	fileStore := &FileStore{
		FilePath: config.Location,
	}

	return fileStore, nil
}

func (fs *FileStore) GetFile(guid string) (string, error) {
	filepath := path.Join(fs.FilePath, guid)
	if _, err := os.Stat(filepath); os.IsNotExist(err) {
		return "", fmt.Errorf("file: %s Not Found", guid)
	}
	return filepath, nil
}
func (fs *FileStore) SaveFile(src multipart.File) (string, error) {
	uuid, err := internal.NewUUID()
	// Destination
	dst, err := os.Create(uuid)
	if err != nil {
		return "", err
	}
	defer dst.Close()

	// Copy
	if _, err = io.Copy(dst, src); err != nil {
		return "", err
	}

	return uuid, nil
}
func (fs *FileStore) DeleteFile(guid string) error {
	filepath := path.Join(fs.FilePath, guid)
	if _, err := os.Stat(filepath); os.IsNotExist(err) {
		return fmt.Errorf("file: %s Not Found", guid)
	}

	err := os.Remove(filepath)

	if err != nil {
		return fmt.Errorf("error deleteing file: %s %e", guid, err)
	}

	return nil
}
