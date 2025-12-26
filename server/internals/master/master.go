package master

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mini-dropbox/internals/common"
	"mini-dropbox/internals/seedream"
	"net/http"
	"sync"
)

type Master struct {
	Port  string
	Files map[string]*common.FileInfo
	Nodes map[string]*common.NodeInfo
	mutex sync.RWMutex
}

func StartMaster(ctx context.Context, port string) {

	master := &Master{
		Port:  port,
		Files: make(map[string]*common.FileInfo),
		Nodes: make(map[string]*common.NodeInfo),
	}

	// Regsiter known nodes

	master.registerNode("8001")
	master.registerNode("8002")

	mux := http.NewServeMux()
	mux.HandleFunc("/register", master.handleRegister)
	mux.HandleFunc("/get", master.handleGet)
	mux.HandleFunc("/list", master.handleList)
	mux.HandleFunc("/seedream/generate", master.handleSeedreamGenerate)

	fmt.Printf("Starting master on port %s\n", port)

	server := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	go server.ListenAndServe()

	// Wait for context cancellation
	<-ctx.Done()
	server.Shutdown(context.Background())
}

func (m *Master) registerNode(address string) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	m.Nodes[address] = &common.NodeInfo{
		Address: address,
		Status:  "active",
	}

}

func (m *Master) handleRegister(w http.ResponseWriter, r *http.Request) {
	// handle file from nodes

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var fileInfo common.FileInfo

	if fileInfo.Size == 0 {
		http.Error(w, "Size is required", http.StatusBadRequest)
		return
	}
	if err := json.NewDecoder(r.Body).Decode(&fileInfo); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	m.mutex.Lock()
	m.Files[fileInfo.Hash] = &fileInfo
	m.mutex.Unlock()

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "File %s registered successfully", fileInfo.Hash)
}

func (m *Master) handleGet(w http.ResponseWriter, r *http.Request) {
	hash := r.URL.Query().Get("hash")

	if hash == "" {
		http.Error(w, "Hash parameter is required", http.StatusBadRequest)
		return
	}

	m.mutex.RLock()
	fileInfo, exists := m.Files[hash]

	m.mutex.RUnlock()

	if !exists {
		http.Error(w, "File not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(fileInfo)
}

func (m *Master) handleList(w http.ResponseWriter, r *http.Request) {
	m.mutex.RLock()
	files := make([]*common.FileInfo, 0, len(m.Files))
	for _, file := range m.Files {
		files = append(files, file)
	}
	m.mutex.RUnlock()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(files)
}

// handleSeedreamGenerate:
// POST multipart/form-data with fields:
//  - prompt: string (required)
//  - image1: file (optional)
//  - image2: file (optional)
// Responds with JSON containing job_id/images as returned by Seedream.
func (m *Master) handleSeedreamGenerate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	if err := r.ParseMultipartForm(20 << 20); err != nil {
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}
	prompt := r.FormValue("prompt")
	if prompt == "" {
		http.Error(w, "prompt is required", http.StatusBadRequest)
		return
	}
	var images [][]byte
	if f, _, err := r.FormFile("image1"); err == nil && f != nil {
		defer f.Close()
		if b, readErr := io.ReadAll(f); readErr == nil {
			images = append(images, b)
		}
	}
	if f, _, err := r.FormFile("image2"); err == nil && f != nil {
		defer f.Close()
		if b, readErr := io.ReadAll(f); readErr == nil {
			images = append(images, b)
		}
	}
	ctx := r.Context()
	resp, err := seedream.GenerateImages(ctx, "", prompt, images)
	if err != nil {
		http.Error(w, fmt.Sprintf("generation failed: %v", err), http.StatusBadGateway)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
