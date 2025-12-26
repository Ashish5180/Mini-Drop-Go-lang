package common

import "errors"

const (
	// Hash format constants
	MD5HashLength = 32
)

var (
	// Pre-computed errors for zero-allocation error returns
	ErrHashRequired = errors.New("hash is required")
	ErrSizeInvalid  = errors.New("size must be positive")
	ErrHashInvalid  = errors.New("invalid hash format")
)

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
	if len(f.Hash) != MD5HashLength {
		return ErrHashInvalid
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
