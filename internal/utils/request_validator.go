package utils

func ValidateParams(params ...string) bool {
	for _, param := range params {
		if param == "" {
			return false
		}
	}

	return true
}

func ValidateMaxRetries(maxRetries int) bool {
	return maxRetries >= 1
}
