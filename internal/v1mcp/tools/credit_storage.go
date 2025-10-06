package tools

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// CreditAllocation represents user's credit allocation
type CreditAllocation struct {
	Service         string    `json:"service"`
	AllocatedCredits float64  `json:"allocatedCredits"`
	UsedCredits     float64   `json:"usedCredits"`
	LastUpdated     time.Time `json:"lastUpdated"`
	Notes           string    `json:"notes,omitempty"`
}

// UserCreditData represents the user's complete credit configuration
type UserCreditData struct {
	TotalCredits     float64            `json:"totalCredits"`
	Allocations      []CreditAllocation `json:"allocations"`
	LastUpdated      time.Time          `json:"lastUpdated"`
	BillingCycle     string             `json:"billingCycle"`
	Organization     string             `json:"organization,omitempty"`
	ContractType     string             `json:"contractType,omitempty"`
	CostPerCredit    float64            `json:"costPerCredit,omitempty"`
}

// CreditStorage manages persistent storage of user credit data
type CreditStorage struct {
	dataPath string
}

// NewCreditStorage creates a new credit storage instance
func NewCreditStorage() (*CreditStorage, error) {
	// Get user's home directory
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get home directory: %v", err)
	}

	// Create .vision-one directory if it doesn't exist
	dataDir := filepath.Join(homeDir, ".vision-one")
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create data directory: %v", err)
	}

	return &CreditStorage{
		dataPath: filepath.Join(dataDir, "credits.json"),
	}, nil
}

// Load reads the user's credit data from disk
func (cs *CreditStorage) Load() (*UserCreditData, error) {
	data, err := os.ReadFile(cs.dataPath)
	if err != nil {
		if os.IsNotExist(err) {
			// Return empty data if file doesn't exist
			return &UserCreditData{
				Allocations:  []CreditAllocation{},
				LastUpdated:  time.Now(),
				BillingCycle: "monthly",
			}, nil
		}
		return nil, fmt.Errorf("failed to read credit data: %v", err)
	}

	var creditData UserCreditData
	if err := json.Unmarshal(data, &creditData); err != nil {
		return nil, fmt.Errorf("failed to parse credit data: %v", err)
	}

	return &creditData, nil
}

// Save writes the user's credit data to disk
func (cs *CreditStorage) Save(data *UserCreditData) error {
	data.LastUpdated = time.Now()

	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to serialize credit data: %v", err)
	}

	if err := os.WriteFile(cs.dataPath, jsonData, 0644); err != nil {
		return fmt.Errorf("failed to write credit data: %v", err)
	}

	return nil
}

// UpdateAllocation updates or adds a service allocation
func (cs *CreditStorage) UpdateAllocation(service string, allocated, used float64, notes string) error {
	data, err := cs.Load()
	if err != nil {
		return err
	}

	found := false
	for i, allocation := range data.Allocations {
		if allocation.Service == service {
			data.Allocations[i].AllocatedCredits = allocated
			data.Allocations[i].UsedCredits = used
			data.Allocations[i].Notes = notes
			data.Allocations[i].LastUpdated = time.Now()
			found = true
			break
		}
	}

	if !found {
		data.Allocations = append(data.Allocations, CreditAllocation{
			Service:          service,
			AllocatedCredits: allocated,
			UsedCredits:      used,
			Notes:            notes,
			LastUpdated:      time.Now(),
		})
	}

	return cs.Save(data)
}

// GetAllocation retrieves allocation for a specific service
func (cs *CreditStorage) GetAllocation(service string) (*CreditAllocation, error) {
	data, err := cs.Load()
	if err != nil {
		return nil, err
	}

	for _, allocation := range data.Allocations {
		if allocation.Service == service {
			return &allocation, nil
		}
	}

	return nil, fmt.Errorf("no allocation found for service: %s", service)
}

// CalculateUtilization calculates the utilization percentage for each service
func (cs *CreditStorage) CalculateUtilization() (map[string]float64, error) {
	data, err := cs.Load()
	if err != nil {
		return nil, err
	}

	utilization := make(map[string]float64)
	for _, allocation := range data.Allocations {
		if allocation.AllocatedCredits > 0 {
			utilization[allocation.Service] = (allocation.UsedCredits / allocation.AllocatedCredits) * 100
		}
	}

	return utilization, nil
}

