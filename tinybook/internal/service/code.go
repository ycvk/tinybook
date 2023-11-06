package service

import (
	"context"
	"errors"
	"fmt"
	"geek_homework/tinybook/internal/repository"
	"geek_homework/tinybook/internal/service/sms"
	regexp "github.com/wasilibs/go-re2"
	"math/rand"
)

type CodeService interface {
	Send(ctx context.Context, biz, phone, timeInterval string) error
	Verify(ctx context.Context, biz, phone, code string) (bool, error)
}

type codeService struct {
	repo       repository.CodeRepository
	smsService sms.Service
}

func NewCodeService(repo repository.CodeRepository, smsService sms.Service) CodeService {
	return &codeService{
		repo:       repo,
		smsService: smsService,
	}
}

func (codeService *codeService) Send(ctx context.Context, biz, phone, timeInterval string) error {
	// 校验手机号码
	if !codeService.validatePhoneNum(phone) {
		return errors.New(fmt.Sprintf("手机号码格式不正确: %s", phone))
	}
	code := codeService.generateCode()
	err := codeService.repo.Set(ctx, biz, phone, code, timeInterval)
	if err != nil {
		return err
	}
	// 发送短信
	const tplId = "123" // 验证码模板id
	return codeService.smsService.Send(ctx, tplId, []string{code}, phone)
}

func (codeService *codeService) Verify(ctx context.Context, biz, phone, code string) (bool, error) {
	// 校验手机号码
	if !codeService.validatePhoneNum(phone) {
		return false, errors.New(fmt.Sprintf("手机号码格式不正确: %s", phone))
	}
	return codeService.repo.Verify(ctx, biz, phone, code)
}

func (codeService *codeService) generateCode() string {
	// 生成6位随机数
	code := rand.Intn(100_0000)
	return fmt.Sprintf("%06d", code)
}

func (codeService *codeService) validatePhoneNum(phone string) bool {
	compile := regexp.MustCompile(`^(\+86)?1[3-9]\d{9}$`)
	return compile.MatchString(phone)
}
