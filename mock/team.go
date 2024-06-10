package mock

func TeamName() string {
	return RandomString(40)
}

func TeamNameShort() string {
	return RandomString(1)
}

func TeamNameLong() string {
	return RandomString(100)
}

func TeamDescription() string {
	return RandomString(200)
}

func TeamDescriptionShort() string {
	return "a"
}

func TeamDescriptionLong() string {
	return RandomString(500)
}
