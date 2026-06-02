package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestCopy(t *testing.T) {
	tempDir := t.TempDir() //удаляется после завершения теста

	testData := []struct {
		name         string
		offset       int64
		limit        int64
		expectedData string
		expectError  bool
	}{
		{
			name:         "Взять первые 3 байта",
			offset:       0,
			limit:        2,
			expectedData: "Go",
			expectError:  false,
		},
		{
			name:         "Пропустить 3 байта и взять 4",
			offset:       3,
			limit:        4,
			expectedData: "Docu",
			expectError:  false,
		},
		{
			name:         "Limit больше, чем размер файла",
			offset:       1000,
			limit:        100,
			expectedData: ", 386\tDebian GNU/kFreeBSD not supported\nLinux 2.6.23 or later with glibc\tamd64, 386, arm, arm64,\ns39",
			expectError:  false,
		},
	}

	for _, tc := range testData {
		t.Run(tc.name, func(t *testing.T) {
			dstPath := filepath.Join(tempDir, "out_"+tc.name+".txt")

			err := Copy("testdata/input.txt", dstPath, tc.offset, tc.limit)

			if (err != nil) != tc.expectError {
				t.Fatalf("Ожидали ошибку: %v, получили: %v", tc.expectError, err)
			}

			if tc.expectError {
				return
			}

			resultData, err := os.ReadFile(dstPath)
			if err != nil {
				t.Fatalf("Не удалось прочитать файл результата: %v", err)
			}

			if string(resultData) != tc.expectedData {
				t.Errorf("Результат не совпал!\nОжидали: %q\nПолучили: %q", tc.expectedData, string(resultData))
			}
		})
	}
}
