package integrationapi

import (
	"github.com/saasus-platform/saasus-sdk-go/tests/testlib"
)

// GetIntegrationStories returns all integration API test stories
func GetIntegrationStories() []testlib.Story {
	return []testlib.Story{
		GetPostmanStoryStandardMethods(),
		GetPostmanStoryStandardWithBodyMethods(),
		GetPostmanStoryWithResponseMethods(),
		GetPostmanStoryWithBodyWithResponseMethods(),
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
				Name:           "Pre_GetEventBridgeSettings",
				ClientMethod:   "GetEventBridgeSettings",
				Parameters:     nil,
				ExpectedStatus: 200,
				Validation: func(response interface{}) error {
					return validateEventBridgeSettingsResponse(response, false)
				},
				StateUpdate: func(response interface{}, variables map[string]interface{}) error {
					variables["aws_account_id"] = getTestAWSAccountID()
					variables["aws_region"] = getTestAWSRegion()
					variables["story_name"] = "Postman Collection Story - Standard Methods"
					return nil
				},
			},
			{
				Name:           "SaveEventBridgeSettings",
				ClientMethod:   "SaveEventBridgeSettings",
				Parameters:     createSaveEventBridgeSettingsParams(),
				ExpectedStatus: 200,
				Validation: func(response interface{}) error {
					return validateEventBridgeSettingsResponse(response, true)
				},
			},
			{
				Name:           "GetEventBridgeSettings_AfterSave",
				ClientMethod:   "GetEventBridgeSettings",
				Parameters:     nil,
				ExpectedStatus: 200,
				Validation: func(response interface{}) error {
					return validateEventBridgeSettingsResponse(response, true)
				},
			},
			{
				Name:           "DeleteEventBridgeSettings",
				ClientMethod:   "DeleteEventBridgeSettings",
				Parameters:     nil,
				ExpectedStatus: 200,
				Validation: func(response interface{}) error {
					return nil
				},
			},
			{
				Name:           "GetEventBridgeSettings_AfterDelete",
				ClientMethod:   "GetEventBridgeSettings",
				Parameters:     nil,
				ExpectedStatus: 200,
				Validation: func(response interface{}) error {
					return validateEventBridgeSettingsResponse(response, false)
				},
			},
			{
				Name:           "GetEventBridgeSettings_Setup",
				ClientMethod:   "GetEventBridgeSettings",
				Parameters:     nil,
				ExpectedStatus: 200,
				Validation: func(response interface{}) error {
					return validateEventBridgeSettingsResponse(response, false)
				},
				StateUpdate: func(response interface{}, variables map[string]interface{}) error {
					variables["aws_account_id"] = getTestAWSAccountID()
					variables["aws_region"] = getTestAWSRegion()
					return nil
				},
			},
			{
				Name:           "SaveEventBridgeSettings_ForTest",
				ClientMethod:   "SaveEventBridgeSettings",
				Parameters:     createSaveEventBridgeSettingsParams(),
				ExpectedStatus: 200,
				Validation: func(response interface{}) error {
					return validateEventBridgeSettingsResponse(response, true)
				},
			},
			{
				Name:           "CreateEventBridgeTestEvent",
				ClientMethod:   "CreateEventBridgeTestEvent",
				Parameters:     nil,
				ExpectedStatus: 201,
				Validation: func(response interface{}) error {
					return nil
				},
			},
			{
				Name:           "CreateEventBridgeEvent",
				ClientMethod:   "CreateEventBridgeEvent",
				Parameters:     createEventBridgeEventParams(),
				ExpectedStatus: 501,
				Validation: func(response interface{}) error {
					return nil
				},
			},
			{
				Name:           "DeleteEventBridgeSettings_Cleanup",
				ClientMethod:   "DeleteEventBridgeSettings",
				Parameters:     nil,
				ExpectedStatus: 200,
				Validation: func(response interface{}) error {
					return nil
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
			return cleanupEventBridgeSettings()
		},
		Cleanup: func() error {
			return cleanupEventBridgeSettings()
		},
	}
}

