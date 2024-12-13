package utils

import (
	"fmt"
	"strings"

	"github.com/SawitProRecruitment/UserService/generated"
)

func IsRegisterUserPayloadValid(payload generated.RegisterUserJSONRequestBody) (bool, string) {
	isPayloadValid := true
	errorMessages := make([]string, 0)

	if isValid := IsStartWithCountryCode(payload.PhoneNumber, "+62"); !isValid {
		isPayloadValid = false
		errorMessages = append(errorMessages, "phone_number field must start with +62")
	}

	minPhoneLen, maxPhoneLen := 10, 13
	if isValid := IsLengthBetweenRange(payload.PhoneNumber, minPhoneLen, maxPhoneLen); !isValid {
		isPayloadValid = false
		errorMessages = append(errorMessages, fmt.Sprintf("phone_number must be between %d to %d characters long", minPhoneLen, maxPhoneLen))
	}

	minNameLen, maxNameLen := 3, 60
	if isValid := IsLengthBetweenRange(payload.FullName, minNameLen, maxNameLen); !isValid {
		isPayloadValid = false
		errorMessages = append(errorMessages, fmt.Sprintf("full_name must be between %d to %d characters long", minNameLen, maxNameLen))
	}

	minPasswordLen, maxPasswordLen := 6, 64
	if isValid := IsLengthBetweenRange(payload.Password, minPasswordLen, maxPasswordLen); !isValid {
		isPayloadValid = false
		errorMessages = append(errorMessages, fmt.Sprintf("password must be between %d to %d characters long", minPasswordLen, maxPasswordLen))
	}

	if isValid := ContainsUpperCase(payload.Password); !isValid {
		isPayloadValid = false
		errorMessages = append(errorMessages, "password must contain 1 upper case")
	}

	if isValid := ContainsNumber(payload.Password); !isValid {
		isPayloadValid = false
		errorMessages = append(errorMessages, "password must contain 1 number")
	}

	if isValid := ContainsSpecialCharacter(payload.Password); !isValid {
		isPayloadValid = false
		errorMessages = append(errorMessages, "password must contain 1 special character")
	}

	return isPayloadValid, strings.Join(errorMessages, ", ")
}
