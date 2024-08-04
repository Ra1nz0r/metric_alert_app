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
		name string
		args args
	}{
		{
			name: "test1",
			args: args{
				ms: storage.New(),
				w:  httptest.NewRecorder(),
				r:  httptest.NewRequest(http.MethodGet, "/", nil),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mock := mocks.NewMockMetricService(ctrl)

			nh := NewHandlers(mock)

			mock.EXPECT().MakeStorageCopy().Times(1)

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
		name       string
		w          *httptest.ResponseRecorder
		req        *http.Request
		router     *chi.Mux
		args       args
		reqURL     string
		wantCode   int
		wantError  error
		wantHeader string
	}{
		{
			name:   "Test 1. Gauge metric.",
			w:      httptest.NewRecorder(),
			router: chi.NewRouter(),
			args: args{
				mType: "gauge",
				mName: "Alloc",
			},
			wantCode:   200,
			wantError:  nil,
			wantHeader: "text/plain; charset=UTF-8",
		},
		{
			name:   "Test 2. Counter metric.",
			w:      httptest.NewRecorder(),
			router: chi.NewRouter(),
			args: args{
				mType: "counter",
				mName: "PollCount",
			},
			wantCode:   200,
			wantError:  nil,
			wantHeader: "text/plain; charset=UTF-8",
		},
		{
			name:   "Test 3. Incorrect metric type.",
			w:      httptest.NewRecorder(),
			router: chi.NewRouter(),
			args: args{
				mType: "incType",
				mName: "PollCount",
			},
			wantCode:   404,
			wantError:  fmt.Errorf("type not found"),
			wantHeader: "application/json; charset=UTF-8",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mock := mocks.NewMockMetricService(ctrl)

			nh := NewHandlers(mock)

			mock.EXPECT().GetMetricVal(tt.args.mType, tt.args.mName).Return(nil, tt.wantError).Times(1)

			tt.router.Get("/value/{type}/{name}", nh.GetMetricByName)

			reqURL := fmt.Sprintf("/value/%s/%s", tt.args.mType, tt.args.mName)

			tt.req = httptest.NewRequest(http.MethodGet, reqURL, nil)

			tt.router.ServeHTTP(tt.w, tt.req)

			result := tt.w.Result()

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
		name     string
		w        *httptest.ResponseRecorder
		req      *http.Request
		router   *chi.Mux
		args     args
		reqURL   string
		wantCode int
	}{
		{
			name:   "Test 1. Correct gauge metric.",
			w:      httptest.NewRecorder(),
			router: chi.NewRouter(),
			args: args{
				mType:  "gauge",
				mName:  "Alloc",
				mValue: float64(4.51),
			},
			wantCode: 200,
		},
		{
			name:   "Test 2. Correct counter metric.",
			w:      httptest.NewRecorder(),
			router: chi.NewRouter(),
			args: args{
				mType:  "counter",
				mName:  "PollCount",
				mValue: int64(23),
			},
			wantCode: 200,
		},
		{
			name:   "Test 3. Empty metric name.",
			w:      httptest.NewRecorder(),
			router: chi.NewRouter(),
			args: args{
				mType:  "gauge",
				mName:  "",
				mValue: 23,
			},
			wantCode: 404,
		},
		{
			name:   "Test 4. Empty metric type.",
			w:      httptest.NewRecorder(),
			router: chi.NewRouter(),
			args: args{
				mType:  "",
				mName:  "PollCount",
				mValue: 23,
			},
			wantCode: 400,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mock := mocks.NewMockMetricService(ctrl)
			nh := NewHandlers(mock)

			switch {
			case tt.args.mType == "gauge" && tt.args.mName != "":
				mock.EXPECT().UpdateGauge(tt.args.mName, tt.args.mValue).Times(1)

			case tt.args.mType == "counter" && tt.args.mName != "":
				mock.EXPECT().UpdateCounter(tt.args.mName, tt.args.mValue).Times(1)
			}

			tt.router.Post("/update/{type}/{name}/{value}", nh.UpdateMetrics)

			tt.req = httptest.NewRequest(http.MethodPost, makeReqURL(tt.args.mType, tt.args.mName, tt.args.mValue), nil)

			tt.router.ServeHTTP(tt.w, tt.req)

			result := tt.w.Result()

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
