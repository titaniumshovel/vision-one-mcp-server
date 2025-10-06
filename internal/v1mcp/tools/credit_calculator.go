package tools

import (
	"fmt"
	"math"
)

// CreditCalculator provides credit calculation and conversion utilities
// Based on official Trend Micro documentation (credit-conversion.pdf, Sept 29, 2025)
type CreditCalculator struct {
	// Credit conversion rates from official Trend Vision One documentation
	ConversionRates map[string]float64
}

// NewCreditCalculator creates a new credit calculator with official Trend Micro conversion rates
// Source: Trend Vision One Credit Conversion Rates & Requirements Report (Sept 29, 2025)
func NewCreditCalculator() *CreditCalculator {
	return &CreditCalculator{
		ConversionRates: map[string]float64{
			// Endpoint Security (per deployment/month)
			"endpoint_core":            45.0,  // Trend Vision One Endpoint Security - Core
			"endpoint_essentials":      65.0,  // Trend Vision One Endpoint Security - Essentials
			"endpoint_pro":             300.0, // Trend Vision One Endpoint Security - Pro
			"endpoint_sensor_xdr":      20.0,  // Endpoint Sensor detection and response (XDR add-on)
			"server_workload_advanced": 235.0, // Advanced Server/Workload Protection

			// Sandbox Analysis
			"sandbox_reserved_daily": 50.0, // Per daily reserved submission

			// CREM - Cyber Risk Exposure Management (Desktop/Server per month)
			"crem_core_asset":       20.0, // Per assessed desktop/server (Core)
			"crem_essentials_asset": 50.0, // Per assessed desktop/server (Essentials)

			// CREM - Cloud Account Assessment (annual, tiered by resource count)
			"crem_cloud_tier1": 1000.0, // Up to 500 resources
			"crem_cloud_tier2": 2000.0, // 501-1,000 resources
			"crem_cloud_tier3": 3000.0, // 1,001-1,500 resources
			"crem_cloud_tier4": 4000.0, // 1,501-2,000 resources
			"crem_cloud_tier5": 5000.0, // 2,001-2,500 resources
			"crem_cloud_tier6": 6000.0, // 2,501-3,000 resources
			"crem_cloud_tier7": 7000.0, // 3,001-3,500 resources
			"crem_cloud_tier8": 8000.0, // 3,501+ resources (maximum)

			// Email Security (per user account per month)
			"email_sensor":          5.0,  // Email Sensor (detection only)
			"email_collab_core":     25.0, // Cloud Email and Collaboration Protection (core)
			"email_collab_advanced": 25.0, // Advanced Protection add-on (additional)
			"email_gateway_core":    25.0, // Cloud Email Gateway Protection (core)
			"email_gateway_advanced": 25.0, // Advanced Gateway Protection add-on (additional)

			// Container Security
			"container_k8s_node":       1100.0, // Protected Kubernetes node or Amazon ECS instance
			"container_serverless_pod": 110.0,  // Protected serverless container pod or task

			// Data Lake / Agentic SIEM
			"datalake_ingestion_analytic_gb": 3.0,  // Third-party data ingestion (analytic) per GB
			"datalake_ingestion_archival_gb": 1.0,  // Third-party data ingestion (archival) per GB
			"datalake_retention_analytic_gb": 0.2,  // Per GB per month (analytic)
			"datalake_retention_archival_gb": 0.05, // Per GB per month (archival)
			"data_transfer_tb":               800.0, // Per TB to third-party platforms

			// File Security
			"file_security_per_scan": 0.01,    // 5,000 credits per 500,000 scans = 0.01 per scan
			"file_security_bucket":   9636.0, // Per bucket (unlimited scans)

			// Network Security
			"network_sensor_gbps": 25000.0, // XDR Network Sensor per 1 Gbps of traffic

			// Forensics
			"forensics_evidence_gb": 400.0, // Data allowance for evidence collection per GB

			// Mobile Inventory
			"mobile_sensor_device": 5.0, // Mobile Sensor per device enrollment

			// Threat Intelligence
			"threat_insights_user":   50000.0,  // Per user account
			"threat_intel_feed":      150000.0, // Per year
			"threat_intel_feed_mssp": 500000.0, // For service providers per year

			// NOTE: Alert investigation and workbench costs are NOT fixed per-alert
			// They consume data lake search credits based on investigation complexity
			// No per-alert conversion rate exists in official Trend Micro documentation
		},
	}
}

