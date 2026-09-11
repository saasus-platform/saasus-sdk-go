package pricingapi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"os"
	"path/filepath"
	"time"

	"github.com/saasus-platform/saasus-sdk-go/generated/pricingapi"
)

// generateUniqueName creates a unique name by appending timestamp and random number
func generateUniqueName(baseName interface{}) string {
	return fmt.Sprintf("%v_%d_%d", baseName, time.Now().UnixNano()/1000000, rand.Intn(10000))
}

// TestParams represents the structure of test parameters
type TestParams struct {
	CreateTenant struct {
		Name                 string         `json:"name"`
		BackOfficeStaffEmail string         `json:"back_office_staff_email"`
		Attributes           map[string]any `json:"attributes"`
	} `json:"CreateTenant"`

	CreatePricingUnit struct {
		CreateParams map[string]any `json:"createParams"`
	} `json:"CreatePricingUnit"`

	UpdatePricingUnit struct {
		UpdateParams map[string]any `json:"updateParams"`
	} `json:"UpdatePricingUnit"`

	CreatePricingMenu struct {
		CreateParams map[string]any `json:"createParams"`
	} `json:"CreatePricingMenu"`

	UpdatePricingMenu struct {
		UpdateParams map[string]any `json:"updateParams"`
	} `json:"UpdatePricingMenu"`

	CreatePricingPlan struct {
		CreateParams map[string]any `json:"createParams"`
	} `json:"CreatePricingPlan"`

	UpdatePricingPlan struct {
		UpdateParams map[string]any `json:"updateParams"`
	} `json:"UpdatePricingPlan"`

	UpdatePricingPlansUsed struct {
		UpdateParams map[string]any `json:"updateParams"`
	} `json:"UpdatePricingPlansUsed"`

	CreateMeteringUnit struct {
		CreateParams map[string]any `json:"createParams"`
	} `json:"CreateMeteringUnit"`

	UpdateMeteringUnitByID struct {
		UpdateParams map[string]any `json:"updateParams"`
	} `json:"UpdateMeteringUnitByID"`

	UpdateMeteringUnitTimestampCount struct {
		UpdateParams map[string]any `json:"updateParams"`
	} `json:"UpdateMeteringUnitTimestampCount"`

	UpdateMeteringUnitTimestampCountNow struct {
		UpdateParams map[string]any `json:"updateParams"`
	} `json:"UpdateMeteringUnitTimestampCountNow"`

	CreateTaxRate struct {
		CreateParams map[string]any `json:"createParams"`
	} `json:"CreateTaxRate"`

	UpdateTaxRate struct {
		UpdateParams map[string]any `json:"updateParams"`
	} `json:"UpdateTaxRate"`
}

var cachedParams *TestParams

// loadTestParams loads test parameters from testdata/test_params.json
func loadTestParams() *TestParams {
	if cachedParams != nil {
		return cachedParams
	}

	data, err := os.ReadFile(filepath.Join("testdata", "test_params.json"))
	if err != nil {
		panic(fmt.Sprintf("failed to read test params: %v", err))
	}

	var params TestParams
	if err := json.Unmarshal(data, &params); err != nil {
		panic(fmt.Sprintf("failed to unmarshal test params: %v", err))
	}

	cachedParams = &params
	return cachedParams
}

func mustGetStringVariable(variables map[string]any, key string) string {
	if value, ok := variables[key].(string); ok && value != "" {
		return value
	}
	panic(fmt.Sprintf("required variable %q is missing; ensure previous steps stored it", key))
}

func resolveTenantID(variables map[string]any) string {
	if tenantID, ok := variables["tenant_id"].(string); ok && tenantID != "" {
		return tenantID
	}
	tenantID := getTestTenantID()
	if tenantID == "" {
		panic("tenant ID is not configured; set TEST_TENANT_ID or capture it from a preceding CreateTenant step")
	}
	return tenantID
}

// Helper function to get test tenant ID
func getTestTenantID() string {
	tenantID := os.Getenv("TEST_TENANT_ID")
	if tenantID == "" {
		tenantID = "test-tenant-id"
	}
	return tenantID
}

// ============================================================================
// Tenant Parameter Functions (Auth API)
// ============================================================================

// createTenantParams returns parameters for CreateTenant
func createTenantParams(variables map[string]any) any {
	params := loadTestParams()
	createParams := params.CreateTenant

	body := map[string]any{
		"name":                    generateUniqueName(createParams.Name),
		"back_office_staff_email": createParams.BackOfficeStaffEmail,
		"attributes":              createParams.Attributes,
	}

	bodyBytes, _ := json.Marshal(body)

	return struct {
		ContentType string
		Body        io.Reader
	}{
		ContentType: "application/json",
		Body:        bytes.NewReader(bodyBytes),
	}
}

