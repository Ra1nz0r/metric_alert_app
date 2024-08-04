package storage

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetMetricVal(t *testing.T) {
	type args struct {
		mType string
		mName string
	}
	tests := []struct {
		name    string
		ms      *MemStorage
		args    args
		want    any
		wantErr bool
	}{
		{
			name: "Test 1. Correct gauge.",
			ms: &MemStorage{
				gauge: map[string]float64{
					"Alloc": 4.51,
				},
			},
			args: args{
				mType: "gauge",
				mName: "Alloc",
			},
			want:    4.51,
			wantErr: false,
		},
		{
			name: "Test 2. Correct counter.",
			ms: &MemStorage{
				counter: map[string]int64{
					"PollCount": 22,
				},
			},
			args: args{
				mType: "counter",
				mName: "PollCount",
			},
			want:    int64(22),
			wantErr: false,
		},
		{
			name: "Test 3. Incorrect name.",
			ms: &MemStorage{
				gauge: map[string]float64{
					"Alloc": 4.51,
				},
			},
			args: args{
				mType: "gauge",
				mName: "incAlloc",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Test 4. Incorrect type.",
			ms: &MemStorage{
				gauge: map[string]float64{
					"Alloc": 4.51,
				},
			},
			args: args{
				mType: "incGauge",
				mName: "Alloc",
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.ms.GetMetricVal(tt.args.mType, tt.args.mName)
			if (err != nil) != tt.wantErr {
				t.Errorf("MemStorage.GetMetricVal() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MemStorage.GetMetricVal() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMakeStorageCopy(t *testing.T) {
	t.Run("Test copy.", func(t *testing.T) {
		ms := &MemStorage{
			gauge: map[string]float64{
				"Alloc":       4.51,
				"LastGC":      62.2,
				"RandomValue": 32.938,
			},
			counter: map[string]int64{
				"PollCount": 782,
			},
		}

		gauge, counter := ms.MakeStorageCopy()

		require.NotNil(t, gauge)
		require.NotNil(t, counter)

		(*gauge)["Alloc"] = 8.888
		(*counter)["PollCount"] = 238

		assert.NotEqual(t, ms.gauge, gauge)
		assert.NotEqual(t, &ms.gauge, &gauge)

		assert.NotEqual(t, ms.gauge, counter)
		assert.NotEqual(t, &ms.counter, &counter)
	})

}

func TestUpdateGauge(t *testing.T) {
	type args struct {
		name  string
		value float64
	}
	tests := []struct {
		name string
		ms   *MemStorage
		args args
	}{
		{
			name: "Test 1. Correct value.",
			ms: &MemStorage{
				gauge: map[string]float64{
					"Alloc": 4.51,
				},
			},
			args: args{
				name:  "Alloc",
				value: 65.234,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.ms.UpdateGauge(tt.args.name, tt.args.value)

			for k, v := range tt.ms.gauge {
				assert.Equal(t, tt.args.name, k)
				assert.Equal(t, tt.args.value, v)
			}
		})
	}
}

func TestUpdateCounter(t *testing.T) {
	type args struct {
		name  string
		value int64
	}
	tests := []struct {
		name string
		ms   *MemStorage
		args args
		want int64
	}{
		{
			name: "Test 1. Correct value.",
			ms: &MemStorage{
				counter: map[string]int64{
					"PollCount": 23,
				},
			},
			args: args{
				name:  "PollCount",
				value: 27,
			},
			want: 50,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.ms.UpdateCounter(tt.args.name, tt.args.value)

			for k, v := range tt.ms.counter {
				assert.Equal(t, tt.args.name, k)
				assert.Equal(t, tt.want, v)
			}
		})
	}
}
