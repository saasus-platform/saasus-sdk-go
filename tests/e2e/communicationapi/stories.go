package communicationapi

import (
	"github.com/saasus-platform/saasus-sdk-go/tests/testlib"
)

// GetCommunicationStories returns all communication API test stories
func GetCommunicationStories() []testlib.Story {
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
		Description: "フィードバック管理の完全なライフサイクルをStandardメソッドでテスト",
		Variables: map[string]interface{}{
			"feedback_id": "",
			"comment_id":  "",
			"user_id":     "",
		},
		Steps: []testlib.Step{
			{
				Name:           "GetFeedbacks",
				ClientMethod:   "GetFeedbacks",
				Parameters:     nil,
				ExpectedStatus: 200,
				Validation: func(response interface{}) error {
					return validateFeedbacksListResponse(response)
				},
				StateUpdate: func(response interface{}, variables map[string]interface{}) error {
					variables["story_name"] = "Postman Collection Story - Standard Methods"
					variables["user_id"] = getTestUserID()
					return nil
				},
			},
			{
				Name:           "CreateFeedback",
				ClientMethod:   "CreateFeedback",
				Parameters:     createTestFeedbackParams,
				ExpectedStatus: 201,
				Validation: func(response interface{}) error {
					return validateFeedbackResponse(response)
				},
				StateUpdate: func(response interface{}, variables map[string]interface{}) error {
					feedbackID, err := extractFeedbackID(response)
					if err != nil {
						return err
					}
					variables["feedback_id"] = feedbackID
					return nil
				},
			},
			{
				Name:           "GetFeedback",
				ClientMethod:   "GetFeedback",
				Parameters:     getFeedbackParams,
				ExpectedStatus: 200,
				Validation: func(response interface{}) error {
					return validateFeedbackResponse(response)
				},
			},
			{
				Name:           "UpdateFeedback",
				ClientMethod:   "UpdateFeedback",
				Parameters:     updateFeedbackParams,
				ExpectedStatus: 200,
				Validation: func(response interface{}) error {
					return validateFeedbackResponse(response)
				},
			},
			{
				Name:           "UpdateFeedbackStatus",
				ClientMethod:   "UpdateFeedbackStatus",
				Parameters:     updateFeedbackStatusParams,
				ExpectedStatus: 200,
				Validation: func(response interface{}) error {
					return validateUpdateFeedbackStatusResponse(response)
				},
			},
			{
				Name:           "CreateFeedbackComment",
				ClientMethod:   "CreateFeedbackComment",
				Parameters:     createFeedbackCommentParams,
				ExpectedStatus: 201,
				Validation: func(response interface{}) error {
					return validateFeedbackCommentResponse(response)
				},
				StateUpdate: func(response interface{}, variables map[string]interface{}) error {
					commentID, err := extractCommentID(response)
					if err != nil {
						return err
					}
					variables["comment_id"] = commentID
					return nil
				},
			},
			{
				Name:           "GetFeedbackComment",
				ClientMethod:   "GetFeedbackComment",
				Parameters:     getFeedbackCommentParams,
				ExpectedStatus: 200,
				Validation: func(response interface{}) error {
					return validateFeedbackCommentResponse(response)
				},
			},
			{
				Name:           "UpdateFeedbackComment",
				ClientMethod:   "UpdateFeedbackComment",
				Parameters:     updateFeedbackCommentParams,
				ExpectedStatus: 200,
				Validation: func(response interface{}) error {
					return validateFeedbackCommentResponse(response)
				},
			},
			{
				Name:           "CreateVoteUser",
				ClientMethod:   "CreateVoteUser",
				Parameters:     createVoteUserParams,
				ExpectedStatus: 201,
				Validation: func(response interface{}) error {
					return validateCreateVoteUserResponse(response)
				},
			},
			{
				Name:           "DeleteVoteForFeedback",
				ClientMethod:   "DeleteVoteForFeedback",
				Parameters:     deleteVoteParams,
				ExpectedStatus: 200,
				Validation: func(response interface{}) error {
					return validateDeleteVoteResponse(response)
				},
			},
			{
				Name:           "DeleteFeedbackComment",
				ClientMethod:   "DeleteFeedbackComment",
				Parameters:     deleteFeedbackCommentParams,
				ExpectedStatus: 200,
				Validation: func(response interface{}) error {
					return nil
				},
			},
			{
				Name:           "DeleteFeedback",
				ClientMethod:   "DeleteFeedback",
				Parameters:     deleteFeedbackParams,
				ExpectedStatus: 200,
				Validation: func(response interface{}) error {
					return nil
				},
				StateUpdate: func(response interface{}, variables map[string]interface{}) error {
					// クリーンアップ: 変数をクリア
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
			return cleanupTestFeedbacks()
		},
		Cleanup: func() error {
			return cleanupTestFeedbacks()
		},
	}
}

// GetPostmanStoryStandardWithBodyMethods returns the Postman collection story using standard WithBody methods
func GetPostmanStoryStandardWithBodyMethods() testlib.Story {
	return testlib.Story{
		Name:        "Postman Collection Story - Standard WithBody Methods",
		Description: "フィードバック管理の完全なライフサイクルをStandard WithBodyメソッドでテスト",
		Variables: map[string]interface{}{
			"feedback_id": "",
			"comment_id":  "",
			"user_id":     "",
		},
		Steps: []testlib.Step{
			{
				Name:           "GetFeedbacks",
				ClientMethod:   "GetFeedbacks",
				Parameters:     nil,
				ExpectedStatus: 200,
				Validation: func(response interface{}) error {
					return validateFeedbacksListResponse(response)
				},
				StateUpdate: func(response interface{}, variables map[string]interface{}) error {
					variables["story_name"] = "Postman Collection Story - Standard WithBody Methods"
					variables["user_id"] = getTestUserID()
					return nil
				},
			},
			{
				Name:           "CreateFeedbackWithBody",
				ClientMethod:   "CreateFeedbackWithBody",
				Parameters:     createTestFeedbackWithBodyParams,
				ExpectedStatus: 201,
				Validation: func(response interface{}) error {
					return validateFeedbackResponse(response)
				},
				StateUpdate: func(response interface{}, variables map[string]interface{}) error {
					feedbackID, err := extractFeedbackID(response)
					if err != nil {
						return err
					}
					variables["feedback_id"] = feedbackID
					return nil
				},
			},
			{
				Name:           "GetFeedback",
				ClientMethod:   "GetFeedback",
				Parameters:     getFeedbackParams,
				ExpectedStatus: 200,
				Validation: func(response interface{}) error {
					return validateFeedbackResponse(response)
				},
			},
			{
				Name:           "UpdateFeedbackWithBody",
				ClientMethod:   "UpdateFeedbackWithBody",
				Parameters:     updateFeedbackWithBodyParams,
				ExpectedStatus: 200,
				Validation: func(response interface{}) error {
					return validateFeedbackResponse(response)
				},
			},
			{
				Name:           "UpdateFeedbackStatusWithBody",
				ClientMethod:   "UpdateFeedbackStatusWithBody",
				Parameters:     updateFeedbackStatusWithBodyParams,
				ExpectedStatus: 200,
				Validation: func(response interface{}) error {
					return validateUpdateFeedbackStatusResponse(response)
				},
			},
			{
				Name:           "CreateFeedbackCommentWithBody",
				ClientMethod:   "CreateFeedbackCommentWithBody",
				Parameters:     createFeedbackCommentWithBodyParams,
				ExpectedStatus: 201,
				Validation: func(response interface{}) error {
					return validateFeedbackCommentResponse(response)
				},
				StateUpdate: func(response interface{}, variables map[string]interface{}) error {
					commentID, err := extractCommentID(response)
					if err != nil {
						return err
					}
					variables["comment_id"] = commentID
					return nil
				},
			},
			{
				Name:           "GetFeedbackComment",
				ClientMethod:   "GetFeedbackComment",
				Parameters:     getFeedbackCommentParams,
				ExpectedStatus: 200,
				Validation: func(response interface{}) error {
					return validateFeedbackCommentResponse(response)
				},
			},
			{
				Name:           "UpdateFeedbackCommentWithBody",
				ClientMethod:   "UpdateFeedbackCommentWithBody",
				Parameters:     updateFeedbackCommentWithBodyParams,
				ExpectedStatus: 200,
				Validation: func(response interface{}) error {
					return validateFeedbackCommentResponse(response)
				},
			},
			{
				Name:           "CreateVoteUserWithBody",
				ClientMethod:   "CreateVoteUserWithBody",
				Parameters:     createVoteUserWithBodyParams,
				ExpectedStatus: 201,
				Validation: func(response interface{}) error {
					return validateCreateVoteUserResponse(response)
				},
			},
			{
				Name:           "DeleteVoteForFeedback",
				ClientMethod:   "DeleteVoteForFeedback",
				Parameters:     deleteVoteParams,
				ExpectedStatus: 200,
				Validation: func(response interface{}) error {
					return validateDeleteVoteResponse(response)
				},
			},
			{
				Name:           "DeleteFeedbackComment",
				ClientMethod:   "DeleteFeedbackComment",
				Parameters:     deleteFeedbackCommentParams,
				ExpectedStatus: 200,
				Validation: func(response interface{}) error {
					return nil
				},
			},
			{
				Name:           "DeleteFeedback",
				ClientMethod:   "DeleteFeedback",
				Parameters:     deleteFeedbackParams,
				ExpectedStatus: 200,
				Validation: func(response interface{}) error {
					return nil
				},
				StateUpdate: func(response interface{}, variables map[string]interface{}) error {
					// クリーンアップ: 変数をクリア
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
			return cleanupTestFeedbacks()
		},
		Cleanup: func() error {
			return cleanupTestFeedbacks()
		},
	}
}

// GetPostmanStoryWithResponseMethods returns the Postman collection story using WithResponse methods
func GetPostmanStoryWithResponseMethods() testlib.Story {
	return testlib.Story{
		Name:        "Postman Collection Story - WithResponse Methods",
		Description: "フィードバック管理の完全なライフサイクルをWithResponseメソッドでテスト",
		Variables: map[string]interface{}{
			"feedback_id": "",
			"comment_id":  "",
			"user_id":     "",
		},
		Steps: []testlib.Step{
			{
				Name:           "GetFeedbacks",
				ClientMethod:   "GetFeedbacksWithResponse",
				Parameters:     nil,
				ExpectedStatus: 200,
				Validation: func(response interface{}) error {
					return validateFeedbacksListWithResponse(response)
				},
				StateUpdate: func(response interface{}, variables map[string]interface{}) error {
					variables["story_name"] = "Postman Collection Story - WithResponse Methods"
					variables["user_id"] = getTestUserID()
					return nil
				},
			},
			{
				Name:           "CreateFeedback",
				ClientMethod:   "CreateFeedbackWithResponse",
				Parameters:     createTestFeedbackParams,
				ExpectedStatus: 201,
				Validation: func(response interface{}) error {
					return validateFeedbackWithResponse(response)
				},
				StateUpdate: func(response interface{}, variables map[string]interface{}) error {
					feedbackID, err := extractFeedbackIDFromWithResponse(response)
					if err != nil {
						return err
					}
					variables["feedback_id"] = feedbackID
					return nil
				},
			},
			{
				Name:           "GetFeedback",
				ClientMethod:   "GetFeedbackWithResponse",
				Parameters:     getFeedbackParams,
				ExpectedStatus: 200,
				Validation: func(response interface{}) error {
					return validateFeedbackWithResponse(response)
				},
			},
			{
				Name:           "UpdateFeedback",
				ClientMethod:   "UpdateFeedbackWithResponse",
				Parameters:     updateFeedbackParams,
				ExpectedStatus: 200,
				Validation: func(response interface{}) error {
					return validateFeedbackWithResponse(response)
				},
			},
			{
				Name:           "UpdateFeedbackStatus",
				ClientMethod:   "UpdateFeedbackStatusWithResponse",
				Parameters:     updateFeedbackStatusParams,
				ExpectedStatus: 200,
				Validation: func(response interface{}) error {
					return validateUpdateFeedbackStatusWithResponse(response)
				},
			},
			{
				Name:           "CreateFeedbackComment",
				ClientMethod:   "CreateFeedbackCommentWithResponse",
				Parameters:     createFeedbackCommentParams,
				ExpectedStatus: 201,
				Validation: func(response interface{}) error {
					return validateFeedbackCommentWithResponse(response)
				},
				StateUpdate: func(response interface{}, variables map[string]interface{}) error {
					commentID, err := extractCommentIDFromWithResponse(response)
					if err != nil {
						return err
					}
					variables["comment_id"] = commentID
					return nil
				},
			},
			{
				Name:           "GetFeedbackComment",
				ClientMethod:   "GetFeedbackCommentWithResponse",
				Parameters:     getFeedbackCommentParams,
				ExpectedStatus: 200,
				Validation: func(response interface{}) error {
					return validateFeedbackCommentWithResponse(response)
				},
			},
			{
				Name:           "UpdateFeedbackComment",
				ClientMethod:   "UpdateFeedbackCommentWithResponse",
				Parameters:     updateFeedbackCommentParams,
				ExpectedStatus: 200,
				Validation: func(response interface{}) error {
					return validateFeedbackCommentWithResponse(response)
				},
			},
			{
				Name:           "CreateVoteUser",
				ClientMethod:   "CreateVoteUserWithResponse",
				Parameters:     createVoteUserParams,
				ExpectedStatus: 201,
				Validation: func(response interface{}) error {
					return validateCreateVoteUserWithResponse(response)
				},
			},
			{
				Name:           "DeleteVoteForFeedback",
				ClientMethod:   "DeleteVoteForFeedbackWithResponse",
				Parameters:     deleteVoteParams,
				ExpectedStatus: 200,
				Validation: func(response interface{}) error {
					return validateDeleteVoteWithResponse(response)
				},
			},
			{
				Name:           "DeleteFeedbackComment",
				ClientMethod:   "DeleteFeedbackCommentWithResponse",
				Parameters:     deleteFeedbackCommentParams,
				ExpectedStatus: 200,
				Validation: func(response interface{}) error {
					return validateDeleteFeedbackCommentWithResponse(response)
				},
			},
			{
				Name:           "DeleteFeedback",
				ClientMethod:   "DeleteFeedbackWithResponse",
				Parameters:     deleteFeedbackParams,
				ExpectedStatus: 200,
				Validation: func(response interface{}) error {
					return validateDeleteFeedbackWithResponse(response)
				},
				StateUpdate: func(response interface{}, variables map[string]interface{}) error {
					// クリーンアップ: 変数をクリア
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
			return cleanupTestFeedbacks()
		},
		Cleanup: func() error {
			return cleanupTestFeedbacks()
		},
	}
}

// GetPostmanStoryWithBodyWithResponseMethods returns the Postman collection story using WithBodyWithResponse methods
func GetPostmanStoryWithBodyWithResponseMethods() testlib.Story {
	return testlib.Story{
		Name:        "Postman Collection Story - WithBodyWithResponse Methods",
		Description: "フィードバック管理の完全なライフサイクルをWithBodyWithResponseメソッドでテスト",
		Variables: map[string]interface{}{
			"feedback_id": "",
			"comment_id":  "",
			"user_id":     "",
		},
		Steps: []testlib.Step{
			{
				Name:           "GetFeedbacks",
				ClientMethod:   "GetFeedbacksWithResponse",
				Parameters:     nil,
				ExpectedStatus: 200,
				Validation: func(response interface{}) error {
					return validateFeedbacksListWithResponse(response)
				},
				StateUpdate: func(response interface{}, variables map[string]interface{}) error {
					variables["story_name"] = "Postman Collection Story - WithBodyWithResponse Methods"
					variables["user_id"] = getTestUserID()
					return nil
				},
			},
			{
				Name:           "CreateFeedbackWithBody",
				ClientMethod:   "CreateFeedbackWithBodyWithResponse",
				Parameters:     createTestFeedbackWithBodyParams,
				ExpectedStatus: 201,
				Validation: func(response interface{}) error {
					return validateFeedbackWithResponse(response)
				},
				StateUpdate: func(response interface{}, variables map[string]interface{}) error {
					feedbackID, err := extractFeedbackIDFromWithResponse(response)
					if err != nil {
						return err
					}
					variables["feedback_id"] = feedbackID
					return nil
				},
			},
			{
				Name:           "GetFeedback",
				ClientMethod:   "GetFeedbackWithResponse",
				Parameters:     getFeedbackParams,
				ExpectedStatus: 200,
				Validation: func(response interface{}) error {
					return validateFeedbackWithResponse(response)
				},
			},
			{
				Name:           "UpdateFeedbackWithBody",
				ClientMethod:   "UpdateFeedbackWithBodyWithResponse",
				Parameters:     updateFeedbackWithBodyParams,
				ExpectedStatus: 200,
				Validation: func(response interface{}) error {
					return validateFeedbackWithResponse(response)
				},
			},
			{
				Name:           "UpdateFeedbackStatusWithBody",
				ClientMethod:   "UpdateFeedbackStatusWithBodyWithResponse",
				Parameters:     updateFeedbackStatusWithBodyParams,
				ExpectedStatus: 200,
				Validation: func(response interface{}) error {
					return validateUpdateFeedbackStatusWithResponse(response)
				},
			},
			{
				Name:           "CreateFeedbackCommentWithBody",
				ClientMethod:   "CreateFeedbackCommentWithBodyWithResponse",
				Parameters:     createFeedbackCommentWithBodyParams,
				ExpectedStatus: 201,
				Validation: func(response interface{}) error {
					return validateFeedbackCommentWithResponse(response)
				},
				StateUpdate: func(response interface{}, variables map[string]interface{}) error {
					commentID, err := extractCommentIDFromWithResponse(response)
					if err != nil {
						return err
					}
					variables["comment_id"] = commentID
					return nil
				},
			},
			{
				Name:           "GetFeedbackComment",
				ClientMethod:   "GetFeedbackCommentWithResponse",
				Parameters:     getFeedbackCommentParams,
				ExpectedStatus: 200,
				Validation: func(response interface{}) error {
					return validateFeedbackCommentWithResponse(response)
				},
			},
			{
				Name:           "UpdateFeedbackCommentWithBody",
				ClientMethod:   "UpdateFeedbackCommentWithBodyWithResponse",
				Parameters:     updateFeedbackCommentWithBodyParams,
				ExpectedStatus: 200,
				Validation: func(response interface{}) error {
					return validateFeedbackCommentWithResponse(response)
				},
			},
			{
				Name:           "CreateVoteUserWithBody",
				ClientMethod:   "CreateVoteUserWithBodyWithResponse",
				Parameters:     createVoteUserWithBodyParams,
				ExpectedStatus: 201,
				Validation: func(response interface{}) error {
					return validateCreateVoteUserWithResponse(response)
				},
			},
			{
				Name:           "DeleteVoteForFeedback",
				ClientMethod:   "DeleteVoteForFeedbackWithResponse",
				Parameters:     deleteVoteParams,
				ExpectedStatus: 200,
				Validation: func(response interface{}) error {
					return validateDeleteVoteWithResponse(response)
				},
			},
			{
				Name:           "DeleteFeedbackComment",
				ClientMethod:   "DeleteFeedbackCommentWithResponse",
				Parameters:     deleteFeedbackCommentParams,
				ExpectedStatus: 200,
				Validation: func(response interface{}) error {
					return validateDeleteFeedbackCommentWithResponse(response)
				},
			},
			{
				Name:           "DeleteFeedback",
				ClientMethod:   "DeleteFeedbackWithResponse",
				Parameters:     deleteFeedbackParams,
				ExpectedStatus: 200,
				Validation: func(response interface{}) error {
					return validateDeleteFeedbackWithResponse(response)
				},
				StateUpdate: func(response interface{}, variables map[string]interface{}) error {
					// クリーンアップ: 変数をクリア
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
			return cleanupTestFeedbacks()
		},
		Cleanup: func() error {
			return cleanupTestFeedbacks()
		},
	}
}

// GetCommunicationMethods returns all communication API client methods
func GetCommunicationMethods() []string {
	return []string{
		// Standard methods
		"GetFeedbacks",
		"CreateFeedback",
		"CreateFeedbackWithBody",
		"GetFeedback",
		"UpdateFeedback",
		"UpdateFeedbackWithBody",
		"DeleteFeedback",
		"CreateFeedbackComment",
		"CreateFeedbackCommentWithBody",
		"GetFeedbackComment",
		"UpdateFeedbackComment",
		"UpdateFeedbackCommentWithBody",
		"DeleteFeedbackComment",
		"UpdateFeedbackStatus",
		"UpdateFeedbackStatusWithBody",
		"CreateVoteUser",
		"CreateVoteUserWithBody",
		"DeleteVoteForFeedback",
		// WithResponse methods
		"GetFeedbacksWithResponse",
		"CreateFeedbackWithResponse",
		"CreateFeedbackWithBodyWithResponse",
		"GetFeedbackWithResponse",
		"UpdateFeedbackWithResponse",
		"UpdateFeedbackWithBodyWithResponse",
		"DeleteFeedbackWithResponse",
		"CreateFeedbackCommentWithResponse",
		"CreateFeedbackCommentWithBodyWithResponse",
		"GetFeedbackCommentWithResponse",
		"UpdateFeedbackCommentWithResponse",
		"UpdateFeedbackCommentWithBodyWithResponse",
		"DeleteFeedbackCommentWithResponse",
		"UpdateFeedbackStatusWithResponse",
		"UpdateFeedbackStatusWithBodyWithResponse",
		"CreateVoteUserWithResponse",
		"CreateVoteUserWithBodyWithResponse",
		"DeleteVoteForFeedbackWithResponse",
	}
}
