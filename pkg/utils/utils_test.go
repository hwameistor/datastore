package utils

import (
	"k8s.io/apimachinery/pkg/api/resource"
	"testing"
)

func TestCapacityRoundUp(t *testing.T) {
	testCases := []struct {
		desc         string
		originalSize float64
		expectedSize string
	}{
		{
			desc:         "bytes to str 1KB",
			originalSize: 1024,
			expectedSize: "1Ki",
		},
		{
			desc:         "bytes to str 1.7KB",
			originalSize: 1.7 * 1024,
			expectedSize: "2Ki",
		},
		{
			desc:         "bytes to str 1MB",
			originalSize: 1024 * 1024,
			expectedSize: "1Mi",
		},
		{
			desc:         "bytes to str 1.5MB",
			originalSize: 1.5 * 1024 * 1024,
			expectedSize: "2Mi",
		},
		{
			desc:         "bytes to str 1GB",
			originalSize: 1024 * 1024 * 1024,
			expectedSize: "1Gi",
		},
		{
			desc:         "bytes to str 1.3GB",
			originalSize: 1.3 * 1024 * 1024 * 1024,
			expectedSize: "2Gi",
		},
		{
			desc:         "bytes to str 1023.8 GB",
			originalSize: 1023.8 * 1024 * 1024 * 1024,
			expectedSize: "1Ti",
		},
		{
			desc:         "bytes to str 1TB",
			originalSize: 1024 * 1024 * 1024 * 1024,
			expectedSize: "1Ti",
		},
		{
			desc:         "bytes to str 1.3TB - don't round up",
			originalSize: 1.5 * 1024 * 1024 * 1024 * 1024,
			expectedSize: "1536Gi",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.desc, func(t *testing.T) {
			roundedCapacity := CapacityRoundUp(int64(testCase.originalSize))
			if roundedCapacity != int64(testCase.originalSize) {
				t.Logf("origin %d, roundedUp %d", int64(testCase.originalSize), roundedCapacity)
			}

			got := resource.NewQuantity(roundedCapacity, resource.BinarySI).String()
			if got != testCase.expectedSize {
				t.Fatalf("expected %s, got %s", testCase.expectedSize, got)
			}
		})
	}
}
