# Mini-Dropbox

A distributed file storage system built in Go that implements a master-node architecture for file storage and retrieval.

## Project Structure

```
Mini-dropbox/
├── cmd/
│   └── main.go              # Application entry point
├── internals/
│   ├── common/
│   │   └── types.go         # Shared data structures
│   ├── master/
│   │   └── master.go        # Master node implementation
│   ├── node/
│   │   └── node.go          # Storage node implementation
│   └── storage/
│       └── storage.go       # File storage operations
├── go.mod                   # Go module definition
└── README.md                # This file
```

## Architecture

This project implements a **distributed file storage system** with the following components:

### 1. Master Node (Port 9000)
- **Purpose**: Central coordinator that manages file metadata and node information
- **Responsibilities**:
  - File registration and metadata management
  - Node health monitoring
  - File location tracking
  - API endpoints for file operations

### 2. Storage Nodes (Ports 8001, 8002)
- **Purpose**: Individual storage units that store actual file data
- **Responsibilities**:
  - File upload and storage
  - File retrieval
  - Health monitoring
  - Local file management

### 3. Storage Layer
- **Purpose**: Low-level file operations and hash generation
- **Responsibilities**:
  - File I/O operations
  - MD5 hash generation
  - Directory management

## Features

- **Distributed Storage**: Files are distributed across multiple nodes
- **Hash-based Storage**: Files are identified by MD5 hash for deduplication
- **HTTP API**: RESTful endpoints for file operations
- **Concurrent Operations**: Go routines for handling multiple requests
- **File Deduplication**: Same content files get the same hash

## API Endpoints

### Master Node (Port 9000)
- `POST /register` - Register a file with metadata
- `GET /get?hash=<hash>` - Get file information by hash
- `GET /list` - List all registered files

### Storage Nodes (Ports 8001, 8002)
- `POST /upload` - Upload a file (multipart form)
- `GET /retrieve?hash=<hash>` - Retrieve a file by hash
- `GET /health` - Health check endpoint

### Seedream Proxy (Master, Port 9000)
- `POST /seedream/generate` - Proxy to Seedream API for image generation using 0/1/2 input images (multipart)
  - Form fields:
    - `prompt` (string, required)
    - `image1` (file, optional)
    - `image2` (file, optional)

## Getting Started

1. **Prerequisites**: Go 1.23 or higher
2. **Run the application**:
   ```bash
   go run cmd/main.go
   ```

This will start:
- 2 storage nodes on ports 8001 and 8002
- 1 master node on port 9000

### Seedream API Setup

Set your API key before using the Seedream proxy endpoint (recommended). If not set, the server uses a built-in development fallback key:
```bash
export SEEDREAM_API_KEY="<your_seedream_api_key>"
```

#### Generate with 0/1/2 input images

```bash
# 0 images
curl -X POST http://localhost:9000/seedream/generate \
  -F 'prompt=Ultra-realistic portrait of a smiling person'

# 1 image
curl -X POST http://localhost:9000/seedream/generate \
  -F 'prompt=Create a stylized variant of this person' \
  -F "image1=@/path/to/ref1.jpg"

# 2 images
curl -X POST http://localhost:9000/seedream/generate \
  -F 'prompt=Blend features of two references into one portrait' \
  -F "image1=@/path/to/ref1.jpg" \
  -F "image2=@/path/to/ref2.jpg"
```

## 🧪 Testing

### **Quick Start Testing**

1. **Start the application**:
   ```bash
   go run cmd/main.go
   ```

2. **Run complete test suite**:
   ```bash
   chmod +x test_*.sh
   ./test_complete.sh
   ```

### **Individual Test Scripts**

- **`test_health.sh`** - Test health endpoints of all nodes
- **`test_master.sh`** - Test master node API endpoints
- **`test_upload.sh`** - Test file upload to storage nodes
- **`test_retrieve.sh <hash>`** - Test file retrieval by hash
- **`test_complete.sh`** - Run all tests in sequence

### **Manual Testing with curl**

#### **Test File Upload**
```bash
# Create test file
echo "Hello Mini-Dropbox!" > test.txt

# Upload to Node 8001
curl -X POST -F "file=@test.txt" http://localhost:8001/upload

# Upload to Node 8002
curl -X POST -F "file=@test.txt" http://localhost:8002/upload
```

#### **Test Master Node API**
```bash
# List all files
curl http://localhost:9000/list

# Get file by hash
curl "http://localhost:9000/get?hash=<file_hash>"

# Register file manually
curl -X POST http://localhost:9000/register \
  -H "Content-Type: application/json" \
  -d '{"hash":"test123","name":"test.txt","size":1024,"replicas":["8001"]}'
```

#### **Test File Retrieval**
```bash
# Download file from Node 8001
curl "http://localhost:8001/retrieve?hash=<file_hash>" -o downloaded_file.txt

# Download file from Node 8002
curl "http://localhost:8002/retrieve?hash=<file_hash>" -o downloaded_file.txt
```

### **Expected Test Results**

- **Health Checks**: Storage nodes return "OK", Master node may not have health endpoint
- **File Upload**: Returns JSON with success status and file hash
- **File Retrieval**: Downloads file content to local system
- **Master API**: Manages file metadata and provides file information

## Data Flow

1. **File Upload**: Client uploads file to a storage node
2. **Hash Generation**: Node generates MD5 hash of file content
3. **Local Storage**: File is stored locally with hash as filename
4. **Metadata Registration**: Node registers file metadata with master
5. **File Retrieval**: Client requests file by hash from master, then retrieves from appropriate node

## Security Features

- File size limit (10MB per upload)
- Input validation and error handling
- Proper HTTP status codes
- File permission management (0644 for files, 0755 for directories)

## Future Enhancements

- File replication across nodes
- Load balancing
- Authentication and authorization
- File compression
- Backup and recovery mechanisms
