# Image to PDF Converter - Backend

A Go-based backend service that converts images to PDF files following clean architecture principles.

## Architecture

The backend follows Go best practices with a layered architecture:

```
backend/
├── cmd/
│   └── main.go              # Application entry point
├── internal/
│   ├── config/              # Configuration management
│   │   └── config.go
│   ├── handlers/            # HTTP handlers
│   │   │── base.go
│   │   │── download.go
│   │   │── health.go
│   │   └── upload.go
│   ├── services/            # Business logic
│   │   ├── pdf_service.go   # PDF conversion logic
│   │   └── file_service.go  # File operations
│   ├── models/              # Data structures
│   │   └── models.go
│   └── utils/               # Utility functions
│       └── file_utils.go
├── pkg/                     # Public packages (if any)
├── temp/                    # Temporary files directory
├── uploads/                 # Upload directory
├── output/                  # Generated PDF output
├── go.mod
├── go.sum
├── Dockerfile
└── README.md
```

## Features

- **Clean Architecture**: Separated concerns with clear layer boundaries
- **Configuration Management**: Environment-based configuration
- **File Validation**: Size and type validation for uploaded files
- **PDF Generation**: High-quality PDF conversion with aspect ratio preservation
- **CORS Support**: Cross-origin resource sharing for frontend integration
- **Health Checks**: Service health monitoring
- **Logging**: Comprehensive request and error logging
- **Docker Support**: Containerized deployment

## Dependencies

- **HTTP Router**: `github.com/go-chi/chi/v5` - Fast HTTP router
- **PDF Generation**: `github.com/jung-kurt/gofpdf` - PDF creation library
- **CORS**: `github.com/rs/cors` - Cross-origin resource sharing

## Configuration

The application uses environment variables for configuration:

| Variable | Default | Description |
|----------|---------|-------------|
| `PORT` | `8080` | Server port |
| `HOST` | `localhost` | Server host |
| `DEBUG` | `true` | Debug mode |
| `FRONTEND_URL` | `http://localhost:3000` | Frontend URL for CORS |
| `MAX_FILE_SIZE` | `10485760` | Max file size in bytes (10MB) |
| `MAX_FILES` | `10` | Maximum number of files per upload |
| `TEMP_DIR` | `./temp` | Temporary files directory |
| `UPLOAD_DIR` | `./uploads` | Upload directory |
| `PDF_OUTPUT_DIR` | `./output` | PDF output directory |

## API Endpoints

### Upload Images
- **POST** `/upload`
- **Content-Type**: `multipart/form-data`
- **Form Field**: `files` (multiple files)
- **Response**: JSON with PDF filename

### Download PDF
- **GET** `/download?file={filename}`
- **Response**: PDF file download

### Health Check
- **GET** `/health`
- **Response**: JSON with service status

### Root
- **GET** `/`
- **Response**: JSON with API information

## Running the Application

### Development
```bash
# Install dependencies
go mod tidy

# Run the application
go run cmd/main.go
```

### Production
```bash
# Build the binary
go build -o img-to-pdf-converter cmd/main.go

# Run the binary
./img-to-pdf-converter
```

### Docker
```bash
# Build the image
docker build -t img-to-pdf-converter .

# Run the container
docker run -p 8080:8080 img-to-pdf-converter
```

## Testing

```bash
# Run tests
go test ./...

# Run tests with coverage
go test -cover ./...
```

## Project Structure Explanation

- **`cmd/`**: Contains the main application entry points
- **`internal/`**: Private application code that shouldn't be imported by other applications
- **`pkg/`**: Public packages that can be imported by other applications
- **`config/`**: Configuration management and environment variable handling
- **`handlers/`**: HTTP request handlers (controllers in MVC terms)
- **`services/`**: Business logic layer containing core application functionality
- **`models/`**: Data structures and interfaces
- **`utils/`**: Utility functions and helpers

This structure follows Go community standards and makes the codebase maintainable and testable.
