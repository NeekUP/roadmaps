package core

import (
	"regexp"
	"strings"
	"unicode/utf8"

	"github.com/gosimple/slug"

	"github.com/google/uuid"
)

func IsValidEmail(email string) bool {

	re := regexp.MustCompile("^[\\w0-9.!#$%&'*+/=?^_`{|}~-]+@[\\w0-9](?:[\\w0-9-]{0,61}[\\w0-9])?(?:\\.[\\w0-9](?:[\\w0-9-]{0,61}[\\w0-9])?)*\\.[\\w]{2,15}$$")
	return re.MatchString(email)
}

func IsValidPassword(pass string) bool {
	negativePattern, _ := regexp.Compile("[\\s\"|*']")

	l := utf8.RuneCountInString(pass)
	return !(l == 0 || l > 64 || negativePattern.MatchString(pass))
}

func IsValidUserName(name string) bool {
	positivePattern, _ := regexp.Compile("^[a-zA-Z0-9_-]+$")

	l := utf8.RuneCountInString(name)
	return !(l < 2 || l > 32 || !positivePattern.MatchString(name))
}

func IsValidTokenFormat(token string) bool {
	return len(token) > 1 && strings.Contains(token, ".")
}

func IsValidFingerprint(fp string) bool {
	return len(fp) > 1
}

func IsValidUserAgent(useragent string) bool {
	return len(useragent) > 1
}

func IsValidUserID(id string) bool {
	_, err := uuid.Parse(id)
	return err == nil
}

// func IsJson(str string) bool {
// 	// TODO: add more tests
// 	var js map[string]interface{}
// 	return json.Unmarshal([]byte(str), &js) == nil
// }

func IsValidDscription(desc string) bool {
	l := len(desc)
	return l < 8000
}

func IsValidTopicTitle(title string) bool {
	l := utf8.RuneCountInString(title)
	return l > 0 && l < 100
}

func IsValidTopicName(name string) bool {
	l := utf8.RuneCountInString(name)
	return l > 0 && l < 100 && slug.IsSlug(name)
}

func IsValidPlanTitle(title string) bool {
	r := regexp.MustCompile("[<>;\"'/\\.]")
	doubleSpace := regexp.MustCompile("[\\s]{2,}")
	startSpace := regexp.MustCompile("^\\s")
	endSpace := regexp.MustCompile("\\s$")

	l := utf8.RuneCountInString(title)

	return l > 1 &&
		l < 100 &&
		!r.MatchString(title) &&
		!doubleSpace.MatchString(title) &&
		!startSpace.MatchString(title) &&
		!endSpace.MatchString(title)
}

func IsValidCommentText(text string) bool {
	length := utf8.RuneCountInString(text)
	return length > 0 && length < 1000
}

func IsValidCommentTitle(text string) bool {
	length := utf8.RuneCountInString(text)
	return length < 256
}
