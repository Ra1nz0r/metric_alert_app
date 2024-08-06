package storage

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

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
