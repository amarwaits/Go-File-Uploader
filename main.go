package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"text/template"
	"time"
)

const MAX_UPLOAD_SIZE = 8 * 1024 * 1024 // 8MB

var TOKEN = os.Getenv("token")

// Compile templates on app start
var templates = template.Must(template.ParseFiles("public/upload.html"))

// Display named template
func displayForm(w http.ResponseWriter, page string, data interface{}) {
	templates.ExecuteTemplate(w, page+".html", data)
}

// Handler for file upload
func uploadFile(w http.ResponseWriter, r *http.Request) {
	token := r.FormValue("token")
	if token != TOKEN {
		http.Error(w, "Invalid token", http.StatusForbidden)
		return
	}

	// if r.ContentLength > MAX_UPLOAD_SIZE {
	// 	http.Error(w, "Request too large", http.StatusForbidden)
	// 	return
	// }

	r.Body = http.MaxBytesReader(w, r.Body, MAX_UPLOAD_SIZE)
	if err := r.ParseMultipartForm(MAX_UPLOAD_SIZE); err != nil {
		// TODO Debug this - wasn't working
		http.Error(w, "The uploaded file is too big. File should be less than 8MB in size", http.StatusForbidden)
		return
	}

	// Get file header for filename, size and headers
	file, fileHeader, err := r.FormFile("file")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer file.Close()

	// TODO Save this information in DB
	tempFileName := fmt.Sprintf("%d%s", time.Now().UnixNano(), filepath.Ext(fileHeader.Filename))
	fmt.Printf("Original file name: %+v\n", fileHeader.Filename)
	fmt.Printf("File Size: %+v\n", fileHeader.Size)
	fmt.Printf("MIME Header: %+v\n", fileHeader.Header)
	fmt.Printf("Temp file name: %+v\n", tempFileName)

	if fileHeader.Size > MAX_UPLOAD_SIZE {
		http.Error(w, "The uploaded file is too big. File should be less than 8MB in size", http.StatusForbidden)
		return
	}

	// Read first 512 bytes
	buff := make([]byte, 512)
	_, err = file.Read(buff)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Check for content type
	filetype := http.DetectContentType(buff)
	switch filetype {
	case "image/jpeg", "image/jpg", "image/png", "image/gif":
	default:
		http.Error(w, "File format is not allowed. Please upload a JPG/JPEG, PNG or GIF image", http.StatusForbidden)
		return
	}

	_, err = file.Seek(0, io.SeekStart)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Create file
	destination, err := os.Create(tempFileName)
	defer destination.Close()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Copy the uploaded file to the created file on the filesystem
	if _, err := io.Copy(destination, file); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "Upload successful")
}

func uploadRouter(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		displayForm(w, "upload", TOKEN)
	case "POST":
		uploadFile(w, r)
	}
}

func main() {
	// Upload route
	http.HandleFunc("/upload", uploadRouter)

	// Listen on port 8080
	http.ListenAndServe(":8080", nil)
}
