package communicationapi

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"reflect"
	"strings"
	"time"

	"github.com/saasus-platform/saasus-sdk-go/generated/communicationapi"
)

// ============================================================================
// Feedback Validation Functions (Subtask 2.1)
// ============================================================================

// validateFeedbackResponse validates standard Feedback response
func validateFeedbackResponse(response interface{}) error {
	var feedback *communicationapi.Feedback

	// Check if response is an HTTP response with status code
	if httpResp, ok := response.(*http.Response); ok {
		if httpResp.StatusCode == 401 {
			return fmt.Errorf("authentication failed: received 401 Unauthorized - check SAASUS_SECRET_KEY, SAASUS_API_KEY, and SAASUS_SAAS_ID")
		}
		if httpResp.StatusCode < 200 || httpResp.StatusCode >= 300 {
			return fmt.Errorf("unexpected HTTP status: %d", httpResp.StatusCode)
		}

		// Decode body
		bodyBytes, err := io.ReadAll(httpResp.Body)
		if err != nil {
			// Handle context canceled errors gracefully
			if strings.Contains(err.Error(), "context canceled") {
				// Treat as empty body
				bodyBytes = []byte{}
			} else {
				return fmt.Errorf("failed to read response body: %w", err)
			}
		}
		httpResp.Body.Close()
		httpResp.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

		if len(bodyBytes) == 0 {
			// Empty body, skip decoding
			feedback = nil
		} else {
			var f communicationapi.Feedback
			if err := json.Unmarshal(bodyBytes, &f); err != nil {
				return fmt.Errorf("failed to decode feedback response: %w", err)
			}
			feedback = &f
		}
	} else if f, ok := response.(*communicationapi.Feedback); ok {
		feedback = f
	}

	if feedback != nil {
		if feedback.Id == "" {
			return fmt.Errorf("feedback ID is empty")
		}
		if feedback.FeedbackTitle == "" {
			return fmt.Errorf("feedback title is empty")
		}
		if feedback.FeedbackDescription == "" {
			return fmt.Errorf("feedback description is empty")
		}
		if feedback.UserId == "" {
			return fmt.Errorf("feedback user ID is empty")
		}
	}

	return nil
}

// validateFeedbackWithResponse validates WithResponse Feedback response
func validateFeedbackWithResponse(response interface{}) error {
	// Check if response has StatusCode method (WithResponse types)
	respValue := reflect.ValueOf(response)
	if respValue.Kind() == reflect.Ptr && !respValue.IsNil() {
		statusCodeMethod := respValue.MethodByName("StatusCode")
		if statusCodeMethod.IsValid() {
			results := statusCodeMethod.Call([]reflect.Value{})
			if len(results) > 0 {
				if statusCode, ok := results[0].Interface().(int); ok {
					if statusCode == 401 {
						return fmt.Errorf("authentication failed: received 401 Unauthorized - check SAASUS_SECRET_KEY, SAASUS_API_KEY, and SAASUS_SAAS_ID")
					}
					if statusCode < 200 || statusCode >= 300 {
						return fmt.Errorf("unexpected HTTP status: %d", statusCode)
					}
				}
			}
		}

		// Try to get JSON200 field for body validation
		json200Field := respValue.Elem().FieldByName("JSON200")
		if json200Field.IsValid() && !json200Field.IsNil() {
			if feedback, ok := json200Field.Interface().(*communicationapi.Feedback); ok {
				if feedback.Id == "" {
					return fmt.Errorf("feedback ID is empty")
				}
				if feedback.FeedbackTitle == "" {
					return fmt.Errorf("feedback title is empty")
				}
				if feedback.FeedbackDescription == "" {
					return fmt.Errorf("feedback description is empty")
				}
				if feedback.UserId == "" {
					return fmt.Errorf("feedback user ID is empty")
				}
			}
		}
	}

	return nil
}

