package tools

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/mark3labs/mcp-go/mcp"
	mcpserver "github.com/mark3labs/mcp-go/server"
	"github.com/trendmicro/vision-one-mcp-server/internal/v1client"
)

var ToolsetsReadOnlyCredits = []func(client *v1client.V1ApiClient) mcpserver.ServerTool{
	toolCreditsEndpointSecurityAnalysis,
	toolCreditsWorkbenchAlertsAnalysis,
	toolCreditsSandboxUsageAnalysis,
	toolCreditsOATDetectionsAnalysis,
	toolCreditsSearchStatisticsAnalysis,
	toolCreditsComprehensiveAnalysis,
	toolCreditCalculator,
	toolCreditsUserInput,
	toolCreditsUtilizationReport,
}

func toolCreditsEndpointSecurityAnalysis(client *v1client.V1ApiClient) mcpserver.ServerTool {
	return mcpserver.ServerTool{
		Tool: mcp.NewTool(
			"credits_endpoint_security_analysis",
			mcp.WithDescription("Analyze endpoint security credit usage including Pro licenses and security features"),
			mcp.WithToolAnnotation(mcp.ToolAnnotation{
				ReadOnlyHint: toPtr(true),
			}),
			mcp.WithString("filter", mcp.Description("Filter to apply to endpoint query")),
			mcp.WithString("top",
				mcp.Description("Number of endpoints to analyze (default: 100, use 'all' for comprehensive analysis)"),
				mcp.Enum("10", "50", "100", "500", "1000", "all"),
			),
			mcp.WithString("fetchAll",
				mcp.Description("Whether to fetch all endpoints (may be slow for large environments)"),
				mcp.Enum("true", "false"),
			),
		),
		Handler: func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			top, err := optionalStrInt("top", request.Params.Arguments)
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}

			filter, err := optionalValue[string]("filter", request.Params.Arguments)
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}

			queryParams := v1client.QueryParameters{
				Top: top,
			}

			response, err := client.CreditsListEndpoints(filter, queryParams)
			if err != nil {
				return mcp.NewToolResultError(fmt.Sprintf("Failed to retrieve endpoints: %v", err)), nil
			}

			return handleStatusResponse(response, nil, http.StatusOK, "failed to analyze endpoint security credits")
		},
	}
}


func toolCreditsOATDetectionsAnalysis(client *v1client.V1ApiClient) mcpserver.ServerTool {
	return mcpserver.ServerTool{
		Tool: mcp.NewTool(
			"credits_oat_detections_analysis",
			mcp.WithDescription("Analyze Observed Attack Techniques (OAT) detections for credit usage"),
			mcp.WithToolAnnotation(mcp.ToolAnnotation{
				ReadOnlyHint: toPtr(true),
			}),
			mcp.WithString("filter", mcp.Description("Filter to apply to OAT detections query")),
			mcp.WithString("top",
				mcp.Description("Number of detections to analyze"),
				mcp.Enum("10", "50", "100", "500"),
			),
			mcp.WithString("startDateTime",
				mcp.Description("Start time for detection analysis in ISO 8601 format"),
			),
			mcp.WithString("endDateTime",
				mcp.Description("End time for detection analysis in ISO 8601 format"),
			),
		),
		Handler: func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			top, err := optionalStrInt("top", request.Params.Arguments)
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}

			filter, err := optionalValue[string]("filter", request.Params.Arguments)
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}

			startDateTime, err := optionalTimeValue("startDateTime", request.Params.Arguments)
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}

			endDateTime, err := optionalTimeValue("endDateTime", request.Params.Arguments)
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}

			queryParams := v1client.QueryParameters{
				Top:           top,
				StartDateTime: startDateTime,
				EndDateTime:   endDateTime,
			}

			response, err := client.CreditsListOATDetections(filter, queryParams)
			if err != nil {
				return mcp.NewToolResultError(fmt.Sprintf("Failed to retrieve OAT detections: %v", err)), nil
			}

			return handleStatusResponse(response, nil, http.StatusOK, "failed to analyze OAT detections")
		},
	}
}