// CalculateEndpointCredits calculates credits for endpoint deployment (monthly)
// Uses official Trend Micro conversion rates
func (c *CreditCalculator) CalculateEndpointCredits(coreCount, essentialsCount, proCount, serverCount int) float64 {
	credits := float64(coreCount)*c.ConversionRates["endpoint_core"] +
		float64(essentialsCount)*c.ConversionRates["endpoint_essentials"] +
		float64(proCount)*c.ConversionRates["endpoint_pro"] +
		float64(serverCount)*c.ConversionRates["server_workload_advanced"]

	return credits
}

// CalculateWorkbenchCredits provides qualitative analysis, not fixed per-alert costs
// Workbench alert investigations consume data lake search credits based on complexity
// NOTE: There is NO fixed per-alert credit rate in official Trend Micro documentation
func (c *CreditCalculator) CalculateWorkbenchCredits(alertCount int, avgDataSearchedGB float64) float64 {
	// Estimate based on data lake search usage
	// This is an approximation since actual costs vary by investigation complexity
	estimatedCredits := avgDataSearchedGB * c.ConversionRates["datalake_ingestion_analytic_gb"]

	return estimatedCredits
}

// CalculateSandboxCredits calculates credits for sandbox usage (reserved submissions)
func (c *CreditCalculator) CalculateSandboxCredits(dailyReservedSubmissions int) float64 {
	monthlyCredits := float64(dailyReservedSubmissions*30) * c.ConversionRates["sandbox_reserved_daily"]
	return monthlyCredits
}

// CalculateDataLakeCredits calculates credits for data lake usage
// Uses official Trend Micro conversion rates for ingestion and retention
func (c *CreditCalculator) CalculateDataLakeCredits(analyticGB, archivalGB, retentionMonths int) float64 {
	ingestionCredits := float64(analyticGB)*c.ConversionRates["datalake_ingestion_analytic_gb"] +
		float64(archivalGB)*c.ConversionRates["datalake_ingestion_archival_gb"]

	retentionCredits := (float64(analyticGB)*c.ConversionRates["datalake_retention_analytic_gb"] +
		float64(archivalGB)*c.ConversionRates["datalake_retention_archival_gb"]) *
		float64(retentionMonths)

	return ingestionCredits + retentionCredits
}

// CalculateCREMCredits calculates credits for CREM features (Desktop/Server assessment)
// Uses official Trend Micro conversion rates
func (c *CreditCalculator) CalculateCREMCredits(coreAssets, essentialsAssets int) float64 {
	credits := float64(coreAssets)*c.ConversionRates["crem_core_asset"] +
		float64(essentialsAssets)*c.ConversionRates["crem_essentials_asset"]

	return credits
}

// CalculateCREMCloudCredits calculates credits for CREM cloud account assessment
// Tiered pricing based on resource count (annual credits)
func (c *CreditCalculator) CalculateCREMCloudCredits(resourceCount int) float64 {
	if resourceCount <= 500 {
		return c.ConversionRates["crem_cloud_tier1"]
	} else if resourceCount <= 1000 {
		return c.ConversionRates["crem_cloud_tier2"]
	} else if resourceCount <= 1500 {
		return c.ConversionRates["crem_cloud_tier3"]
	} else if resourceCount <= 2000 {
		return c.ConversionRates["crem_cloud_tier4"]
	} else if resourceCount <= 2500 {
		return c.ConversionRates["crem_cloud_tier5"]
	} else if resourceCount <= 3000 {
		return c.ConversionRates["crem_cloud_tier6"]
	} else if resourceCount <= 3500 {
		return c.ConversionRates["crem_cloud_tier7"]
	}
	return c.ConversionRates["crem_cloud_tier8"] // 3501+ resources
}

// NOTE: OAT (Observed Attack Techniques) detections do not have standalone credit costs
// OAT is included as part of XDR sensor deployment (20 credits per endpoint)
// Individual OAT detections do not consume additional credits

// EstimateMonthlyCost provides a rough cost estimate in USD based on credits
// NOTE: Actual credit pricing varies significantly by contract, volume, and region
// Contact Trend Micro sales for accurate pricing
func (c *CreditCalculator) EstimateMonthlyCost(totalCredits float64) float64 {
	// This is a ROUGH estimate only - actual pricing varies by customer contract
	// Typical range is $0.10 - $1.00 per credit depending on volume and commitment
	estimatedCostPerCredit := 0.50 // Mid-range estimate
	return totalCredits * estimatedCostPerCredit
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
		report += fmt.Sprintf("  %-20s: %8.0f credits (%5.1f%%) - $%.2f/month\n",
			service, credits, percentage, cost)
	}

	report += fmt.Sprintf("\nTotal Monthly Credits: %.0f\n", total)
	report += fmt.Sprintf("Estimated Monthly Cost: $%.2f\n", c.EstimateMonthlyCost(total))

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