// validateFeedbacksListResponse validates standard Feedbacks list response
func validateFeedbacksListResponse(response interface{}) error {
	var feedbacks *communicationapi.Feedbacks

	// Check if response is an HTTP response with status code
	if httpResp, ok := response.(*http.Response); ok {
		if httpResp.StatusCode == 401 {
			return fmt.Errorf("authentication failed: received 401 Unauthorized - check SAASUS_SECRET_KEY, SAASUS_API_KEY, and SAASUS_SAAS_ID")
		}
		if httpResp.StatusCode < 200 || httpResp.StatusCode >= 300 {
			return fmt.Errorf("unexpected HTTP status: %d", httpResp.StatusCode)
		}

		// Decode body
		bodyBytes, err := io.ReadAll(httpResp.Body)
		if err != nil {
			// Handle context canceled errors gracefully
			if strings.Contains(err.Error(), "context canceled") {
				// Treat as empty body
				bodyBytes = []byte{}
			} else {
				return fmt.Errorf("failed to read response body: %w", err)
			}
		}
		httpResp.Body.Close()
		httpResp.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

		if len(bodyBytes) == 0 {
			// Empty body, skip decoding
			feedbacks = nil
		} else {
			var f communicationapi.Feedbacks
			if err := json.Unmarshal(bodyBytes, &f); err != nil {
				return fmt.Errorf("failed to decode feedbacks list response: %w", err)
			}
			feedbacks = &f
		}
	} else if f, ok := response.(*communicationapi.Feedbacks); ok {
		feedbacks = f
	}

	if feedbacks != nil {
		// List can be empty, so just check the structure exists
		_ = feedbacks
	}

	return nil
}

// validateFeedbacksListWithResponse validates WithResponse Feedbacks list response
func validateFeedbacksListWithResponse(response interface{}) error {
	// Check if response has StatusCode method (WithResponse types)
	respValue := reflect.ValueOf(response)
	if respValue.Kind() == reflect.Ptr && !respValue.IsNil() {
		statusCodeMethod := respValue.MethodByName("StatusCode")
		if statusCodeMethod.IsValid() {
			results := statusCodeMethod.Call([]reflect.Value{})
			if len(results) > 0 {
				if statusCode, ok := results[0].Interface().(int); ok {
					if statusCode == 401 {
						return fmt.Errorf("authentication failed: received 401 Unauthorized - check SAASUS_SECRET_KEY, SAASUS_API_KEY, and SAASUS_SAAS_ID")
					}
					if statusCode < 200 || statusCode >= 300 {
						return fmt.Errorf("unexpected HTTP status: %d", statusCode)
					}
				}
			}
		}
	}

	return nil
}

// ============================================================================
// Comment Validation Functions (Subtask 2.2)
// ============================================================================

// validateFeedbackCommentResponse validates standard Comment response
func validateFeedbackCommentResponse(response interface{}) error {
	var comment *communicationapi.Comment

	// Check if response is an HTTP response with status code
	if httpResp, ok := response.(*http.Response); ok {
		if httpResp.StatusCode == 401 {
			return fmt.Errorf("authentication failed: received 401 Unauthorized - check SAASUS_SECRET_KEY, SAASUS_API_KEY, and SAASUS_SAAS_ID")
		}
		if httpResp.StatusCode < 200 || httpResp.StatusCode >= 300 {
			return fmt.Errorf("unexpected HTTP status: %d", httpResp.StatusCode)
		}

		// Decode body
		bodyBytes, err := io.ReadAll(httpResp.Body)
		if err != nil {
			// Handle context canceled errors gracefully
			if strings.Contains(err.Error(), "context canceled") {
				// Treat as empty body
				bodyBytes = []byte{}
			} else {
				return fmt.Errorf("failed to read response body: %w", err)
			}
		}
		httpResp.Body.Close()
		httpResp.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

		if len(bodyBytes) == 0 {
			// Empty body, skip decoding
			comment = nil
		} else {
			var c communicationapi.Comment
			if err := json.Unmarshal(bodyBytes, &c); err != nil {
				return fmt.Errorf("failed to decode comment response: %w", err)
			}
			comment = &c
		}
	} else if c, ok := response.(*communicationapi.Comment); ok {
		comment = c
	}

	if comment != nil {
		if comment.Id == "" {
			return fmt.Errorf("comment ID is empty")
		}
		if comment.Body == "" {
			return fmt.Errorf("comment body is empty")
		}
		if comment.CreatedAt == 0 {
			return fmt.Errorf("comment created_at is zero")
		}
	}

	return nil
}

