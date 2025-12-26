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
		Addr:    ":" + port,
		Handler: mux,
	}
	
	go server.ListenAndServe()
	
	// Wait for context cancellation
	<-ctx.Done()
	server.Shutdown(context.Background())
}

func (n *Node) handleUpload(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	err := r.ParseMultipartForm(10 << 20) // limit to 10MB

	if err != nil {
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}

	file, header, err := r.FormFile("file")

	if err != nil {
		http.Error(w, "Failed to get file from form", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Read the file

	data, err := io.ReadAll(file)
	if err != nil {
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
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)

	fmt.Printf("Stored file %s with hash %s\n", header.Filename, hash)

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

	w.Header().Set("Content-Type", "application/octet-stream")
	w.Write(data)

	fmt.Printf("Retrieved file with hash %s\n", hash)

}

func (n *Node) handleHealth(w http.ResponseWriter, r *http.Request) {
	// Handle health check
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}
