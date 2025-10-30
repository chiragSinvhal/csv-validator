package services

import (
	"bytes"
	"mime/multipart"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFileService_SaveFile(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "file-test")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	fs := NewFileService(tempDir, tempDir+"-downloads")

	// Make a fake file upload
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, _ := writer.CreateFormFile("file", "test.csv")
	part.Write([]byte("name,email\nJohn,Chirag@test.com"))
	writer.Close()

	req := bytes.NewReader(body.Bytes())
	reader := multipart.NewReader(req, writer.Boundary())
	form, _ := reader.ReadForm(1024)
	file := form.File["file"][0]

	savedPath, err := fs.SaveFile(file, "test-job-123")
	require.NoError(t, err)
	assert.Contains(t, savedPath, "test-job-123")
	assert.FileExists(t, savedPath)

	// Read back and verify
	content, _ := os.ReadFile(savedPath)
	assert.Contains(t, string(content), "Chirag,Chirag@test.com")
}

func TestFileService_GetFile(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "file-test")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	fs := NewFileService(tempDir, tempDir+"-downloads")

	// Create a test file directly
	testFile := filepath.Join(tempDir, "existing.csv")
	os.WriteFile(testFile, []byte("test"), 0644)

	// Should find it
	path, err := fs.GetFile("existing.csv")
	assert.NoError(t, err)
	assert.Equal(t, testFile, path)

	// Should fail on missing file
	_, err = fs.GetFile("nope.csv")
	assert.Error(t, err)
}

func TestFileService_ValidateFile_SizeCheck(t *testing.T) {
	tempDir, _ := os.MkdirTemp("", "file-test")
	defer os.RemoveAll(tempDir)

	fs := NewFileService(tempDir, tempDir+"-downloads")

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, _ := writer.CreateFormFile("file", "big.csv")
	part.Write(make([]byte, 2048)) // 2KB file
	writer.Close()

	req := bytes.NewReader(body.Bytes())
	reader := multipart.NewReader(req, writer.Boundary())
	form, _ := reader.ReadForm(4096)
	file := form.File["file"][0]

	// Should pass with 10KB limit
	info, err := fs.ValidateFile(file, 10*1024)
	assert.NoError(t, err)
	assert.Equal(t, "big.csv", info.Filename)

	// Should fail with 1KB limit
	_, err = fs.ValidateFile(file, 1024)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "exceeds maximum")
}

func TestFileService_ValidateFile_ExtensionCheck(t *testing.T) {
	tempDir, _ := os.MkdirTemp("", "file-test")
	defer os.RemoveAll(tempDir)

	fs := NewFileService(tempDir, tempDir+"-downloads")

	makeFile := func(name string) *multipart.FileHeader {
		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)
		part, _ := writer.CreateFormFile("file", name)
		part.Write([]byte("name,email\ntest,test@example.com"))
		writer.Close()

		req := bytes.NewReader(body.Bytes())
		reader := multipart.NewReader(req, writer.Boundary())
		form, _ := reader.ReadForm(1024)
		return form.File["file"][0]
	}

	// CSV should pass
	csvFile := makeFile("good.csv")
	_, err := fs.ValidateFile(csvFile, 10*1024)
	assert.NoError(t, err)

	// TXT should fail
	txtFile := makeFile("bad.txt")
	_, err = fs.ValidateFile(txtFile, 10*1024)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "CSV")
}

func TestFileService_DeleteFile(t *testing.T) {
	tempDir, _ := os.MkdirTemp("", "file-test")
	defer os.RemoveAll(tempDir)

	fs := NewFileService(tempDir, tempDir+"-downloads")

	testFile := filepath.Join(tempDir, "delete-me.csv")
	os.WriteFile(testFile, []byte("temp"), 0644)

	err := fs.DeleteFile("delete-me.csv")
	assert.NoError(t, err)
	assert.NoFileExists(t, testFile)

	// Deleting again should error
	err = fs.DeleteFile("delete-me.csv")
	assert.Error(t, err)
}

func TestIsTextFile(t *testing.T) {
	// Normal CSV content
	csvData := []byte("name,email,age\nJohn,Chirag@test.com,30")
	assert.True(t, isTextFile(csvData))

	// Binary data (lots of null bytes)
	binaryData := make([]byte, 100)
	binaryData[0] = 0xFF
	binaryData[1] = 0x00
	binaryData[2] = 0xFE
	assert.False(t, isTextFile(binaryData))

	// Empty data
	assert.False(t, isTextFile([]byte{}))

	// Mix of text and control chars
	mixedData := []byte("name,email\n\r\tJohn,test@example.com")
	assert.True(t, isTextFile(mixedData))
}
