// Package apilogapi provides test stories and method definitions for Apilog API E2E tests.
package apilogapi

import (
	"github.com/saasus-platform/saasus-sdk-go/tests/testlib"
)

// GetApilogStories returns all apilog API test stories
func GetApilogStories() []testlib.Story {
	return []testlib.Story{
		GetPostmanStoryStandardMethods(),
		GetPostmanStoryWithResponseMethods(),
	}
}

// GetPostmanStoryStandardMethods returns the Postman collection story using standard methods
func GetPostmanStoryStandardMethods() testlib.Story {
	return testlib.Story{
		Name:        "Postman Collection Story - Standard Methods",
		Description: "Reproduces the Postman collection test flow using standard client methods",
		Variables:   map[string]interface{}{},
		Steps: []testlib.Step{
			{
				Name:           "Pre_GetApiLogs",
				ClientMethod:   "GetLogs",
				Parameters:     getLogsParams,
				ExpectedStatus: 200,
				Validation: func(response interface{}) error {
					return validateGetLogsResponse(response)
				},
				StateUpdate: func(response interface{}, variables map[string]interface{}) error {
					return extractLogsVariables(response, variables)
				},
			},
			{
				Name:           "GetApiLogs",
				ClientMethod:   "GetLogs",
				Parameters:     getLogsParams,
				ExpectedStatus: 200,
				Validation: func(response interface{}) error {
					return validateGetLogsResponse(response)
				},
			},
			{
				Name:           "GetApiLogs With QueryParameters",
				ClientMethod:   "GetLogs",
				Parameters:     getLogsWithQueryParams,
				ExpectedStatus: 200,
				Validation: func(response interface{}) error {
					return validateGetLogsResponse(response)
				},
			},
			{
				Name:           "GetApiLog",
				ClientMethod:   "GetLog",
				Parameters:     getLogParams,
				ExpectedStatus: 200,
				Validation: func(response interface{}) error {
					return validateGetLogResponse(response)
				},
			},
		},
	}
}

// GetPostmanStoryWithResponseMethods returns the Postman collection story using WithResponse methods
func GetPostmanStoryWithResponseMethods() testlib.Story {
	return testlib.Story{
		Name:        "Postman Collection Story - WithResponse Methods",
		Description: "Reproduces the Postman collection test flow using WithResponse client methods",
		Variables:   map[string]interface{}{},
		Steps: []testlib.Step{
			{
				Name:           "Pre_GetApiLogs",
				ClientMethod:   "GetLogsWithResponse",
				Parameters:     getLogsParams,
				ExpectedStatus: 200,
				Validation: func(response interface{}) error {
					return validateGetLogsWithResponse(response)
				},
				StateUpdate: func(response interface{}, variables map[string]interface{}) error {
					return extractLogsVariablesWithResponse(response, variables)
				},
			},
			{
				Name:           "GetApiLogs",
				ClientMethod:   "GetLogsWithResponse",
				Parameters:     getLogsParams,
				ExpectedStatus: 200,
				Validation: func(response interface{}) error {
					return validateGetLogsWithResponse(response)
				},
			},
			{
				Name:           "GetApiLogs With QueryParameters",
				ClientMethod:   "GetLogsWithResponse",
				Parameters:     getLogsWithQueryParams,
				ExpectedStatus: 200,
				Validation: func(response interface{}) error {
					return validateGetLogsWithResponse(response)
				},
			},
			{
				Name:           "GetApiLog",
				ClientMethod:   "GetLogWithResponse",
				Parameters:     getLogParams,
				ExpectedStatus: 200,
				Validation: func(response interface{}) error {
					return validateGetLogWithResponse(response)
				},
			},
		},
	}
}

// GetApilogMethods returns all apilog API client methods
func GetApilogMethods() []string {
	return []string{
		"GetLogs",
		"GetLog",
		"GetLogsWithResponse",
		"GetLogWithResponse",
	}
}
