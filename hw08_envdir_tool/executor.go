package main

import (
	"errors"
	"os"
	"os/exec"
	"strings"
)

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmd []string, env Environment) (returnCode int) {
	if len(cmd) == 0 {
		return 1 // Нет команды для запуска
	}

	//создаем команду
	command := exec.Command(cmd[0], cmd[1:]...)

	//если это не сделать, логи будут писаться в stdout самой программы, в терминале мы их не увидим
	command.Stdin = os.Stdin
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr

	//нельзя не "копировать" остальное окружение, если брать только то, что пришло из аргумента функции,
	//то у нас не будет ли текущей директории, ни кодировки и прочее
	rawEnv := os.Environ()

	tempMap := make(map[string]string, len(rawEnv)+len(env))

	for _, kv := range rawEnv {
		k, v, _ := strings.Cut(kv, "=")
		tempMap[k] = v
	}

	//добавляем/удаляем из tempMap переменные из аргументов
	for k, v := range env {
		if v.NeedRemove {
			delete(tempMap, k)
		} else {
			tempMap[k] = v.Value
		}
	}

	//обратно всё склеиваем
	finalEnv := make([]string, 0, len(tempMap))
	for k, v := range tempMap {
		finalEnv = append(finalEnv, k+"="+v)
	}

	//отдаем изолированное окружение команде
	command.Env = finalEnv

	//запуск
	err := command.Run()
	if err != nil {
		var exitError *exec.ExitError
		//ищем первую попавшуюся системную ошибку
		if errors.As(err, &exitError) {
			//код выхода утилиты должен совпадать с кодом выхода программы
			return exitError.ExitCode()
		}
		return 1
	}

	return 0
}
