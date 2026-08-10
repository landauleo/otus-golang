package hw10programoptimization

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"strings"
)

type User struct {
	ID       int
	Name     string
	Username string
	Email    string
	Phone    string
	Password string
	Address  string
}

type DomainStat map[string]int

func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {
	emails, err := getUsers(r)
	if err != nil {
		return nil, fmt.Errorf("get users error: %w", err)
	}
	return countDomains(emails, domain)
}

type emails [100_000]string //это массив!!!

// ВАЖНО: когда даешь имена возвращаемым значениям,
// язык создает их в начале вызова функции и инициализирует нулевыми значениями
func getUsers(r io.Reader) (result emails, err error) {
	scanner := bufio.NewScanner(r)
	var i = 0
	var user User

	for scanner.Scan() {
		if err := json.Unmarshal(scanner.Bytes(), &user); err != nil {
			continue
		}
		result[i] = user.Email
		i++
	}
	return
}

func countDomains(e emails, domain string) (DomainStat, error) {
	result := make(DomainStat)

	for _, email := range e {
		//at - это то, как символ @ читается на ангельском
		atIndex := strings.Index(email, "@")
		if atIndex == -1 {
			continue
		}
		domainIndex := strings.LastIndex(email, domain)
		if domainIndex == -1 || domainIndex <= atIndex {
			continue
		}

		key := strings.ToLower(strings.TrimSuffix(email[atIndex+1:], "."))

		num := result[key]
		num++
		result[key] = num
	}
	return result, nil
}
