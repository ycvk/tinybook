package domain

import (
	regexp "github.com/wasilibs/go-re2"
	"strings"
	"time"
)

type User struct {
	Id       int64  `json:"id"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Nickname string `json:"nickname"`
	Birthday string `json:"birthday"`
	AboutMe  string `json:"aboutMe"`
}

func (user *User) ValidateEmail() bool {
	user.Email = strings.TrimSpace(user.Email)
	return newUserRegexPattern().EmailRegex.MatchString(user.Email)
}

func (user *User) ValidatePassword() bool {
	user.Password = strings.TrimSpace(user.Password)
	return newUserRegexPattern().PasswordRegex.MatchString(user.Password)
}

func (user *User) ValidateNickname() bool {
	user.Nickname = strings.TrimSpace(user.Nickname)
	return len(user.Nickname) > 0 && len(user.Nickname) <= 10
}

func (user *User) ValidateBirthday() bool {
	birthday := user.Birthday
	// 正则检查时间格式
	match := newUserRegexPattern().BirthdayRegex.MatchString(birthday)
	if !match {
		return false
	}

	// 转换字符串为时间
	inputTime, err := time.Parse("2006-01-01", birthday)
	if err != nil {
		return false
	}

	// 检查时间是否早于当前时间
	currentTime := time.Now()
	return inputTime.Before(currentTime)
}

func (user *User) ValidateAboutMe() bool {
	user.AboutMe = strings.TrimSpace(user.AboutMe)
	return len(user.AboutMe) <= 100
}

type userRegexPattern struct {
	EmailRegex    *regexp.Regexp
	PasswordRegex *regexp.Regexp
	BirthdayRegex *regexp.Regexp
}

func newUserRegexPattern() *userRegexPattern {
	return &userRegexPattern{
		EmailRegex:    regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`),
		PasswordRegex: regexp.MustCompile(`^[\w-]{6,16}$`),
		BirthdayRegex: regexp.MustCompile(`\d{4}-\d{2}-\d{2}`),
	}
}