// ============================================================================
// Pricing Units Parameter Functions
// ============================================================================

// createPricingUnitParams returns parameters for CreatePricingUnit
func createPricingUnitParams(variables map[string]any) any {
	params := loadTestParams()
	createParams := params.CreatePricingUnit.CreateParams
	meteringUnitName := mustGetStringVariable(variables, "metering_unit_name")

	// Generate unique name and store it for later use
	unitName := generateUniqueName(createParams["name"])
	variables["pricing_unit_name"] = unitName

	body := map[string]any{
		"type":               createParams["type"],
		"name":               unitName,
		"display_name":       createParams["display_name"],
		"description":        createParams["description"],
		"currency":           createParams["currency"],
		"upper_count":        int(createParams["upper_count"].(float64)),
		"metering_unit_name": meteringUnitName,
		"aggregate_usage":    createParams["aggregate_usage"],
		"tiers": []map[string]any{
			{
				"up_to":       int(createParams["up_to"].(float64)),
				"unit_amount": int(createParams["unit_amount"].(float64)),
				"flat_amount": int(createParams["flat_amount"].(float64)),
				"inf":         createParams["inf"],
			},
		},
	}

	bodyBytes, _ := json.Marshal(body)

	return struct {
		ContentType string
		Body        io.Reader
	}{
		ContentType: "application/json",
		Body:        bytes.NewReader(bodyBytes),
	}
}

// getPricingUnitParams returns parameters for GetPricingUnit
func getPricingUnitParams(variables map[string]any) any {
	return mustGetStringVariable(variables, "pricing_unit_id")
}

// updatePricingUnitParams returns parameters for UpdatePricingUnit
func updatePricingUnitParams(variables map[string]any) any {
	pricingUnitID := mustGetStringVariable(variables, "pricing_unit_id")
	unitName := mustGetStringVariable(variables, "pricing_unit_name")
	meteringUnitName := mustGetStringVariable(variables, "metering_unit_name")
	params := loadTestParams()

	// Get original creation params to maintain structural fields
	createParams := params.CreatePricingUnit.CreateParams
	updateParams := params.UpdatePricingUnit.UpdateParams

	// Build body with original structural fields but updated display fields
	body := map[string]any{
		"type":               createParams["type"],
		"name":               unitName,                     // Use original name from creation
		"display_name":       updateParams["display_name"], // Updated
		"description":        updateParams["description"],  // Updated
		"currency":           createParams["currency"],
		"upper_count":        int(createParams["upper_count"].(float64)),
		"metering_unit_name": meteringUnitName,
		"aggregate_usage":    createParams["aggregate_usage"],
		"tiers": []map[string]any{
			{
				"up_to":       int(createParams["up_to"].(float64)),
				"unit_amount": int(createParams["unit_amount"].(float64)),
				"flat_amount": int(createParams["flat_amount"].(float64)),
				"inf":         createParams["inf"],
			},
		},
	}

	bodyBytes, _ := json.Marshal(body)

	return struct {
		PricingUnitId string
		ContentType   string
		Body          io.Reader
	}{
		PricingUnitId: pricingUnitID,
		ContentType:   "application/json",
		Body:          bytes.NewReader(bodyBytes),
	}
}

// deletePricingUnitParams returns parameters for DeletePricingUnit
func deletePricingUnitParams(variables map[string]any) any {
	return mustGetStringVariable(variables, "pricing_unit_id")
}

// ============================================================================
// Pricing Menus Parameter Functions
// ============================================================================

// createPricingMenuParams returns parameters for CreatePricingMenu
func createPricingMenuParams(variables map[string]any) any {
	params := loadTestParams()
	createParams := params.CreatePricingMenu.CreateParams

	unitID := mustGetStringVariable(variables, "pricing_unit_id")

	// Generate unique name and store it for later use
	menuName := generateUniqueName(createParams["menu_name"])
	variables["pricing_menu_name"] = menuName

	body := map[string]any{
		"name":         menuName,
		"display_name": createParams["menu_display_name"],
		"description":  createParams["menu_description"],
		"unit_ids":     []string{unitID},
	}

	bodyBytes, _ := json.Marshal(body)

	return struct {
		ContentType string
		Body        io.Reader
	}{
		ContentType: "application/json",
		Body:        bytes.NewReader(bodyBytes),
	}
}

// getPricingMenuParams returns parameters for GetPricingMenu
func getPricingMenuParams(variables map[string]any) any {
	return mustGetStringVariable(variables, "pricing_menu_id")
}

