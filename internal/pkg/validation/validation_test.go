package validation

import (
	"fmt"
	"os"
	"testing"

	"github.com/nortoo/usms/internal/pkg/etc"
	"go.uber.org/zap"
)

var s Service

func TestMain(m *testing.M) {
	cfg, err := etc.Load("testdata/app.yml")
	if err != nil {
		fmt.Println("Failed to load config file", zap.Error(err))
		os.Exit(1)
	}
	s = New(cfg)

	m.Run()
}

func TestIsValidUsername(t *testing.T) {
	samples := map[string]bool{
		// Valid usernames
		"abcde":               true,
		"A1234":               true,
		"Z_987":               true,
		"z_987":               true,
		"user_name":           true,
		"A1_bc":               true,
		"a123456789012345678": true, // 20 chars
		"Z_1234567890123456":  true, // 20 chars
		"Alice_1":             true,
		"Bob_123":             true,
		"Charlie_2023":        true,
		"David_The_2nd":       true,
		"Eve_007":             true,
		"Franklin_12345":      true,
		"Grace_Hopper":        true,
		"Helen_2022":          true,
		"Isaac_Newton":        true,
		"JackieChan_1":        true,
		"K_L_M_N_O":           true,
		"Liam_underscore":     true,
		"Mike123456789012":    true,
		"Nina_1234567890":     true,
		"Oscar_The_Best":      true,
		"Paul_123456789":      true,
		"QueenBee_2020":       true,
		"Robert_2021":         true,
		"Steve_12345":         true,
		"Tommy_underscore":    true,

		// Invalid usernames
		"A123":                  false, // too short
		"Z_9":                   false, // too short
		"A1_":                   false, // too short
		"1abc":                  false, // starts with digit
		"_abc":                  false, // starts with underscore
		"ab":                    false, // too short
		"a":                     false, // too short
		"a12345678901234567890": false, // 21 chars, too long
		"a-bc":                  false, // invalid character '-'
		"a.bc":                  false, // invalid character '.'
		"a bc":                  false, // contains space
		"":                      false, // empty
		"A!bc":                  false, // invalid character '!'
	}

	for username, expected := range samples {
		allowed, _ := s.IsValidUsername(username)
		if allowed != expected {
			t.Errorf("Username: %s, Expected: %v, Got: %v", username, expected, allowed)
		}
	}
}

func TestIsValidPassword(t *testing.T) {
	samples := map[string]bool{
		// Valid passwords
		"Password1!":   true, // 10 chars, has uppercase, lowercase, digit, special char
		"MyPass#123":   true, // 10 chars, has uppercase, lowercase, digit, special char
		"Secure$Pass1": true, // 12 chars, has uppercase, lowercase, digit, special char
		"Test%123":     true, // 8 chars, has uppercase, lowercase, digit, special char
		"Valid^Pass9":  true, // 11 chars, has uppercase, lowercase, digit, special char
		"valid^Pass9":  true, // 11 chars, has uppercase, lowercase, digit, special char

		// Valid password - contains spaces
		"Password 1!": true, // 11 chars, has space
		"Test 123@":   true, // 9 chars, has space

		// Invalid passwords - too short
		"Abc123@": false, // 7 chars, has uppercase, lowercase, digit, special char
		"Pass1!":  false, // 6 chars (less than 8)
		"Abc1@":   false, // 5 chars (less than 8)
		"Test#":   false, // 5 chars (less than 8)

		// Invalid passwords - missing lowercase
		"PASSWORD1!": false, // 10 chars, no lowercase
		"TEST123@":   false, // 8 chars, no lowercase
		"SECURE#123": false, // 10 chars, no lowercase

		// Invalid passwords - missing uppercase
		"password1!": false, // 10 chars, no uppercase
		"test123@":   false, // 8 chars, no uppercase
		"secure#123": false, // 10 chars, no uppercase

		// Invalid passwords - missing digit
		"Password!":   false, // 9 chars, no digit
		"TestPass@":   false, // 9 chars, no digit
		"SecurePass#": false, // 12 chars, no digit

		// Invalid passwords - missing special character
		"Password1":   false, // 9 chars, no special char
		"TestPass123": false, // 11 chars, no special char
		"SecurePass1": false, // 12 chars, no special char

		// Invalid passwords - empty or whitespace only
		"":     false, // empty string
		"   ":  false, // whitespace only
		"\t\n": false, // whitespace characters

		// Invalid passwords - only digits and special chars
		"12345678!":  false, // 9 chars, only digits and special char
		"123456789@": false, // 10 chars, only digits and special char

		// Invalid passwords - only letters and special chars
		"OnlyLetters!": false, // 12 chars, only letters and special char (no digit)
		"JustWords@":   false, // 10 chars, only letters and special char (no digit)

		// Invalid passwords - mixed but missing requirements
		"pass123!":  false, // 8 chars, missing uppercase
		"PASS123@":  false, // 8 chars, missing lowercase
		"Passwor#":  false, // 8 chars, missing digit
		"TestPass1": false, // 9 chars, missing special char
	}

	for password, expected := range samples {
		allowed, _ := s.IsValidPassword(password)
		if allowed != expected {
			t.Errorf("Password: %s, Expected: %v, Got: %v", password, expected, allowed)
		}
	}
}
