package main

import (
	"errors"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCopy(tt *testing.T) {
	tt.Run("ошибка вида ErrUnsupportedFile", func(tt *testing.T) {
		err := Copy("", "test.txt", 0, 0)
		require.True(tt, errors.Is(err, ErrUnsupportedFile), "Текущая ошибка: ", ErrUnsupportedFile)

		file, errCreate := os.Create("empty.txt")
		file.Close()
		errCopy := Copy("empty.txt", "test.txt", 0, 0)
		errDel := os.Remove("empty.txt")
		require.NoError(tt, errCreate)
		require.NoError(tt, errDel)
		require.True(tt, errors.Is(errCopy, ErrUnsupportedFile), "Текущая ошибка: ", ErrUnsupportedFile)
	})

	tt.Run("ошибка вида ErrOffsetExceedsFileSize", func(tt *testing.T) {
		file, err := os.Stat("testdata/input.txt")
		require.NoError(tt, err)

		err = Copy("testdata/input.txt", "test.txt", file.Size()+1, 0)
		require.True(tt, errors.Is(err, ErrOffsetExceedsFileSize), "Текущая ошибка: ", ErrOffsetExceedsFileSize)
	})

	tt.Run("полное копирование с отрицательными параметрами", func(tt *testing.T) {
		err := Copy("testdata/input.txt", "test.txt", -1, -5)
		require.NoError(tt, err)

		fileSrc, errSrc := os.Stat("testdata/input.txt")
		require.NoError(tt, errSrc)

		fileDest, errDest := os.Stat("test.txt")
		require.NoError(tt, errDest)

		require.Equal(tt, fileSrc.Size(), fileDest.Size())

		errDel := os.Remove("test.txt")
		require.NoError(tt, errDel)
	})

	tt.Run("копирование части файла", func(tt *testing.T) {
		err := Copy("testdata/input.txt", "test.txt", 10, 100)
		require.NoError(tt, err)

		fileDest, errDest := os.Stat("test.txt")
		require.NoError(tt, errDest)

		require.Equal(tt, int64(100), fileDest.Size())

		errDel := os.Remove("test.txt")
		require.NoError(tt, errDel)
	})
}
