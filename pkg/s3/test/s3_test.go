package test

import (
	"context"
	"github.com/PosokhovVadim/stawberry/internal/config"
	"github.com/PosokhovVadim/stawberry/pkg/s3"
	"github.com/stretchr/testify/assert"
	"io"
	"log"
	"os"
	"testing"
)

func TestAWSUpload(t *testing.T) {
	ctx := context.Background()
	objectKey := "fortest"
	file, err := os.Open("79135.png")
	assert.NoError(t, err)
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(file)

	cfg := config.LoadConfig()
	bucket := objectstorage.ObjectStorageConn(cfg)

	err = bucket.UploadFileWithPresignedURL(ctx, objectKey, file)
	assert.NoError(t, err)
}

func TestAWSDownload(t *testing.T) {
	ctx := context.Background()
	objectKey := "fortest"

	cfg := config.LoadConfig()
	bucket := objectstorage.ObjectStorageConn(cfg)

	expectedFile, err := os.Open("79135.png")
	defer func(expectedFile *os.File) {
		err := expectedFile.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(expectedFile)
	assert.NoError(t, err)

	bytesExpectedFile, err := io.ReadAll(expectedFile)
	assert.NoError(t, err)

	file, err := bucket.DownloadFile(ctx, objectKey)
	assert.NoError(t, err)
	assert.Equal(t, bytesExpectedFile, file)
}
