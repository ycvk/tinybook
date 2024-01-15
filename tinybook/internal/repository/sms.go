package repository

import (
	"context"
	"github.com/bytedance/sonic"
	"tinybook/tinybook/internal/repository/dao"
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
	smsDaos := make([]*dao.SMS, 0)
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
	smsDaos := make([]*dao.SMS, 0)
	for _, number := range numbers {
		smsDaos = append(smsDaos, &dao.SMS{
			TplId: tplId,
			Code:  marshalString,
			Phone: number,
		})
	}
	return g.dao.Delete(ctx, smsDaos...)
}