// GetUnderutilizedServices returns services using less than the threshold percentage
func (cs *CreditStorage) GetUnderutilizedServices(threshold float64) ([]CreditAllocation, error) {
	data, err := cs.Load()
	if err != nil {
		return nil, err
	}

	var underutilized []CreditAllocation
	for _, allocation := range data.Allocations {
		if allocation.AllocatedCredits > 0 {
			utilization := (allocation.UsedCredits / allocation.AllocatedCredits) * 100
			if utilization < threshold {
				underutilized = append(underutilized, allocation)
			}
		}
	}

	return underutilized, nil
}

// GetOverutilizedServices returns services using more than the threshold percentage
func (cs *CreditStorage) GetOverutilizedServices(threshold float64) ([]CreditAllocation, error) {
	data, err := cs.Load()
	if err != nil {
		return nil, err
	}

	var overutilized []CreditAllocation
	for _, allocation := range data.Allocations {
		if allocation.AllocatedCredits > 0 {
			utilization := (allocation.UsedCredits / allocation.AllocatedCredits) * 100
			if utilization > threshold {
				overutilized = append(overutilized, allocation)
			}
		}
	}

	return overutilized, nil
}

// GenerateUtilizationReport creates a formatted utilization report
func (cs *CreditStorage) GenerateUtilizationReport() (string, error) {
	data, err := cs.Load()
	if err != nil {
		return "", err
	}

	report := "Credit Utilization Report\n"
	report += "========================\n\n"
	report += fmt.Sprintf("Organization: %s\n", data.Organization)
	report += fmt.Sprintf("Total Credits: %.0f\n", data.TotalCredits)
	report += fmt.Sprintf("Billing Cycle: %s\n", data.BillingCycle)
	report += fmt.Sprintf("Last Updated: %s\n\n", data.LastUpdated.Format("2006-01-02 15:04:05"))

	totalAllocated := 0.0
	totalUsed := 0.0

	report += "Service Allocations:\n"
	report += "-------------------\n"
	for _, allocation := range data.Allocations {
		utilization := 0.0
		if allocation.AllocatedCredits > 0 {
			utilization = (allocation.UsedCredits / allocation.AllocatedCredits) * 100
		}

		status := "🟢" // Green
		if utilization > 90 {
			status = "🔴" // Red
		} else if utilization > 75 {
			status = "🟡" // Yellow
		}

		report += fmt.Sprintf("%s %-20s: %8.0f / %8.0f credits (%5.1f%%)",
			status, allocation.Service, allocation.UsedCredits, allocation.AllocatedCredits, utilization)

		if allocation.Notes != "" {
			report += fmt.Sprintf(" - %s", allocation.Notes)
		}
		report += "\n"

		totalAllocated += allocation.AllocatedCredits
		totalUsed += allocation.UsedCredits
	}

	report += "\n"
	report += fmt.Sprintf("Total Allocated: %.0f credits (%.1f%% of total)\n",
		totalAllocated, (totalAllocated/data.TotalCredits)*100)
	report += fmt.Sprintf("Total Used: %.0f credits (%.1f%% of allocated)\n",
		totalUsed, (totalUsed/totalAllocated)*100)

	unallocated := data.TotalCredits - totalAllocated
	if unallocated > 0 {
		report += fmt.Sprintf("Unallocated: %.0f credits (%.1f%% of total)\n",
			unallocated, (unallocated/data.TotalCredits)*100)
	}

	// Add cost information if available
	if data.CostPerCredit > 0 {
		monthlyCost := totalUsed * data.CostPerCredit
		report += fmt.Sprintf("\nEstimated Monthly Cost: $%,.2f\n", monthlyCost)
	}

	return report, nil
}

// ImportFromCSV imports credit data from a CSV file
func (cs *CreditStorage) ImportFromCSV(csvPath string) error {
	// This would be implemented to import data from CSV format
	// For now, returning a placeholder
	return fmt.Errorf("CSV import not yet implemented")
}

// ExportToCSV exports credit data to a CSV file
func (cs *CreditStorage) ExportToCSV(csvPath string) error {
	// This would be implemented to export data to CSV format
	// For now, returning a placeholder
	return fmt.Errorf("CSV export not yet implemented")
}