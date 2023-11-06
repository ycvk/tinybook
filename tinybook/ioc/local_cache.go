package ioc

import (
	"github.com/Yiling-J/theine-go"
	"sync"
)

var (
	localOnce sync.Once
	client    *theine.Cache[string, any]
)

func InitLocalCache() *theine.Cache[string, any] {
	localOnce.Do(func() {
		var err error
		client, err = theine.NewBuilder[string, any](1000).Build()
		if err != nil {
			panic(err)
		}
	})
	return client
}
