package storage

import (
	"crypto/md5"
	"fmt"
	"hash"
	"os"
	"path/filepath"
	"sync"
)

var (
	// Buffer pool for efficient memory reuse
	bufferPool = sync.Pool{
		New: func() interface{} {
			return make([]byte, 32*1024) // 32KB buffers
		},
	}
	// Hash pool for MD5 computation reuse
	hashPool = sync.Pool{
		New: func() interface{} {
			return md5.New()
		},
	}
)

type FileStorage struct {
	StorageDir string
	mu         sync.RWMutex
	cache      map[string]bool // Simple existence cache
	cacheMu    sync.RWMutex
}

func (fs *FileStorage) GetFile(hash string) ([]byte, error) {
	fs.mu.RLock()
	defer fs.mu.RUnlock()
	filePath := filepath.Join(fs.StorageDir, hash)
	return os.ReadFile(filePath)
}

func NewFileStorage(dir string) *FileStorage {
	os.MkdirAll(dir, 0755)
	return &FileStorage{
		StorageDir: dir,
		cache:      make(map[string]bool, 100), // Pre-allocate cache
	}
}

func (fs *FileStorage) StoreFile(filename string, data []byte) (string, error) {
	// Compute hash with pooled hasher for efficiency
	hasher := hashPool.Get().(hash.Hash)
	hasher.Reset()
	hasher.Write(data)
	hashStr := fmt.Sprintf("%x", hasher.Sum(nil))
	hashPool.Put(hasher)

	// Check cache first (fast path)
	fs.cacheMu.RLock()
	exists := fs.cache[hashStr]
	fs.cacheMu.RUnlock()

	if exists {
		return hashStr, nil
	}

	fs.mu.Lock()
	defer fs.mu.Unlock()

	filePath := filepath.Join(fs.StorageDir, hashStr)

	// Double-check file existence
	if _, err := os.Stat(filePath); err == nil {
		fs.updateCache(hashStr, true)
		return hashStr, nil
	}

	// Write file atomically
	err := os.WriteFile(filePath, data, 0644)
	if err == nil {
		fs.updateCache(hashStr, true)
	}
	return hashStr, err
}

func (fs *FileStorage) RetrieveFile(hash string) ([]byte, error) {
	fs.mu.RLock()
	defer fs.mu.RUnlock()
	filePath := filepath.Join(fs.StorageDir, hash)
	return os.ReadFile(filePath)
}

func (fs *FileStorage) DeleteFile(hash string) error {
	fs.mu.Lock()
	defer fs.mu.Unlock()
	filePath := filepath.Join(fs.StorageDir, hash)
	err := os.Remove(filePath)
	if err == nil {
		fs.updateCache(hash, false)
	}
	return err
}

func (fs *FileStorage) FileExists(hash string) bool {
	// Check cache first
	fs.cacheMu.RLock()
	exists, cached := fs.cache[hash]
	fs.cacheMu.RUnlock()

	if cached {
		return exists
	}

	// Check filesystem
	fs.mu.RLock()
	defer fs.mu.RUnlock()
	filePath := filepath.Join(fs.StorageDir, hash)
	_, err := os.Stat(filePath)
	exists = err == nil
	fs.updateCache(hash, exists)
	return exists
}

// updateCache safely updates the file existence cache
func (fs *FileStorage) updateCache(hash string, exists bool) {
	fs.cacheMu.Lock()
	defer fs.cacheMu.Unlock()
	fs.cache[hash] = exists
}
