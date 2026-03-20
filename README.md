# FileStore Server

A lightweight file storage server built with Go, supporting file upload, download, user management, and chunked uploads.

## 🚀 Features

- 📁 **File Management**
  - Upload files with metadata
  - Download files
  - Update file metadata
  - Delete files
  - Query file information

- 🔐 **User Authentication**
  - User registration (signup)
  - User login (signin)
  - User information retrieval
  - JWT-based authentication

- ⚡ **Chunked Upload**
  - Support for large file uploads
  - Chunk status checking
  - Automatic chunk merging

- 🗄️ **Storage Backend**
  - MySQL database for metadata
  - Redis for session/cache management
  - Local file system storage

## 🏗️ Architecture

```
filestore-server/
├── main.go              # Entry point and HTTP routing
├── db/                  # Database operations
│   ├── mysql/conn.go    # MySQL connection
│   ├── file.go          # File-related DB operations
│   └── user.go          # User-related DB operations
├── handler/             # HTTP request handlers
│   ├── auth.go          # Authentication middleware
│   ├── handler.go       # File upload/download handlers
│   └── user.go          # User management handlers
├── meta/                # File metadata management
│   └── filemeta.go      # File metadata structure
├── rd/                  # Redis operations
│   └── redis.go         # Redis connection and operations
├── util/                # Utility functions
│   └── chunk.go         # Chunk upload utilities
├── static/              # Static files (frontend assets)
├── uploads/             # Uploaded files storage
└── go.mod               # Go module definition
```

## 🛠️ Technology Stack

- **Language**: Go 1.24.0
- **Database**: MySQL
- **Cache**: Redis
- **Web Framework**: net/http (standard library)
- **Authentication**: JWT

## 📋 API Endpoints

### File Operations

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/file/upload` | Upload a file |
| GET | `/file/meta` | Get file metadata |
| GET | `/file/query` | Query files |
| GET | `/file/download` | Download a file |
| POST | `/file/update` | Update file metadata |
| POST | `/file/delete` | Delete a file |

### Chunk Upload

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/file/upload/chunk` | Upload file chunk |
| GET | `/file/upload/status` | Check chunk upload status |
| POST | `/file/upload/merge` | Merge uploaded chunks |

### User Operations

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/user/signup` | User registration |
| POST | `/user/signin` | User login |
| GET | `/user/info` | Get user information |

## 🚦 Getting Started

### Prerequisites

- Go 1.24.0 or higher
- MySQL database
- Redis server

### Installation

1. **Clone the repository**
   ```bash
   git clone <repository-url>
   cd filestore-server
   ```

2. **Install dependencies**
   ```bash
   go mod tidy
   ```

3. **Configure database connections**
   - Update MySQL connection settings in `db/mysql/conn.go`
   - Update Redis connection settings in `rd/redis.go`

4. **Create database tables**
   ```sql
   -- Create users table
   CREATE TABLE users (
       id INT AUTO_INCREMENT PRIMARY KEY,
       username VARCHAR(50) UNIQUE NOT NULL,
       password VARCHAR(100) NOT NULL,
       email VARCHAR(100),
       create_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP
   );

   -- Create files table
   CREATE TABLE file_meta (
       id INT AUTO_INCREMENT PRIMARY KEY,
       file_hash VARCHAR(100) NOT NULL,
       file_name VARCHAR(255) NOT NULL,
       file_size BIGINT DEFAULT 0,
       file_path VARCHAR(255) NOT NULL,
       create_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
       update_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
       status INT DEFAULT 0
   );
   ```

5. **Run the server**
   ```bash
   go run main.go
   ```

6. **Access the server**
   - Server will start on `http://localhost:8080`
   - Static files served from `/static/`

## 📝 Usage Examples

### Upload a File
```bash
curl -X POST -F "file=@/path/to/your/file.txt" http://localhost:8080/file/upload
```

### Download a File
```bash
curl -X GET "http://localhost:8080/file/download?filehash=abc123" --output file.txt
```

### User Registration
```bash
curl -X POST -H "Content-Type: application/json" \
  -d '{"username":"testuser","password":"password123"}' \
  http://localhost:8080/user/signup
```

### User Login
```bash
curl -X POST -H "Content-Type: application/json" \
  -d '{"username":"testuser","password":"password123"}' \
  http://localhost:8080/user/signin
```

## ⚙️ Configuration

### Database Configuration
Update the connection strings in the respective files:

**MySQL** (`db/mysql/conn.go`):
```go
db, err := sql.Open("mysql", "username:password@tcp(localhost:3306)/database_name")
```

**Redis** (`rd/redis.go`):
```go
client := redis.NewClient(&redis.Options{
    Addr:     "localhost:6379",
    Password: "",
    DB:       0,
})
```

### Server Configuration
The server runs on port 8080 by default. To change the port, modify `main.go`:
```go
err := http.ListenAndServe(":8080", nil)  // Change 8080 to desired port
```

## 🔒 Security Considerations

- Passwords should be hashed before storing (implement in user.go)
- Use HTTPS in production
- Validate file types and sizes
- Implement rate limiting for uploads
- Sanitize file names and paths

## 🚧 Development Status

This project is under active development. Features may change and APIs are subject to modification.

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests if applicable
5. Submit a pull request

## 📄 License

This project is for educational and learning purposes.

## 🆘 Support

For issues and questions, please open an issue in the repository.