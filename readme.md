# Image to PDF Converter (Go)

A simple web application backend written in Go that converts uploaded images into a single PDF file.

## 📂 Project Structure
```
img-to-pdf-backend/
├── cmd/                # Application entry point
│   └── server/
│       └── main.go
├── internal/
│   ├── api/            # API layer (router + handlers)
│   │   ├── router.go
│   │   └── handlers/
│   │       ├── health.go
│   │       └── convert.go
│   ├── services/       # Business logic (PDF conversion)
│   │   └── pdf_service.go
│   └── storage/        # (optional) for file storage logic
├── uploads/            # Temporary file storage (ignored in git)
├── go.mod
├── go.sum
└── README.md
```

## 🚀 Running the Server

1. Clone the repo:
   ```bash
   git clone https://github.com/keshupandre/img-to-pdf-backend.git
   cd img-to-pdf-backend
   ```

2. Install dependencies:
   ```bash
   go mod tidy
   ```

3. Run the server:
   ```bash
   go run cmd/server/main.go
   ```

4. Test endpoints:
   - Health check: [http://localhost:8080/api/health](http://localhost:8080/api/health)
   - Convert images to PDF: POST request with `form-data` key `images` → returns PDF URL

## 🛠️ Tech Stack
- [Go](https://go.dev/) (Golang)
- [Chi Router](https://github.com/go-chi/chi)
- [gofpdf](https://github.com/jung-kurt/gofpdf)

## 📌 Notes
- The `uploads/` folder is automatically created to store images and generated PDFs.
- PDFs are overwritten each time (can be extended to use unique IDs).
- `.env` file can be added for configuration.

