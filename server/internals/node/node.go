package node

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mini-dropbox/internals/common"
	"mini-dropbox/internals/storage"
	"net/http"
	"path/filepath"
	"time"
)

type Node struct {
	Port    string
	Storage *storage.FileStorage
}

func StartNode(ctx context.Context, port string) {
	node := &Node{
		Port:    port,
		Storage: storage.NewFileStorage(filepath.Join("data", "node_"+port)),
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/upload", node.handleUpload)
	mux.HandleFunc("/retrieve", node.handleRetrieve)
	mux.HandleFunc("/health", node.handleHealth)

	fmt.Printf("Node starting on port %s\n", port)

	// Use context for graceful shutdown
	server := &http.Server{
		Addr:         ":" + port,
		Handler:      mux,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	go server.ListenAndServe()

	// Wait for context cancellation
	<-ctx.Done()
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	server.Shutdown(shutdownCtx)
}

func (n *Node) handleUpload(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if err := r.ParseMultipartForm(10 << 20); err != nil {
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "File not found in form", http.StatusBadRequest)
		return
	}
	defer file.Close()

	data, err := io.ReadAll(io.LimitReader(file, 10<<20))
	if err != nil || len(data) == 0 {
		http.Error(w, "Failed to read file", http.StatusInternalServerError)
		return
	}

	// store file
	hash, err := n.Storage.StoreFile(header.Filename, data)
	if err != nil {
		http.Error(w, "Failed to store file", http.StatusInternalServerError)
		return
	}

	response := common.UploadResponse{
		Success: true,
		Hash:    hash,
		Message: "File Upload Successful",
		Size:    int64(len(data)),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)

	fmt.Printf("Stored file %s with hash %s (%d bytes)\n", header.Filename, hash, len(data))

}

func (n *Node) handleRetrieve(w http.ResponseWriter, r *http.Request) {
	hash := r.URL.Query().Get("hash")

	if hash == "" {
		http.Error(w, "Missing hash parameter", http.StatusBadRequest)
		return
	}

	data, err := n.Storage.GetFile(hash)
	if err != nil {
		http.Error(w, "File not found", http.StatusNotFound)
		return
	}

	// Set headers for efficient transfer
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Length", fmt.Sprintf("%d", len(data)))
	w.Header().Set("Cache-Control", "public, max-age=31536000") // Cache for 1 year (immutable content)
	w.WriteHeader(http.StatusOK)
	w.Write(data)

	fmt.Printf("Retrieved file with hash %s (%d bytes)\n", hash, len(data))

}

func (n *Node) handleHealth(w http.ResponseWriter, r *http.Request) {
	// Fast path for health checks with caching
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "no-cache")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"healthy","port":"` + n.Port + `"}`))
}
