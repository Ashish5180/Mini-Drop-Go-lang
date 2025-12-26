package storage

import (
	"crypto/md5"
	"fmt"
	"hash"
	"os"
	"path/filepath"
	"sync"
)

const (
	// Buffer and cache sizes
	defaultBufferSize = 32 * 1024 // 32KB
	defaultCacheSize  = 100
)

var (
	// Buffer pool for efficient memory reuse
	bufferPool = sync.Pool{
		New: func() interface{} {
			return make([]byte, defaultBufferSize)
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

// GetFile retrieves file by hash (alias for RetrieveFile)
func (fs *FileStorage) GetFile(hash string) ([]byte, error) {
	return fs.RetrieveFile(hash)
}

func NewFileStorage(dir string) *FileStorage {
	os.MkdirAll(dir, 0755)
	return &FileStorage{
		StorageDir: dir,
		cache:      make(map[string]bool, defaultCacheSize),
	}
}

func (fs *FileStorage) StoreFile(filename string, data []byte) (string, error) {
	// Compute hash with pooled hasher
	hasher := hashPool.Get().(hash.Hash)
	hasher.Reset()
	hasher.Write(data)
	hashStr := fmt.Sprintf("%x", hasher.Sum(nil))
	hashPool.Put(hasher)

	filePath := filepath.Join(fs.StorageDir, hashStr)

	// Fast path: check cache first
	fs.cacheMu.RLock()
	if exists := fs.cache[hashStr]; exists {
		fs.cacheMu.RUnlock()
		return hashStr, nil
	}
	fs.cacheMu.RUnlock()

	// Check and write with single lock
	fs.mu.Lock()
	defer fs.mu.Unlock()

	if _, err := os.Stat(filePath); err == nil {
		fs.setCached(hashStr, true)
		return hashStr, nil
	}

	if err := os.WriteFile(filePath, data, 0644); err != nil {
		return hashStr, err
	}

	fs.setCached(hashStr, true)
	return hashStr, nil
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
	if err := os.Remove(filePath); err != nil {
		return err
	}

	fs.setCached(hash, false)
	return nil
}

func (fs *FileStorage) FileExists(hash string) bool {
	fs.cacheMu.RLock()
	exists, cached := fs.cache[hash]
	fs.cacheMu.RUnlock()

	if cached {
		return exists
	}

	filePath := filepath.Join(fs.StorageDir, hash)
	_, err := os.Stat(filePath)
	exists = err == nil
	fs.setCached(hash, exists)
	return exists
}

// setCached updates cache without lock (caller must hold lock or handle sync)
func (fs *FileStorage) setCached(hash string, exists bool) {
	fs.cacheMu.Lock()
	fs.cache[hash] = exists
	fs.cacheMu.Unlock()
}
