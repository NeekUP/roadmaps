package core

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
)

func IsValidEmail(email string) (bool, ErrorCode) {

	re := regexp.MustCompile("^[\\w0-9.!#$%&'*+/=?^_`{|}~-]+@[\\w0-9](?:[\\w0-9-]{0,61}[\\w0-9])?(?:\\.[\\w0-9](?:[\\w0-9-]{0,61}[\\w0-9])?)*\\.[\\w]{2,15}$$")
	if re.MatchString(email) {
		return true, None
	}

	return false, BadEmail
}

func IsValidPassword(pass string) (bool, ErrorCode) {
	negativePattern, _ := regexp.Compile("[\\s\"|*']")

	l := len(pass)
	if l == 0 || l > 64 || negativePattern.MatchString(pass) {
		return false, BadPassword
	}

	return true, None
}

func IsValidUserName(name string) (bool, ErrorCode) {
	positivePattern, _ := regexp.Compile("^[a-zA-Z0-9_-]+$")

	l := len(name)
	if l < 2 || l > 32 || !positivePattern.MatchString(name) {
		return false, BadUserName
	}

	return true, None
}

func IsValidTokenFormat(token string) (bool, error) {
	result := len(token) > 1 && strings.Contains(token, ".")
	if result {
		return true, nil
	}

	return false, fmt.Errorf("Token not valid [%s]", token)
}

func IsValidFingerprint(fp string) (bool, error) {
	result := len(fp) > 1
	if result {
		return true, nil
	}

	return false, fmt.Errorf("Fingerprint not valid [%s]", fp)
}

func IsValidUserAgent(useragent string) (bool, error) {
	result := len(useragent) > 1
	if result {
		return true, nil
	}

	return false, fmt.Errorf("Useragent not valid [%s]", useragent)
}

func IsJson(str string) bool {
	// TODO: add more tests
	var js map[string]interface{}
	return json.Unmarshal([]byte(str), &js) == nil
}

func IsValidSourceTitle(title string) (bool, error) {
	r := regexp.MustCompile("^\\W$")
	l := len(title)
	if l == 0 || l > 100 || r.MatchString(title) {
		return false, fmt.Errorf("Source title contain not valid value [%s]", title)
	}

	return true, nil
}
