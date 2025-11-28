package objects

import (
	"bytes"
	"compress/zlib"
	"crypto/sha1"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

const RootDir = ".gogit"

func SaveObject(objType string, data []byte) (string, error) {
	header := fmt.Sprintf("%s %d\x00", objType, len(data))

	storeContent := append([]byte(header), data...)
	hash := fmt.Sprintf("%x", sha1.Sum(storeContent))

	dirPath := filepath.Join(RootDir, "objects", hash[:2])
	filePath := filepath.Join(dirPath, hash[2:])

	if _, err := os.Stat(filePath); err == nil {
		return hash, nil
	}

	if err := os.MkdirAll(dirPath, 0755); err != nil {
		return "", err
	}

	f, err := os.Create(filePath)
	if err != nil {
		return "", err
	}
	defer f.Close()

	zw := zlib.NewWriter(f)
	if _, err := zw.Write(storeContent); err != nil {
		zw.Close()
		return "", err
	}
	zw.Close()

	return hash, nil
}

func ReadObject(hash string) (string, []byte, error) {
	if len(hash) < 4 {
		return "", nil, fmt.Errorf("короткий хеш")
	}
	path := filepath.Join(RootDir, "objects", hash[:2], hash[2:])

	f, err := os.Open(path)
	if err != nil {
		return "", nil, err
	}
	defer f.Close()

	zr, err := zlib.NewReader(f)
	if err != nil {
		return "", nil, err
	}
	defer zr.Close()

	content, err := io.ReadAll(zr)
	if err != nil {
		return "", nil, err
	}

	nullByteIdx := bytes.IndexByte(content, 0)
	if nullByteIdx == -1 {
		return "", nil, fmt.Errorf("неверный формат объекта")
	}

	header := string(content[:nullByteIdx])
	var objType string
	var size int
	fmt.Sscanf(header, "%s %d", &objType, &size)

	data := content[nullByteIdx+1:]
	return objType, data, nil
}

func HashFile(path string) (string, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	header := fmt.Sprintf("blob %d\x00", len(content))
	storeContent := append([]byte(header), content...)
	return fmt.Sprintf("%x", sha1.Sum(storeContent)), nil
}