func toolCreditsSandboxUsageAnalysis(client *v1client.V1ApiClient) mcpserver.ServerTool {
	return mcpserver.ServerTool{
		Tool: mcp.NewTool(
			"credits_sandbox_usage_analysis",
			mcp.WithDescription("Analyze sandbox submission usage for credit consumption"),
			mcp.WithToolAnnotation(mcp.ToolAnnotation{
				ReadOnlyHint: toPtr(true),
			}),
			mcp.WithString("filter", mcp.Description("Filter to apply to sandbox submissions query")),
			mcp.WithString("top",
				mcp.Description("Number of submissions to analyze"),
				mcp.Enum("10", "50", "100", "500"),
			),
			mcp.WithString("startDateTime",
				mcp.Description("Start time for submission analysis in ISO 8601 format"),
			),
			mcp.WithString("endDateTime",
				mcp.Description("End time for submission analysis in ISO 8601 format"),
			),
		),
		Handler: func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			top, err := optionalStrInt("top", request.Params.Arguments)
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}

			filter, err := optionalValue[string]("filter", request.Params.Arguments)
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}

			startDateTime, err := optionalTimeValue("startDateTime", request.Params.Arguments)
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}

			endDateTime, err := optionalTimeValue("endDateTime", request.Params.Arguments)
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}

			queryParams := v1client.QueryParameters{
				Top:           top,
				StartDateTime: startDateTime,
				EndDateTime:   endDateTime,
			}

			response, err := client.CreditsListSandboxSubmissions(filter, queryParams)
			if err != nil {
				return mcp.NewToolResultError(fmt.Sprintf("Failed to retrieve sandbox submissions: %v", err)), nil
			}

			return handleStatusResponse(response, nil, http.StatusOK, "failed to analyze sandbox usage")
		},
	}
}

func toolCreditsWorkbenchAlertsAnalysis(client *v1client.V1ApiClient) mcpserver.ServerTool {
	return mcpserver.ServerTool{
		Tool: mcp.NewTool(
			"credits_workbench_alerts_analysis",
			mcp.WithDescription("Analyze workbench alert investigation activity for credits"),
			mcp.WithToolAnnotation(mcp.ToolAnnotation{
				ReadOnlyHint: toPtr(true),
			}),
			mcp.WithString("filter", mcp.Description("Filter to apply to workbench alerts query")),
			mcp.WithString("top",
				mcp.Description("Number of alerts to analyze"),
				mcp.Enum("10", "50", "100", "500"),
			),
			mcp.WithString("startDateTime",
				mcp.Description("Start time for alert analysis in ISO 8601 format"),
			),
			mcp.WithString("endDateTime",
				mcp.Description("End time for alert analysis in ISO 8601 format"),
			),
		),
		Handler: func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			top, err := optionalStrInt("top", request.Params.Arguments)
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}

			filter, err := optionalValue[string]("filter", request.Params.Arguments)
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}

			startDateTime, err := optionalTimeValue("startDateTime", request.Params.Arguments)
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}

			endDateTime, err := optionalTimeValue("endDateTime", request.Params.Arguments)
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}

			queryParams := v1client.QueryParameters{
				Top:           top,
				StartDateTime: startDateTime,
				EndDateTime:   endDateTime,
			}

			response, err := client.CreditsListWorkbenchAlerts(filter, queryParams)
			if err != nil {
				return mcp.NewToolResultError(fmt.Sprintf("Failed to retrieve workbench alerts: %v", err)), nil
			}

			return handleStatusResponse(response, nil, http.StatusOK, "failed to analyze workbench alerts")
		},
	}
}

func toolCreditsSearchStatisticsAnalysis(client *v1client.V1ApiClient) mcpserver.ServerTool {
	return mcpserver.ServerTool{
		Tool: mcp.NewTool(
			"credits_search_statistics_analysis",
			mcp.WithDescription("Analyze search activity and sensor statistics for credit usage"),
			mcp.WithToolAnnotation(mcp.ToolAnnotation{
				ReadOnlyHint: toPtr(true),
			}),
			mcp.WithString("filter", mcp.Description("Filter to apply to search statistics query")),
			mcp.WithString("startDateTime",
				mcp.Description("Start time for search statistics analysis in ISO 8601 format"),
			),
			mcp.WithString("endDateTime",
				mcp.Description("End time for search statistics analysis in ISO 8601 format"),
			),
		),
		Handler: func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			filter, err := optionalValue[string]("filter", request.Params.Arguments)
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}

			startDateTime, err := optionalTimeValue("startDateTime", request.Params.Arguments)
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}

			endDateTime, err := optionalTimeValue("endDateTime", request.Params.Arguments)
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}

			queryParams := v1client.QueryParameters{
				StartDateTime: startDateTime,
				EndDateTime:   endDateTime,
			}

			response, err := client.CreditsGetSearchStatistics(filter, queryParams)
			if err != nil {
				return mcp.NewToolResultError(fmt.Sprintf("Failed to retrieve search statistics: %v", err)), nil
			}

			return handleStatusResponse(response, nil, http.StatusOK, "failed to analyze search statistics")
		},
	}
}

