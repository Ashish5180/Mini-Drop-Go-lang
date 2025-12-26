package common

type FileInfo struct {
	Hash     string   `json:"hash"`
	Name     string   `json:"name"`
	Size     int64    `json:"size"`
	Replicas []string `json:"replicas"` // node addresses
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
