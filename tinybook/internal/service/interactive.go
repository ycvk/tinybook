package service

import (
	"context"
	"geek_homework/tinybook/internal/repository"
)

type InteractiveService interface {
	IncreaseReadCount(ctx context.Context, biz string, bizId int64) error
	Like(ctx context.Context, biz string, id int64, uid int64) error
	Unlike(ctx context.Context, biz string, id int64, uid int64) error
	Collect(ctx context.Context, biz string, id int64, cid int64, uid int64) error
}

type interactiveService struct {
	repo repository.InteractiveRepository
}

func (i *interactiveService) Collect(ctx context.Context, biz string, id int64, cid int64, uid int64) error {
	return i.repo.Collect(ctx, biz, id, cid, uid)
}

func (i *interactiveService) Like(ctx context.Context, biz string, id int64, uid int64) error {
	return i.repo.IncreaseLikeCount(ctx, biz, id, uid)
}

func (i *interactiveService) Unlike(ctx context.Context, biz string, id int64, uid int64) error {
	return i.repo.DecreaseLikeCount(ctx, biz, id, uid)
}

func NewInteractiveService(repo repository.InteractiveRepository) InteractiveService {
	return &interactiveService{repo: repo}
}

func (i *interactiveService) IncreaseReadCount(ctx context.Context, biz string, bizId int64) error {
	return i.repo.IncreaseReadCount(ctx, biz, bizId)
}
