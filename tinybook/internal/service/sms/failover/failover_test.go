package failover

import (
	"context"
	"geek_homework/tinybook/internal/service/sms"
	smsmocks "geek_homework/tinybook/internal/service/sms/mocks"
	"github.com/cockroachdb/errors"
	"github.com/zeebo/assert"
	"go.uber.org/mock/gomock"
	"testing"
)

func TestFailoverSMSService_Send(t *testing.T) {
	testCases := []struct {
		name        string
		mock        func(controller *gomock.Controller) []sms.Service
		expectedErr error
	}{
		{
			name: "first services success",
			mock: func(controller *gomock.Controller) []sms.Service {
				service1 := smsmocks.NewMockService(controller)
				service1.EXPECT().Send(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
				service2 := smsmocks.NewMockService(controller)
				return []sms.Service{service1, service2}
			},
			expectedErr: nil,
		},
		{
			name: "second services success",
			mock: func(controller *gomock.Controller) []sms.Service {
				service1 := smsmocks.NewMockService(controller)
				service1.EXPECT().Send(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New("first services failed"))
				service2 := smsmocks.NewMockService(controller)
				service2.EXPECT().Send(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
				return []sms.Service{service1, service2}
			},
			expectedErr: nil,
		},
		{
			name: "all services failed",
			mock: func(controller *gomock.Controller) []sms.Service {
				service1 := smsmocks.NewMockService(controller)
				service1.EXPECT().Send(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New("first services failed"))
				service2 := smsmocks.NewMockService(controller)
				service2.EXPECT().Send(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New("second services failed"))
				return []sms.Service{service1, service2}
			},
			expectedErr: ErrAllServicesFailed,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			services := testCase.mock(ctrl)
			smsService := NewFailoverSMSService(services...)
			err := smsService.Send(context.Background(), "test", []string{"test"}, "12345678901")
			if err != nil {
				assert.Equal(t, testCase.expectedErr.Error(), err.Error())
			} else {
				assert.Equal(t, testCase.expectedErr, err)
			}
		})
	}
}