func toolCreditsComprehensiveAnalysis(client *v1client.V1ApiClient) mcpserver.ServerTool {
	return mcpserver.ServerTool{
		Tool: mcp.NewTool(
			"credits_comprehensive_analysis",
			mcp.WithDescription("Run comprehensive credit usage analysis across all Vision One modules"),
			mcp.WithToolAnnotation(mcp.ToolAnnotation{
				ReadOnlyHint: toPtr(true),
			}),
			mcp.WithString("sampleSize",
				mcp.Description("Number of items to analyze per module"),
				mcp.Enum("10", "50", "100", "500"),
			),
			mcp.WithString("startDateTime",
				mcp.Description("Start time for analysis in ISO 8601 format"),
			),
			mcp.WithString("endDateTime",
				mcp.Description("End time for analysis in ISO 8601 format"),
			),
		),
		Handler: func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			sampleSizeStr, _ := optionalValue[string]("sampleSize", request.Params.Arguments)
			if sampleSizeStr == "" {
				sampleSizeStr = "100"
			}
			sampleSize := 100
			if size, err := strconv.Atoi(sampleSizeStr); err == nil {
				sampleSize = size
			}

			startTime, _ := optionalTimeValue("startDateTime", request.Params.Arguments)
			endTime, _ := optionalTimeValue("endDateTime", request.Params.Arguments)

			// Set default time range if not provided
			if endTime.IsZero() {
				endTime = time.Now()
			}

			if startTime.IsZero() {
				startTime = endTime.AddDate(0, 0, -30)
			}

			return performComprehensiveAnalysis(client, sampleSize, startTime, endTime)
		},
	}
}

