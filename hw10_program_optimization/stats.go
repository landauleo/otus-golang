package hw10programoptimization

import (
	"bufio"
	"encoding/json"
	"io"
	"strings"
)

type User struct {
	Email string `json:"email"` //вот это имба
}

type DomainStat map[string]int

func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {
	return getEmails(r, domain)
}

// ВАЖНО: когда даешь имена возвращаемым значениям,
// язык создает их в начале вызова функции и инициализирует нулевыми значениями
func getEmails(r io.Reader, domain string) (DomainStat, error) {
	result := make(DomainStat)
	scanner := bufio.NewScanner(r)
	var user User

	for scanner.Scan() {
		if err := json.Unmarshal(scanner.Bytes(), &user); err != nil {
			continue
		}
		email := user.Email
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
