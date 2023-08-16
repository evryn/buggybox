package common_test

import (
	"kermoo/modules/common"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetValue(t *testing.T) {
	t.Run("value is exactly set", func(t *testing.T) {
		val := float32(42)
		s := &common.SingleFloat{
			Exactly: &val,
		}

		got, err := s.ToFloat()
		assert.NoError(t, err)
		assert.Equal(t, val, got)
	})

	t.Run("value is a range", func(t *testing.T) {
		min := float32(10)
		max := float32(20)
		s := &common.SingleFloat{
			Between: []float32{min, max},
		}

		got, err := s.ToFloat()
		assert.NoError(t, err)
		assert.GreaterOrEqual(t, got, min)
		assert.LessOrEqual(t, got, max)
	})
}

func TestToSingleValues(t *testing.T) {
	t.Run("value is exactly set", func(t *testing.T) {
		val := float32(42)
		v := common.MultiFloat{
			SingleFloat: common.SingleFloat{
				Exactly: &val,
			},
		}

		got, err := v.ToSingleFloats()
		assert.NoError(t, err)
		assert.Len(t, got, 1)
		assert.Equal(t, val, *got[0].Exactly)
	})

	t.Run("value is a chart", func(t *testing.T) {
		bars := []float32{1, 2, 3}

		v := common.MultiFloat{
			Chart: &common.FloatChart{Bars: bars},
		}

		got, err := v.ToSingleFloats()
		assert.NoError(t, err)
		assert.Len(t, got, len(bars))
		for i, bar := range bars {
			assert.Equal(t, bar, *got[i].Exactly)
		}
	})

	t.Run("value is a range", func(t *testing.T) {
		min := float32(10)
		max := float32(20)
		v := common.MultiFloat{
			SingleFloat: common.SingleFloat{
				Between: []float32{min, max},
			},
		}

		got, err := v.ToSingleFloats()
		assert.NoError(t, err)
		assert.Len(t, got, 1)
		assert.Equal(t, min, got[0].Between[0])
		assert.Equal(t, max, got[0].Between[1])
	})
}
