package wasmtooci

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"io"
	"os"
	"path/filepath"
)

//TODO: detlete old
func StoreFileAsCAS(src, dst string) (string, int64, error) {
	// Open the source file
	file, err := os.Open(src)
	if err != nil {
		return "", 0, err
	}
	defer file.Close()

	// Create a SHA-256 hash
	hash := sha256.New()

	// Copy the file contents to the hash
	if _, err := io.Copy(hash, file); err != nil {
		return "", 0, err
	}

	// Compute the hash string
	hashString := hex.EncodeToString(hash.Sum(nil))

	// Create the destination file path
	casFilePath := filepath.Join(dst, hashString)
	casFile, err := os.Create(casFilePath)
	if err != nil {
		return "", 0, err
	}
	defer casFile.Close()

	// Reset the file pointer to the beginning of the source file
	if _, err := file.Seek(0, 0); err != nil {
		return "", 0, err
	}

	// Copy the file contents to the destination file
	size, err := io.Copy(casFile, file)
	if err != nil {
		return "", 0, err
	}

	_, fileName := filepath.Split(casFile.Name())

	return fileName, size, nil
}

func StoreAsCAS(src io.Reader, dst string) (string, int64, error) {
	// Create a buffer to store the entire JSON data
	var buf bytes.Buffer
	tee := io.TeeReader(src, &buf)

	// Create a SHA-256 hash
	hash := sha256.New()
	if _, err := io.Copy(hash, tee); err != nil {
		return "", 0, err
	}

	// Compute the hash string
	hashString := hex.EncodeToString(hash.Sum(nil))

	// Create the destination file path
	casFilePath := filepath.Join(dst, hashString)
	casFile, err := os.Create(casFilePath)
	if err != nil {
		return "", 0, err
	}
	defer casFile.Close()

	// Write the buffered JSON data to the destination file
	size, err := io.Copy(casFile, &buf)
	if err != nil {
		return "", 0, err
	}

	_, fileName := filepath.Split(casFile.Name())

	return fileName, size, nil
}

func MarshalToBuffer(data interface{}) (*bytes.Buffer, error) {
	jsonData,err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	return bytes.NewBuffer(jsonData), nil
}
