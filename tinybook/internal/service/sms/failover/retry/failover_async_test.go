package retry

import (
	"context"
	"geek_homework/tinybook/internal/repository"
	repomocks "geek_homework/tinybook/internal/repository/mocks"
	"geek_homework/tinybook/internal/service/sms"
	smsmocks "geek_homework/tinybook/internal/service/sms/mocks"
	"geek_homework/tinybook/pkg/limiter"
	limitermocks "geek_homework/tinybook/pkg/limiter/mocks"
	"github.com/cockroachdb/errors"
	"github.com/zeebo/assert"
	"go.uber.org/mock/gomock"
	"testing"
	"time"
)

func TestAsyncFailoverSMSService_Send(t *testing.T) {
	testCases := []struct {
		name        string
		mock        func(controller *gomock.Controller) (sms.Service, repository.SMSRepository, limiter.Limiter, AsyncRetry)
		expectedErr error
	}{
		{
			name: "未限流，未失败，直接发送成功",
			mock: func(controller *gomock.Controller) (sms.Service, repository.SMSRepository, limiter.Limiter, AsyncRetry) {
				smsServ := smsmocks.NewMockService(controller)
				mockLimiter := limitermocks.NewMockLimiter(controller)
				smsRepository := repomocks.NewMockSMSRepository(controller)
				smsServ.EXPECT().Send(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
				mockLimiter.EXPECT().Limit(gomock.Any(), gomock.Any()).Return(false, nil)
				return smsServ, smsRepository, mockLimiter, nil
			},
			expectedErr: nil,
		},
		{
			name: "未限流，发送失败，重试成功",
			mock: func(controller *gomock.Controller) (sms.Service, repository.SMSRepository, limiter.Limiter, AsyncRetry) {
				smsServ := smsmocks.NewMockService(controller)
				mockLimiter := limitermocks.NewMockLimiter(controller)
				smsRepository := repomocks.NewMockSMSRepository(controller)
				retry := NewMockAsyncRetry(controller)
				mockLimiter.EXPECT().Limit(gomock.Any(), gomock.Any()).Return(false, nil)
				smsServ.EXPECT().Send(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New("发送失败"))
				retry.EXPECT().StartRetryLoop(gomock.Any()).Return(true, nil)
				return smsServ, smsRepository, mockLimiter, retry
			},
			expectedErr: nil,
		},
		{
			name: "未限流，发送失败，重试失败",
			mock: func(controller *gomock.Controller) (sms.Service, repository.SMSRepository, limiter.Limiter, AsyncRetry) {
				smsServ := smsmocks.NewMockService(controller)
				mockLimiter := limitermocks.NewMockLimiter(controller)
				smsRepository := repomocks.NewMockSMSRepository(controller)
				retry := NewMockAsyncRetry(controller)
				mockLimiter.EXPECT().Limit(gomock.Any(), gomock.Any()).Return(false, nil)
				smsServ.EXPECT().Send(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New("发送失败"))
				retry.EXPECT().StartRetryLoop(gomock.Any()).Return(false, errors.New("重试失败"))
				return smsServ, smsRepository, mockLimiter, retry
			},
			expectedErr: errors.New("重试失败"),
		},
		{
			name: "未限流，发送失败，重试成功，删除数据库记录失败",
			mock: func(controller *gomock.Controller) (sms.Service, repository.SMSRepository, limiter.Limiter, AsyncRetry) {
				smsServ := smsmocks.NewMockService(controller)
				mockLimiter := limitermocks.NewMockLimiter(controller)
				smsRepository := repomocks.NewMockSMSRepository(controller)
				retry := NewMockAsyncRetry(controller)
				mockLimiter.EXPECT().Limit(gomock.Any(), gomock.Any()).Return(false, nil)
				smsServ.EXPECT().Send(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New("发送失败"))
				retry.EXPECT().StartRetryLoop(gomock.Any()).Return(true, nil)
				smsRepository.EXPECT().Delete(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New("删除数据库记录失败"))
				return smsServ, smsRepository, mockLimiter, retry
			},
			expectedErr: errors.New("删除数据库记录失败"),
		},
		{
			name: "未限流，发送失败，重试成功，删除数据库记录成功",
			mock: func(controller *gomock.Controller) (sms.Service, repository.SMSRepository, limiter.Limiter, AsyncRetry) {
				smsServ := smsmocks.NewMockService(controller)
				mockLimiter := limitermocks.NewMockLimiter(controller)
				smsRepository := repomocks.NewMockSMSRepository(controller)
				retry := NewMockAsyncRetry(controller)
				mockLimiter.EXPECT().Limit(gomock.Any(), gomock.Any()).Return(false, nil)
				smsServ.EXPECT().Send(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New("发送失败"))
				retry.EXPECT().StartRetryLoop(gomock.Any()).Return(true, nil)
				smsRepository.EXPECT().Delete(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
				return smsServ, smsRepository, mockLimiter, retry
			},
			expectedErr: errors.New("发送失败"),
		},
		{
			name: "限流，存储数据库失败",
			mock: func(controller *gomock.Controller) (sms.Service, repository.SMSRepository, limiter.Limiter, AsyncRetry) {
				smsServ := smsmocks.NewMockService(controller)
				mockLimiter := limitermocks.NewMockLimiter(controller)
				smsRepository := repomocks.NewMockSMSRepository(controller)
				mockLimiter.EXPECT().Limit(gomock.Any(), gomock.Any()).Return(true, nil)
				smsRepository.EXPECT().Save(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New("存储数据库失败"))
				return smsServ, smsRepository, mockLimiter, nil
			},
			expectedErr: errors.New("存储数据库失败"),
		},
		{
			name: "限流，存储数据库成功，重试成功，删除数据库记录失败",
			mock: func(controller *gomock.Controller) (sms.Service, repository.SMSRepository, limiter.Limiter, AsyncRetry) {
				smsServ := smsmocks.NewMockService(controller)
				mockLimiter := limitermocks.NewMockLimiter(controller)
				smsRepository := repomocks.NewMockSMSRepository(controller)
				retry := NewMockAsyncRetry(controller)
				mockLimiter.EXPECT().Limit(gomock.Any(), gomock.Any()).Return(true, nil)
				smsRepository.EXPECT().Save(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
				retry.EXPECT().StartRetryLoop(gomock.Any()).Return(true, nil)
				smsRepository.EXPECT().Delete(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New("删除数据库记录失败"))
				smsServ.EXPECT().Send(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New("删除数据库记录失败"))
				return smsServ, smsRepository, mockLimiter, retry
			},
			expectedErr: errors.New("删除数据库记录失败"),
		},
		{
			name: "限流，重试失败",
			mock: func(controller *gomock.Controller) (sms.Service, repository.SMSRepository, limiter.Limiter, AsyncRetry) {
				smsServ := smsmocks.NewMockService(controller)
				mockLimiter := limitermocks.NewMockLimiter(controller)
				smsRepository := repomocks.NewMockSMSRepository(controller)
				retry := NewMockAsyncRetry(controller)
				mockLimiter.EXPECT().Limit(gomock.Any(), gomock.Any()).Return(true, nil)
				smsRepository.EXPECT().Save(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
				retry.EXPECT().StartRetryLoop(gomock.Any()).Return(false, errors.New("重试失败"))
				return smsServ, smsRepository, mockLimiter, retry
			},
			expectedErr: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			smsService, smsrepo, l, retry := tc.mock(ctrl)
			// 错误率监控器
			monitor := NewErrorRateMonitor(0.3, 0.5, 30*time.Second)
			// 重试任务
			//retryTask := NewRetryTask(3)
			// 异步重试服务
			failoverSMSService := NewAsyncFailoverSMSService(l, smsService, smsrepo, monitor, retry)

			err := failoverSMSService.Send(context.Background(), "test", []string{"778899"}, "13011223344")
			if err != nil {
				assert.Equal(t, tc.expectedErr.Error(), err.Error())
			} else {
				assert.Equal(t, tc.expectedErr, err)
			}
		})
	}
}