// validateFeedbackCommentWithResponse validates WithResponse Comment response
func validateFeedbackCommentWithResponse(response interface{}) error {
	// Check if response has StatusCode method (WithResponse types)
	respValue := reflect.ValueOf(response)
	if respValue.Kind() == reflect.Ptr && !respValue.IsNil() {
		statusCodeMethod := respValue.MethodByName("StatusCode")
		if statusCodeMethod.IsValid() {
			results := statusCodeMethod.Call([]reflect.Value{})
			if len(results) > 0 {
				if statusCode, ok := results[0].Interface().(int); ok {
					if statusCode == 401 {
						return fmt.Errorf("authentication failed: received 401 Unauthorized - check SAASUS_SECRET_KEY, SAASUS_API_KEY, and SAASUS_SAAS_ID")
					}
					if statusCode < 200 || statusCode >= 300 {
						return fmt.Errorf("unexpected HTTP status: %d", statusCode)
					}
				}
			}
		}

		// Try to get JSON200 or JSON201 field for body validation
		json200Field := respValue.Elem().FieldByName("JSON200")
		json201Field := respValue.Elem().FieldByName("JSON201")

		var comment *communicationapi.Comment
		if json200Field.IsValid() && !json200Field.IsNil() {
			comment, _ = json200Field.Interface().(*communicationapi.Comment)
		} else if json201Field.IsValid() && !json201Field.IsNil() {
			comment, _ = json201Field.Interface().(*communicationapi.Comment)
		}

		if comment != nil {
			if comment.Id == "" {
				return fmt.Errorf("comment ID is empty")
			}
			if comment.Body == "" {
				return fmt.Errorf("comment body is empty")
			}
			if comment.CreatedAt == 0 {
				return fmt.Errorf("comment created_at is zero")
			}
		}
	}

	return nil
}

// ============================================================================
// Other Validation Functions (Subtask 2.3)
// ============================================================================

// validateUpdateFeedbackStatusResponse validates standard status update response
func validateUpdateFeedbackStatusResponse(response interface{}) error {
	// Check if response is an HTTP response with status code
	if httpResp, ok := response.(*http.Response); ok {
		if httpResp.StatusCode == 401 {
			return fmt.Errorf("authentication failed: received 401 Unauthorized - check SAASUS_SECRET_KEY, SAASUS_API_KEY, and SAASUS_SAAS_ID")
		}
		if httpResp.StatusCode < 200 || httpResp.StatusCode >= 300 {
			return fmt.Errorf("unexpected HTTP status: %d", httpResp.StatusCode)
		}
	}
	return nil
}

// validateUpdateFeedbackStatusWithResponse validates WithResponse status update response
func validateUpdateFeedbackStatusWithResponse(response interface{}) error {
	// Check if response has StatusCode method (WithResponse types)
	respValue := reflect.ValueOf(response)
	if respValue.Kind() == reflect.Ptr && !respValue.IsNil() {
		statusCodeMethod := respValue.MethodByName("StatusCode")
		if statusCodeMethod.IsValid() {
			results := statusCodeMethod.Call([]reflect.Value{})
			if len(results) > 0 {
				if statusCode, ok := results[0].Interface().(int); ok {
					if statusCode == 401 {
						return fmt.Errorf("authentication failed: received 401 Unauthorized - check SAASUS_SECRET_KEY, SAASUS_API_KEY, and SAASUS_SAAS_ID")
					}
					if statusCode < 200 || statusCode >= 300 {
						return fmt.Errorf("unexpected HTTP status: %d", statusCode)
					}
				}
			}
		}
	}
	return nil
}

// validateCreateVoteUserResponse validates standard vote creation response
func validateCreateVoteUserResponse(response interface{}) error {
	// Check if response is an HTTP response with status code
	if httpResp, ok := response.(*http.Response); ok {
		if httpResp.StatusCode == 401 {
			return fmt.Errorf("authentication failed: received 401 Unauthorized - check SAASUS_SECRET_KEY, SAASUS_API_KEY, and SAASUS_SAAS_ID")
		}
		if httpResp.StatusCode < 200 || httpResp.StatusCode >= 300 {
			return fmt.Errorf("unexpected HTTP status: %d", httpResp.StatusCode)
		}
	}
	return nil
}

