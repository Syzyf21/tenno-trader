package internal

import (
	"bytes"
	"compress/zlib"
	"fmt"
	"io"
	"os"
)

func DecryptAlecaFrame(filePath string) ([]byte, error) {
	encryptedData, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	const xorKey byte = 0x5A
	decryptedBytes := make([]byte, len(encryptedData))
	for i := 0; i < len(encryptedData); i++ {
		decryptedBytes[i] = encryptedData[i] ^ xorKey
	}

	b := bytes.NewReader(decryptedBytes)
	zlibReader, err := zlib.NewReader(b)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize zlib reader: %w", err)
	}
	defer zlibReader.Close()

	var resultBuffer bytes.Buffer
	_, err = io.Copy(&resultBuffer, zlibReader)
	if err != nil {
		return nil, fmt.Errorf("failed during decompression: %w", err)
	}

	return resultBuffer.Bytes(), nil
}

func main() {
	filePath := "path/to/lastData.dat"

	jsonData, err := DecryptAlecaFrame(filePath)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Println(string(jsonData))
}