// updatePricingMenuParams returns parameters for UpdatePricingMenu
func updatePricingMenuParams(variables map[string]any) any {
	pricingMenuID := mustGetStringVariable(variables, "pricing_menu_id")
	menuName := mustGetStringVariable(variables, "pricing_menu_name")
	params := loadTestParams()
	updateParams := params.UpdatePricingMenu.UpdateParams

	unitID := mustGetStringVariable(variables, "pricing_unit_id")

	body := map[string]any{
		"name":         menuName, // Use original name from creation
		"display_name": updateParams["menu_display_name"],
		"description":  updateParams["menu_description"],
		"unit_ids":     []string{unitID},
	}

	bodyBytes, _ := json.Marshal(body)

	return struct {
		MenuId      string
		ContentType string
		Body        io.Reader
	}{
		MenuId:      pricingMenuID,
		ContentType: "application/json",
		Body:        bytes.NewReader(bodyBytes),
	}
}

// deletePricingMenuParams returns parameters for DeletePricingMenu
func deletePricingMenuParams(variables map[string]any) any {
	return mustGetStringVariable(variables, "pricing_menu_id")
}

// ============================================================================
// Pricing Plans Parameter Functions
// ============================================================================

// createPricingPlanParams returns parameters for CreatePricingPlan
func createPricingPlanParams(variables map[string]any) any {
	params := loadTestParams()
	createParams := params.CreatePricingPlan.CreateParams

	menuID := mustGetStringVariable(variables, "pricing_menu_id")

	// Generate unique name and store it for later use
	planName := generateUniqueName(createParams["plan_name"])
	variables["pricing_plan_name"] = planName

	body := map[string]any{
		"name":         planName,
		"display_name": createParams["plan_display_name"],
		"description":  createParams["plan_description"],
		"menu_ids":     []string{menuID},
	}

	bodyBytes, _ := json.Marshal(body)

	return struct {
		ContentType string
		Body        io.Reader
	}{
		ContentType: "application/json",
		Body:        bytes.NewReader(bodyBytes),
	}
}

// getPricingPlanParams returns parameters for GetPricingPlan
func getPricingPlanParams(variables map[string]any) any {
	return mustGetStringVariable(variables, "pricing_plan_id")
}

// updatePricingPlanParams returns parameters for UpdatePricingPlan
func updatePricingPlanParams(variables map[string]any) any {
	pricingPlanID := mustGetStringVariable(variables, "pricing_plan_id")
	planName := mustGetStringVariable(variables, "pricing_plan_name")
	params := loadTestParams()
	updateParams := params.UpdatePricingPlan.UpdateParams

	menuID := mustGetStringVariable(variables, "pricing_menu_id")

	body := map[string]any{
		"name":         planName, // Use original name from creation
		"display_name": updateParams["plan_display_name"],
		"description":  updateParams["plan_description"],
		"menu_ids":     []string{menuID},
	}

	bodyBytes, _ := json.Marshal(body)

	return struct {
		PlanId      string
		ContentType string
		Body        io.Reader
	}{
		PlanId:      pricingPlanID,
		ContentType: "application/json",
		Body:        bytes.NewReader(bodyBytes),
	}
}

// updatePricingPlansUsedParams returns parameters for UpdatePricingPlansUsed
func updatePricingPlansUsedParams(variables map[string]any) any {
	planID := mustGetStringVariable(variables, "pricing_plan_id")

	body := map[string]any{
		"plan_ids": []string{planID},
	}

	bodyBytes, _ := json.Marshal(body)

	return struct {
		ContentType string
		Body        io.Reader
	}{
		ContentType: "application/json",
		Body:        bytes.NewReader(bodyBytes),
	}
}

// deletePricingPlanParams returns parameters for DeletePricingPlan
func deletePricingPlanParams(variables map[string]any) any {
	return mustGetStringVariable(variables, "pricing_plan_id")
}

// ============================================================================
// Metering Parameter Functions
// ============================================================================

// createMeteringUnitParams returns parameters for CreateMeteringUnit
func createMeteringUnitParams(variables map[string]any) any {
	params := loadTestParams()
	createParams := params.CreateMeteringUnit.CreateParams

	uniqueUnitName := generateUniqueName(createParams["unit_name"])

	body := map[string]any{
		"unit_name":       uniqueUnitName,
		"display_name":    createParams["display_name"],
		"description":     createParams["description"],
		"aggregate_usage": createParams["aggregate_usage"],
	}

	// make unit name available for subsequent steps before the API response arrives
	variables["metering_unit_name"] = body["unit_name"]

	bodyBytes, _ := json.Marshal(body)

	return struct {
		ContentType string
		Body        io.Reader
	}{
		ContentType: "application/json",
		Body:        bytes.NewReader(bodyBytes),
	}
}

