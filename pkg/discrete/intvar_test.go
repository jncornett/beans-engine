package discrete

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFromPMF(t *testing.T) {
	const (
		sampleSize        = 100
		allowedErrorRatio = 0.1 // 10% variation
	)
	tests := []struct {
		name         string
		pmf          []IntVarPoint
		wantMean     float64
		wantVariance float64
	}{
		{
			name: "empty",
		},
		{
			name: "zero freq",
			pmf:  []IntVarPoint{{X: Const(42), Y: 0}},
		},
		{
			name:     "one",
			pmf:      []IntVarPoint{{X: Const(42), Y: 1}},
			wantMean: 42,
		},
		{
			name: "two",
			pmf: []IntVarPoint{
				{X: Const(42), Y: 3},
				{X: Const(101), Y: 2},
			},
			wantMean:     65,
			wantVariance: 789,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			iv := FromPMF(tt.pmf)
			samples := SampleK(sampleSize, iv)
			mean, variance := summarize(samples)
			meanError := 2 * allowedErrorRatio * mean
			varianceError := 2 * allowedErrorRatio * variance
			assert.InDelta(t, tt.wantMean, mean, meanError)
			assert.InDelta(t, tt.wantVariance, variance, varianceError)
		})
	}
}

func TestRange(t *testing.T) {
	const (
		sampleSize        = 100
		allowedErrorRatio = 0.1 // 10% variation
	)
	tests := []struct {
		name         string
		min, max     int64
		wantMean     float64
		wantVariance float64
	}{
		{
			name: "zero",
		},
		{
			name: "zero size",
			min:  5,
			max:  5,
		},
		{
			name: "negative size",
			min:  10,
			max:  5,
		},
		{
			name:     "one",
			min:      101,
			max:      102,
			wantMean: 101,
		},
		{
			name:         "many",
			min:          0,
			max:          100,
			wantMean:     50,
			wantVariance: 850,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			iv := Range(tt.min, tt.max)
			samples := SampleK(sampleSize, iv)
			mean, variance := summarize(samples)
			meanError := 2 * allowedErrorRatio * mean
			varianceError := 2 * allowedErrorRatio * variance
			assert.InDelta(t, tt.wantMean, mean, meanError)
			assert.InDelta(t, tt.wantVariance, variance, varianceError)
		})
	}

}

func summarize(samples []int64) (mean, variance float64) {
	if len(samples) == 0 {
		return 0, 0
	}
	if len(samples) == 1 {
		return float64(samples[0]), 0
	}
	var sum float64
	for _, s := range samples {
		sum += float64(s)
	}
	mean = sum / float64(len(samples))
	var sse float64
	for _, s := range samples {
		r := float64(s) - mean
		sse += r * r
	}
	return mean, sse / float64(len(samples)-1)
}
