package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	vision "cloud.google.com/go/vision/apiv1"
	"google.golang.org/api/option"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.Use(cors.Default())

	r.POST("/upload", uploadHandler)

	r.Run(":8080")
}

func uploadHandler(c *gin.Context) {
	file, err := c.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid file"})
		return
	}

	// ファイルの保存
	filePath := fmt.Sprintf("./uploads/%s", file.Filename)
	if err := c.SaveUploadedFile(file, filePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
		return
	}

	// 画像からテキストを抽出
	text, err := extractText(filePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to extract text"})
		return
	}

	// 結果をJSONで返す
	c.JSON(http.StatusOK, gin.H{"text": text, "file": file.Filename})
}

func extractText(filePath string) (string, error) {
	ctx := context.Background()

	client, err := vision.NewImageAnnotatorClient(ctx, option.WithCredentialsFile("path/to/your/service-account-file.json"))
	if err != nil {
		return "", err
	}
	defer client.Close()

	fileBytes, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}

	image := vision.NewImageFromBytes(fileBytes)
	annotation, err := client.DetectDocumentText(ctx, image, nil)
	if err != nil {
		return "", err
	}

	return annotation.Text, nil
}
