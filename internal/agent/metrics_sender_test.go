package agent

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/ra1nz0r/metric_alert_app/internal/storage/mocks"
)

func UpdateMetrics(t *testing.T) {
	t.Run("Test 1. Update metrics.", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		mock := mocks.NewMockMetricService(mockCtrl)

		ss := NewSender(mock)

		mock.EXPECT().UpdateGauge(gomock.Any().String(), gomock.Any()).Times(28)
		mock.EXPECT().UpdateCounter(gomock.Any().String(), gomock.Any()).Times(1)

		ss.UpdateMetrics()
	})
}