// validateCreateVoteUserWithResponse validates WithResponse vote creation response
func validateCreateVoteUserWithResponse(response interface{}) error {
	// Check if response has StatusCode method (WithResponse types)
	respValue := reflect.ValueOf(response)
	if respValue.Kind() == reflect.Ptr && !respValue.IsNil() {
		statusCodeMethod := respValue.MethodByName("StatusCode")
		if statusCodeMethod.IsValid() {
			results := statusCodeMethod.Call([]reflect.Value{})
			if len(results) > 0 {
				if statusCode, ok := results[0].Interface().(int); ok {
					if statusCode == 401 {
						return fmt.Errorf("authentication failed: received 401 Unauthorized - check SAASUS_SECRET_KEY, SAASUS_API_KEY, and SAASUS_SAAS_ID")
					}
					if statusCode < 200 || statusCode >= 300 {
						return fmt.Errorf("unexpected HTTP status: %d", statusCode)
					}
				}
			}
		}
	}
	return nil
}

// validateDeleteVoteResponse validates standard vote deletion response
func validateDeleteVoteResponse(response interface{}) error {
	// Check if response is an HTTP response with status code
	if httpResp, ok := response.(*http.Response); ok {
		if httpResp.StatusCode == 401 {
			return fmt.Errorf("authentication failed: received 401 Unauthorized - check SAASUS_SECRET_KEY, SAASUS_API_KEY, and SAASUS_SAAS_ID")
		}
		if httpResp.StatusCode < 200 || httpResp.StatusCode >= 300 {
			return fmt.Errorf("unexpected HTTP status: %d", httpResp.StatusCode)
		}
	}
	return nil
}

// validateDeleteVoteWithResponse validates WithResponse vote deletion response
func validateDeleteVoteWithResponse(response interface{}) error {
	// Check if response has StatusCode method (WithResponse types)
	respValue := reflect.ValueOf(response)
	if respValue.Kind() == reflect.Ptr && !respValue.IsNil() {
		statusCodeMethod := respValue.MethodByName("StatusCode")
		if statusCodeMethod.IsValid() {
			results := statusCodeMethod.Call([]reflect.Value{})
			if len(results) > 0 {
				if statusCode, ok := results[0].Interface().(int); ok {
					if statusCode == 401 {
						return fmt.Errorf("authentication failed: received 401 Unauthorized - check SAASUS_SECRET_KEY, SAASUS_API_KEY, and SAASUS_SAAS_ID")
					}
					if statusCode < 200 || statusCode >= 300 {
						return fmt.Errorf("unexpected HTTP status: %d", statusCode)
					}
				}
			}
		}
	}
	return nil
}

// ============================================================================
// Test Data Generation Functions (Subtask 2.4)
// ============================================================================

// createTestFeedbackParams creates parameters for CreateFeedback
func createTestFeedbackParams(variables map[string]interface{}) interface{} {
	userID := getTestUserID()
	if uid, ok := variables["user_id"].(string); ok && uid != "" {
		userID = uid
	}

	return communicationapi.CreateFeedbackJSONRequestBody{
		UserId:              userID,
		FeedbackTitle:       "Test Feedback Title",
		FeedbackDescription: "Test Feedback Description for E2E testing",
	}
}

// createTestFeedbackWithBodyParams creates parameters for CreateFeedbackWithBody
func createTestFeedbackWithBodyParams(variables map[string]interface{}) interface{} {
	userID := getTestUserID()
	if uid, ok := variables["user_id"].(string); ok && uid != "" {
		userID = uid
	}

	param := communicationapi.CreateFeedbackParam{
		UserId:              userID,
		FeedbackTitle:       "Test Feedback Title (WithBody)",
		FeedbackDescription: "Test Feedback Description for E2E testing (WithBody)",
	}
	body, _ := json.Marshal(param)
	return struct {
		ContentType string
		Body        io.Reader
	}{
		ContentType: "application/json",
		Body:        bytes.NewReader(body),
	}
}

