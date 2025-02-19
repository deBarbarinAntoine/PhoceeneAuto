package validator

import (
	"encoding/json"
	"fmt"
	"regexp"
	"slices"
	"strings"
	"time"
	"unicode"
)

const (
	MinPasswordLength = 8
	MaxPasswordLength = 72
)

// Validator represents a validation instance with non-field errors and field-specific errors.
type Validator struct {
	NonFieldErrors []string          `json:"non_field_errors"`
	FieldErrors    map[string]string `json:"field_errors"`
}

var (
	EmailRX     = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	PrintableRX = regexp.MustCompile("^[[:print:]]+$")
	FileRX      = regexp.MustCompile(`^[^\0/\\\s]+$`)
)

// New creates a new Validator instance.
//
// Returns:
//
//	*Validator - A pointer to the newly created Validator
func New() *Validator {
	return &Validator{
		FieldErrors: make(map[string]string),
	}
}

// Valid checks if there are any non-field or field-specific errors.
//
// Returns:
//
//	bool - True if there are no errors, false otherwise
func (v *Validator) Valid() bool {
	return len(v.NonFieldErrors) == 0 && len(v.FieldErrors) == 0
}

// AddNonFieldError adds a non-field error message to the Validator.
//
// Parameters:
//
//	message - The error message to add
func (v *Validator) AddNonFieldError(message string) {
	v.NonFieldErrors = append(v.NonFieldErrors, message)
}

// AddFieldError adds a field-specific error message to the Validator.
//
// Parameters:
//
//	key - The field name
//	message - The error message to add
func (v *Validator) AddFieldError(key, message string) {
	if _, exists := v.FieldErrors[key]; !exists {
		v.FieldErrors[key] = message
	}
}

// Errors returns the validation errors as a JSON-encoded byte slice.
//
// Returns:
//
//	[]byte - The JSON-encoded byte slice of errors
func (v *Validator) Errors() []byte {
	errors, err := json.Marshal(v)
	if err != nil {
		panic(err) // FIXME -> maybe change that to handle errors better if necessary
	}
	return errors
}

// Check adds a field-specific error message if the condition is false.
//
// Parameters:
//
//	ok - The condition to check
//	key - The field name
//	message - The error message to add if the condition is false
func (v *Validator) Check(ok bool, key, message string) {
	if !ok {
		v.AddFieldError(key, message)
	}
}

// ValidateDate checks if a date string is in the correct format.
//
// Parameters:
//
//	date - The date string to validate
//	fieldName - The field name for error messages
func (v *Validator) ValidateDate(date, fieldName string) {
	_, err := time.Parse("01/02/2006", date)
	if err != nil {
		v.AddFieldError(fieldName, "invalid date")
	}
}

// CheckID checks if an integer ID is greater than 0.
//
// Parameters:
//
//	id - The integer ID to check
//	fieldName - The field name for error messages
func (v *Validator) CheckID(id int, fieldName string) {
	v.Check(id > 0, fieldName, "ID must be greater than 0")
}

// ValidateEmail checks if an email string is valid.
//
// Parameters:
//
//	email - The email string to validate
func (v *Validator) ValidateEmail(email string) {
	v.StringCheck(email, 5, 150, true, "email")
	v.Check(Matches(email, EmailRX), "email", "must be a valid email address")
}

// StringCheck checks if a string meets the specified length and mandatory requirements.
//
// Parameters:
//
//	str - The string to check
//	min - The minimum length of the string
//	max - The maximum length of the string
//	isMandatory - Whether the string must be provided
//	key - The field name for error messages
func (v *Validator) StringCheck(str string, min, max int, isMandatory bool, key string) {
	if isMandatory {
		v.Check(str != "", key, "must be provided")
	}
	v.Check(len(str) >= min, key, fmt.Sprintf("must be minimum %d bytes long", min))
	v.Check(len(str) <= max, key, fmt.Sprintf("must not be more than %d bytes long", max))
}