// getMeteringUnitDateCountParams returns parameters for GetMeteringUnitDateCountByTenantIdAndUnitNameAndDate
func getMeteringUnitDateCountParams(variables map[string]any) any {
	tenantID := resolveTenantID(variables)
	meteringUnitName := mustGetStringVariable(variables, "metering_unit_name")
	date := "2024-01-01"

	return struct {
		TenantId         string
		MeteringUnitName string
		Date             string
	}{
		TenantId:         tenantID,
		MeteringUnitName: meteringUnitName,
		Date:             date,
	}
}

// updateMeteringUnitTimestampCountParams returns parameters for UpdateMeteringUnitTimestampCount
func updateMeteringUnitTimestampCountParams(variables map[string]any) any {
	tenantID := resolveTenantID(variables)
	meteringUnitName := mustGetStringVariable(variables, "metering_unit_name")
	timestamp := 1640995200

	body := map[string]any{
		"method": "add",
		"count":  10,
	}

	bodyBytes, _ := json.Marshal(body)

	return struct {
		TenantId         string
		MeteringUnitName string
		Timestamp        int
		ContentType      string
		Body             io.Reader
	}{
		TenantId:         tenantID,
		MeteringUnitName: meteringUnitName,
		Timestamp:        timestamp,
		ContentType:      "application/json",
		Body:             bytes.NewReader(bodyBytes),
	}
}

// deleteMeteringUnitTimestampCountParams returns parameters for DeleteMeteringUnitTimestampCount
func deleteMeteringUnitTimestampCountParams(variables map[string]any) any {
	tenantID := resolveTenantID(variables)
	meteringUnitName := mustGetStringVariable(variables, "metering_unit_name")
	timestamp := 1640995200

	return struct {
		TenantId         string
		MeteringUnitName string
		Timestamp        int
	}{
		TenantId:         tenantID,
		MeteringUnitName: meteringUnitName,
		Timestamp:        timestamp,
	}
}

// updateMeteringUnitTimestampCountNowParams returns parameters for UpdateMeteringUnitTimestampCountNow
func updateMeteringUnitTimestampCountNowParams(variables map[string]any) any {
	tenantID := resolveTenantID(variables)
	meteringUnitName := mustGetStringVariable(variables, "metering_unit_name")

	body := map[string]any{
		"method": "add",
		"count":  5,
	}

	bodyBytes, _ := json.Marshal(body)

	return struct {
		TenantId         string
		MeteringUnitName string
		ContentType      string
		Body             io.Reader
	}{
		TenantId:         tenantID,
		MeteringUnitName: meteringUnitName,
		ContentType:      "application/json",
		Body:             bytes.NewReader(bodyBytes),
	}
}

// getMeteringUnitDateCountTodayParams returns parameters for GetMeteringUnitDateCountByTenantIdAndUnitNameToday
func getMeteringUnitDateCountTodayParams(variables map[string]any) any {
	tenantID := resolveTenantID(variables)
	meteringUnitName := mustGetStringVariable(variables, "metering_unit_name")

	return struct {
		TenantId         string
		MeteringUnitName string
	}{
		TenantId:         tenantID,
		MeteringUnitName: meteringUnitName,
	}
}

// getMeteringUnitMonthCountThisMonthParams returns parameters for GetMeteringUnitMonthCountByTenantIdAndUnitNameThisMonth
func getMeteringUnitMonthCountThisMonthParams(variables map[string]any) any {
	tenantID := resolveTenantID(variables)
	meteringUnitName := mustGetStringVariable(variables, "metering_unit_name")

	return struct {
		TenantId         string
		MeteringUnitName string
	}{
		TenantId:         tenantID,
		MeteringUnitName: meteringUnitName,
	}
}

// getMeteringUnitMonthCountParams returns parameters for GetMeteringUnitMonthCountByTenantIdAndUnitNameAndMonth
func getMeteringUnitMonthCountParams(variables map[string]any) any {
	tenantID := resolveTenantID(variables)
	meteringUnitName := mustGetStringVariable(variables, "metering_unit_name")
	month := "2024-01"

	return struct {
		TenantId         string
		MeteringUnitName string
		Month            string
	}{
		TenantId:         tenantID,
		MeteringUnitName: meteringUnitName,
		Month:            month,
	}
}