// createFeedbackCommentParams creates parameters for CreateFeedbackComment
func createFeedbackCommentParams(variables map[string]interface{}) interface{} {
	feedbackID, _ := variables["feedback_id"].(string)
	return struct {
		FeedbackId string
		Body       communicationapi.CreateFeedbackCommentJSONRequestBody
	}{
		FeedbackId: feedbackID,
		Body: communicationapi.CreateFeedbackCommentParam{
			Body: "Test comment body",
		},
	}
}

// createFeedbackCommentWithBodyParams creates parameters for CreateFeedbackCommentWithBody
func createFeedbackCommentWithBodyParams(variables map[string]interface{}) interface{} {
	feedbackID, _ := variables["feedback_id"].(string)
	param := communicationapi.CreateFeedbackCommentParam{
		Body: "Test comment body (WithBody)",
	}
	body, _ := json.Marshal(param)
	return struct {
		FeedbackId  string
		ContentType string
		Body        io.Reader
	}{
		FeedbackId:  feedbackID,
		ContentType: "application/json",
		Body:        bytes.NewReader(body),
	}
}

// createVoteUserParams creates parameters for CreateVoteUser
func createVoteUserParams(variables map[string]interface{}) interface{} {
	feedbackID, _ := variables["feedback_id"].(string)
	userID := getTestUserID()
	if uid, ok := variables["user_id"].(string); ok && uid != "" {
		userID = uid
	}

	return struct {
		FeedbackId string
		Body       communicationapi.CreateVoteUserJSONRequestBody
	}{
		FeedbackId: feedbackID,
		Body: communicationapi.CreateVoteUserParam{
			UserId: userID,
		},
	}
}

// createVoteUserWithBodyParams creates parameters for CreateVoteUserWithBody
func createVoteUserWithBodyParams(variables map[string]interface{}) interface{} {
	feedbackID, _ := variables["feedback_id"].(string)
	userID := getTestUserID()
	if uid, ok := variables["user_id"].(string); ok && uid != "" {
		userID = uid
	}

	param := communicationapi.CreateVoteUserParam{
		UserId: userID,
	}
	body, _ := json.Marshal(param)
	return struct {
		FeedbackId  string
		ContentType string
		Body        io.Reader
	}{
		FeedbackId:  feedbackID,
		ContentType: "application/json",
		Body:        bytes.NewReader(body),
	}
}

// updateFeedbackParams creates parameters for UpdateFeedback
func updateFeedbackParams(variables map[string]interface{}) interface{} {
	feedbackID, _ := variables["feedback_id"].(string)
	return struct {
		FeedbackId string
		Body       communicationapi.UpdateFeedbackJSONRequestBody
	}{
		FeedbackId: feedbackID,
		Body: communicationapi.UpdateFeedbackParam{
			FeedbackTitle:       "Updated Feedback Title",
			FeedbackDescription: "Updated Feedback Description",
		},
	}
}

// updateFeedbackWithBodyParams creates parameters for UpdateFeedbackWithBody
func updateFeedbackWithBodyParams(variables map[string]interface{}) interface{} {
	feedbackID, _ := variables["feedback_id"].(string)
	param := communicationapi.UpdateFeedbackParam{
		FeedbackTitle:       "Updated Feedback Title (WithBody)",
		FeedbackDescription: "Updated Feedback Description (WithBody)",
	}
	body, _ := json.Marshal(param)
	return struct {
		FeedbackId  string
		ContentType string
		Body        io.Reader
	}{
		FeedbackId:  feedbackID,
		ContentType: "application/json",
		Body:        bytes.NewReader(body),
	}
}

// updateFeedbackCommentParams creates parameters for UpdateFeedbackComment
func updateFeedbackCommentParams(variables map[string]interface{}) interface{} {
	feedbackID, _ := variables["feedback_id"].(string)
	commentID, _ := variables["comment_id"].(string)
	return struct {
		FeedbackId string
		CommentId  string
		Body       communicationapi.UpdateFeedbackCommentJSONRequestBody
	}{
		FeedbackId: feedbackID,
		CommentId:  commentID,
		Body: communicationapi.UpdateFeedbackCommentParam{
			Body: "Updated comment body",
		},
	}
}

