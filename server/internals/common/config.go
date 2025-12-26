package common

import "time"

// ServerConfig contains optimized server configuration
type ServerConfig struct {
	// HTTP Server timeouts
	ReadTimeout    time.Duration
	WriteTimeout   time.Duration
	IdleTimeout    time.Duration
	MaxHeaderBytes int

	// Rate limiting
	MaxRequestsPerSec int
	MaxConcurrent     int

	// File handling
	MaxFileSize int64
	BufferSize  int

	// Cache settings
	EnableCache bool
	CacheSize   int
	CacheTTL    time.Duration
}

// DefaultMasterConfig returns optimized config for master node
func DefaultMasterConfig() *ServerConfig {
	return &ServerConfig{
		ReadTimeout:       15 * time.Second,
		WriteTimeout:      15 * time.Second,
		IdleTimeout:       60 * time.Second,
		MaxHeaderBytes:    1 << 20, // 1 MB
		MaxRequestsPerSec: 1000,
		MaxConcurrent:     100,
		MaxFileSize:       10 << 20,  // 10 MB
		BufferSize:        32 * 1024, // 32 KB
		EnableCache:       true,
		CacheSize:         1000,
		CacheTTL:          5 * time.Minute,
	}
}

// DefaultNodeConfig returns optimized config for storage node
func DefaultNodeConfig() *ServerConfig {
	return &ServerConfig{
		ReadTimeout:       30 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       120 * time.Second,
		MaxHeaderBytes:    1 << 20, // 1 MB
		MaxRequestsPerSec: 500,
		MaxConcurrent:     50,
		MaxFileSize:       10 << 20,  // 10 MB
		BufferSize:        64 * 1024, // 64 KB
		EnableCache:       true,
		CacheSize:         500,
		CacheTTL:          10 * time.Minute,
	}
}
