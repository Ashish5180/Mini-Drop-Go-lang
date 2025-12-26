package storage

import (
	"crypto/md5"
	"fmt"
	"os"
	"path/filepath"
	"sync"
)

type FileStorage struct {
	StorageDir string
	mu         sync.RWMutex
}

func (fs *FileStorage) GetFile(hash string) ([]byte, error) {
	fs.mu.RLock()
	defer fs.mu.RUnlock()
	filePath := filepath.Join(fs.StorageDir, hash)
	return os.ReadFile(filePath)
}

func NewFileStorage(dir string) *FileStorage {
	os.MkdirAll(dir, 0755)
	return &FileStorage{StorageDir: dir}
}

func (fs *FileStorage) StoreFile(filename string, data []byte) (string, error) {
	fs.mu.Lock()
	defer fs.mu.Unlock()

	// Generate file hash
	hash := fmt.Sprintf("%x", md5.Sum(data))
	filePath := filepath.Join(fs.StorageDir, hash)

	// Check if file already exists (deduplication)
	if _, err := os.Stat(filePath); err == nil {
		return hash, nil
	}

	return hash, os.WriteFile(filePath, data, 0644)
}

func (fs *FileStorage) RetrieveFile(hash string) ([]byte, error) {
	filePath := filepath.Join(fs.StorageDir, hash)
	return os.ReadFile(filePath)
}

func (fs *FileStorage) DeleteFile(hash string) error {
	filePath := filepath.Join(fs.StorageDir, hash)
	return os.Remove(filePath)
}

func (fs *FileStorage) FileExists(hash string) bool {
	filePath := filepath.Join(fs.StorageDir, hash)
	_, err := os.Stat(filePath)
	return err == nil

}
