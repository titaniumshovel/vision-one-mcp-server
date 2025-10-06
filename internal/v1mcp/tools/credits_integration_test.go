package tools

import (
	"os"
	"testing"
	"time"

	"github.com/trendmicro/vision-one-mcp-server/internal/v1client"
)

// TestCreditManagementIntegration tests all credit management tools with a real Vision One endpoint
// Run with: go test -v -run TestCreditManagementIntegration
// Set V1_API_KEY and V1_REGION environment variables before running
func TestCreditManagementIntegration(t *testing.T) {
	apiKey := os.Getenv("V1_API_KEY")
	region := os.Getenv("V1_REGION")

	if apiKey == "" {
		t.Skip("Skipping integration test: V1_API_KEY not set")
	}

	if region == "" {
		region = "us" // Default to US region
	}

	client, err := v1client.NewV1ApiClient(v1client.ClientOptions{
		ApiKey: apiKey,
		Region: region,
	})
	if err != nil {
		t.Fatalf("Failed to create V1 API client: %v", err)
	}

	t.Run("EndpointSecurityAnalysis", func(t *testing.T) {
		testEndpointSecurityAnalysis(t, client)
	})

	t.Run("WorkbenchAlertsAnalysis", func(t *testing.T) {
		testWorkbenchAlertsAnalysis(t, client)
	})

	t.Run("SandboxUsageAnalysis", func(t *testing.T) {
		testSandboxUsageAnalysis(t, client)
	})

	t.Run("OATDetectionsAnalysis", func(t *testing.T) {
		testOATDetectionsAnalysis(t, client)
	})

	t.Run("CreditCalculator", func(t *testing.T) {
		testCreditCalculator(t)
	})

	t.Run("CreditStorage", func(t *testing.T) {
		testCreditStorage(t)
	})

	t.Run("ComprehensiveAnalysis", func(t *testing.T) {
		testComprehensiveAnalysis(t, client)
	})
}

func testEndpointSecurityAnalysis(t *testing.T, client *v1client.V1ApiClient) {
	t.Log("Testing endpoint security credit analysis...")

	analysis, err := analyzeEndpointCredits(client, "", 10)
	if err != nil {
		t.Logf("Warning: Failed to analyze endpoints: %v", err)
		return
	}

	if analysis == nil {
		t.Error("Expected analysis result, got nil")
		return
	}

	t.Logf("Module: %s", analysis.Module)
	t.Logf("Total Endpoints: %v", analysis.Usage["totalEndpoints"])
	t.Logf("Pro Licenses: %v", analysis.Usage["proLicenses"])
	t.Logf("Estimated Credits: %d", analysis.EstimatedCredits)

	if len(analysis.Recommendations) > 0 {
		t.Logf("Recommendations:")
		for _, rec := range analysis.Recommendations {
			t.Logf("  - %s", rec)
		}
	}

	// Verify the structure
	if analysis.Module != "Endpoint Security" {
		t.Errorf("Expected module 'Endpoint Security', got '%s'", analysis.Module)
	}

	if analysis.Usage == nil {
		t.Error("Expected usage data, got nil")
	}
}

func testWorkbenchAlertsAnalysis(t *testing.T, client *v1client.V1ApiClient) {
	t.Log("Testing workbench alerts credit analysis...")

	endTime := time.Now()
	startTime := endTime.AddDate(0, 0, -7)

	analysis, err := analyzeWorkbenchCredits(client, "", startTime, endTime)
	if err != nil {
		t.Logf("Warning: Failed to analyze workbench alerts: %v", err)
		return
	}

	if analysis == nil {
		t.Error("Expected analysis result, got nil")
		return
	}

	t.Logf("Module: %s", analysis.Module)
	t.Logf("Total Alerts: %v", analysis.Usage["totalAlerts"])
	t.Logf("High Severity: %v", analysis.Usage["highSeverityAlerts"])
	t.Logf("Estimated Credits: %d", analysis.EstimatedCredits)

	if len(analysis.Recommendations) > 0 {
		t.Logf("Recommendations:")
		for _, rec := range analysis.Recommendations {
			t.Logf("  - %s", rec)
		}
	}

	if analysis.Module != "Workbench Investigations" {
		t.Errorf("Expected module 'Workbench Investigations', got '%s'", analysis.Module)
	}
}

func testSandboxUsageAnalysis(t *testing.T, client *v1client.V1ApiClient) {
	t.Log("Testing sandbox usage credit analysis...")

	endTime := time.Now()
	startTime := endTime.AddDate(0, 0, -7)

	analysis, err := analyzeSandboxCredits(client, "", startTime, endTime)
	if err != nil {
		t.Logf("Warning: Failed to analyze sandbox usage: %v", err)
		return
	}

	if analysis == nil {
		t.Error("Expected analysis result, got nil")
		return
	}

	t.Logf("Module: %s", analysis.Module)
	t.Logf("Total Submissions: %v", analysis.Usage["totalSubmissions"])
	t.Logf("File Submissions: %v", analysis.Usage["fileSubmissions"])
	t.Logf("Estimated Credits: %d", analysis.EstimatedCredits)

	if len(analysis.Recommendations) > 0 {
		t.Logf("Recommendations:")
		for _, rec := range analysis.Recommendations {
			t.Logf("  - %s", rec)
		}
	}

	if analysis.Module != "Sandbox Analysis" {
		t.Errorf("Expected module 'Sandbox Analysis', got '%s'", analysis.Module)
	}
}

