package common

import "errors"

type FileInfo struct {
	Hash     string   `json:"hash"`
	Name     string   `json:"name"`
	Size     int64    `json:"size"`
	Replicas []string `json:"replicas"` // node addresses
}

// Validate checks if FileInfo has required fields
func (f *FileInfo) Validate() error {
	if f.Hash == "" {
		return errors.New("hash is required")
	}
	if f.Size <= 0 {
		return errors.New("size must be positive")
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
}
