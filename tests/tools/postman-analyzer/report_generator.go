package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

// ReportGenerator handles report generation and file export
type ReportGenerator struct {
	analysis *PostmanAnalysis
}

// NewReportGenerator creates a new report generator
func NewReportGenerator(analysis *PostmanAnalysis) *ReportGenerator {
	return &ReportGenerator{
		analysis: analysis,
	}
}

// GenerateConsoleReport prints analysis to console
func (g *ReportGenerator) GenerateConsoleReport() {
	fmt.Printf("🔍 %s - Postman Collection Analysis\n", g.analysis.CollectionName)
	fmt.Println(strings.Repeat("=", 60))

	g.printOverview()
	g.printTestFlow()
	g.printMethodCoverage()
	g.printGeneratedMethods()

	fmt.Println("\n✅ Postman Collection Analysis Complete")
}

// ExportFiles exports analysis to JSON and text files
func (g *ReportGenerator) ExportFiles(outputDir string) error {
	fmt.Println("\n=== File Export ===")

	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	if err := g.exportJSON(outputDir); err != nil {
		return err
	}

	if err := g.exportTextReport(outputDir); err != nil {
		return err
	}

	return nil
}

// printOverview prints collection overview
func (g *ReportGenerator) printOverview() {
	fmt.Printf("Collection: %s\n", g.analysis.CollectionName)
	fmt.Printf("Total Test Steps: %d\n", len(g.analysis.TestFlow))
	fmt.Printf("Total Client Methods: %d\n", len(g.analysis.ClientMethods))
}

// printTestFlow prints test flow details
func (g *ReportGenerator) printTestFlow() {
	fmt.Println("\n=== Test Flow ===")
	for i, step := range g.analysis.TestFlow {
		fmt.Printf("%d. %s\n", i+1, step.Name)
		fmt.Printf("   Method: %s %s\n", step.Method, step.Endpoint)

		// モジュール情報を追加
		moduleName := step.Module
		if moduleName == "" {
			moduleName = "unknown"
		}
		fmt.Printf("   Module: %s\n", moduleName)

		// 型情報を追加（要件10.5）
		if step.OperationID != "" {
			fmt.Printf("   Operation ID: %s\n", step.OperationID)
		}
		if step.ClientMethodName != "" {
			fmt.Printf("   Client Method: %s\n", step.ClientMethodName)
		}
		if step.ParameterTypeName != "" {
			fmt.Printf("   Parameter Type: %s (%d fields)\n", step.ParameterTypeName, len(step.ParameterTypeFields))
		}

		fmt.Printf("   Expected Status: %d\n", step.ExpectedStatus)
		if len(step.Tests) > 0 {
			fmt.Printf("   Tests: %d test assertions\n", len(step.Tests))
		}
		fmt.Println()
	}
}

// printMethodCoverage prints method coverage analysis
func (g *ReportGenerator) printMethodCoverage() {
	coverage := g.calculateCoverage()

	fmt.Println("=== Method Coverage Analysis ===")
	fmt.Printf("Total Client Methods: %d\n", len(g.analysis.ClientMethods))
	fmt.Printf("Standard Methods: %d\n", coverage["standard"])
	fmt.Printf("WithResponse Methods: %d\n", coverage["with_response"])
	fmt.Printf("WithBody Methods: %d\n", coverage["with_body"])
	fmt.Printf("WithBodyWithResponse Methods: %d\n", coverage["with_body_with_response"])
}

// printGeneratedMethods prints list of generated methods
func (g *ReportGenerator) printGeneratedMethods() {
	fmt.Println("\n=== Generated Methods ===")
	for _, method := range g.analysis.ClientMethods {
		paramInfo := ""
		if len(method.ParamTypes) > 0 {
			paramInfo = fmt.Sprintf(" [Params: %s]", strings.Join(method.ParamTypes, ", "))
		}
		fmt.Printf("✅ %s (%s)%s\n", method.Name, method.Type, paramInfo)
	}
}

// calculateCoverage calculates method coverage by type
func (g *ReportGenerator) calculateCoverage() map[string]int {
	coverage := map[string]int{
		"standard":                0,
		"with_response":           0,
		"with_body":               0,
		"with_body_with_response": 0,
	}

	for _, method := range g.analysis.ClientMethods {
		coverage[method.Type]++
	}

	return coverage
}

// exportJSON exports analysis as JSON
func (g *ReportGenerator) exportJSON(outputDir string) error {
	jsonData, err := json.MarshalIndent(g.analysis, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}

	jsonFile := filepath.Join(outputDir, "postman_analysis.json")
	if err := os.WriteFile(jsonFile, jsonData, 0644); err != nil {
		return fmt.Errorf("failed to write JSON file: %w", err)
	}

	fmt.Printf("📄 JSON analysis exported to %s\n", jsonFile)
	return nil
}

