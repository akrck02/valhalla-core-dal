package mock

import "math/rand"

func RandomString(length int) string {

	str := ""
	for i := 0; i < length; i++ {
		str += string(rune(97 + rand.Intn(26)))
	}

	return str
}
