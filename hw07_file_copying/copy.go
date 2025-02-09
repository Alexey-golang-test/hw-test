package main

import (
	"errors"
	"io"
	"os"

	"github.com/cheggaaa/pb/v3"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
)

const bufferSize int = 10

func Copy(fromPath, toPath string, offset, limit int64) error {
	// отступ в источнике, по умолчанию - 0
	if offset < 0 {
		offset = 0
	}

	// количество копируемых байт, по умолчанию - 0 (весь файл)
	if limit < 0 {
		limit = 0
	}

	// Проверка на zero-value
	if fromPath == "" || toPath == "" {
		return ErrUnsupportedFile
	}

	// Размер файла для проверок на ошибки
	info, err := os.Stat(fromPath)
	if err != nil {
		return err
	}
	fileSize := info.Size()

	// offset больше, чем размер файла - невалидная ситуация
	if offset > fileSize {
		return ErrOffsetExceedsFileSize
	}

	// программа может НЕ обрабатывать файлы, у которых неизвестна длина
	if fileSize <= 0 {
		return ErrUnsupportedFile
	}

	// Коррекция limit, чтобы точнее определить кол-во байтов для progress bar
	if limit == 0 {
		limit = fileSize

		if offset > 0 {
			limit = limit - offset + 1
		}
	}

	// необходимо выводить в консоль прогресс копирования в процентах (%)
	count := limit / int64(bufferSize)
	if limit > count*int64(bufferSize) {
		count++
	}
	bar := pb.Start64(count)

	srcFile, err := os.Open(fromPath)
	if err != nil {
		return err
	}
	defer srcFile.Close()
	buffer := make([]byte, bufferSize)

	destFile, err := os.Create(toPath)
	if err != nil {
		return err
	}
	defer destFile.Close()

	for ii := int64(0); ii < count; ii++ {
		nn, err := srcFile.ReadAt(buffer, offset+ii*int64(bufferSize))
		if !errors.Is(err, io.EOF) && err != nil {
			return err
		}

		_, err = destFile.Write(buffer[:nn])
		if err != nil {
			return err
		}

		bar.Increment()
	}

	bar.Finish()

	return nil
}
