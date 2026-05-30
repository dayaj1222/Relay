package uploads

import (
	"crypto/rand"
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

var port = os.Getenv("PORT")

const MaxUploadSize = 50 * 1024 * 1024
const UploadDir = "./uploads"

func UploadHandler(c *gin.Context) {
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, MaxUploadSize)

	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid file key or size limit exceeded"})
		return
	}

	b := make([]byte, 8)
	rand.Read(b)
	uniqueName := fmt.Sprintf("%x%s", b, filepath.Ext(file.Filename))

	if err := os.MkdirAll(UploadDir, os.ModePerm); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create directory"})
		return
	}

	dst := filepath.Join(UploadDir, uniqueName)
	if err := c.SaveUploadedFile(file, dst); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
		return
	}

	if port == "" {
		port = "8000"
	}
	fileURL := fmt.Sprintf("http://localhost:%s/uploads/%s", port, uniqueName)
	c.JSON(http.StatusOK, gin.H{
		"fileUrl":  fileURL,
		"fileName": file.Filename,
	})
}
