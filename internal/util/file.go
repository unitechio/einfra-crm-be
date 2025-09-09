
package util

import (
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

// SaveUploadedFile saves an uploaded file to the specified path.
// It creates the destination directory if it does not exist.
func SaveUploadedFile(file *multipart.FileHeader, dst string) error {
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	// Create the destination directory if it doesn't exist
	if err := os.MkdirAll(filepath.Dir(dst), 0750); err != nil {
		return err
	}

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, src)
	return err
}

// ServeFiles creates a Gin handler that serves files from a directory.
func ServeFiles(dir string) gin.HandlerFunc {
	return func(c *gin.Context) {
		fileName := c.Param("filename")
		filePath := filepath.Join(dir, fileName)

		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			c.String(http.StatusNotFound, "File not found")
			return
		}

		c.File(filePath)
	}
}
