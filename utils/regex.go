package utils

import (
	"fmt"
	"regexp"
)

func IsStartWithCountryCode(str string, countryCode string) bool {
	regexSyntax := fmt.Sprintf(`^\%s`, countryCode)
	regex := regexp.MustCompile(regexSyntax)
	return regex.MatchString(str)
}

func IsLengthBetweenRange(str string, min, max int) bool {
	regexSyntax := fmt.Sprintf(`^.{%d,%d}$`, min, max)
	regex := regexp.MustCompile(regexSyntax)
	return regex.MatchString(str)
}

func ContainsUpperCase(str string) bool {
	uppercaseRegex := regexp.MustCompile(`[A-Z]`)
	return uppercaseRegex.MatchString(str)
}

func ContainsNumber(str string) bool {
	numberRegex := regexp.MustCompile(`\d`)
	return numberRegex.MatchString(str)
}

func ContainsSpecialCharacter(str string) bool {
	specialCharRegex := regexp.MustCompile(`[^a-zA-Z\d]`)
	return specialCharRegex.MatchString(str)
}
