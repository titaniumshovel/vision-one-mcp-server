package tools

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/trendmicro/vision-one-mcp-server/internal/v1client"
)

// CreditsAnalysisData represents analyzed credit usage data
type CreditsAnalysisData struct {
	Module           string                 `json:"module"`
	Usage            map[string]interface{} `json:"usage"`
	Recommendations  []string               `json:"recommendations"`
	EstimatedCredits int                    `json:"estimatedCredits,omitempty"`
	Efficiency       string                 `json:"efficiency,omitempty"`
}

// analyzeEndpointCredits analyzes endpoint security credit usage
func analyzeEndpointCredits(client *v1client.V1ApiClient, filter string, top int) (*CreditsAnalysisData, error) {
	queryParams := v1client.QueryParameters{
		Top: top,
	}

	response, err := client.CreditsListEndpoints(filter, queryParams)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve endpoints: %v", err)
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %v", err)
	}

	var data map[string]interface{}
	if err := json.Unmarshal(body, &data); err != nil {
		return nil, fmt.Errorf("failed to parse response: %v", err)
	}

	// Analyze endpoint data for credit usage patterns
	analysis := &CreditsAnalysisData{
		Module: "Endpoint Security",
		Usage:  make(map[string]interface{}),
	}

	if items, ok := data["items"].([]interface{}); ok {
		totalEndpoints := len(items)
		proLicenses := 0
		standardLicenses := 0

		for _, item := range items {
			if endpoint, ok := item.(map[string]interface{}); ok {
				// Check for Pro license features
				if licenses, ok := endpoint["creditAllocatedLicenses"].([]interface{}); ok {
					for _, license := range licenses {
						if licStr, ok := license.(string); ok {
							if strings.Contains(licStr, "Pro") || strings.Contains(licStr, "Advanced") {
								proLicenses++
							} else {
								standardLicenses++
							}
						}
					}
				}
			}
		}

		analysis.Usage["totalEndpoints"] = totalEndpoints
		analysis.Usage["proLicenses"] = proLicenses
		analysis.Usage["standardLicenses"] = standardLicenses

		// Calculate estimated credits (Pro licenses consume more credits)
		analysis.EstimatedCredits = (proLicenses * 10) + (standardLicenses * 2)

		// Generate recommendations
		if proLicenses > 0 {
			utilization := float64(proLicenses) / float64(totalEndpoints) * 100
			if utilization < 30 {
				analysis.Recommendations = append(analysis.Recommendations,
					fmt.Sprintf("Pro license utilization is low (%.1f%%). Consider reducing Pro licenses to save credits.", utilization))
			}
			analysis.Usage["proUtilization"] = fmt.Sprintf("%.1f%%", utilization)
		}

		if totalEndpoints > 1000 {
			analysis.Recommendations = append(analysis.Recommendations,
				"High endpoint count detected. Consider implementing endpoint grouping to optimize credit usage.")
		}
	}

	return analysis, nil
}

// analyzeWorkbenchCredits analyzes workbench investigation credit usage
func analyzeWorkbenchCredits(client *v1client.V1ApiClient, filter string, startTime, endTime time.Time) (*CreditsAnalysisData, error) {
	queryParams := v1client.QueryParameters{
		Top:           100,
		StartDateTime: startTime,
		EndDateTime:   endTime,
	}

	response, err := client.CreditsListWorkbenchAlerts(filter, queryParams)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve workbench alerts: %v", err)
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %v", err)
	}

	var data map[string]interface{}
	if err := json.Unmarshal(body, &data); err != nil {
		return nil, fmt.Errorf("failed to parse response: %v", err)
	}

	analysis := &CreditsAnalysisData{
		Module: "Workbench Investigations",
		Usage:  make(map[string]interface{}),
	}

	if items, ok := data["items"].([]interface{}); ok {
		totalAlerts := len(items)
		highSeverity := 0
		investigated := 0

		for _, item := range items {
			if alert, ok := item.(map[string]interface{}); ok {
				// Count severity levels
				if severity, ok := alert["severity"].(string); ok {
					if severity == "high" || severity == "critical" {
						highSeverity++
					}
				}

				// Count investigated alerts
				if status, ok := alert["investigationStatus"].(string); ok {
					if status != "New" {
						investigated++
					}
				}
			}
		}

		analysis.Usage["totalAlerts"] = totalAlerts
		analysis.Usage["highSeverityAlerts"] = highSeverity
		analysis.Usage["investigatedAlerts"] = investigated

		// Calculate daily average
		days := endTime.Sub(startTime).Hours() / 24
		if days > 0 {
			dailyAvg := float64(totalAlerts) / days
			analysis.Usage["dailyAverage"] = fmt.Sprintf("%.1f", dailyAvg)

			// Estimate credits based on investigation activity
			analysis.EstimatedCredits = investigated * 5 // Each investigation consumes ~5 credits

			if dailyAvg > 20 {
				analysis.Recommendations = append(analysis.Recommendations,
					fmt.Sprintf("High alert volume (%.1f/day). Consider tuning detection rules to reduce noise.", dailyAvg))
			}
		}

		if highSeverity > totalAlerts/3 {
			analysis.Recommendations = append(analysis.Recommendations,
				"High percentage of critical alerts. Review detection thresholds to focus on true positives.")
		}
	}

	return analysis, nil
}

