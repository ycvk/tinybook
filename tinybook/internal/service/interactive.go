package service

import (
	"context"
	"geek_homework/tinybook/internal/domain"
	"geek_homework/tinybook/internal/repository"
	"golang.org/x/sync/errgroup"
)

type InteractiveService interface {
	IncreaseReadCount(ctx context.Context, biz string, bizId int64) error
	Like(ctx context.Context, biz string, id int64, uid int64) error
	Unlike(ctx context.Context, biz string, id int64, uid int64) error
	Collect(ctx context.Context, biz string, id int64, cid int64, uid int64) error
	GetInteractive(ctx context.Context, biz string, id int64, uid int64) (domain.Interactive, error)
}

type interactiveService struct {
	repo repository.InteractiveRepository
}

func (i *interactiveService) GetInteractive(ctx context.Context, biz string, id int64, uid int64) (domain.Interactive, error) {
	interactive, err := i.repo.GetInteractive(ctx, biz, id)
	if err != nil {
		return domain.Interactive{}, err
	}
	var eg errgroup.Group
	eg.Go(func() error {
		var er error
		interactive.Liked, er = i.repo.Liked(ctx, biz, id, uid)
		return er
	})
	eg.Go(func() error {
		var er error
		interactive.Collected, er = i.repo.Collected(ctx, biz, id, uid)
		return er
	})
	return interactive, eg.Wait()
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
