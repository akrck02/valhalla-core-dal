package mock

import "math/rand"

func randomString(length int) string {

	str := ""
	for i := 0; i < length; i++ {
		str += string(rune(97 + rand.Intn(26)))
	}

	return str
}

func Email() string {

	//Return a random character string
	return randomString(10) + "@valhalla.org"
}

func EmailNotDot() string {

	return "thelegendof@lumberjackcom"
}

func EmailNotAt() string {
	return "thelegendoflumberjack.com"
}

func EmailShort() string {
	return "a@s."
}

func Password() string {
	return "PasswordPassword1#"
}

func PasswordShort() string {
	return "Pass1#"
}

func PasswordNotUpperCase() string {
	return "passwordpassword1#"
}

func PasswordNotLowerCase() string {
	return "PASSWORDPASSWORD1#"
}

func PasswordNotNumber() string {
	return "PasswordPassword#"
}

func PasswordNotSpecialChar() string {
	return "PasswordPassword1"
}

func Username() string {
	return randomString(10)
}
