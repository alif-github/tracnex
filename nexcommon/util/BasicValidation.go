package util

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode"
)

func IsMacAddress(input string) (output bool) {
	temp := strings.Split(input, ":")
	if len(temp) == 6 {
		output = true
		for i := 0; i < len(temp); i++ {
			if len(temp[i]) != 2 {
				output = false
			}
		}
		return
	}
	return false
}

func IsEmailAddress(input string) (output bool) {
	emailRegexp := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	return emailRegexp.MatchString(input)
}

func IsPhoneNumber(input string) (number int, isValid bool) {
	number, isValid = IsNumeric(input)
	return number, len(input) >= 9 && len(input) <= 13 && isValid
}
func IsPhoneNumberWithCountryCode(input string) bool {
	phoneNumberRegexp := regexp.MustCompile("[+][0-9]+[-][1-9][0-9]{8,12}")
	return phoneNumberRegexp.MatchString(input)
}
func IsCountryCode(input string) bool {
	countryCodeRegexp := regexp.MustCompile("^[+][0-9]{1,3}$")
	return countryCodeRegexp.MatchString(input)
}
func IsIPPrivate(input string) (output bool) {
	ipPrivateRegexp := regexp.MustCompile("^https?://localhost.*")
	if ipPrivateRegexp.MatchString(input) {
		return true
	}
	ipPrivateRegexp = regexp.MustCompile("^https?://192.168..*")
	if ipPrivateRegexp.MatchString(input) {
		return true
	}

	ipPrivateRegexp = regexp.MustCompile("^https?://fd[0-9a-z]+:.*")
	return ipPrivateRegexp.MatchString(input)
}

func IsNexsoftPasswordStandardValid(password string) (bool, string, string) {
	if len(password) < 8 {
		return false, "NEED_MORE_THAN", "8"
	} else if len(password) >= 8 && len(password) <= 50 {
	next:
		for name, classes := range map[string][]*unicode.RangeTable{
			"UPPERCASE": {unicode.Upper, unicode.Title},
			"LOWERCASE": {unicode.Lower},
			"NUMERIC":   {unicode.Number, unicode.Digit},
			"SPECIAL":   {unicode.Space, unicode.Symbol, unicode.Punct, unicode.Mark},
		} {
			for _, r := range password {
				if unicode.IsOneOf(classes, r) {
					continue next
				}
			}
			return false, name, ""
		}
	} else {
		return false, "NEED_LESS_THAN", "50"
	}
	return true, "", ""
}

func IsNexsoftUsernameStandardValid(username string) (bool, string, string) {
	if len(username) < 6 {
		return false, "NEED_MORE_THAN", "6"
	} else if len(username) > 20 {
		return false, "NEED_LESS_THAN", "20"
	} else {
		usernameRegex := regexp.MustCompile("^[a-z][a-z0-9_.]+$")
		return usernameRegex.MatchString(username), "USERNAME_REGEX_MESSAGE", ""
	}
}

const NameStandardRegex = "^[A-Z][a-z]+(([ ][A-Z][a-z])?[a-z]*)*$"

func IsNexsoftNameStandardValid(username string) (bool, string, string) {
	usernameRegex := regexp.MustCompile(NameStandardRegex)
	return usernameRegex.MatchString(username), "NAME_REGEX_MESSAGE", ""
}

func IsNexsoftAdditionalInformationKeyStandardValid(username string) (bool, string) {
	usernameRegex := regexp.MustCompile("^[a-z][a-z0-9_.-]+$")
	return usernameRegex.MatchString(username), "ADDITIONAL_INFO_REGEX"
}

func IsOnlyContainLowerCase(username string) (bool, string) {
	usernameRegex := regexp.MustCompile("^[a-z]+$")
	return usernameRegex.MatchString(username), "LOWERCASE_REGEX"
}

func IsOnlyContainLowerCaseAndNumber(username string) (bool, string) {
	usernameRegex := regexp.MustCompile("^[a-z][a-z0-9]+$")
	return usernameRegex.MatchString(username), "LOWERCASE_AND_NUMBER_REGEX"
}

func IsStringEmpty(input string) bool {
	return input == ""
}

func IsTimestampValid(input string) (bool, string) {
	format := "2006-01-02T15:04:05.999999999"
	timestamp, err := time.Parse(format, input)

	if err != nil {
		return false, ""
	} else {
		return true, strings.Replace(timestamp.UTC().Format(time.RFC3339Nano), "Z", "", -1)
	}
}

func IsNumeric(input string) (int, bool) {
	result, err := strconv.Atoi(input)
	if err != nil {
		return -1, false
	} else {
		return result, true
	}
}

func IsNexsoftPermissionStandardValid(permission string) (bool, string) {
	permissionRegex := regexp.MustCompile("^[a-z]+[a-z._-]+[a-z1-9]+[:][a-z]+[a-z-][a-z]+$")
	return permissionRegex.MatchString(permission), "PERMISSION_REGEX_MESSAGE"
}

func IsNPWPValid(npwp string) bool {
	npwpRegex := regexp.MustCompile(fmt.Sprintf(`^(?:(0[0-9]{15})|([1-9]{1}[0-9]{5}((?:(0[1-9]|[12]\d|[4-6]\d|[37][01]))((0[1-9]|[1][0-2]))[0-9]{6}))|(?:[0-9]{15}))$`))
	return npwpRegex.MatchString(npwp)
}

func IsNIKValid(nik string) bool {
	nikRegex := regexp.MustCompile("^[0-9]{16,20}$")
	return nikRegex.MatchString(nik)
}

func IsFacsimileValid(fax string) bool {
	facsimileRegex := regexp.MustCompile("^[0-9]{5,25}$")
	return facsimileRegex.MatchString(fax)
}

const ProfileNameStandardRegex = "^[A-Z0-9](?:|(?:[a-z0-9]+|(?:[a-z0-9]|[a-z0-9])(?:([_-]|)[a-z0-9])+)|[ ]([A-Z0-9](?:|(?:[a-z0-9]+|(?:[a-z0-9]|[a-z0-9])(?:([_-]|)[a-z0-9])+))+)+)+$"

func IsNexsoftProfileNameStandardValid(profileName string) (bool, string) {
	NameOrTitle := regexp.MustCompile(ProfileNameStandardRegex)
	return NameOrTitle.MatchString(profileName), "PROFILE_NAME_REGEX_MESSAGE"
}

func IsNameWithUppercaseValid(input string) bool {
	UppercaseRegex := regexp.MustCompile("^[A-Z]+$")
	return UppercaseRegex.MatchString(input)
}

const LongNumericRegex = "^[0-9]+$"
const AlphanumericRegex = "^[A-Za-z0-9]+$"

func IsLongNumeric(input string) bool {
	longNumericRegex := regexp.MustCompile(LongNumericRegex)
	return longNumericRegex.MatchString(input)
}
func IsDataScopeValid(input string) (bool, string) {
	dataScopeRegexp := regexp.MustCompile("^(nexsoft[.]([a-z]+[a-z._-]+[a-z1-9]))[:]([1-9][0-9]*|all)$")
	return dataScopeRegexp.MatchString(input), "DATA_SCOPE"
}

const DirectoryNameStandardRegex = "^(([a-z0-9](?:([_]|)[a-z0-9])+))$"

func IsNexsoftDirectoryNameStandardValid(profileName string) (bool, string) {
	NameOrTitle := regexp.MustCompile(DirectoryNameStandardRegex)
	return NameOrTitle.MatchString(profileName), "DIRECTORY_NAME_REGEX_MESSAGE"
}