// GetPostmanStoryStandardWithBodyMethods returns the Postman collection story using standard WithBody methods
func GetPostmanStoryStandardWithBodyMethods() testlib.Story {
	return testlib.Story{
		Name:        "Postman Collection Story - Standard Methods With Body",
		Description: "Reproduces the Postman collection test flow using standard WithBody client methods",
		Variables:   map[string]interface{}{},
		Steps: []testlib.Step{
			{
				Name:           "Pre_GetEventBridgeSettings",
				ClientMethod:   "GetEventBridgeSettings",
				Parameters:     nil,
				ExpectedStatus: 200,
				Validation: func(response interface{}) error {
					return validateEventBridgeSettingsResponse(response, false)
				},
				StateUpdate: func(response interface{}, variables map[string]interface{}) error {
					variables["aws_account_id"] = getTestAWSAccountID()
					variables["aws_region"] = getTestAWSRegion()
					variables["story_name"] = "Postman Collection Story - Standard Methods With Body"
					return nil
				},
			},
			{
				Name:           "SaveEventBridgeSettingsWithBody",
				ClientMethod:   "SaveEventBridgeSettingsWithBody",
				Parameters:     createSaveEventBridgeSettingsBodyParams(),
				ExpectedStatus: 200,
				Validation: func(response interface{}) error {
					return validateEventBridgeSettingsResponse(response, true)
				},
			},
			{
				Name:           "GetEventBridgeSettings_AfterSave",
				ClientMethod:   "GetEventBridgeSettings",
				Parameters:     nil,
				ExpectedStatus: 200,
				Validation: func(response interface{}) error {
					return validateEventBridgeSettingsResponse(response, true)
				},
			},
			{
				Name:           "DeleteEventBridgeSettings",
				ClientMethod:   "DeleteEventBridgeSettings",
				Parameters:     nil,
				ExpectedStatus: 200,
				Validation: func(response interface{}) error {
					return nil
				},
			},
			{
				Name:           "GetEventBridgeSettings_AfterDelete",
				ClientMethod:   "GetEventBridgeSettings",
				Parameters:     nil,
				ExpectedStatus: 200,
				Validation: func(response interface{}) error {
					return validateEventBridgeSettingsResponse(response, false)
				},
			},
			{
				Name:           "GetEventBridgeSettings_Setup",
				ClientMethod:   "GetEventBridgeSettings",
				Parameters:     nil,
				ExpectedStatus: 200,
				Validation: func(response interface{}) error {
					return validateEventBridgeSettingsResponse(response, false)
				},
				StateUpdate: func(response interface{}, variables map[string]interface{}) error {
					variables["aws_account_id"] = getTestAWSAccountID()
					variables["aws_region"] = getTestAWSRegion()
					return nil
				},
			},
			{
				Name:           "SaveEventBridgeSettingsWithBody_ForTest",
				ClientMethod:   "SaveEventBridgeSettingsWithBody",
				Parameters:     createSaveEventBridgeSettingsBodyParams(),
				ExpectedStatus: 200,
				Validation: func(response interface{}) error {
					return validateEventBridgeSettingsResponse(response, true)
				},
			},
			{
				Name:           "CreateEventBridgeTestEvent",
				ClientMethod:   "CreateEventBridgeTestEvent",
				Parameters:     nil,
				ExpectedStatus: 201,
				Validation: func(response interface{}) error {
					return nil
				},
			},
			{
				Name:           "CreateEventBridgeEventWithBody",
				ClientMethod:   "CreateEventBridgeEventWithBody",
				Parameters:     createEventBridgeEventBodyParams(),
				ExpectedStatus: 501,
				Validation: func(response interface{}) error {
					return nil
				},
			},
			{
				Name:           "DeleteEventBridgeSettings_Cleanup",
				ClientMethod:   "DeleteEventBridgeSettings",
				Parameters:     nil,
				ExpectedStatus: 200,
				Validation: func(response interface{}) error {
					return nil
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
			return cleanupEventBridgeSettings()
		},
		Cleanup: func() error {
			return cleanupEventBridgeSettings()
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
				Name:           "Pre_GetEventBridgeSettingsWithResponse",
				ClientMethod:   "GetEventBridgeSettingsWithResponse",
				Parameters:     nil,
				ExpectedStatus: 200,
				Validation: func(response interface{}) error {
					return validateEventBridgeSettingsWithResponse(response, false)
				},
				StateUpdate: func(response interface{}, variables map[string]interface{}) error {
					variables["aws_account_id"] = getTestAWSAccountID()
					variables["aws_region"] = getTestAWSRegion()
					variables["story_name"] = "Postman Collection Story - WithResponse Methods"
					return nil
				},
			},
			{
				Name:           "SaveEventBridgeSettingsWithResponse",
				ClientMethod:   "SaveEventBridgeSettingsWithResponse",
				Parameters:     createSaveEventBridgeSettingsParams(),
				ExpectedStatus: 200,
				Validation: func(response interface{}) error {
					return validateEventBridgeSettingsWithResponse(response, true)
				},
			},
			{
				Name:           "GetEventBridgeSettingsWithResponse_AfterSave",
				ClientMethod:   "GetEventBridgeSettingsWithResponse",
				Parameters:     nil,
				ExpectedStatus: 200,
				Validation: func(response interface{}) error {
					return validateEventBridgeSettingsWithResponse(response, true)
				},
			},
			{
				Name:           "DeleteEventBridgeSettingsWithResponse",
				ClientMethod:   "DeleteEventBridgeSettingsWithResponse",
				Parameters:     nil,
				ExpectedStatus: 200,
				Validation: func(response interface{}) error {
					return validateEventBridgeSettingsWithResponse(response, true)
				},
			},
			{
				Name:           "GetEventBridgeSettingsWithResponse_AfterDelete",
				ClientMethod:   "GetEventBridgeSettingsWithResponse",
				Parameters:     nil,
				ExpectedStatus: 200,
				Validation: func(response interface{}) error {
					return validateEventBridgeSettingsWithResponse(response, false)
				},
			},
			{
				Name:           "GetEventBridgeSettingsWithResponse_Setup",
				ClientMethod:   "GetEventBridgeSettingsWithResponse",
				Parameters:     nil,
				ExpectedStatus: 200,
				Validation: func(response interface{}) error {
					return validateEventBridgeSettingsWithResponse(response, false)
				},
				StateUpdate: func(response interface{}, variables map[string]interface{}) error {
					variables["aws_account_id"] = getTestAWSAccountID()
					variables["aws_region"] = getTestAWSRegion()
					return nil
				},
			},
			{
				Name:           "SaveEventBridgeSettingsWithResponse_ForTest",
				ClientMethod:   "SaveEventBridgeSettingsWithResponse",
				Parameters:     createSaveEventBridgeSettingsParams(),
				ExpectedStatus: 200,
				Validation: func(response interface{}) error {
					return validateEventBridgeSettingsWithResponse(response, true)
				},
			},
			{
				Name:           "CreateEventBridgeTestEventWithResponse",
				ClientMethod:   "CreateEventBridgeTestEventWithResponse",
				Parameters:     nil,
				ExpectedStatus: 201,
				Validation: func(response interface{}) error {
					return validateEventBridgeSettingsWithResponse(response, true)
				},
			},
			{
				Name:           "CreateEventBridgeEventWithResponse",
				ClientMethod:   "CreateEventBridgeEventWithResponse",
				Parameters:     createEventBridgeEventParams(),
				ExpectedStatus: 501,
				Validation:     nil,
			},
			{
				Name:           "DeleteEventBridgeSettingsWithResponse_Cleanup",
				ClientMethod:   "DeleteEventBridgeSettingsWithResponse",
				Parameters:     nil,
				ExpectedStatus: 200,
				Validation: func(response interface{}) error {
					return validateEventBridgeSettingsWithResponse(response, true)
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
			return cleanupEventBridgeSettings()
		},
		Cleanup: func() error {
			return cleanupEventBridgeSettings()
		},
	}
}

// GetPostmanStoryWithBodyWithResponseMethods returns the Postman collection story using WithBodyWithResponse methods
func GetPostmanStoryWithBodyWithResponseMethods() testlib.Story {
	return testlib.Story{
		Name:        "Postman Collection Story - WithResponse Methods With Body",
		Description: "Reproduces the Postman collection test flow using WithBodyWithResponse client methods",
		Variables:   map[string]interface{}{},
		Steps: []testlib.Step{
			{
				Name:           "Pre_GetEventBridgeSettingsWithResponse",
				ClientMethod:   "GetEventBridgeSettingsWithResponse",
				Parameters:     nil,
				ExpectedStatus: 200,
				Validation: func(response interface{}) error {
					return validateEventBridgeSettingsWithResponse(response, false)
				},
				StateUpdate: func(response interface{}, variables map[string]interface{}) error {
					variables["aws_account_id"] = getTestAWSAccountID()
					variables["aws_region"] = getTestAWSRegion()
					variables["story_name"] = "Postman Collection Story - WithResponse Methods With Body"
					return nil
				},
			},
			{
				Name:           "SaveEventBridgeSettingsWithBodyWithResponse",
				ClientMethod:   "SaveEventBridgeSettingsWithBodyWithResponse",
				Parameters:     createSaveEventBridgeSettingsBodyParams(),
				ExpectedStatus: 200,
				Validation: func(response interface{}) error {
					return validateEventBridgeSettingsWithResponse(response, true)
				},
			},
			{
				Name:           "GetEventBridgeSettingsWithResponse_AfterSave",
				ClientMethod:   "GetEventBridgeSettingsWithResponse",
				Parameters:     nil,
				ExpectedStatus: 200,
				Validation: func(response interface{}) error {
					return validateEventBridgeSettingsWithResponse(response, true)
				},
			},
			{
				Name:           "DeleteEventBridgeSettingsWithResponse",
				ClientMethod:   "DeleteEventBridgeSettingsWithResponse",
				Parameters:     nil,
				ExpectedStatus: 200,
				Validation: func(response interface{}) error {
					return validateEventBridgeSettingsWithResponse(response, true)
				},
			},
			{
				Name:           "GetEventBridgeSettingsWithResponse_AfterDelete",
				ClientMethod:   "GetEventBridgeSettingsWithResponse",
				Parameters:     nil,
				ExpectedStatus: 200,
				Validation: func(response interface{}) error {
					return validateEventBridgeSettingsWithResponse(response, false)
				},
			},
			{
				Name:           "GetEventBridgeSettingsWithResponse_Setup",
				ClientMethod:   "GetEventBridgeSettingsWithResponse",
				Parameters:     nil,
				ExpectedStatus: 200,
				Validation: func(response interface{}) error {
					return validateEventBridgeSettingsWithResponse(response, false)
				},
				StateUpdate: func(response interface{}, variables map[string]interface{}) error {
					variables["aws_account_id"] = getTestAWSAccountID()
					variables["aws_region"] = getTestAWSRegion()
					return nil
				},
			},
			{
				Name:           "SaveEventBridgeSettingsWithBodyWithResponse_ForTest",
				ClientMethod:   "SaveEventBridgeSettingsWithBodyWithResponse",
				Parameters:     createSaveEventBridgeSettingsBodyParams(),
				ExpectedStatus: 200,
				Validation: func(response interface{}) error {
					return validateEventBridgeSettingsWithResponse(response, true)
				},
			},
			{
				Name:           "CreateEventBridgeTestEventWithResponse",
				ClientMethod:   "CreateEventBridgeTestEventWithResponse",
				Parameters:     nil,
				ExpectedStatus: 201,
				Validation: func(response interface{}) error {
					return validateEventBridgeSettingsWithResponse(response, true)
				},
			},
			{
				Name:           "CreateEventBridgeEventWithBodyWithResponse",
				ClientMethod:   "CreateEventBridgeEventWithBodyWithResponse",
				Parameters:     createEventBridgeEventBodyParams(),
				ExpectedStatus: 501,
				Validation:     nil,
			},
			{
				Name:           "DeleteEventBridgeSettingsWithResponse_Cleanup",
				ClientMethod:   "DeleteEventBridgeSettingsWithResponse",
				Parameters:     nil,
				ExpectedStatus: 200,
				Validation: func(response interface{}) error {
					return validateEventBridgeSettingsWithResponse(response, true)
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
			return cleanupEventBridgeSettings()
		},
		Cleanup: func() error {
			return cleanupEventBridgeSettings()
		},
	}
}

// GetIntegrationMethods returns all integration API client methods
func GetIntegrationMethods() []string {
	return []string{
		// Standard methods
		"GetEventBridgeSettings",
		"SaveEventBridgeSettings",
		"DeleteEventBridgeSettings",
		"CreateEventBridgeEvent",
		"CreateEventBridgeTestEvent",
		// WithBody methods
		"SaveEventBridgeSettingsWithBody",
		"CreateEventBridgeEventWithBody",
		// WithResponse methods
		"GetEventBridgeSettingsWithResponse",
		"SaveEventBridgeSettingsWithResponse",
		"DeleteEventBridgeSettingsWithResponse",
		"CreateEventBridgeEventWithResponse",
		"CreateEventBridgeTestEventWithResponse",
		// WithBodyWithResponse methods
		"SaveEventBridgeSettingsWithBodyWithResponse",
		"CreateEventBridgeEventWithBodyWithResponse",
	}
}