// updateFeedbackCommentWithBodyParams creates parameters for UpdateFeedbackCommentWithBody
func updateFeedbackCommentWithBodyParams(variables map[string]interface{}) interface{} {
	feedbackID, _ := variables["feedback_id"].(string)
	commentID, _ := variables["comment_id"].(string)
	param := communicationapi.UpdateFeedbackCommentParam{
		Body: "Updated comment body (WithBody)",
	}
	body, _ := json.Marshal(param)
	return struct {
		FeedbackId  string
		CommentId   string
		ContentType string
		Body        io.Reader
	}{
		FeedbackId:  feedbackID,
		CommentId:   commentID,
		ContentType: "application/json",
		Body:        bytes.NewReader(body),
	}
}

// updateFeedbackStatusParams creates parameters for UpdateFeedbackStatus
func updateFeedbackStatusParams(variables map[string]interface{}) interface{} {
	feedbackID, _ := variables["feedback_id"].(string)
	return struct {
		FeedbackId string
		Body       communicationapi.UpdateFeedbackStatusJSONRequestBody
	}{
		FeedbackId: feedbackID,
		Body: communicationapi.UpdateFeedbackStatusParam{
			Status: 1, // Status 1 for testing
		},
	}
}

// updateFeedbackStatusWithBodyParams creates parameters for UpdateFeedbackStatusWithBody
func updateFeedbackStatusWithBodyParams(variables map[string]interface{}) interface{} {
	feedbackID, _ := variables["feedback_id"].(string)
	param := communicationapi.UpdateFeedbackStatusParam{
		Status: 1, // Status 1 for testing
	}
	body, _ := json.Marshal(param)
	return struct {
		FeedbackId  string
		ContentType string
		Body        io.Reader
	}{
		FeedbackId:  feedbackID,
		ContentType: "application/json",
		Body:        bytes.NewReader(body),
	}
}

// getFeedbackParams creates parameters for GetFeedback
func getFeedbackParams(variables map[string]interface{}) interface{} {
	feedbackID, _ := variables["feedback_id"].(string)
	return feedbackID
}

// getFeedbackCommentParams creates parameters for GetFeedbackComment
func getFeedbackCommentParams(variables map[string]interface{}) interface{} {
	feedbackID, _ := variables["feedback_id"].(string)
	commentID, _ := variables["comment_id"].(string)
	return struct {
		FeedbackId string
		CommentId  string
	}{
		FeedbackId: feedbackID,
		CommentId:  commentID,
	}
}

// deleteFeedbackParams creates parameters for DeleteFeedback
func deleteFeedbackParams(variables map[string]interface{}) interface{} {
	feedbackID, _ := variables["feedback_id"].(string)
	return struct {
		FeedbackId string
	}{
		FeedbackId: feedbackID,
	}
}

// deleteFeedbackCommentParams creates parameters for DeleteFeedbackComment
func deleteFeedbackCommentParams(variables map[string]interface{}) interface{} {
	feedbackID, _ := variables["feedback_id"].(string)
	commentID, _ := variables["comment_id"].(string)
	return struct {
		FeedbackId string
		CommentId  string
	}{
		FeedbackId: feedbackID,
		CommentId:  commentID,
	}
}

// deleteVoteParams creates parameters for DeleteVoteForFeedback
func deleteVoteParams(variables map[string]interface{}) interface{} {
	feedbackID, _ := variables["feedback_id"].(string)
	userID := getTestUserID()
	if uid, ok := variables["user_id"].(string); ok && uid != "" {
		userID = uid
	}

	return struct {
		FeedbackId string
		UserId     string
	}{
		FeedbackId: feedbackID,
		UserId:     userID,
	}
}

// ============================================================================
// Utility Functions (Subtask 2.5)
// ============================================================================

// cleanupTestFeedbacks cleans up any test feedbacks
func cleanupTestFeedbacks() error {
	// Implementation will clean up any existing test feedbacks
	// This is a placeholder for the actual cleanup logic
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Try to delete any existing test feedbacks
	// This will be implemented in the actual test execution engine
	_ = ctx
	return nil
}

// getTestUserID returns test user ID from environment or default
func getTestUserID() string {
	return getEnvWithDefault("TEST_USER_ID", "00000000-0000-0000-0000-000000000000")
}