func testOATDetectionsAnalysis(t *testing.T, client *v1client.V1ApiClient) {
	t.Log("Testing OAT detections credit analysis...")

	endTime := time.Now()
	startTime := endTime.AddDate(0, 0, -7)

	analysis, err := analyzeOATCredits(client, "", startTime, endTime)
	if err != nil {
		t.Logf("Warning: Failed to analyze OAT detections: %v", err)
		return
	}

	if analysis == nil {
		t.Error("Expected analysis result, got nil")
		return
	}

	t.Logf("Module: %s", analysis.Module)
	t.Logf("Total Detections: %v", analysis.Usage["totalDetections"])
	t.Logf("Unique Endpoints: %v", analysis.Usage["uniqueEndpoints"])
	t.Logf("Estimated Credits: %d", analysis.EstimatedCredits)

	if len(analysis.Recommendations) > 0 {
		t.Logf("Recommendations:")
		for _, rec := range analysis.Recommendations {
			t.Logf("  - %s", rec)
		}
	}

	if analysis.Module != "OAT Detections" {
		t.Errorf("Expected module 'OAT Detections', got '%s'", analysis.Module)
	}
}

func testCreditCalculator(t *testing.T) {
	t.Log("Testing credit calculator...")

	calc := NewCreditCalculator()
	if calc == nil {
		t.Fatal("Failed to create credit calculator")
	}

	// Test endpoint credit calculation
	endpointCredits := calc.CalculateEndpointCredits(100, 50, 10)
	t.Logf("Endpoint credits (100 standard, 50 pro, 10 server): %.0f", endpointCredits)

	if endpointCredits <= 0 {
		t.Error("Expected positive endpoint credits")
	}

	// Test workbench credit calculation
	workbenchCredits := calc.CalculateWorkbenchCredits(20, 5, 10)
	t.Logf("Workbench credits (20 alerts/day, 5 hunting/day, 10 rules): %.0f", workbenchCredits)

	if workbenchCredits <= 0 {
		t.Error("Expected positive workbench credits")
	}

	// Test sandbox credit calculation
	sandboxCredits := calc.CalculateSandboxCredits(50, 10, 5)
	t.Logf("Sandbox credits (50 files/day, 10 URLs/day, 5 priority/day): %.0f", sandboxCredits)

	if sandboxCredits <= 0 {
		t.Error("Expected positive sandbox credits")
	}

	// Test cost estimation
	totalCredits := endpointCredits + workbenchCredits + sandboxCredits
	monthlyCost := calc.EstimateMonthlyCost(totalCredits)
	t.Logf("Total credits: %.0f, Estimated monthly cost: $%.2f", totalCredits, monthlyCost)

	// Test optimization recommendations
	usage := map[string]float64{
		"endpoint_security": endpointCredits,
		"workbench":         workbenchCredits,
		"sandbox":           sandboxCredits,
	}

	recommendations := calc.GetOptimizationRecommendations(usage)
	if len(recommendations) > 0 {
		t.Log("Optimization recommendations:")
		for _, rec := range recommendations {
			t.Logf("  - %s", rec)
		}
	}

	// Test report generation
	report := calc.FormatCreditReport(usage)
	if report == "" {
		t.Error("Expected non-empty credit report")
	}
	t.Logf("Credit Report:\n%s", report)
}

func testCreditStorage(t *testing.T) {
	t.Log("Testing credit storage...")

	storage, err := NewCreditStorage()
	if err != nil {
		t.Fatalf("Failed to create credit storage: %v", err)
	}

	// Test saving allocation
	err = storage.UpdateAllocation("endpoint_security", 10000, 6500, "Test allocation")
	if err != nil {
		t.Errorf("Failed to update allocation: %v", err)
	}

	// Test retrieving allocation
	allocation, err := storage.GetAllocation("endpoint_security")
	if err != nil {
		t.Errorf("Failed to get allocation: %v", err)
	} else {
		t.Logf("Retrieved allocation: %s - %.0f/%.0f credits",
			allocation.Service, allocation.UsedCredits, allocation.AllocatedCredits)
	}

	// Test utilization calculation
	utilization, err := storage.CalculateUtilization()
	if err != nil {
		t.Errorf("Failed to calculate utilization: %v", err)
	} else {
		for service, util := range utilization {
			t.Logf("Service %s utilization: %.1f%%", service, util)
		}
	}

	// Test underutilized services
	underutilized, err := storage.GetUnderutilizedServices(50)
	if err != nil {
		t.Errorf("Failed to get underutilized services: %v", err)
	} else {
		t.Logf("Found %d underutilized services", len(underutilized))
	}

	// Test report generation
	report, err := storage.GenerateUtilizationReport()
	if err != nil {
		t.Errorf("Failed to generate utilization report: %v", err)
	} else if report == "" {
		t.Error("Expected non-empty utilization report")
	} else {
		t.Logf("Utilization Report:\n%s", report)
	}
}

func testComprehensiveAnalysis(t *testing.T, client *v1client.V1ApiClient) {
	t.Log("Testing comprehensive credit analysis...")

	endTime := time.Now()
	startTime := endTime.AddDate(0, 0, -7)

	result, err := performComprehensiveAnalysis(client, 10, startTime, endTime)
	if err != nil {
		t.Fatalf("Failed to perform comprehensive analysis: %v", err)
	}

	if result == nil {
		t.Fatal("Expected result, got nil")
	}

	// Check that we got content
	if len(result.Content) == 0 {
		t.Error("Expected result content, got empty")
	} else {
		for _, content := range result.Content {
			if textContent, ok := content.(interface{ GetText() string }); ok {
				text := textContent.GetText()
				if text != "" {
					t.Logf("Comprehensive Analysis Result:\n%s", text)
				}
			}
		}
	}
}
