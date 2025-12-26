package common

import (
	"errors"
	"sync"
)

var (
	// Pre-computed error messages for efficiency
	ErrHashRequired = errors.New("hash is required")
	ErrSizeInvalid  = errors.New("size must be positive")
	ErrNameRequired = errors.New("name is required")

	// String interning pool for hash deduplication
	hashInternPool = sync.Map{}
)

// InternHash returns a canonical representation of the hash string
// to reduce memory usage when same hashes appear multiple times
func InternHash(hash string) string {
	if actual, loaded := hashInternPool.LoadOrStore(hash, hash); loaded {
		return actual.(string)
	}
	return hash
}

type FileInfo struct {
	Hash     string   `json:"hash"`
	Name     string   `json:"name"`
	Size     int64    `json:"size"`
	Replicas []string `json:"replicas"` // node addresses
}

// Validate checks if FileInfo has required fields
func (f *FileInfo) Validate() error {
	if f.Hash == "" {
		return ErrHashRequired
	}
	if len(f.Hash) != 32 { // MD5 hash is always 32 hex chars
		return errors.New("invalid hash format")
	}
	if f.Size <= 0 {
		return ErrSizeInvalid
	}
	return nil
}

type NodeInfo struct {
	Address string `json:"address"`
	Port    int    `json:"port"`
	Status  string `json:"status"` // e.g., "active", "inactive"
}

type UploadResponse struct {
	Success bool   `json:"success"`
	Hash    string `json:"hash"`
	Message string `json:"message"`
	Size    int64  `json:"size,omitempty"` // File size in bytes
}