// getMeteringUnitDateCountsByDateParams returns parameters for GetMeteringUnitDateCountsByTenantIdAndDate
func getMeteringUnitDateCountsByDateParams(variables map[string]any) any {
	tenantID := resolveTenantID(variables)
	date := "2024-01-01"

	return struct {
		TenantId string
		Date     string
	}{
		TenantId: tenantID,
		Date:     date,
	}
}

// getMeteringUnitMonthCountsByMonthParams returns parameters for GetMeteringUnitMonthCountsByTenantIdAndMonth
func getMeteringUnitMonthCountsByMonthParams(variables map[string]any) any {
	tenantID := resolveTenantID(variables)
	month := "2024-01"

	return struct {
		TenantId string
		Month    string
	}{
		TenantId: tenantID,
		Month:    month,
	}
}

// getMeteringUnitDateCountByDatePeriodParams returns parameters for GetMeteringUnitDateCountByTenantIdAndUnitNameAndDatePeriod
func getMeteringUnitDateCountByDatePeriodParams(variables map[string]any) any {
	tenantID := resolveTenantID(variables)
	meteringUnitName := mustGetStringVariable(variables, "metering_unit_name")

	// Create empty params struct for optional query parameters
	params := &pricingapi.GetMeteringUnitDateCountByTenantIdAndUnitNameAndDatePeriodParams{}

	return struct {
		TenantId         string
		MeteringUnitName string
		Params           *pricingapi.GetMeteringUnitDateCountByTenantIdAndUnitNameAndDatePeriodParams
	}{
		TenantId:         tenantID,
		MeteringUnitName: meteringUnitName,
		Params:           params,
	}
}

// updateMeteringUnitByIDParams returns parameters for UpdateMeteringUnitByID
func updateMeteringUnitByIDParams(variables map[string]any) any {
	meteringUnitID := mustGetStringVariable(variables, "metering_unit_id")
	meteringUnitName := mustGetStringVariable(variables, "metering_unit_name")
	params := loadTestParams()
	updateParams := params.UpdateMeteringUnitByID.UpdateParams

	// Use original unit_name from variables (cannot change name or aggregate_usage of unit in use)
	body := map[string]any{
		"unit_name":    meteringUnitName,
		"display_name": updateParams["display_name"],
		"description":  updateParams["description"],
		// aggregate_usage must stay the same as when created
		"aggregate_usage": "max",
	}

	bodyBytes, _ := json.Marshal(body)

	return struct {
		MeteringUnitId string
		ContentType    string
		Body           io.Reader
	}{
		MeteringUnitId: meteringUnitID,
		ContentType:    "application/json",
		Body:           bytes.NewReader(bodyBytes),
	}
}

// deleteMeteringUnitByIDParams returns parameters for DeleteMeteringUnitByID
func deleteMeteringUnitByIDParams(variables map[string]any) any {
	return mustGetStringVariable(variables, "metering_unit_id")
}

// ============================================================================
// Tax Rate Parameter Functions
// ============================================================================

// createTaxRateParams returns parameters for CreateTaxRate
func createTaxRateParams(variables map[string]any) any {
	params := loadTestParams()
	createParams := params.CreateTaxRate.CreateParams

	body := map[string]any{
		"name":         generateUniqueName(createParams["name"]),
		"display_name": createParams["display_name"],
		"description":  createParams["description"],
		"percentage":   createParams["percentage"],
		"inclusive":    createParams["inclusive"],
		"country":      createParams["country"],
	}

	bodyBytes, _ := json.Marshal(body)

	return struct {
		ContentType string
		Body        io.Reader
	}{
		ContentType: "application/json",
		Body:        bytes.NewReader(bodyBytes),
	}
}

// updateTaxRateParams returns parameters for UpdateTaxRate
func updateTaxRateParams(variables map[string]any) any {
	taxRateID, _ := variables["tax_rate_id"].(string)
	params := loadTestParams()
	updateParams := params.UpdateTaxRate.UpdateParams

	body := map[string]any{
		"display_name": updateParams["display_name"],
		"description":  updateParams["description"],
	}

	bodyBytes, _ := json.Marshal(body)

	return struct {
		TaxRateId   string
		ContentType string
		Body        io.Reader
	}{
		TaxRateId:   taxRateID,
		ContentType: "application/json",
		Body:        bytes.NewReader(bodyBytes),
	}
}

// Parameter functions for non-WithBody variants (they take typed structs instead of io.Reader)