// getEnvWithDefault returns environment variable or default value
func getEnvWithDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// extractFeedbackID extracts feedback ID from standard response
func extractFeedbackID(response interface{}) (string, error) {
	if id, ok := extractFeedbackIDFromMapResponse(response); ok {
		return id, nil
	}

	if feedback, ok := response.(*communicationapi.Feedback); ok {
		if feedback.Id == "" {
			return "", fmt.Errorf("feedback ID is empty in response")
		}
		return feedback.Id, nil
	}

	// Try to extract from HTTP response body
	if httpResp, ok := response.(*http.Response); ok {
		defer httpResp.Body.Close()
		var feedback communicationapi.Feedback
		if err := json.NewDecoder(httpResp.Body).Decode(&feedback); err != nil {
			return "", fmt.Errorf("failed to decode feedback response: %w", err)
		}
		if feedback.Id == "" {
			return "", fmt.Errorf("feedback ID is empty in response")
		}
		return feedback.Id, nil
	}

	return "", fmt.Errorf("unable to extract feedback ID from response type: %T", response)
}

// extractFeedbackIDFromWithResponse extracts feedback ID from WithResponse response
func extractFeedbackIDFromWithResponse(response interface{}) (string, error) {
	if id, ok := extractFeedbackIDFromMapResponse(response); ok {
		return id, nil
	}

	respValue := reflect.ValueOf(response)
	if respValue.Kind() == reflect.Ptr && !respValue.IsNil() {
		// Try to get JSON200 or JSON201 field
		json200Field := respValue.Elem().FieldByName("JSON200")
		json201Field := respValue.Elem().FieldByName("JSON201")

		var feedback *communicationapi.Feedback
		if json200Field.IsValid() && !json200Field.IsNil() {
			feedback, _ = json200Field.Interface().(*communicationapi.Feedback)
		} else if json201Field.IsValid() && !json201Field.IsNil() {
			feedback, _ = json201Field.Interface().(*communicationapi.Feedback)
		}

		if feedback != nil && feedback.Id != "" {
			return feedback.Id, nil
		}
	}

	return "", fmt.Errorf("unable to extract feedback ID from WithResponse type: %T", response)
}

func extractFeedbackIDFromMapResponse(response interface{}) (string, bool) {
	data, ok := response.(map[string]interface{})
	if !ok {
		return "", false
	}

	if id, ok := extractFeedbackIDFromMap(data); ok {
		return id, true
	}

	return "", false
}

func extractFeedbackIDFromMap(data map[string]interface{}) (string, bool) {
	if id, ok := getStringField(data, "id"); ok {
		return id, true
	}

	if id, ok := extractFeedbackIDFromNestedMap(data, "JSON200"); ok {
		return id, true
	}
	if id, ok := extractFeedbackIDFromNestedMap(data, "JSON201"); ok {
		return id, true
	}

	return "", false
}

func extractFeedbackIDFromNestedMap(data map[string]interface{}, key string) (string, bool) {
	raw, ok := data[key]
	if !ok {
		return "", false
	}

	nested, ok := raw.(map[string]interface{})
	if !ok {
		return "", false
	}

	return getStringField(nested, "id")
}

func getStringField(data map[string]interface{}, key string) (string, bool) {
	value, ok := data[key]
	if !ok {
		return "", false
	}

	id, ok := value.(string)
	if !ok || id == "" {
		return "", false
	}

	return id, true
}

// extractCommentID extracts comment ID from standard response
func extractCommentID(response interface{}) (string, error) {
	if id, ok := extractCommentIDFromMapResponse(response); ok {
		return id, nil
	}

	if comment, ok := response.(*communicationapi.Comment); ok {
		if comment.Id == "" {
			return "", fmt.Errorf("comment ID is empty in response")
		}
		return comment.Id, nil
	}

	// Try to extract from HTTP response body
	if httpResp, ok := response.(*http.Response); ok {
		defer httpResp.Body.Close()
		var comment communicationapi.Comment
		if err := json.NewDecoder(httpResp.Body).Decode(&comment); err != nil {
			return "", fmt.Errorf("failed to decode comment response: %w", err)
		}
		if comment.Id == "" {
			return "", fmt.Errorf("comment ID is empty in response")
		}
		return comment.Id, nil
	}

	return "", fmt.Errorf("unable to extract comment ID from response type: %T", response)
}

