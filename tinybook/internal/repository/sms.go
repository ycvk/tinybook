package repository

import (
	"context"
	"geek_homework/tinybook/internal/repository/dao"
	"github.com/bytedance/sonic"
)

type SMSRepository interface {
	Save(ctx context.Context, tplId string, args []string, numbers ...string) error
	Delete(ctx context.Context, tplId string, args []string, numbers ...string) error
}

type GormSMSRepository struct {
	dao dao.SMSDAO
}

func NewGormSMSRepository(dao dao.SMSDAO) SMSRepository {
	return &GormSMSRepository{dao: dao}
}

func (g GormSMSRepository) Save(ctx context.Context, tplId string, args []string, numbers ...string) error {
	marshalString, err := sonic.MarshalString(args)
	if err != nil {
		return err
	}
	smsDaos := make([]*dao.SMS, len(numbers))
	for _, number := range numbers {
		smsDaos = append(smsDaos, &dao.SMS{
			TplId: tplId,
			Code:  marshalString,
			Phone: number,
		})
	}
	return g.dao.Insert(ctx, smsDaos...)
}

func (g GormSMSRepository) Delete(ctx context.Context, tplId string, args []string, numbers ...string) error {
	marshalString, err := sonic.MarshalString(args)
	if err != nil {
		return err
	}
	smsDaos := make([]*dao.SMS, len(numbers))
	for _, number := range numbers {
		smsDaos = append(smsDaos, &dao.SMS{
			TplId: tplId,
			Code:  marshalString,
			Phone: number,
		})
	}
	return g.dao.Delete(ctx, smsDaos...)
}