// createMeteringUnitStructParams returns typed parameters for CreateMeteringUnit (non-WithBody)
func createMeteringUnitStructParams(variables map[string]any) any {
	params := loadTestParams()
	createParams := params.CreateMeteringUnit.CreateParams

	aggregateUsageStr, ok := createParams["aggregate_usage"].(string)
	if !ok {
		aggregateUsageStr = "max"
	}
	var aggregateUsage *pricingapi.AggregateUsage
	if aggregateUsageStr != "" {
		au := pricingapi.AggregateUsage(aggregateUsageStr)
		aggregateUsage = &au
	}

	// Return the body directly, not wrapped in a struct
	return pricingapi.MeteringUnitProps{
		UnitName:       generateUniqueName(createParams["unit_name"]),
		DisplayName:    createParams["display_name"].(string),
		Description:    createParams["description"].(string),
		AggregateUsage: aggregateUsage,
	}
}

// updateMeteringUnitByIDStructParams returns typed parameters for UpdateMeteringUnit (non-WithBody)
func updateMeteringUnitByIDStructParams(variables map[string]any) any {
	meteringUnitID := mustGetStringVariable(variables, "metering_unit_id_2")
	meteringUnitName := mustGetStringVariable(variables, "metering_unit_name_2")
	params := loadTestParams()
	updateParams := params.UpdateMeteringUnitByID.UpdateParams

	aggregateUsage := pricingapi.Max
	body := pricingapi.MeteringUnitProps{
		UnitName:       meteringUnitName, // Cannot change name
		DisplayName:    updateParams["display_name"].(string),
		Description:    updateParams["description"].(string),
		AggregateUsage: &aggregateUsage,
	}

	return struct {
		MeteringUnitId string
		Body           pricingapi.UpdateMeteringUnitByIDJSONRequestBody
	}{
		MeteringUnitId: meteringUnitID,
		Body:           body,
	}
}

// updateMeteringUnitTimestampCountStructParams returns typed parameters for UpdateMeteringUnitTimestampCount (non-WithBody)
func updateMeteringUnitTimestampCountStructParams(variables map[string]any) any {
	tenantID := resolveTenantID(variables)
	meteringUnitName := mustGetStringVariable(variables, "metering_unit_name_2")
	params := loadTestParams()
	updateParams := params.UpdateMeteringUnitTimestampCount.UpdateParams

	body := pricingapi.UpdateMeteringUnitTimestampCountParam{
		Count:  int(updateParams["count"].(float64)),
		Method: pricingapi.Add,
	}

	return struct {
		TenantId         string
		MeteringUnitName string
		Timestamp        int
		Body             pricingapi.UpdateMeteringUnitTimestampCountJSONRequestBody
	}{
		TenantId:         tenantID,
		MeteringUnitName: meteringUnitName,
		Timestamp:        int(time.Now().Unix()),
		Body:             body,
	}
}

// updateMeteringUnitTimestampCountNowStructParams returns typed parameters for UpdateMeteringUnitTimestampCountNow (non-WithBody)
func updateMeteringUnitTimestampCountNowStructParams(variables map[string]any) any {
	tenantID := resolveTenantID(variables)
	meteringUnitName := mustGetStringVariable(variables, "metering_unit_name_2")
	params := loadTestParams()
	updateParams := params.UpdateMeteringUnitTimestampCountNow.UpdateParams

	body := pricingapi.UpdateMeteringUnitTimestampCountNowParam{
		Count:  int(updateParams["count"].(float64)),
		Method: pricingapi.Add,
	}

	return struct {
		TenantId         string
		MeteringUnitName string
		Body             pricingapi.UpdateMeteringUnitTimestampCountNowJSONRequestBody
	}{
		TenantId:         tenantID,
		MeteringUnitName: meteringUnitName,
		Body:             body,
	}
}

// createPricingUnitStructParams returns typed parameters for CreatePricingUnit (non-WithBody)
func createPricingUnitStructParams(variables map[string]any) any {
	params := loadTestParams()
	createParams := params.CreatePricingUnit.CreateParams

	// Generate unique name and store it for later use
	unitName := generateUniqueName(createParams["name1"])
	variables["pricing_unit_name_2"] = unitName

	// Create fixed unit
	fixedUnit := pricingapi.PricingFixedUnitForSave{
		Name:              unitName,
		DisplayName:       createParams["display_name1"].(string),
		Description:       createParams["description1"].(string),
		Type:              pricingapi.Fixed,
		UnitAmount:        uint64(createParams["unit_amount1"].(float64)),
		Currency:          pricingapi.Currency(createParams["currency1"].(string)),
		RecurringInterval: pricingapi.RecurringInterval(createParams["recurring_interval1"].(string)),
	}

	// Use FromPricingFixedUnitForSave method to create union type
	var body pricingapi.PricingUnitForSave
	_ = body.FromPricingFixedUnitForSave(fixedUnit)
	return body
}

