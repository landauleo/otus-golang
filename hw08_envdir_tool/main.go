package main

import "os"

// проверка:
// ./go-envdir ./my_env_dir env
// env просто печатает в консоль все переменные окружения, которые ей доступны)
// в ./my_env_dir  разные файлики с корнер-кейсами из ТЗ
func main() {
	// os.Args[0] — имя программы
	// os.Args[1] — путь к папке
	// os.Args[2] — сама команда
	if len(os.Args) < 3 {
		os.Exit(1)
	}

	dir := os.Args[1]
	commandWithArgs := os.Args[2:]

	env, err := ReadDir(dir)
	if err != nil {
		os.Exit(1)
	}

	exitCode := RunCmd(commandWithArgs, env)

	// 3. Завершаем утилиту с тем же кодом
	os.Exit(exitCode)
}