// extractCommentIDFromWithResponse extracts comment ID from WithResponse response
func extractCommentIDFromWithResponse(response interface{}) (string, error) {
	if id, ok := extractCommentIDFromMapResponse(response); ok {
		return id, nil
	}

	respValue := reflect.ValueOf(response)
	if respValue.Kind() == reflect.Ptr && !respValue.IsNil() {
		// Try to get JSON200 or JSON201 field
		json200Field := respValue.Elem().FieldByName("JSON200")
		json201Field := respValue.Elem().FieldByName("JSON201")

		var comment *communicationapi.Comment
		if json200Field.IsValid() && !json200Field.IsNil() {
			comment, _ = json200Field.Interface().(*communicationapi.Comment)
		} else if json201Field.IsValid() && !json201Field.IsNil() {
			comment, _ = json201Field.Interface().(*communicationapi.Comment)
		}

		if comment != nil && comment.Id != "" {
			return comment.Id, nil
		}
	}

	return "", fmt.Errorf("unable to extract comment ID from WithResponse type: %T", response)
}

func extractCommentIDFromMapResponse(response interface{}) (string, bool) {
	data, ok := response.(map[string]interface{})
	if !ok {
		return "", false
	}

	if id, ok := extractCommentIDFromMap(data); ok {
		return id, true
	}

	return "", false
}

func extractCommentIDFromMap(data map[string]interface{}) (string, bool) {
	if id, ok := getStringField(data, "id"); ok {
		return id, true
	}

	if id, ok := extractCommentIDFromNestedMap(data, "JSON200"); ok {
		return id, true
	}
	if id, ok := extractCommentIDFromNestedMap(data, "JSON201"); ok {
		return id, true
	}

	return "", false
}

func extractCommentIDFromNestedMap(data map[string]interface{}, key string) (string, bool) {
	raw, ok := data[key]
	if !ok {
		return "", false
	}

	nested, ok := raw.(map[string]interface{})
	if !ok {
		return "", false
	}

	return getStringField(nested, "id")
}

// validateDeleteFeedbackCommentWithResponse validates WithResponse comment deletion response
func validateDeleteFeedbackCommentWithResponse(response interface{}) error {
	// Check if response has StatusCode method (WithResponse types)
	respValue := reflect.ValueOf(response)
	if respValue.Kind() == reflect.Ptr && !respValue.IsNil() {
		statusCodeMethod := respValue.MethodByName("StatusCode")
		if statusCodeMethod.IsValid() {
			results := statusCodeMethod.Call([]reflect.Value{})
			if len(results) > 0 {
				if statusCode, ok := results[0].Interface().(int); ok {
					if statusCode == 401 {
						return fmt.Errorf("authentication failed: received 401 Unauthorized - check SAASUS_SECRET_KEY, SAASUS_API_KEY, and SAASUS_SAAS_ID")
					}
					if statusCode < 200 || statusCode >= 300 {
						return fmt.Errorf("unexpected HTTP status: %d", statusCode)
					}
				}
			}
		}
	}
	return nil
}

// validateDeleteFeedbackWithResponse validates WithResponse feedback deletion response
func validateDeleteFeedbackWithResponse(response interface{}) error {
	// Check if response has StatusCode method (WithResponse types)
	respValue := reflect.ValueOf(response)
	if respValue.Kind() == reflect.Ptr && !respValue.IsNil() {
		statusCodeMethod := respValue.MethodByName("StatusCode")
		if statusCodeMethod.IsValid() {
			results := statusCodeMethod.Call([]reflect.Value{})
			if len(results) > 0 {
				if statusCode, ok := results[0].Interface().(int); ok {
					if statusCode == 401 {
						return fmt.Errorf("authentication failed: received 401 Unauthorized - check SAASUS_SECRET_KEY, SAASUS_API_KEY, and SAASUS_SAAS_ID")
					}
					if statusCode < 200 || statusCode >= 300 {
						return fmt.Errorf("unexpected HTTP status: %d", statusCode)
					}
				}
			}
		}
	}
	return nil
}
