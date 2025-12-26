package common

import "time"

const (
	// Common timeout values
	DefaultReadTimeout  = 15 * time.Second
	DefaultWriteTimeout = 15 * time.Second
	DefaultIdleTimeout  = 60 * time.Second

	// File size limits
	DefaultMaxFileSize = 10 << 20 // 10 MB
)

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
		ReadTimeout:   DefaultReadTimeout,
		WriteTimeout:  DefaultWriteTimeout,
		IdleTimeout:   DefaultIdleTimeout,
		MaxFileSize:   DefaultMaxFileSize,
		MaxConcurrent: 100,
		CacheSize:     1000,
	}
}

// DefaultNodeConfig returns storage node configuration
func DefaultNodeConfig() *ServerConfig {
	return &ServerConfig{
		ReadTimeout:   2 * DefaultReadTimeout,
		WriteTimeout:  2 * DefaultWriteTimeout,
		IdleTimeout:   2 * DefaultIdleTimeout,
		MaxFileSize:   DefaultMaxFileSize,
		MaxConcurrent: 50,
		CacheSize:     500,
	}
}