// updatePricingUnitStructParams returns typed parameters for UpdatePricingUnit (non-WithBody)
func updatePricingUnitStructParams(variables map[string]any) any {
	pricingUnitID := mustGetStringVariable(variables, "pricing_unit_id_2")
	unitName := mustGetStringVariable(variables, "pricing_unit_name_2")
	params := loadTestParams()
	createParams := params.CreatePricingUnit.CreateParams
	updateParams := params.UpdatePricingUnit.UpdateParams

	// Create fixed unit with complete fields - use original name from variables
	fixedUnit := pricingapi.PricingFixedUnitForSave{
		Name:              unitName, // Use original name from creation
		DisplayName:       updateParams["display_name"].(string),
		Description:       updateParams["description"].(string),
		Type:              pricingapi.Fixed,
		UnitAmount:        uint64(createParams["unit_amount1"].(float64)),
		Currency:          pricingapi.Currency(createParams["currency1"].(string)),
		RecurringInterval: pricingapi.RecurringInterval(createParams["recurring_interval1"].(string)),
	}

	// Use FromPricingFixedUnitForSave method to create union type
	var body pricingapi.PricingUnitForSave
	_ = body.FromPricingFixedUnitForSave(fixedUnit)

	return struct {
		PricingUnitId string
		Body          pricingapi.UpdatePricingUnitJSONRequestBody
	}{
		PricingUnitId: pricingUnitID,
		Body:          body,
	}
}

// createPricingMenuStructParams returns typed parameters for CreatePricingMenu (non-WithBody)
func createPricingMenuStructParams(variables map[string]any) any {
	pricingUnitID, _ := variables["pricing_unit_id_2"].(string)
	params := loadTestParams()
	createParams := params.CreatePricingMenu.CreateParams

	// Generate unique name and store it for later use
	menuName := generateUniqueName(createParams["menu_name"])
	variables["pricing_menu_name_2"] = menuName

	// Create the struct directly with proper types
	body := pricingapi.CreatePricingMenuJSONRequestBody{
		Name:        menuName,
		DisplayName: createParams["menu_display_name"].(string),
		Description: createParams["menu_description"].(string),
		UnitIds:     []pricingapi.Uuid{pricingapi.Uuid(pricingUnitID)},
	}

	return body
}

// updatePricingMenuStructParams returns typed parameters for UpdatePricingMenu (non-WithBody)
func updatePricingMenuStructParams(variables map[string]any) any {
	pricingMenuID := mustGetStringVariable(variables, "pricing_menu_id_2")
	menuName := mustGetStringVariable(variables, "pricing_menu_name_2")
	pricingUnitID := mustGetStringVariable(variables, "pricing_unit_id_2")
	params := loadTestParams()
	updateParams := params.UpdatePricingMenu.UpdateParams

	body := pricingapi.UpdatePricingMenuParam{
		Name:        menuName, // Use original name from creation
		DisplayName: updateParams["menu_display_name"].(string),
		Description: updateParams["menu_description"].(string),
		UnitIds:     []pricingapi.Uuid{pricingapi.Uuid(pricingUnitID)},
	}

	return struct {
		MenuId string
		Body   pricingapi.UpdatePricingMenuJSONRequestBody
	}{
		MenuId: pricingMenuID,
		Body:   body,
	}
}

// createPricingPlanStructParams returns typed parameters for CreatePricingPlan (non-WithBody)
func createPricingPlanStructParams(variables map[string]any) any {
	pricingMenuID, _ := variables["pricing_menu_id_2"].(string)
	params := loadTestParams()
	createParams := params.CreatePricingPlan.CreateParams

	// Generate unique name and store it for later use
	planName := generateUniqueName(createParams["plan_name"])
	variables["pricing_plan_name_2"] = planName

	// Create the struct directly with proper types
	body := pricingapi.CreatePricingPlanJSONRequestBody{
		Name:        planName,
		DisplayName: createParams["plan_display_name"].(string),
		Description: createParams["plan_description"].(string),
		MenuIds:     []pricingapi.Uuid{pricingapi.Uuid(pricingMenuID)},
	}

	return body
}

