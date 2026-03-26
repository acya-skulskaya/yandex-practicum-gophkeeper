package secret_data

import (
	"fmt"
	"os"

	"github.com/acya-skulskaya/yandex-practicum-gophkeeper/internal/config"
)

type SecretTypeBinaryData struct {
	OriginalFilePath string `json:"original_file_path"`
	FileName         string `json:"file_name"`
	Metadata         string `json:"metadata"`
	FileSize         int64  `json:"file_size"`
}

func (b *SecretTypeBinaryData) PrintData() error {
	fmt.Printf("ORIGINAL FILE PATH: %s\n", b.OriginalFilePath)
	fmt.Printf("FILE NAME: %s\n", b.FileName)
	fmt.Printf("FILE SIZE: %d\n", b.FileSize)
	fmt.Printf("METADATA: %s\n", b.Metadata)
	return nil
}

func (b *SecretTypeBinaryData) Fill(requireNotEmpty bool) (string, error) {
	var fileInfo os.FileInfo
	var metadata, binaryFilePath string

	if metadata == "" {
		fmt.Print("enter new metadata or leave empty: ")
		//nolint:errcheck,gosec // can be left empty
		fmt.Scanln(&metadata)
	}

	for binaryFilePath == "" {
		fmt.Print("enter new path to the binary file: ")
		_, err := fmt.Scanln(&binaryFilePath)
		if err != nil {
			fmt.Printf("error: could not read path: %s\n", err.Error())
		}
		fileInfo, err = os.Stat(binaryFilePath)
		if err != nil {
			fmt.Printf("error: could not read path: %s\n", err.Error())
			binaryFilePath = ""
		} else if fileInfo.IsDir() {
			fmt.Printf("error: binary file path can't be a directory\n")
			binaryFilePath = ""
		}
	}

	if fileInfo.Size() > config.MaxBinaryFileSize {
		return "", fmt.Errorf("error: binary file size is too large (max size is %dMB)", config.MaxBinaryFileSize/1024/1024)
	}

	b.OriginalFilePath = binaryFilePath
	b.FileName = fileInfo.Name()
	b.FileSize = fileInfo.Size()

	if metadata != "" {
		b.Metadata = metadata
	}

	return binaryFilePath, nil
}
