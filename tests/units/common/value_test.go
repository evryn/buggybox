package values_test

import (
	"kermoo/modules/values"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetValue(t *testing.T) {
	t.Run("percentage is exactly set", func(t *testing.T) {
		val := float32(42)
		s := &values.SingleFloat{
			Exactly: &val,
		}

		got, err := s.ToFloat()
		assert.NoError(t, err)
		assert.Equal(t, val, got)
	})

	t.Run("percentage is a range", func(t *testing.T) {
		min := float32(10)
		max := float32(20)
		s := &values.SingleFloat{
			Between: []float32{min, max},
		}

		got, err := s.ToFloat()
		assert.NoError(t, err)
		assert.GreaterOrEqual(t, got, min)
		assert.LessOrEqual(t, got, max)
	})
}

func TestToSingleValues(t *testing.T) {
	t.Run("percentage is exactly set", func(t *testing.T) {
		val := float32(42)
		v := values.MultiFloat{
			SingleFloat: values.SingleFloat{
				Exactly: &val,
			},
		}

		got, err := v.ToSingleFloats()
		assert.NoError(t, err)
		assert.Len(t, got, 1)
		assert.Equal(t, val, *got[0].Exactly)
	})

	t.Run("percentage is a chart", func(t *testing.T) {
		bars := []float32{1, 2, 3}

		v := values.MultiFloat{
			Chart: &values.FloatChart{Bars: bars},
		}

		got, err := v.ToSingleFloats()
		assert.NoError(t, err)
		assert.Len(t, got, len(bars))
		for i, bar := range bars {
			assert.Equal(t, bar, *got[i].Exactly)
		}
	})

	t.Run("percentage is a range", func(t *testing.T) {
		min := float32(10)
		max := float32(20)
		v := values.MultiFloat{
			SingleFloat: values.SingleFloat{
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