// analyzeSandboxCredits analyzes sandbox submission credit usage
func analyzeSandboxCredits(client *v1client.V1ApiClient, filter string, startTime, endTime time.Time) (*CreditsAnalysisData, error) {
	queryParams := v1client.QueryParameters{
		Top:           100,
		StartDateTime: startTime,
		EndDateTime:   endTime,
	}

	response, err := client.CreditsListSandboxSubmissions(filter, queryParams)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve sandbox submissions: %v", err)
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %v", err)
	}

	var data map[string]interface{}
	if err := json.Unmarshal(body, &data); err != nil {
		return nil, fmt.Errorf("failed to parse response: %v", err)
	}

	analysis := &CreditsAnalysisData{
		Module: "Sandbox Analysis",
		Usage:  make(map[string]interface{}),
	}

	if items, ok := data["items"].([]interface{}); ok {
		totalSubmissions := len(items)
		fileSubmissions := 0
		urlSubmissions := 0

		for _, item := range items {
			if submission, ok := item.(map[string]interface{}); ok {
				if subType, ok := submission["type"].(string); ok {
					if subType == "file" {
						fileSubmissions++
					} else if subType == "url" {
						urlSubmissions++
					}
				}
			}
		}

		analysis.Usage["totalSubmissions"] = totalSubmissions
		analysis.Usage["fileSubmissions"] = fileSubmissions
		analysis.Usage["urlSubmissions"] = urlSubmissions

		// Calculate daily average
		days := endTime.Sub(startTime).Hours() / 24
		if days > 0 {
			dailyAvg := float64(totalSubmissions) / days
			analysis.Usage["dailyAverage"] = fmt.Sprintf("%.1f", dailyAvg)

			// Estimate credits (sandbox submissions are expensive)
			analysis.EstimatedCredits = totalSubmissions * 15

			if dailyAvg > 100 {
				analysis.Recommendations = append(analysis.Recommendations,
					"Very high sandbox submission rate. Implement file type filtering to reduce unnecessary analysis.")
			}
		}

		if fileSubmissions > urlSubmissions*3 {
			analysis.Recommendations = append(analysis.Recommendations,
				"Consider implementing hash-based deduplication for file submissions to avoid re-analyzing known files.")
		}
	}

	return analysis, nil
}

// analyzeOATCredits analyzes OAT detection credit usage
func analyzeOATCredits(client *v1client.V1ApiClient, filter string, startTime, endTime time.Time) (*CreditsAnalysisData, error) {
	queryParams := v1client.QueryParameters{
		Top:           100,
		StartDateTime: startTime,
		EndDateTime:   endTime,
	}

	response, err := client.CreditsListOATDetections(filter, queryParams)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve OAT detections: %v", err)
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %v", err)
	}

	var data map[string]interface{}
	if err := json.Unmarshal(body, &data); err != nil {
		return nil, fmt.Errorf("failed to parse response: %v", err)
	}

	analysis := &CreditsAnalysisData{
		Module: "OAT Detections",
		Usage:  make(map[string]interface{}),
	}

	if items, ok := data["items"].([]interface{}); ok {
		totalDetections := len(items)
		riskLevels := make(map[string]int)
		uniqueEndpoints := make(map[string]bool)

		for _, item := range items {
			if detection, ok := item.(map[string]interface{}); ok {
				// Count risk levels
				if risk, ok := detection["riskLevel"].(string); ok {
					riskLevels[risk]++
				}

				// Track unique endpoints
				if endpoint, ok := detection["endpointName"].(string); ok {
					uniqueEndpoints[endpoint] = true
				}
			}
		}

		analysis.Usage["totalDetections"] = totalDetections
		analysis.Usage["uniqueEndpoints"] = len(uniqueEndpoints)
		analysis.Usage["riskLevels"] = riskLevels

		// Calculate daily average
		days := endTime.Sub(startTime).Hours() / 24
		if days > 0 {
			dailyAvg := float64(totalDetections) / days
			analysis.Usage["dailyAverage"] = fmt.Sprintf("%.1f", dailyAvg)

			// Estimate credits
			analysis.EstimatedCredits = totalDetections * 3

			if dailyAvg > 50 {
				analysis.Recommendations = append(analysis.Recommendations,
					"High OAT detection volume. Review detection rules to reduce false positives.")
			}
		}

		highRisk := riskLevels["high"] + riskLevels["critical"]
		if highRisk > totalDetections/2 {
			analysis.Recommendations = append(analysis.Recommendations,
				"Majority of detections are high-risk. This indicates either severe threats or overly sensitive rules.")
		}
	}

	return analysis, nil
}