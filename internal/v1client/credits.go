package v1client

import (
	"net/http"
)

// CreditsListEndpoints lists endpoints for credit analysis
func (c *V1ApiClient) CreditsListEndpoints(filter string, queryParams QueryParameters) (*http.Response, error) {
	return c.searchAndFilter(
		"v3.0/endpointSecurity/endpoints",
		filter,
		queryParams,
	)
}

// CreditsListDataLakePipelines lists active data lake pipelines consuming credits
func (c *V1ApiClient) CreditsListDataLakePipelines(filter string, queryParams QueryParameters) (*http.Response, error) {
	return c.searchAndFilter(
		"v3.0/datalake/pipelines",
		filter,
		queryParams,
	)
}

// CreditsListOATDetections lists OAT detections for credit analysis
func (c *V1ApiClient) CreditsListOATDetections(filter string, queryParams QueryParameters) (*http.Response, error) {
	return c.searchAndFilter(
		"v3.0/xdr/oat/detections",
		filter,
		queryParams,
	)
}

// CreditsListOATPipelines lists Observed Attack Techniques pipelines
func (c *V1ApiClient) CreditsListOATPipelines(filter string, queryParams QueryParameters) (*http.Response, error) {
	return c.searchAndFilter(
		"v3.0/oat/pipelines",
		filter,
		queryParams,
	)
}

// CreditsListSandboxAnalysisResults lists sandbox analysis results for credit usage
func (c *V1ApiClient) CreditsListSandboxAnalysisResults(filter string, queryParams QueryParameters) (*http.Response, error) {
	return c.searchAndFilter(
		"v3.0/sandbox/analysisResults",
		filter,
		queryParams,
	)
}

// CreditsListSandboxSubmissions lists sandbox submissions for credit usage analysis
func (c *V1ApiClient) CreditsListSandboxSubmissions(filter string, queryParams QueryParameters) (*http.Response, error) {
	return c.searchAndFilter(
		"v3.0/sandbox/submissions",
		filter,
		queryParams,
	)
}

// CreditsListWorkbenchAlerts lists workbench alerts for investigation activity analysis
func (c *V1ApiClient) CreditsListWorkbenchAlerts(filter string, queryParams QueryParameters) (*http.Response, error) {
	return c.searchAndFilter(
		"v3.0/workbench/alerts",
		filter,
		queryParams,
	)
}

// CreditsGetSearchStatistics gets search statistics for credit usage analysis
func (c *V1ApiClient) CreditsGetSearchStatistics(filter string, queryParams QueryParameters) (*http.Response, error) {
	return c.searchAndFilter(
		"v3.0/search/statistics",
		filter,
		queryParams,
	)
}

// CreditsListDataRetentionModels lists data retention models
func (c *V1ApiClient) CreditsListDataRetentionModels(filter string, queryParams QueryParameters) (*http.Response, error) {
	return c.searchAndFilter(
		"v3.0/datalake/dataRetentionModels",
		filter,
		queryParams,
	)
}

// CreditsListCREMHighRiskDevices lists high risk devices from CREM for enhanced analysis
func (c *V1ApiClient) CreditsListCREMHighRiskDevices(filter string, queryParams QueryParameters) (*http.Response, error) {
	return c.searchAndFilter(
		"v3.0/asrm/highRiskDevices",
		filter,
		queryParams,
	)
}

// CreditsListCREMHighRiskUsers lists high risk users from CREM for enhanced analysis
func (c *V1ApiClient) CreditsListCREMHighRiskUsers(filter string, queryParams QueryParameters) (*http.Response, error) {
	return c.searchAndFilter(
		"v3.0/asrm/highRiskUsers",
		filter,
		queryParams,
	)
}

// CreditsListCREMCompromiseIndicators lists compromise indicators from CREM
func (c *V1ApiClient) CreditsListCREMCompromiseIndicators(filter string, queryParams QueryParameters) (*http.Response, error) {
	return c.searchAndFilter(
		"v3.0/asrm/compromiseIndicators",
		filter,
		queryParams,
	)
}

// Note: The following endpoints (allocation, balance, usage statistics, limits) are placeholders
// and don't exist in the Vision One API yet. They have been removed to avoid confusion.
// Real credit data should be obtained through:
// - Search statistics for data lake usage (CreditsGetSearchStatistics)
// - Workbench alerts for investigation activity (CreditsListWorkbenchAlerts)
// - Endpoint inventory for license allocation (CreditsListEndpoints)
// - Sandbox submissions for analysis quotas (CreditsListSandboxSubmissions)