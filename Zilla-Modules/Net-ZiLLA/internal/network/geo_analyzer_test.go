package network

import (
	"context"
	"testing"

	"net-zilla/internal/models"
	"net-zilla/pkg/logger"
)

func TestGeoAnalyzer_AnalyzeRegion(t *testing.T) {
	ga := NewGeoAnalyzer(logger.NewLogger())

	tests := []struct {
		name string
		geo  *models.GeoAnalysis
		want int
	}{
		{
			name: "Clean Region",
			geo:  &models.GeoAnalysis{Country: "USA", IsProxy: false, HostingType: "ISP/Residential"},
			want: 0,
		},
		{
			name: "High Risk Country",
			geo:  &models.GeoAnalysis{Country: "Unknown", IsProxy: false, HostingType: "ISP/Residential"},
			want: 40,
		},
		{
			name: "Proxy Detection",
			geo:  &models.GeoAnalysis{Country: "USA", IsProxy: true, HostingType: "ISP/Residential"},
			want: 30,
		},
		{
			name: "Hosting Provider",
			geo:  &models.GeoAnalysis{Country: "USA", IsProxy: false, HostingType: "Hosting Provider"},
			want: 15,
		},
		{
			name: "Combined Risks",
			geo:  &models.GeoAnalysis{Country: "Unknown", IsProxy: true, HostingType: "Hosting Provider"},
			want: 85, // 40 + 30 + 15
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ga.AnalyzeRegion(context.Background(), tt.geo); got != tt.want {
				t.Errorf("AnalyzeRegion() = %v, want %v", got, tt.want)
			}
		})
	}
}
