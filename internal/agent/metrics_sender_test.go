package agent

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
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

func HelloWorld(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello World")
}

func TestMakeRequest(t *testing.T) {
	type args struct {
		resURL string
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testServer := httptest.NewServer(http.HandlerFunc(HelloWorld))
			defer testServer.Close()
			testClient := testServer.Client()

			resp, err := testClient.Get(testServer.URL)
			if err != nil {
				t.Errorf("Get error: %v", err)
			}
			if resp.StatusCode != http.StatusOK {
				t.Errorf("response code is not 200: %d", resp.StatusCode)
			}
			data, err := io.ReadAll(resp.Body)
			if err != nil {
				t.Errorf("io.ReadAll error: %v", err)
			}
			if string(data) != "Hello World\n" {
				t.Error("response body does not equal to Hello World")
			}
		})
	}
}