// updatePricingPlanStructParams returns typed parameters for UpdatePricingPlan (non-WithBody)
func updatePricingPlanStructParams(variables map[string]any) any {
	pricingPlanID := mustGetStringVariable(variables, "pricing_plan_id_2")
	planName := mustGetStringVariable(variables, "pricing_plan_name_2")
	pricingMenuID := mustGetStringVariable(variables, "pricing_menu_id_2")
	params := loadTestParams()
	updateParams := params.UpdatePricingPlan.UpdateParams

	body := pricingapi.UpdatePricingPlanParam{
		Name:        planName, // Use original name from creation
		DisplayName: updateParams["plan_display_name"].(string),
		Description: updateParams["plan_description"].(string),
		MenuIds:     []pricingapi.Uuid{pricingapi.Uuid(pricingMenuID)},
	}

	return struct {
		PlanId string
		Body   pricingapi.UpdatePricingPlanJSONRequestBody
	}{
		PlanId: pricingPlanID,
		Body:   body,
	}
}

// updatePricingPlansUsedStructParams returns typed parameters for UpdatePricingPlansUsed (non-WithBody)
func updatePricingPlansUsedStructParams(variables map[string]any) any {
	pricingPlanID, _ := variables["pricing_plan_id_2"].(string)

	bodyData := map[string]any{
		"plan_ids": []string{pricingPlanID},
	}
	bodyBytes, _ := json.Marshal(bodyData)

	// Use the json.RawMessage approach since this is a union type
	var body pricingapi.UpdatePricingPlansUsedJSONRequestBody
	_ = json.Unmarshal(bodyBytes, &body)

	return body
}

// createTaxRateStructParams returns typed parameters for CreateTaxRate (non-WithBody)
func createTaxRateStructParams(variables map[string]any) any {
	params := loadTestParams()
	createParams := params.CreateTaxRate.CreateParams

	// Return body directly for non-WithBody Create methods
	body := pricingapi.TaxRateProps{
		Name:        generateUniqueName(createParams["name"]),
		DisplayName: createParams["display_name"].(string),
		Description: createParams["description"].(string),
		Percentage:  createParams["percentage"].(float64),
		Inclusive:   createParams["inclusive"].(bool),
		Country:     createParams["country"].(string),
	}

	return body
}

// updateTaxRateStructParams returns typed parameters for UpdateTaxRate (non-WithBody)
func updateTaxRateStructParams(variables map[string]any) any {
	taxRateID, _ := variables["tax_rate_id_2"].(string)
	params := loadTestParams()
	updateParams := params.UpdateTaxRate.UpdateParams

	body := pricingapi.UpdateTaxRateParam{
		DisplayName: updateParams["display_name"].(string),
		Description: updateParams["description"].(string),
	}

	return struct {
		TaxRateId string
		Body      pricingapi.UpdateTaxRateJSONRequestBody
	}{
		TaxRateId: taxRateID,
		Body:      body,
	}
}

// linkPlanToStripeParams returns parameters for LinkPlanToStripe
func linkPlanToStripeParams(variables map[string]any) any {
	pricingPlanID, _ := variables["pricing_plan_id_2"].(string)

	return struct {
		PlanId string
	}{
		PlanId: pricingPlanID,
	}
}

// deleteStripePlanParams returns parameters for DeleteStripePlan
func deleteStripePlanParams(variables map[string]any) any {
	pricingPlanID, _ := variables["pricing_plan_id_2"].(string)

	return struct {
		PlanId string
	}{
		PlanId: pricingPlanID,
	}
}

// deletePricingPlan2Params returns parameters for DeletePricingPlan (using pricing_plan_id_2)
func deletePricingPlan2Params(variables map[string]any) any {
	pricingPlanID, _ := variables["pricing_plan_id_2"].(string)

	return struct {
		PricingPlanId string
	}{
		PricingPlanId: pricingPlanID,
	}
}

// deletePricingMenu2Params returns parameters for DeletePricingMenu (using pricing_menu_id_2)
func deletePricingMenu2Params(variables map[string]any) any {
	pricingMenuID, _ := variables["pricing_menu_id_2"].(string)

	return struct {
		PricingMenuId string
	}{
		PricingMenuId: pricingMenuID,
	}
}

// deletePricingUnit2Params returns parameters for DeletePricingUnit (using pricing_unit_id_2)
func deletePricingUnit2Params(variables map[string]any) any {
	pricingUnitID, _ := variables["pricing_unit_id_2"].(string)

	return struct {
		PricingUnitId string
	}{
		PricingUnitId: pricingUnitID,
	}
}

// deleteMeteringUnitByID2Params returns parameters for DeleteMeteringUnitByID (using metering_unit_id_2)
func deleteMeteringUnitByID2Params(variables map[string]any) any {
	meteringUnitID, _ := variables["metering_unit_id_2"].(string)

	return struct {
		MeteringUnitId string
	}{
		MeteringUnitId: meteringUnitID,
	}
}
