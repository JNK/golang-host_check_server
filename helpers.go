package main

func checkError(e error) {
	if e != nil {
		panic(e)
	}
}

func CheckHealth(config Config) bool {
	for _, check := range config.Checks {
		if !check.Validate() {
			return false
		}
	}
	return true
}

