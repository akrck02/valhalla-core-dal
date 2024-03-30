package mock

func ProjectName() string {
	return randomString(10)
}

func ProjectNameShort() string {
	return "a"
}

func ProjectNameLong() string {
	return "This is a really long project name that is longer than 50 characters"
}

func ProjectDescription() string {
	return randomString(40)
}
