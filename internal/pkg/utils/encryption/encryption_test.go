package encryption

import "testing"

func TestEncryptPassword(t *testing.T) {
	samples := []string{
		// Simple passwords
		"password123",
		"123456789",
		"qwerty",
		"abc123",

		// Complex passwords with special characters
		"P@ssw0rd!",
		"Str0ng#P@ss",
		"Secure$Pass1",
		"MyP@ssw0rd!",
		"T3st#P@ss",

		// Long passwords
		"ThisIsAVeryLongPassword123!@#",
		"SuperSecurePasswordWithManyCharacters2023!",
		"AVeryLongPasswordThatExceedsNormalLength123456789!@#$%",

		// Short passwords
		"a1!",
		"b2@",
		"c3#",

		// Passwords with spaces
		"My Password 123",
		"Secure Pass Word",
		"Test Password With Spaces",

		// Passwords with unicode characters
		"P@ssw0rd中文",
		"MotDePasse123!",
		"Contraseña2023",
		"パスワード123!",

		// Edge cases
		"",           // empty string
		" ",          // single space
		"   ",        // multiple spaces
		"\t\n\r",     // whitespace characters
		"!@#$%^&*()", // only special characters
		"1234567890", // only digits
		"abcdefghij", // only lowercase
		"ABCDEFGHIJ", // only uppercase

		// Common weak passwords
		"password",
		"admin",
		"root",
		"user",
		"guest",
		"test",

		// Passwords with repeated characters
		"aaaaaaaa",
		"11111111",
		"!!!!!!!!",
		"AAAAAAA1",

		// Passwords with patterns
		"12345678901234567890",
		"abcdefghijklmnopqrst",
		"ABCDEFGHIJKLMNOPQRST",
		"!@#$%^&*()!@#$%^&*()",
	}

	for _, s := range samples {
		encryptedPassword, err := EncryptPassword(s)
		if err != nil {
			t.Fatalf("failed to encrypt password: %v", err)
		}

		if !ComparePassword(encryptedPassword, s) {
			t.Fatal("failed to encrypt password, the encrypted password cannot be used to verify the origin password.")
		}
	}
}

func TestEncryptEmailAddress(t *testing.T) {
	validEmail := "user@example.com"
	encryptedEmail := EncryptEmailAddress(validEmail)
	if encryptedEmail != "u***@***.com" {
		t.Fatalf(`expected "u***@***.com", got %q`, encryptedEmail)
	}

	invalidEmailAddress := "@google.com"
	encryptedEmail = EncryptEmailAddress(invalidEmailAddress)
	if encryptedEmail != "" {
		t.Fatalf(`expected empty string, got %q`, encryptedEmail)
	}
}

func TestEncryptMobileNumber(t *testing.T) {
	validMobileNumber := "+12025550100"
	encryptedMobileNumber := EncryptMobileNumber(validMobileNumber)
	if encryptedMobileNumber != "+12****0100" {
		t.Fatalf(`expected "+120****0100", got %q`, encryptedMobileNumber)
	}

	invalidMobileNumber := "+15550101"
	encryptedMobileNumber = EncryptMobileNumber(invalidMobileNumber)
	if encryptedMobileNumber != "" {
		t.Fatalf(`expected empty string, got %q`, encryptedMobileNumber)
	}
}
