package visualization

import (
	"encoding/json"
	"fmt"
	"net-zilla/internal/models"
)

// ExportFormatter handles data serialization into different formats.
type ExportFormatter struct{}

func NewExportFormatter() *ExportFormatter {
	return &ExportFormatter{}
}

// FormatJSON serializes the advanced report into JSON.
func (ef *ExportFormatter) FormatJSON(report *models.AdvancedReport) ([]byte, error) {
	return json.MarshalIndent(report, "", "  ")
}

// FormatCSV creates a simple CSV representation of basic indicators.
func (ef *ExportFormatter) FormatCSV(report *models.AdvancedReport) string {
	header := "Indicator,Value,Type\n"
	row := fmt.Sprintf("Target,%s,URL\n", report.Target)
	if report.BasicAnalysis != nil && report.BasicAnalysis.GeoAnalysis != nil {
		row += fmt.Sprintf("IP,%s,Network\n", report.BasicAnalysis.GeoAnalysis.IP)
	}
	return header + row
}
