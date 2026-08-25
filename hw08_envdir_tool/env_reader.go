package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type Environment map[string]EnvValue

// EnvValue helps to distinguish between empty files and files with the first empty line.
type EnvValue struct {
	Value      string
	NeedRemove bool
}

// ReadDir reads a specified directory and returns map of env variables.
// Variables represented as files where filename is name of variable, file first line is a value.
func ReadDir(dir string) (Environment, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("failed to read env dir %q: %w", dir, err)
	}

	env := make(Environment, len(entries))

	//проходимся по всему содержанию директории
	for _, entry := range entries {
		//нужны только файлы (не папки/симлинки и что-то подобное)
		if !entry.Type().IsRegular() {
			continue
		}

		name := entry.Name()
		if strings.Contains(name, "=") {
			//имя S не должно содержать '='
			continue
		}

		info, err := entry.Info()
		if err != nil {
			return nil, fmt.Errorf("failed to get file info %q: %w", filepath.Join(dir, name), err)
		}

		if info.Size() == 0 {
			//если файл пустой, то переменной нет, записываем её в словарь и помечаем на удаление
			env[name] = EnvValue{NeedRemove: true}
			continue
		}

		//читаем содержимое файла целиком (так как файлы конфигурации обычно маленькие)
		content, err := os.ReadFile(filepath.Join(dir, name))
		if err != nil {
			return nil, fmt.Errorf("failed to read env file %q: %w", filepath.Join(dir, name), err)
		}

		//читаем первую строку
		firstLine, _, _ := strings.Cut(string(content), "\n")

		//пробелы и табуляция в конце T удаляются
		firstLine = strings.TrimRight(firstLine, " \t\r")

		//терминальные нули в T (0x00) заменяются на перевод строки (\n);
		firstLine = strings.ReplaceAll(firstLine, "\x00", "\n")

		env[name] = EnvValue{
			Value:      firstLine,
			NeedRemove: false,
		}
	}

	return env, nil
}
