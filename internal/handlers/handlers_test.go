package handlers

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/golang/mock/gomock"
	"github.com/ra1nz0r/metric_alert_app/internal/storage"
	"github.com/ra1nz0r/metric_alert_app/internal/storage/mocks"
	"github.com/stretchr/testify/assert"
)

func TestGetAllMetrics(t *testing.T) {
	type args struct {
		ms storage.MetricService
		w  *httptest.ResponseRecorder
		r  *http.Request
	}
	tests := []struct {
		name        string
		args        args
		tGaugeStr   *map[string]float64
		tCounterStr *map[string]int64
	}{
		{
			name: "Test 1.",
			args: args{
				ms: storage.New(),
				w:  httptest.NewRecorder(),
				r:  httptest.NewRequest(http.MethodGet, "/", nil),
			},
			tGaugeStr:   &map[string]float64{"Alloc": 4.51},
			tCounterStr: &map[string]int64{"PollCount": 73},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mock := mocks.NewMockMetricService(ctrl)

			nh := NewHandlers(mock)

			mock.EXPECT().MakeStorageCopy().Return(tt.tGaugeStr, tt.tCounterStr).Times(1)

			h := http.HandlerFunc(nh.GetAllMetrics)

			h(tt.args.w, tt.args.r)

			result := tt.args.w.Result()

			assert.Equal(t, http.StatusOK, result.StatusCode)
			assert.Equal(t, "application/json; charset=UTF-8", result.Header.Get("Content-Type"))

			result.Body.Close()
		})
	}
}

func TestGetMetricsByName(t *testing.T) {
	type args struct {
		mType string
		mName string
	}
	tests := []struct {
		name        string
		tGaugeStr   *map[string]float64
		tCounterStr *map[string]int64
		args        args
		wantCode    int
		wantError   error
		wantHeader  string
	}{
		{
			name:      "Test 1. Gauge metric.",
			tGaugeStr: &map[string]float64{"Alloc": 4.51},
			args: args{
				mType: "gauge",
				mName: "Alloc",
			},
			wantCode:   200,
			wantError:  nil,
			wantHeader: "text/plain; charset=UTF-8",
		},
		{
			name:        "Test 2. Counter metric.",
			tCounterStr: &map[string]int64{"PollCount": 73},
			args: args{
				mType: "counter",
				mName: "PollCount",
			},
			wantCode:   200,
			wantError:  nil,
			wantHeader: "text/plain; charset=UTF-8",
		},
		{
			name:        "Test 3. Incorrect metric type.",
			tCounterStr: &map[string]int64{"PollCount": 73},
			args: args{
				mType: "incType",
				mName: "PollCount",
			},
			wantCode:   404,
			wantError:  fmt.Errorf("type not found"),
			wantHeader: "application/json; charset=UTF-8",
		},
		{
			name:      "Test 4. Incorrect gauge metric name.",
			tGaugeStr: &map[string]float64{"Alloc": 7.3},
			args: args{
				mType: "gauge",
				mName: "incAlloc",
			},
			wantCode:   404,
			wantError:  fmt.Errorf("metric not found"),
			wantHeader: "application/json; charset=UTF-8",
		},
		{
			name:        "Test 5. Incorrect counter metric name.",
			tCounterStr: &map[string]int64{"PollCount": 7},
			args: args{
				mType: "counter",
				mName: "incCount",
			},
			wantCode:   404,
			wantError:  fmt.Errorf("metric not found"),
			wantHeader: "application/json; charset=UTF-8",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mock := mocks.NewMockMetricService(ctrl)

			nh := NewHandlers(mock)

			mock.EXPECT().MakeStorageCopy().Return(tt.tGaugeStr, tt.tCounterStr).Times(1)

			router := chi.NewRouter()

			router.Get("/value/{type}/{name}", nh.GetMetricByName)

			reqURL := fmt.Sprintf("/value/%s/%s", tt.args.mType, tt.args.mName)

			req := httptest.NewRequest(http.MethodGet, reqURL, nil)

			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			result := w.Result()

			assert.Equal(t, tt.wantCode, result.StatusCode)
			assert.Equal(t, tt.wantHeader, result.Header.Get("Content-Type"))

			result.Body.Close()
		})
	}
}

func TestUpdateMetrics(t *testing.T) {

	type args struct {
		mType  string
		mName  string
		mValue any
	}
	tests := []struct {
		name        string
		args        args
		reqURL      string
		wantCode    int
		buildEXPECT func(store *mocks.MockMetricService, mName string, mValue any)
	}{
		{
			name: "Test 1. Gauge metric.",
			args: args{
				mType:  "gauge",
				mName:  "Alloc",
				mValue: float64(4.51),
			},
			wantCode: 200,
			buildEXPECT: func(store *mocks.MockMetricService, mName string, mValue any) {
				store.EXPECT().UpdateGauge(mName, mValue).Times(1)
			},
		},
		{
			name: "Test 2. Counter metric.",
			args: args{
				mType:  "counter",
				mName:  "PollCount",
				mValue: int64(23),
			},
			wantCode: 200,
			buildEXPECT: func(store *mocks.MockMetricService, mName string, mValue any) {
				store.EXPECT().UpdateCounter(mName, mValue).Times(1)
			},
		},
		{
			name: "Test 3. Empty metric name.",
			args: args{
				mType:  "gauge",
				mName:  "",
				mValue: 23,
			},
			wantCode:    404,
			buildEXPECT: func(store *mocks.MockMetricService, mName string, mValue any) {},
		},
		{
			name: "Test 4. Empty metric type.",
			args: args{
				mType:  "",
				mName:  "PollCount",
				mValue: 23,
			},
			wantCode:    400,
			buildEXPECT: func(store *mocks.MockMetricService, mName string, mValue any) {},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mock := mocks.NewMockMetricService(ctrl)
			nh := NewHandlers(mock)

			tt.buildEXPECT(mock, tt.args.mName, tt.args.mValue)

			router := chi.NewRouter()

			router.Post("/update/{type}/{name}/{value}", nh.UpdateMetrics)

			req := httptest.NewRequest(http.MethodPost, makeReqURL(tt.args.mType, tt.args.mName, tt.args.mValue), nil)

			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			result := w.Result()

			assert.Equal(t, tt.wantCode, result.StatusCode)

			result.Body.Close()

		})
	}
}

func makeReqURL(mType, mName string, mValue any) string {
	switch mValue.(type) {
	case float64:
		return fmt.Sprintf("/update/%s/%s/%.2f", mType, mName, mValue)
	default:
		return fmt.Sprintf("/update/%s/%s/%d", mType, mName, mValue)
	}
}
