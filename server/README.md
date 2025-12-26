# Mini-Dropbox

A high-performance distributed file storage system built in Go that implements a master-node architecture for scalable file storage and retrieval with built-in AI image generation capabilities.

## ✨ Key Features

- **Distributed Storage**: Horizontal scaling with multiple storage nodes
- **Concurrent Operations**: Efficient goroutine-based request handling
- **File Deduplication**: Content-based addressing using MD5 hashing
- **Graceful Shutdown**: Proper cleanup and timeout handling
- **AI Integration**: Seedream API integration for image generation
- **RESTful API**: Clean HTTP endpoints for all operations
- **Health Monitoring**: Built-in health check endpoints
- **Thread-Safe**: Mutex-protected concurrent file operations
- **Connection Pooling**: Optimized HTTP client with connection reuse

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

- **Distributed Storage**: Files are distributed across multiple nodes for scalability
- **Hash-based Storage**: Content-addressed storage using MD5 for automatic deduplication
- **HTTP API**: RESTful endpoints for all file operations
- **Concurrent Operations**: Efficient Go routines with proper synchronization
- **File Deduplication**: Identical content files share the same hash and storage
- **Graceful Shutdown**: Coordinated shutdown with 5-second timeout
- **Server Timeouts**: Configured read/write/idle timeouts for better resource management
- **Thread-Safe Operations**: Mutex-protected concurrent access to shared resources
- **Connection Pooling**: HTTP client optimization for external API calls
- **AI Image Generation**: Integrated Seedream API with 0-2 reference images support

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

### Prerequisites
- Go 1.23 or higher
- 10MB+ free disk space for storage nodes

### Installation & Running

1. **Clone the repository**:
   ```bash
   git clone <repository-url>
   cd Mini-dropbox/server
   ```

2. **Install dependencies**:
   ```bash
   go mod tidy
   ```

3. **Run the application**:
   ```bash
   go run cmd/main.go
   ```

This will start:
- 2 storage nodes on ports 8001 and 8002
- 1 master node on port 9000

All nodes support graceful shutdown with Ctrl+C.

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

## 🔧 Performance Optimizations

### Recent Improvements (v1.1)

1. **Graceful Shutdown Enhancement**
   - Coordinated WaitGroup-based shutdown
   - 5-second timeout for cleanup operations
   - Proper context cancellation handling

2. **HTTP Server Optimization**
   - Read timeout: 15s (master), 30s (nodes)
   - Write timeout: 15s (master), 30s (nodes)
   - Idle timeout: 60s (master), 120s (nodes)

3. **Concurrency & Thread Safety**
   - RWMutex for optimized read-heavy operations
   - File operation locking in storage layer
   - Validation methods for data integrity

4. **Connection Pooling**
   - HTTP client with connection reuse
   - MaxIdleConns: 10 per host
   - IdleConnTimeout: 90 seconds

5. **Storage Layer Enhancement**
   - Built-in file deduplication check
   - Thread-safe concurrent file access
   - Optimized hash generation

## 📊 System Requirements

- **CPU**: 2+ cores recommended
- **RAM**: 512MB minimum, 1GB recommended
- **Disk**: 10MB+ for storage nodes
- **Network**: Low latency between nodes recommended
- **OS**: Linux, macOS, Windows (with Go support)

## Security Features

- File size limit (10MB per upload)
- Input validation and error handling
- Proper HTTP status codes
- File permission management (0644 for files, 0755 for directories)
- Content-based addressing prevents unauthorized file access

## 🛠️ Development

### Project Structure Best Practices
- `cmd/`: Application entry points
- `internals/`: Internal packages (not importable)
- `internals/common/`: Shared types and utilities
- `internals/master/`: Master node coordination logic
- `internals/node/`: Storage node implementation
- `internals/storage/`: Low-level file operations
- `internals/seedream/`: External API integration

### Adding New Storage Nodes
Edit [cmd/main.go](cmd/main.go) and add:
```go
wg.Add(1)
go func() {
    defer wg.Done()
    node.StartNode(ctx, "8003")
}()
```

## Future Enhancements

- File replication across nodes for redundancy
- Load balancing across storage nodes
- Authentication and authorization
- File compression for storage optimization
- Backup and recovery mechanisms
- Distributed consensus for master node failover

## 📝 License

This project is for educational and development purposes.

## 🤝 Contributing

Contributions are welcome! Please ensure:
- Code follows Go best practices
- All tests pass before submitting
- Proper error handling
- Thread-safe implementations
- Documentation updates for new features

---

**Built with ❤️ using Go | Last Updated: December 2025**