// performComprehensiveAnalysis performs a complete credit usage analysis
func performComprehensiveAnalysis(client *v1client.V1ApiClient, sampleSize int, startTime, endTime time.Time) (*mcp.CallToolResult, error) {
	result := "Comprehensive Vision One Credit Usage Analysis\n"
	result += "==========================================\n\n"
	result += fmt.Sprintf("Analysis Period: %s to %s\n", startTime.Format("2006-01-02"), endTime.Format("2006-01-02"))
	result += fmt.Sprintf("Sample Size: %d items per module\n\n", sampleSize)

	calculator := NewCreditCalculator()
	totalUsage := make(map[string]float64)

	// Analyze Endpoint Security
	endpointAnalysis, err := analyzeEndpointCredits(client, "", sampleSize)
	if err == nil && endpointAnalysis != nil {
		result += fmt.Sprintf("📊 Endpoint Security:\n")
		result += fmt.Sprintf("   Total Endpoints: %v\n", endpointAnalysis.Usage["totalEndpoints"])
		result += fmt.Sprintf("   Pro Licenses: %v\n", endpointAnalysis.Usage["proLicenses"])
		result += fmt.Sprintf("   Estimated Credits: %d\n", endpointAnalysis.EstimatedCredits)
		totalUsage["endpoint_security"] = float64(endpointAnalysis.EstimatedCredits)

		for _, rec := range endpointAnalysis.Recommendations {
			result += fmt.Sprintf("   💡 %s\n", rec)
		}
		result += "\n"
	}

	// Analyze Workbench
	workbenchAnalysis, err := analyzeWorkbenchCredits(client, "", startTime, endTime)
	if err == nil && workbenchAnalysis != nil {
		result += fmt.Sprintf("🔍 Workbench Investigations:\n")
		result += fmt.Sprintf("   Total Alerts: %v\n", workbenchAnalysis.Usage["totalAlerts"])
		result += fmt.Sprintf("   Daily Average: %v\n", workbenchAnalysis.Usage["dailyAverage"])
		result += fmt.Sprintf("   Estimated Credits: %d\n", workbenchAnalysis.EstimatedCredits)
		totalUsage["workbench"] = float64(workbenchAnalysis.EstimatedCredits)

		for _, rec := range workbenchAnalysis.Recommendations {
			result += fmt.Sprintf("   💡 %s\n", rec)
		}
		result += "\n"
	}

	// Analyze Sandbox
	sandboxAnalysis, err := analyzeSandboxCredits(client, "", startTime, endTime)
	if err == nil && sandboxAnalysis != nil {
		result += fmt.Sprintf("🧪 Sandbox Analysis:\n")
		result += fmt.Sprintf("   Total Submissions: %v\n", sandboxAnalysis.Usage["totalSubmissions"])
		result += fmt.Sprintf("   File Submissions: %v\n", sandboxAnalysis.Usage["fileSubmissions"])
		result += fmt.Sprintf("   Estimated Credits: %d\n", sandboxAnalysis.EstimatedCredits)
		totalUsage["sandbox"] = float64(sandboxAnalysis.EstimatedCredits)

		for _, rec := range sandboxAnalysis.Recommendations {
			result += fmt.Sprintf("   💡 %s\n", rec)
		}
		result += "\n"
	}

	// Generate overall recommendations
	result += "📈 Overall Credit Usage Summary:\n"
	result += calculator.FormatCreditReport(totalUsage)
	result += "\n"

	recommendations := calculator.GetOptimizationRecommendations(totalUsage)
	if len(recommendations) > 0 {
		result += "🎯 Optimization Recommendations:\n"
		for _, rec := range recommendations {
			result += fmt.Sprintf("   • %s\n", rec)
		}
	}

	return mcp.NewToolResultText(result), nil
}

// toolCreditCalculator provides credit calculation utilities
func toolCreditCalculator(client *v1client.V1ApiClient) mcpserver.ServerTool {
	return mcpserver.ServerTool{
		Tool: mcp.NewTool(
			"credits_calculator",
			mcp.WithDescription("Calculate credit usage based on deployment metrics"),
			mcp.WithToolAnnotation(mcp.ToolAnnotation{
				ReadOnlyHint: toPtr(true),
			}),
			mcp.WithNumber("standardEndpoints", mcp.Description("Number of standard endpoints")),
			mcp.WithNumber("proEndpoints", mcp.Description("Number of Pro-licensed endpoints")),
			mcp.WithNumber("alertsPerDay", mcp.Description("Average alerts investigated per day")),
			mcp.WithNumber("sandboxFilesPerDay", mcp.Description("Average files submitted to sandbox per day")),
			mcp.WithNumber("searchGBPerDay", mcp.Description("Average GB searched in data lake per day")),
		),
		Handler: func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			standardEndpoints, _ := optionalValue[float64]("standardEndpoints", request.Params.Arguments)
			proEndpoints, _ := optionalValue[float64]("proEndpoints", request.Params.Arguments)
			alertsPerDay, _ := optionalValue[float64]("alertsPerDay", request.Params.Arguments)
			sandboxFilesPerDay, _ := optionalValue[float64]("sandboxFilesPerDay", request.Params.Arguments)
			searchGBPerDay, _ := optionalValue[float64]("searchGBPerDay", request.Params.Arguments)

			calculator := NewCreditCalculator()
			usage := make(map[string]float64)

			// Calculate endpoint credits
			if standardEndpoints > 0 || proEndpoints > 0 {
				usage["endpoint_security"] = calculator.CalculateEndpointCredits(
					int(standardEndpoints), int(proEndpoints), 0)
			}

			// Calculate workbench credits
			if alertsPerDay > 0 {
				usage["workbench"] = calculator.CalculateWorkbenchCredits(int(alertsPerDay), 0, 0)
			}

			// Calculate sandbox credits
			if sandboxFilesPerDay > 0 {
				usage["sandbox"] = calculator.CalculateSandboxCredits(int(sandboxFilesPerDay), 0, 0)
			}

			// Calculate data lake credits
			if searchGBPerDay > 0 {
				usage["data_lake"] = calculator.CalculateDataLakeCredits(int(searchGBPerDay), 0, 0, 0)
			}

			result := calculator.FormatCreditReport(usage)
			result += "\n" + strings.Join(calculator.GetOptimizationRecommendations(usage), "\n")

			return mcp.NewToolResultText(result), nil
		},
	}
}

