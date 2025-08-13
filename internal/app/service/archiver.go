package service

import (
	"archive/zip"
	"fmt"
	"os"
	"path/filepath"
)

// SaveFilesAsZip создает zip-архив из списка файлов
func SaveFilesAsZip(files []DownloadedFile, archivePath string) error {
	f, err := os.Create(archivePath)
	if err != nil {
		return fmt.Errorf("cannot create archive: %w", err)
	}
	defer f.Close()

	zipWriter := zip.NewWriter(f)
	defer zipWriter.Close()

	for _, file := range files {
		if len(file.Data) == 0 {
			continue
		}

		if err := addFileToZip(zipWriter, file); err != nil {
			return fmt.Errorf("cannot add file %s to zip: %w", file.Name, err)
		}
	}

	return nil
}

// addFileToZip добавляет один файл в zip-архив
func addFileToZip(zipWriter *zip.Writer, file DownloadedFile) error {
	w, err := zipWriter.Create(filepath.Base(file.Name))
	if err != nil {
		return fmt.Errorf("cannot create entry in zip: %w", err)
	}

	_, err = w.Write(file.Data)
	if err != nil {
		return fmt.Errorf("cannot write file data: %w", err)
	}

	return nil
}