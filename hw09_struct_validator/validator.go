package hw09structvalidator

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

type ValidationError struct {
	Field string
	Err   error
}

type ValidationErrors []ValidationError

var (
	ErrLen    = errors.New("invalid length")
	ErrRegexp = errors.New("invalid format")
	ErrIn     = errors.New("value not in list")
	ErrMin    = errors.New("value too small")
	ErrMax    = errors.New("value too big")
)

func (v ValidationErrors) Error() string {
	var res []string
	for _, e := range v {
		res = append(res, fmt.Sprintf("%s: %v", e.Field, e.Err))
	}
	return strings.Join(res, ", ")
}

func Validate(v interface{}) error {
	val := reflect.ValueOf(v)
	if val.Kind() == reflect.Ptr { //-> pointer
		val = val.Elem() // разыменовываем указатель, если передали &
	}

	var errs ValidationErrors

	for i := 0; i < val.NumField(); i++ {
		fieldVal := val.Field(i)
		fieldType := val.Type().Field(i) //получаем тип поля (и его теги)

		tag := fieldType.Tag.Get("validate")
		if tag == "" {
			continue
		}

		//разбиваем тег на "правило" и "аргумент" -> "min" и "10"
		parts := strings.SplitN(tag, ":", 2)
		rule, param := parts[0], parts[1]

		check := func(fv reflect.Value) error {
			switch fv.Kind() {
			case reflect.String:
				s := fv.String()
				switch rule {
				case "len":
					var res, _ = strconv.Atoi(param)
					if len(s) != res {
						return ErrLen
					}
				case "regexp":
					if !regexp.MustCompile(param).MatchString(s) {
						return ErrRegexp
					}
				case "in":
					if !isIn(s, strings.Split(param, ",")) {
						return ErrIn
					}
				}
			case reflect.Int:
				n := int(fv.Int())
				switch rule {
				case "min":
					var res, _ = strconv.Atoi(param)
					if n < res {
						return ErrMin
					}
				case "max":
					var res, _ = strconv.Atoi(param)
					if n > res {
						return ErrMax
					}
				case "in":
					if !isIn(strconv.Itoa(n), strings.Split(param, ",")) {
						return ErrIn
					}
				}
			} //свои кейсы реализовывать не стала
			return nil
		}

		//если это слайс — проверяем каждый элемент, иначе — само поле
		if fieldVal.Kind() == reflect.Slice {
			for j := 0; j < fieldVal.Len(); j++ {
				if err := check(fieldVal.Index(j)); err != nil {
					errs = append(errs, ValidationError{Field: fieldType.Name, Err: err})
					break // Достаточно одной ошибки на поле-слайс
				}
			}
		} else {
			if err := check(fieldVal); err != nil {
				errs = append(errs, ValidationError{Field: fieldType.Name, Err: err})
			}
		}
	}

	if len(errs) > 0 {
		return errs
	}
	return nil
}

// условие: строка/число должно входить в множество строк/чисел
func isIn(target string, list []string) bool {
	for _, item := range list {
		if item == target {
			return true
		}
	}
	return false
}
