package common

import "time"

// ServerConfig holds optimized server settings
type ServerConfig struct {
	ReadTimeout   time.Duration
	WriteTimeout  time.Duration
	IdleTimeout   time.Duration
	MaxFileSize   int64
	MaxConcurrent int
	CacheSize     int
}

// DefaultMasterConfig returns master node configuration
func DefaultMasterConfig() *ServerConfig {
	return &ServerConfig{
		ReadTimeout:   15 * time.Second,
		WriteTimeout:  15 * time.Second,
		IdleTimeout:   60 * time.Second,
		MaxFileSize:   10 << 20, // 10 MB
		MaxConcurrent: 100,
		CacheSize:     1000,
	}
}

// DefaultNodeConfig returns storage node configuration
func DefaultNodeConfig() *ServerConfig {
	return &ServerConfig{
		ReadTimeout:   30 * time.Second,
		WriteTimeout:  30 * time.Second,
		IdleTimeout:   120 * time.Second,
		MaxFileSize:   10 << 20, // 10 MB
		MaxConcurrent: 50,
		CacheSize:     500,
	}
}