// exportTextReport exports analysis as text report
func (g *ReportGenerator) exportTextReport(outputDir string) error {
	textReport := g.generateTextReport()
	textFile := filepath.Join(outputDir, "postman_analysis_report.txt")

	if err := os.WriteFile(textFile, []byte(textReport), 0644); err != nil {
		return fmt.Errorf("failed to write text report: %w", err)
	}

	fmt.Printf("📄 Text report exported to %s\n", textFile)
	return nil
}

// generateTextReport creates a detailed text report
func (g *ReportGenerator) generateTextReport() string {
	var report strings.Builder

	g.writeHeader(&report)
	g.writeOverview(&report)
	g.writeTestFlowDetails(&report)
	g.writeGeneratedMethodsList(&report)

	return report.String()
}

// writeHeader writes report header
func (g *ReportGenerator) writeHeader(report *strings.Builder) {
	fmt.Fprintf(report, "%s - Postman Collection Analysis Report\n", g.analysis.CollectionName)
	report.WriteString("Generated at: " + time.Now().Format("2006-01-02 15:04:05") + "\n")
	report.WriteString(strings.Repeat("=", 80) + "\n\n")
}

// writeOverview writes overview section
func (g *ReportGenerator) writeOverview(report *strings.Builder) {
	report.WriteString("POSTMAN COLLECTION OVERVIEW\n")
	report.WriteString(strings.Repeat("-", 40) + "\n")
	fmt.Fprintf(report, "Collection: %s\n", g.analysis.CollectionName)
	fmt.Fprintf(report, "Total Test Steps: %d\n", len(g.analysis.TestFlow))
	fmt.Fprintf(report, "Total Client Methods: %d\n\n", len(g.analysis.ClientMethods))
}

// writeTestFlowDetails writes test flow details
func (g *ReportGenerator) writeTestFlowDetails(report *strings.Builder) {
	report.WriteString("TEST FLOW\n")
	report.WriteString(strings.Repeat("-", 40) + "\n")

	for i, step := range g.analysis.TestFlow {
		fmt.Fprintf(report, "%d. %s\n", i+1, step.Name)
		fmt.Fprintf(report, "   Method: %s %s\n", step.Method, step.Endpoint)

		// モジュール情報を追加
		moduleName := step.Module
		if moduleName == "" {
			moduleName = "unknown"
		}
		fmt.Fprintf(report, "   Module: %s\n", moduleName)

		// 型情報を追加（要件10.5）
		if step.OperationID != "" {
			fmt.Fprintf(report, "   Operation ID: %s\n", step.OperationID)
		}
		if step.ClientMethodName != "" {
			fmt.Fprintf(report, "   Client Method: %s\n", step.ClientMethodName)
		}
		if step.ParameterTypeName != "" {
			fmt.Fprintf(report, "   Parameter Type: %s\n", step.ParameterTypeName)
		}
		if len(step.ParameterTypeFields) > 0 {
			fmt.Fprintf(report, "   Type Fields: %d fields\n", len(step.ParameterTypeFields))
			// フィールド情報の詳細を表示
			var fieldNames []string
			for jsonTag := range step.ParameterTypeFields {
				fieldNames = append(fieldNames, jsonTag)
			}
			sort.Strings(fieldNames)
			for _, jsonTag := range fieldNames {
				field := step.ParameterTypeFields[jsonTag]
				requiredStr := ""
				if field.Required {
					requiredStr = " (required)"
				}
				fmt.Fprintf(report, "     - %s: %s%s\n", jsonTag, field.Type, requiredStr)
			}
		}

		fmt.Fprintf(report, "   Expected Status: %d\n", step.ExpectedStatus)
		fmt.Fprintf(report, "   Tests: %d assertions\n", len(step.Tests))

		if len(step.RequiredVariables) > 0 {
			fmt.Fprintf(report, "   Required Variables: %s\n", strings.Join(step.RequiredVariables, ", "))
		}

		if len(step.Variables) > 0 {
			var setVars []string
			for varName := range step.Variables {
				setVars = append(setVars, varName)
			}
			sort.Strings(setVars)
			fmt.Fprintf(report, "   Sets Variables: %s\n", strings.Join(setVars, ", "))
		}

		if len(step.BodyFields) > 0 {
			fmt.Fprintf(report, "   Body Fields: %s\n", strings.Join(step.BodyFields, ", "))
		}

		report.WriteString("\n")
	}
}

// writeGeneratedMethodsList writes generated methods list
func (g *ReportGenerator) writeGeneratedMethodsList(report *strings.Builder) {
	report.WriteString("GENERATED METHODS\n")
	report.WriteString(strings.Repeat("-", 40) + "\n")

	for _, method := range g.analysis.ClientMethods {
		paramInfo := ""
		if len(method.ParamTypes) > 0 {
			paramInfo = fmt.Sprintf("\n   Params: %s", strings.Join(method.ParamTypes, ", "))
		}
		fmt.Fprintf(report, "✅ %s (%s)%s\n", method.Name, method.Type, paramInfo)
	}
}