// CheckPassword checks if a password string meets the specified criteria.
//
// Parameters:
//
//	password - The password string to check
//	key - The field name for error messages
func (v *Validator) CheckPassword(password, key string) {

	// setting booleans to check the criteria
	var (
		hasUpper   = false
		hasLower   = false
		hasNumber  = false
		hasSpecial = false
	)

	// checking every character
	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsNumber(char):
			hasNumber = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}

	// adding errors if needed
	v.Check(hasUpper, key, "must contain an uppercase character")
	v.Check(hasLower, key, "must contain a lowercase character")
	v.Check(hasNumber, key, "must contain a numeric character")
	v.Check(hasSpecial, key, "must contain a special character")
}

// ValidatePassword checks if a password string is valid.
//
// Parameters:
//
//	password - The password string to validate
func (v *Validator) ValidatePassword(password string) {
	v.StringCheck(password, MinPasswordLength, MaxPasswordLength, true, "password")
	v.CheckPassword(password, "password")
}

// ValidateRegisterPassword checks if a registration password is valid.
//
// Parameters:
//
//	password - The password string to validate
//	confirmationPassword - The confirmation password string to validate
func (v *Validator) ValidateRegisterPassword(password, confirmationPassword string) {
	v.StringCheck(password, MinPasswordLength, MaxPasswordLength, true, "password")
	v.CheckPassword(password, "password")
	v.Check(confirmationPassword != "", "confirm_password", "must be provided")
	v.Check(password == confirmationPassword, "confirm_password", "mismatched passwords!")
}

// ValidateNewPassword checks if a new password is valid.
//
// Parameters:
//
//	newPassword - The new password string to validate
//	confirmationPassword - The confirmation password string to validate
func (v *Validator) ValidateNewPassword(newPassword, confirmationPassword string) {
	v.StringCheck(newPassword, MinPasswordLength, MaxPasswordLength, true, "new_password")
	v.CheckPassword(newPassword, "new_password")
	v.Check(confirmationPassword != "", "confirm_password", "must be provided")
	v.Check(newPassword == confirmationPassword, "confirm_password", "mismatched passwords!")
}

// ValidateToken checks if a token string is valid.
//
// Parameters:
//
//	token - The token string to validate
func (v *Validator) ValidateToken(token string) {
	v.Check(token != "", "token", "you need a special link to access here")
	v.Check(len(token) == 86, "token", "invalid link")
}

// CheckFileName checks if a filename is valid.
//
// Parameters:
//
//	filename - The filename to check
//
// Returns:
//
//	bool - True if the filename is valid, false otherwise
func CheckFileName(filename string) bool {
	if PrintableRX.MatchString(filename) {
		return FileRX.MatchString(filename)
	}
	return false
}

// NotBlank checks if a field name is not blank.
//
// Parameters:
//
//	fieldName - The field name to check
//
// Returns:
//
//	bool - True if the field name is not blank, false otherwise
func NotBlank(fieldName string) bool {
	return strings.TrimSpace(fieldName) != ""
}

// Matches checks if a value matches a regular expression.
//
// Parameters:
//
//	value - The value to check
//	rx - The regular expression to match against
//
// Returns:
//
//	bool - True if the value matches the regular expression, false otherwise
func Matches(value string, rx *regexp.Regexp) bool {
	return rx.MatchString(value)
}

// PermittedValue checks if a value is within a list of permitted values.
//
// Parameters:
//
//	value - The value to check
//	permittedValues - The list of permitted values
//
// Returns:
//
//	bool - True if the value is within the list of permitted values, false otherwise
func PermittedValue[T comparable](value T, permittedValues ...T) bool {
	return slices.Contains(permittedValues, value)
}

// Unique checks if a slice contains only unique values.
//
// Parameters:
//
//	values - The slice to check
//
// Returns:
//
//	bool - True if the slice contains only unique values, false otherwise
func Unique[T comparable](values []T) bool {
	uniqueValues := make(map[T]bool)

	for _, value := range values {
		uniqueValues[value] = true
	}

	return len(values) == len(uniqueValues)
}