// toolCreditsUserInput allows users to input their credit allocations
func toolCreditsUserInput(client *v1client.V1ApiClient) mcpserver.ServerTool {
	return mcpserver.ServerTool{
		Tool: mcp.NewTool(
			"credits_user_input",
			mcp.WithDescription("Store user's credit allocation for analysis"),
			mcp.WithToolAnnotation(mcp.ToolAnnotation{
				ReadOnlyHint: toPtr(false),
			}),
			mcp.WithString("service",
				mcp.Description("Service name"),
				mcp.Required(),
			),
			mcp.WithNumber("allocated",
				mcp.Description("Allocated credits for this service"),
				mcp.Required(),
			),
			mcp.WithNumber("used",
				mcp.Description("Currently used credits"),
			),
			mcp.WithString("notes",
				mcp.Description("Additional notes about this allocation"),
			),
		),
		Handler: func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			service, err := requiredValue[string]("service", request.Params.Arguments)
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}

			allocated, err := requiredValue[float64]("allocated", request.Params.Arguments)
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}

			used, _ := optionalValue[float64]("used", request.Params.Arguments)
			notes, _ := optionalValue[string]("notes", request.Params.Arguments)

			storage, err := NewCreditStorage()
			if err != nil {
				return mcp.NewToolResultError(fmt.Sprintf("Failed to initialize storage: %v", err)), nil
			}

			if err := storage.UpdateAllocation(service, allocated, used, notes); err != nil {
				return mcp.NewToolResultError(fmt.Sprintf("Failed to save allocation: %v", err)), nil
			}

			utilization := 0.0
			if allocated > 0 {
				utilization = (used / allocated) * 100
			}

			result := fmt.Sprintf("Credit allocation saved for %s:\n", service)
			result += fmt.Sprintf("  Allocated: %.0f credits\n", allocated)
			result += fmt.Sprintf("  Used: %.0f credits\n", used)
			result += fmt.Sprintf("  Utilization: %.1f%%\n", utilization)
			if notes != "" {
				result += fmt.Sprintf("  Notes: %s\n", notes)
			}

			return mcp.NewToolResultText(result), nil
		},
	}
}

// toolCreditsUtilizationReport generates a utilization report
func toolCreditsUtilizationReport(client *v1client.V1ApiClient) mcpserver.ServerTool {
	return mcpserver.ServerTool{
		Tool: mcp.NewTool(
			"credits_utilization_report",
			mcp.WithDescription("Generate credit utilization report from stored allocations"),
			mcp.WithToolAnnotation(mcp.ToolAnnotation{
				ReadOnlyHint: toPtr(true),
			}),
		),
		Handler: func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			storage, err := NewCreditStorage()
			if err != nil {
				return mcp.NewToolResultError(fmt.Sprintf("Failed to initialize storage: %v", err)), nil
			}

			report, err := storage.GenerateUtilizationReport()
			if err != nil {
				return mcp.NewToolResultError(fmt.Sprintf("Failed to generate report: %v", err)), nil
			}

			// Add recommendations based on utilization
			underutilized, _ := storage.GetUnderutilizedServices(50)
			overutilized, _ := storage.GetOverutilizedServices(85)

			if len(underutilized) > 0 || len(overutilized) > 0 {
				report += "\n💡 Optimization Opportunities:\n"

				for _, service := range underutilized {
					utilization := (service.UsedCredits / service.AllocatedCredits) * 100
					report += fmt.Sprintf("  • %s is underutilized (%.1f%%). Consider reducing allocation.\n",
						service.Service, utilization)
				}

				for _, service := range overutilized {
					utilization := (service.UsedCredits / service.AllocatedCredits) * 100
					report += fmt.Sprintf("  • %s is approaching limit (%.1f%%). Consider increasing allocation.\n",
						service.Service, utilization)
				}
			}

			return mcp.NewToolResultText(report), nil
		},
	}
}

