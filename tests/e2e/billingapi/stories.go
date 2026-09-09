package billingapi

import (
	"github.com/saasus-platform/saasus-sdk-go/tests/testlib"
)

// GetBillingStories returns all billing API test stories
func GetBillingStories() []testlib.Story {
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
		Variables: map[string]interface{}{
			"secret_key": getTestStripeKey(),
		},
		Steps: []testlib.Step{
			{
				Name:           "Pre_GetStripeConnectionInformation",
				ClientMethod:   "GetStripeInfo",
				Parameters:     nil,
				ExpectedStatus: 200,
				Validation: func(response interface{}) error {
					return validateStripeInfoResponse(response, false)
				},
				StateUpdate: func(response interface{}, variables map[string]interface{}) error {
					variables["secret_key"] = getTestStripeKey()
					variables["story_name"] = "Postman Collection Story - Standard Methods"
					return nil
				},
			},
			{
				Name:           "UpdateStripeConnectionInfo",
				ClientMethod:   "UpdateStripeInfo",
				Parameters:     createUpdateStripeInfoParams(),
				ExpectedStatus: 200,
				Validation: func(response interface{}) error {
					return nil // Validate successful update
				},
			},
			{
				Name:           "UpdateStripeConnectionInfo_WithBody",
				ClientMethod:   "UpdateStripeInfoWithBody",
				Parameters:     createUpdateStripeInfoBodyParams(),
				ExpectedStatus: 200,
				Validation: func(response interface{}) error {
					return nil // Validate successful update with body method
				},
			},
			{
				Name:           "GetStripeConnectionInformation",
				ClientMethod:   "GetStripeInfo",
				Parameters:     nil,
				ExpectedStatus: 200,
				Validation: func(response interface{}) error {
					return validateStripeInfoResponse(response, true)
				},
			},
			{
				Name:           "DeleteStripeConnection",
				ClientMethod:   "DeleteStripeInfo",
				Parameters:     nil,
				ExpectedStatus: 200,
				Validation: func(response interface{}) error {
					return nil // Validate successful deletion
				},
			},
			{
				Name:           "Final_GetStripeConnectionInformation",
				ClientMethod:   "GetStripeInfo",
				Parameters:     nil,
				ExpectedStatus: 200,
				Validation: func(response interface{}) error {
					return validateStripeInfoResponse(response, false)
				},
				StateUpdate: func(response interface{}, variables map[string]interface{}) error {
					// Clear variables (mimicking Postman cleanup)
					for key := range variables {
						if key != "story_name" {
							delete(variables, key)
						}
					}
					return nil
				},
			},
		},
		Setup: func() error {
			return cleanupStripeInfo()
		},
		Cleanup: func() error {
			return cleanupStripeInfo()
		},
	}
}

// GetPostmanStoryWithResponseMethods returns the Postman collection story using WithResponse methods
func GetPostmanStoryWithResponseMethods() testlib.Story {
	return testlib.Story{
		Name:        "Postman Collection Story - WithResponse Methods",
		Description: "Reproduces the Postman collection test flow using WithResponse client methods",
		Variables: map[string]interface{}{
			"secret_key": getTestStripeKey(),
		},
		Steps: []testlib.Step{
			{
				Name:           "Pre_GetStripeConnectionInformation",
				ClientMethod:   "GetStripeInfoWithResponse",
				Parameters:     nil,
				ExpectedStatus: 200,
				Validation: func(response interface{}) error {
					return validateStripeInfoWithResponse(response, false)
				},
				StateUpdate: func(response interface{}, variables map[string]interface{}) error {
					variables["secret_key"] = getTestStripeKey()
					variables["story_name"] = "Postman Collection Story - WithResponse Methods"
					return nil
				},
			},
			{
				Name:           "UpdateStripeConnectionInfo",
				ClientMethod:   "UpdateStripeInfoWithResponse",
				Parameters:     createUpdateStripeInfoParams(),
				ExpectedStatus: 200,
				Validation: func(response interface{}) error {
					return validateUpdateStripeInfoWithResponse(response)
				},
			},
			{
				Name:           "UpdateStripeConnectionInfo_WithBody",
				ClientMethod:   "UpdateStripeInfoWithBodyWithResponse",
				Parameters:     createUpdateStripeInfoBodyParams(),
				ExpectedStatus: 200,
				Validation: func(response interface{}) error {
					return validateUpdateStripeInfoWithResponse(response)
				},
			},
			{
				Name:           "GetStripeConnectionInformation",
				ClientMethod:   "GetStripeInfoWithResponse",
				Parameters:     nil,
				ExpectedStatus: 200,
				Validation: func(response interface{}) error {
					return validateStripeInfoWithResponse(response, true)
				},
			},
			{
				Name:           "DeleteStripeConnection",
				ClientMethod:   "DeleteStripeInfoWithResponse",
				Parameters:     nil,
				ExpectedStatus: 200,
				Validation: func(response interface{}) error {
					return validateDeleteStripeInfoWithResponse(response)
				},
			},
			{
				Name:           "Final_GetStripeConnectionInformation",
				ClientMethod:   "GetStripeInfoWithResponse",
				Parameters:     nil,
				ExpectedStatus: 200,
				Validation: func(response interface{}) error {
					return validateStripeInfoWithResponse(response, false)
				},
				StateUpdate: func(response interface{}, variables map[string]interface{}) error {
					for key := range variables {
						if key != "story_name" {
							delete(variables, key)
						}
					}
					return nil
				},
			},
		},
		Setup: func() error {
			return cleanupStripeInfo()
		},
		Cleanup: func() error {
			return cleanupStripeInfo()
		},
	}
}

// GetBillingMethods returns all billing API client methods
func GetBillingMethods() []string {
	return []string{
		// Standard methods
		"DeleteStripeInfo",
		"GetStripeInfo",
		"UpdateStripeInfoWithBody",
		"UpdateStripeInfo",
		// WithResponse methods
		"DeleteStripeInfoWithResponse",
		"GetStripeInfoWithResponse",
		"UpdateStripeInfoWithBodyWithResponse",
		"UpdateStripeInfoWithResponse",
	}
}
