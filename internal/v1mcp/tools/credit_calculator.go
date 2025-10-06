package tools

import (
	"fmt"
	"math"
)

// CreditCalculator provides credit calculation and conversion utilities
type CreditCalculator struct {
	// Credit conversion rates based on Vision One documentation
	// These are approximate values based on typical usage patterns
	ConversionRates map[string]float64

	// Service-specific multipliers
	ServiceMultipliers map[string]float64
}

// NewCreditCalculator creates a new credit calculator with default rates
func NewCreditCalculator() *CreditCalculator {
	return &CreditCalculator{
		ConversionRates: map[string]float64{
			// Endpoint Security
			"endpoint_standard":    2.0,   // Standard endpoint protection
			"endpoint_pro":        10.0,   // Pro license with advanced features
			"endpoint_server":     15.0,   // Server protection

			// Workbench & Investigations
			"alert_investigation":  5.0,   // Per alert investigation
			"threat_hunting":      10.0,   // Threat hunting query
			"custom_detection":     3.0,   // Custom detection rule

			// Sandbox Analysis
			"sandbox_file":        15.0,   // File analysis
			"sandbox_url":         10.0,   // URL analysis
			"sandbox_priority":    25.0,   // Priority/advanced analysis

			// Data Lake & Search
			"search_gb":            1.0,   // Per GB searched
			"data_retention_day":   0.5,   // Per GB per day retained
			"custom_pipeline":     20.0,   // Custom data pipeline

			// CREM (Cyber Risk Exposure Management)
			"risk_assessment":      8.0,   // Per asset risk assessment
			"attack_surface_scan": 12.0,   // Attack surface discovery scan
			"vulnerability_scan":   5.0,   // Vulnerability assessment

			// OAT (Observed Attack Techniques)
			"oat_detection":        3.0,   // Per detection
			"oat_correlation":      5.0,   // Advanced correlation
			"mitre_mapping":        2.0,   // MITRE ATT&CK mapping
		},

		ServiceMultipliers: map[string]float64{
			"endpoint_security":     1.0,
			"workbench":            1.2,  // Higher multiplier for investigation services
			"sandbox":              1.5,  // Sandbox is resource-intensive
			"data_lake":            0.8,  // Data lake has economies of scale
			"crem":                 1.3,  // CREM requires dedicated allocation
			"oat":                  1.0,
		},
	}
}

// CalculateEndpointCredits calculates credits for endpoint deployment
func (c *CreditCalculator) CalculateEndpointCredits(standardCount, proCount, serverCount int) float64 {
	credits := float64(standardCount)*c.ConversionRates["endpoint_standard"] +
		float64(proCount)*c.ConversionRates["endpoint_pro"] +
		float64(serverCount)*c.ConversionRates["endpoint_server"]

	return credits * c.ServiceMultipliers["endpoint_security"]
}

// CalculateWorkbenchCredits calculates credits for workbench activity
func (c *CreditCalculator) CalculateWorkbenchCredits(alertsPerDay, huntingQueriesPerDay, customRules int) float64 {
	monthlyAlerts := float64(alertsPerDay * 30)
	monthlyHunting := float64(huntingQueriesPerDay * 30)

	credits := monthlyAlerts*c.ConversionRates["alert_investigation"] +
		monthlyHunting*c.ConversionRates["threat_hunting"] +
		float64(customRules)*c.ConversionRates["custom_detection"]

	return credits * c.ServiceMultipliers["workbench"]
}

// CalculateSandboxCredits calculates credits for sandbox usage
func (c *CreditCalculator) CalculateSandboxCredits(filesPerDay, urlsPerDay, priorityPerDay int) float64 {
	monthlyFiles := float64(filesPerDay * 30)
	monthlyUrls := float64(urlsPerDay * 30)
	monthlyPriority := float64(priorityPerDay * 30)

	credits := monthlyFiles*c.ConversionRates["sandbox_file"] +
		monthlyUrls*c.ConversionRates["sandbox_url"] +
		monthlyPriority*c.ConversionRates["sandbox_priority"]

	return credits * c.ServiceMultipliers["sandbox"]
}

// CalculateDataLakeCredits calculates credits for data lake usage
func (c *CreditCalculator) CalculateDataLakeCredits(searchGBPerDay, retentionGB, retentionDays, pipelines int) float64 {
	monthlySearch := float64(searchGBPerDay * 30)
	retentionCredits := float64(retentionGB) * float64(retentionDays) * c.ConversionRates["data_retention_day"]
	pipelineCredits := float64(pipelines) * c.ConversionRates["custom_pipeline"]

	credits := monthlySearch*c.ConversionRates["search_gb"] +
		retentionCredits +
		pipelineCredits

	return credits * c.ServiceMultipliers["data_lake"]
}

// CalculateCREMCredits calculates credits for CREM features
func (c *CreditCalculator) CalculateCREMCredits(assetsScanned, surfaceScansPerMonth, vulnScansPerMonth int) float64 {
	assessmentCredits := float64(assetsScanned) * c.ConversionRates["risk_assessment"]
	surfaceCredits := float64(surfaceScansPerMonth) * c.ConversionRates["attack_surface_scan"]
	vulnCredits := float64(vulnScansPerMonth) * c.ConversionRates["vulnerability_scan"]

	credits := assessmentCredits + surfaceCredits + vulnCredits

	return credits * c.ServiceMultipliers["crem"]
}

