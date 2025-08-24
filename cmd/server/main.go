// cmd/server/main.go
package main

import (
	"fmt"
	"net/http"

	"github.com/keshupandre/img-to-pdf-backend/internal/api"
)

func main() {
	router := api.NewRouter()

	fmt.Println("🚀 Server started at http://localhost:8080")
	http.ListenAndServe(":8080", router)
}
