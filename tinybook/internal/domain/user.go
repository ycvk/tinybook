package domain

import (
	regexp "github.com/wasilibs/go-re2"
	"strings"
	"time"
	"unicode/utf8"
)

type User struct {
	Id       int64  `json:"id"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Nickname string `json:"nickname"`
	Birthday string `json:"birthday"`
	AboutMe  string `json:"aboutMe"`
}

// ValidateEmail 验证邮箱
func (user *User) ValidateEmail() bool {
	user.Email = strings.TrimSpace(user.Email)
	return newUserRegexPattern().EmailRegex.MatchString(user.Email)
}

// ValidatePassword 验证密码
func (user *User) ValidatePassword() bool {
	user.Password = strings.TrimSpace(user.Password)
	return newUserRegexPattern().PasswordRegex.MatchString(user.Password)
}

// ValidateNickname 验证昵称
func (user *User) ValidateNickname() bool {
	user.Nickname = strings.TrimSpace(user.Nickname)
	countInString := utf8.RuneCountInString(user.Nickname)
	return countInString > 0 && countInString <= 10
}

// ValidateBirthday 验证生日
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

// ValidateAboutMe 验证个人简介
func (user *User) ValidateAboutMe() bool {
	user.AboutMe = strings.TrimSpace(user.AboutMe)
	return utf8.RuneCountInString(user.AboutMe) <= 200
}

// userRegexPattern 用户正则表达式
type userRegexPattern struct {
	EmailRegex    *regexp.Regexp
	PasswordRegex *regexp.Regexp
	BirthdayRegex *regexp.Regexp
}

// newUserRegexPattern 编译正则表达式
func newUserRegexPattern() *userRegexPattern {
	return &userRegexPattern{
		EmailRegex:    regexp.MustCompile(`^[a-zA-Z0-9._%+-]{1,64}@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`),
		PasswordRegex: regexp.MustCompile(`^[\w-!@#$%^&*()_+={}\[\]:;"'<>,.?~]{6,16}$`), // 长度6-16位, 只能包含字母、数字、特殊字符
		BirthdayRegex: regexp.MustCompile(`\d{4}-\d{2}-\d{2}`),                          // 格式为yyyy-MM-dd
	}
}