// CalculateOATCredits calculates credits for OAT usage
func (c *CreditCalculator) CalculateOATCredits(detectionsPerDay, correlationsPerDay, mitreMappings int) float64 {
	monthlyDetections := float64(detectionsPerDay * 30)
	monthlyCorrelations := float64(correlationsPerDay * 30)

	credits := monthlyDetections*c.ConversionRates["oat_detection"] +
		monthlyCorrelations*c.ConversionRates["oat_correlation"] +
		float64(mitreMappings)*c.ConversionRates["mitre_mapping"]

	return credits * c.ServiceMultipliers["oat"]
}

// EstimateMonthlyCost estimates the monthly cost in USD based on credits
func (c *CreditCalculator) EstimateMonthlyCost(totalCredits float64) float64 {
	// Approximate cost per credit (varies by contract and volume)
	costPerCredit := 5.0 // $5 per credit is a rough estimate
	return totalCredits * costPerCredit
}

// GetOptimizationRecommendations provides credit optimization recommendations
func (c *CreditCalculator) GetOptimizationRecommendations(usage map[string]float64) []string {
	recommendations := []string{}
	total := 0.0

	for _, credits := range usage {
		total += credits
	}

	if total == 0 {
		return []string{"No usage data available for analysis"}
	}

	// Check for imbalanced usage
	for serviceName, credits := range usage {
		percentage := (credits / total) * 100

		if percentage > 40 {
			recommendations = append(recommendations,
				fmt.Sprintf("%s is consuming %.1f%% of credits. Consider optimization strategies.", serviceName, percentage))
		}

		if percentage < 5 && credits > 0 {
			recommendations = append(recommendations,
				fmt.Sprintf("%s is underutilized (%.1f%%). Consider reducing allocation or increasing usage.", serviceName, percentage))
		}
	}

	// Service-specific recommendations
	if endpointCredits, ok := usage["endpoint_security"]; ok && endpointCredits > 5000 {
		recommendations = append(recommendations,
			"High endpoint security costs. Review Pro license distribution and consider Standard licenses for low-risk endpoints.")
	}

	if sandboxCredits, ok := usage["sandbox"]; ok && sandboxCredits > 3000 {
		recommendations = append(recommendations,
			"High sandbox usage. Implement file reputation checking and hash-based deduplication to reduce submissions.")
	}

	if workbenchCredits, ok := usage["workbench"]; ok && workbenchCredits > 4000 {
		recommendations = append(recommendations,
			"High investigation costs. Consider automated triage and SOAR integration to reduce manual investigation time.")
	}

	if dataLakeCredits, ok := usage["data_lake"]; ok && dataLakeCredits > 2000 {
		recommendations = append(recommendations,
			"High data lake costs. Review retention policies and optimize search queries for efficiency.")
	}

	return recommendations
}

// CalculateSavingsOpportunity estimates potential savings from optimizations
func (c *CreditCalculator) CalculateSavingsOpportunity(currentUsage, optimizedUsage map[string]float64) (float64, float64) {
	currentTotal := 0.0
	optimizedTotal := 0.0

	for _, credits := range currentUsage {
		currentTotal += credits
	}

	for _, credits := range optimizedUsage {
		optimizedTotal += credits
	}

	savingsCredits := currentTotal - optimizedTotal
	savingsDollars := c.EstimateMonthlyCost(savingsCredits)

	return savingsCredits, savingsDollars
}

// FormatCreditReport generates a formatted credit usage report
func (c *CreditCalculator) FormatCreditReport(usage map[string]float64) string {
	report := "Vision One Credit Usage Report\n"
	report += "================================\n\n"

	total := 0.0
	for _, credits := range usage {
		total += credits
	}

	report += "Service Breakdown:\n"
	for service, credits := range usage {
		percentage := 0.0
		if total > 0 {
			percentage = (credits / total) * 100
		}
		cost := c.EstimateMonthlyCost(credits)
		report += fmt.Sprintf("  %-20s: %8.0f credits (%5.1f%%) - $%,.2f/month\n",
			service, credits, percentage, cost)
	}

	report += fmt.Sprintf("\nTotal Monthly Credits: %.0f\n", total)
	report += fmt.Sprintf("Estimated Monthly Cost: $%,.2f\n", c.EstimateMonthlyCost(total))

	// Add efficiency rating
	efficiency := c.calculateEfficiencyRating(usage, total)
	report += fmt.Sprintf("Efficiency Rating: %s\n", efficiency)

	return report
}

// calculateEfficiencyRating determines the overall efficiency of credit usage
func (c *CreditCalculator) calculateEfficiencyRating(usage map[string]float64, total float64) string {
	if total == 0 {
		return "N/A"
	}

	// Calculate variance to determine if usage is balanced
	mean := total / float64(len(usage))
	variance := 0.0

	for _, credits := range usage {
		diff := credits - mean
		variance += diff * diff
	}

	variance = variance / float64(len(usage))
	stdDev := math.Sqrt(variance)
	coefficientOfVariation := stdDev / mean

	// Lower CV means more balanced usage
	if coefficientOfVariation < 0.5 {
		return "Excellent (Well-balanced)"
	} else if coefficientOfVariation < 1.0 {
		return "Good (Mostly balanced)"
	} else if coefficientOfVariation < 1.5 {
		return "Fair (Some imbalance)"
	} else {
		return "Poor (Highly imbalanced)"
	}
}