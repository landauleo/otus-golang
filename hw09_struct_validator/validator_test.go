package hw09structvalidator

import (
	"encoding/json"
	"errors"
	"fmt"
	"testing"
)

type UserRole string

// Test the function on different structures and other types.
type (
	User struct {
		ID     string `json:"id" validate:"len:36"` //всё, что с обратными кавычками - Struct Tags
		Name   string
		Age    int      `validate:"min:18|max:50"`
		Email  string   `validate:"regexp:^\\w+@\\w+\\.\\w+$"`
		Role   UserRole `validate:"in:admin,stuff"`
		Phones []string `validate:"len:11"`
		meta   json.RawMessage
	}

	App struct {
		Version string `validate:"len:5"`
	}

	Token struct {
		Header    []byte
		Payload   []byte
		Signature []byte
	}

	Response struct {
		Code int    `validate:"in:200,404,500"`
		Body string `json:"omitempty"`
	}
)

func TestValidate(t *testing.T) { //главный тест-родитель (внешняя t) -> через это Go сообщает ОС, упал ли весь тест
	tests := []struct {
		in          interface{}
		expectedErr error
	}{
		{
			in: User{
				ID:     "123456789012345678901234567890123456",
				Name:   "Puk",
				Age:    20,
				Email:  "test@example.com",
				Role:   "admin",
				Phones: []string{"79991112233", "79992223344"},
			},
			expectedErr: nil,
		},
		{
			in: App{
				Version: "1",
			},
			expectedErr: errors.New("error"),
		},
		{
			in: Response{
				Code: 200,
			},
			expectedErr: nil,
		},
		{
			in:          Token{Header: []byte("h"), Payload: []byte("p")},
			expectedErr: nil,
		},
	}

	for i, testCase := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) { //тест-ребенок (внутренняя t) -> изолированный пульт-управления
			tc := testCase
			t.Parallel()
			err := Validate(tc.in)

			if tc.expectedErr == nil && err != nil {
				t.Errorf("expected no error, got: %v", err)
			}

			if tc.expectedErr != nil && err == nil {
				t.Errorf("expected error, got nil")
			}
		})
	}
}
