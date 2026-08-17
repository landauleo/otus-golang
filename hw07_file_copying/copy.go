package main

import (
	"errors"
	"fmt"
	"io"
	"os"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
)

func Copy(fromPath, toPath string, offset, limit int64) error {
	//чтобы использовать offset/limit нужно сначала открыть файл
	file, err := os.Open(fromPath)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err) // TODO: в го обычно мы врапаем ошибки,
		// чтобы они содержали больше контекста.
		// Например тут errors.Wrap(err, "failed to open file"). В остальных местах по аналогии
	}

	defer file.Close()

	if offset > 0 {
		fileInfo, fileInfoErr := os.Stat(fromPath)
		if fileInfoErr != nil {
			return err
		}

		if fileInfo.Size() < offset {
			return ErrOffsetExceedsFileSize
		}

		//программа может НЕ обрабатывать файлы, у которых неизвестна длина (например, /dev/urandom)
		if !fileInfo.Mode().IsRegular() {
			return ErrUnsupportedFile
		}

		_, err := file.Seek(offset, io.SeekStart)
		if err != nil {
			return err
		}
	}

	resultFile, err := os.Create(toPath)
	if err != nil {
		return err
	}

	defer resultFile.Close()

	var reader io.Reader
	if limit > 0 {
		reader = io.LimitReader(file, limit)
	} else {
		reader = file
	}

	_, err = io.Copy(resultFile, reader)
	if err != nil {
		return err
	}

	return nil
}
