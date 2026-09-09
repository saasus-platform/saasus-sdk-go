package snapshot

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"time"
)

// StoryComparator handles story-level snapshot comparisons
type StoryComparator struct {
	config *StorySnapshotConfig
	logger *SnapshotLogger
}

// NewStoryComparator creates a new story comparator
func NewStoryComparator(config *StorySnapshotConfig) *StoryComparator {
	// Create logger with appropriate level
	logLevel := LogLevelInfo
	debugMode := false
	if config != nil && config.CaptureLevel == CaptureLevelFull {
		logLevel = LogLevelDebug
		debugMode = true
	}

	logger := NewSnapshotLogger(logLevel, debugMode)

	return &StoryComparator{
		config: config,
		logger: logger,
	}
}

// getSnapshotTag extracts git tag from snapshot metadata
func getSnapshotTag(snapshot *StorySnapshot) string {
	if snapshot != nil && snapshot.Metadata.GitTag != "" {
		return snapshot.Metadata.GitTag
	}
	return "unknown"
}

// CompareStorySnapshots compares two story snapshots
func (sc *StoryComparator) CompareStorySnapshots(oldSnapshot, newSnapshot *StorySnapshot) (*StoryComparison, error) {
	if oldSnapshot == nil || newSnapshot == nil {
		return nil, NewComparisonError(ErrorTypeValidation,
			"cannot compare nil snapshots", nil, "")
	}

	sc.logger.LogPhaseStart(PhaseComparison, fmt.Sprintf("story '%s'", newSnapshot.StoryName))
	sc.logger.SetContext("story", newSnapshot.StoryName)
	defer sc.logger.ClearContext()

	comparison := &StoryComparison{
		StoryName:      newSnapshot.StoryName,
		ComparisonType: "story_level",
		Timestamp:      time.Now(),
		OldSnapshot:    oldSnapshot,
		NewSnapshot:    newSnapshot,
		Differences:    []StoryDifference{},
		CompatibilityReport: StoryCompatibilityReport{
			Level:                  Compatible,
			StoryFlowChanges:       []StoryFlowChange{},
			StepSequenceChanges:    []StepSequenceChange{},
			StateTransitionChanges: []StateTransitionChange{},
			Passed:                 true,
		},
	}

	// Compare story metadata
	sc.compareStoryMetadata(comparison, oldSnapshot, newSnapshot)

	// Compare step sequences
	sc.compareStepSequences(comparison, oldSnapshot.Steps, newSnapshot.Steps)

	// Compare state transitions
	sc.compareStateTransitions(comparison, oldSnapshot.Steps, newSnapshot.Steps)

	// Generate compatibility report
	sc.generateCompatibilityReport(comparison)

	// Generate summary
	sc.generateComparisonSummary(comparison)

	sc.logger.LogComparison(newSnapshot.StoryName,
		getSnapshotTag(oldSnapshot), getSnapshotTag(newSnapshot),
		len(comparison.Differences))

	return comparison, nil
}

// CompareWithPreviousRelease compares current snapshot with previous release
func (sc *StoryComparator) CompareWithPreviousRelease(currentSnapshot *StorySnapshot) (*StoryComparison, error) {
	if currentSnapshot == nil {
		return nil, NewComparisonError(ErrorTypeValidation,
			"current snapshot cannot be nil", nil, "")
	}

	sc.logger.LogInfo("Comparing with previous release for story: %s", currentSnapshot.StoryName)

	// Create file manager to get previous release tag
	var outputDir string
	if sc.config != nil {
		outputDir = sc.config.GetModuleOutputDirectory()
	}
	fileManager := NewStorySnapshotFileManager(outputDir)

	// Get current tag
	currentTag, err := GetGitTagForFileManager()
	if err != nil {
		sc.logger.LogWarning("Failed to get current git tag, using 'current': %v", err)
		currentTag = "current"
	}

	// Get previous release tag
	previousTag, err := fileManager.GetPreviousReleaseTag(currentTag, currentSnapshot.StoryName)
	if err != nil {
		snapErr := NewComparisonError(ErrorTypeFileIO,
			"failed to get previous release tag", err, currentSnapshot.StoryName)
		sc.logger.LogSnapshotError(snapErr)
		return nil, snapErr
	}

	if previousTag == "" {
		// No previous release, create baseline comparison
		sc.logger.LogInfo("No previous release found, creating baseline comparison")
		return sc.createBaselineComparison(currentSnapshot), nil
	}

	sc.logger.LogInfo("Comparing against previous release: %s", previousTag)

	// Load previous snapshot
	previousSnapshot, err := fileManager.LoadStorySnapshotByTag(previousTag, currentSnapshot.StoryName)
	if err != nil {
		snapErr := NewComparisonError(ErrorTypeFileIO,
			"failed to load previous snapshot", err, currentSnapshot.StoryName)
		sc.logger.LogSnapshotError(snapErr)
		return nil, snapErr
	}

	return sc.CompareStorySnapshots(previousSnapshot, currentSnapshot)
}

// LoadStorySnapshot loads a story snapshot from file
func (sc *StoryComparator) LoadStorySnapshot(filePath string) (*StorySnapshot, error) {
	sc.logger.LogFileOperation("read", filePath)

	data, err := os.ReadFile(filePath)
	if err != nil {
		snapErr := NewFileIOError("failed to read snapshot file", err, filePath)
		sc.logger.LogSnapshotError(snapErr)
		return nil, snapErr
	}

	var snapshot StorySnapshot
	if err := json.Unmarshal(data, &snapshot); err != nil {
		snapErr := NewComparisonError(ErrorTypeParsing,
			"failed to unmarshal snapshot", err, "")
		snapErr.Context.FilePath = filePath
		sc.logger.LogSnapshotError(snapErr)
		return nil, snapErr
	}

	sc.logger.LogDebug("Successfully loaded snapshot for story: %s", snapshot.StoryName)
	return &snapshot, nil
}

// SaveComparisonResult saves comparison result to file
func (sc *StoryComparator) SaveComparisonResult(result *StoryComparison, filePath string) error {
	if result == nil {
		return NewComparisonError(ErrorTypeValidation,
			"comparison result cannot be nil", nil, "")
	}

	if filePath != "" {
		// Use provided file path
		sc.logger.LogFileOperation("write", filePath)

		data, err := json.MarshalIndent(result, "", "  ")
		if err != nil {
			snapErr := NewComparisonError(ErrorTypeParsing,
				"failed to marshal comparison result", err, result.StoryName)
			sc.logger.LogSnapshotError(snapErr)
			return snapErr
		}

		if err := os.WriteFile(filePath, data, 0644); err != nil {
			snapErr := NewFileIOError("failed to write comparison file", err, filePath)
			sc.logger.LogSnapshotError(snapErr)
			return snapErr
		}

		sc.logger.LogInfo("Saved comparison result to: %s", filePath)
		return nil
	}

	// Use file manager for automatic naming
	var outputDir string
	if sc.config != nil {
		outputDir = sc.config.GetModuleOutputDirectory()
	}
	fileManager := NewStorySnapshotFileManager(outputDir)
	savedPath, err := fileManager.SaveStoryComparison(result)
	if err != nil {
		snapErr := NewComparisonError(ErrorTypeFileIO,
			"failed to save comparison using file manager", err, result.StoryName)
		sc.logger.LogSnapshotError(snapErr)
		return snapErr
	}

	sc.logger.LogInfo("Saved comparison result to: %s", savedPath)
	return nil
}

// compareStoryMetadata compares story-level metadata
func (sc *StoryComparator) compareStoryMetadata(comparison *StoryComparison, oldSnapshot, newSnapshot *StorySnapshot) {
	// Compare story status
	if oldSnapshot.Status != newSnapshot.Status {
		comparison.Differences = append(comparison.Differences, StoryDifference{
			Type:        "STORY_STATUS",
			Description: fmt.Sprintf("Story status changed from %s to %s", oldSnapshot.Status, newSnapshot.Status),
			OldValue:    oldSnapshot.Status,
			NewValue:    newSnapshot.Status,
			Impact:      sc.determineStatusChangeImpact(oldSnapshot.Status, newSnapshot.Status),
		})
	}

	// Compare story duration (significant changes only)
	oldDuration := oldSnapshot.Duration.Seconds()
	newDuration := newSnapshot.Duration.Seconds()
	if oldDuration > 0 && newDuration > 0 {
		changePercent := (newDuration - oldDuration) / oldDuration * 100
		if changePercent > 50 || changePercent < -50 { // More than 50% change
			comparison.Differences = append(comparison.Differences, StoryDifference{
				Type:        "STORY_DURATION",
				Description: fmt.Sprintf("Story duration changed significantly: %.2fs to %.2fs (%.1f%%)", oldDuration, newDuration, changePercent),
				OldValue:    oldDuration,
				NewValue:    newDuration,
				Impact:      Warning,
			})
		}
	}

	// Compare step count
	if len(oldSnapshot.Steps) != len(newSnapshot.Steps) {
		comparison.Differences = append(comparison.Differences, StoryDifference{
			Type:        "STEP_COUNT",
			Description: fmt.Sprintf("Step count changed from %d to %d", len(oldSnapshot.Steps), len(newSnapshot.Steps)),
			OldValue:    len(oldSnapshot.Steps),
			NewValue:    len(newSnapshot.Steps),
			Impact:      sc.determineStepCountChangeImpact(len(oldSnapshot.Steps), len(newSnapshot.Steps)),
		})
	}
}

// compareStepSequences compares step sequences between snapshots
func (sc *StoryComparator) compareStepSequences(comparison *StoryComparison, oldSteps, newSteps []StepSnapshot) {
	// Create step maps for easier comparison
	oldStepMap := make(map[string]*StepSnapshot)
	newStepMap := make(map[string]*StepSnapshot)

	for i := range oldSteps {
		oldStepMap[oldSteps[i].StepName] = &oldSteps[i]
	}

	for i := range newSteps {
		newStepMap[newSteps[i].StepName] = &newSteps[i]
	}

	// Compare existing steps
	for stepName, newStep := range newStepMap {
		if oldStep, exists := oldStepMap[stepName]; exists {
			// Compare individual step using existing SDK method comparison logic
			stepDiffs := sc.compareStepSnapshots(stepName, oldStep, newStep)
			for _, diff := range stepDiffs {
				comparison.Differences = append(comparison.Differences, diff)
			}
		} else {
			// New step added
			comparison.Differences = append(comparison.Differences, StoryDifference{
				Type:        "STEP_ADDED",
				StepName:    stepName,
				Description: fmt.Sprintf("New step added: %s", stepName),
				NewValue:    newStep,
				Impact:      Compatible,
			})

			comparison.CompatibilityReport.StepSequenceChanges = append(comparison.CompatibilityReport.StepSequenceChanges, StepSequenceChange{
				Type:        "ADDITION",
				StepName:    stepName,
				Description: fmt.Sprintf("Step %s was added to the story", stepName),
				Impact:      Compatible,
			})
		}
	}

	// Check for removed steps
	for stepName, oldStep := range oldStepMap {
		if _, exists := newStepMap[stepName]; !exists {
			comparison.Differences = append(comparison.Differences, StoryDifference{
				Type:        "STEP_REMOVED",
				StepName:    stepName,
				Description: fmt.Sprintf("Step removed: %s", stepName),
				OldValue:    oldStep,
				Impact:      Breaking,
			})

			comparison.CompatibilityReport.StepSequenceChanges = append(comparison.CompatibilityReport.StepSequenceChanges, StepSequenceChange{
				Type:        "REMOVAL",
				StepName:    stepName,
				Description: fmt.Sprintf("Step %s was removed from the story", stepName),
				Impact:      Breaking,
			})
		}
	}

	// Compare step order (if both have same steps)
	sc.compareStepOrder(comparison, oldSteps, newSteps)
}

// compareStepSnapshots compares individual step snapshots using existing SDK method logic
func (sc *StoryComparator) compareStepSnapshots(stepName string, oldStep, newStep *StepSnapshot) []StoryDifference {
	var differences []StoryDifference

	// Compare method name
	if oldStep.Method != newStep.Method {
		differences = append(differences, StoryDifference{
			Type:        "STEP_METHOD",
			StepName:    stepName,
			Field:       "Method",
			Description: fmt.Sprintf("Step method changed from %s to %s", oldStep.Method, newStep.Method),
			OldValue:    oldStep.Method,
			NewValue:    newStep.Method,
			Impact:      Breaking,
		})
	}

	// Compare success status
	if oldStep.Success != newStep.Success {
		differences = append(differences, StoryDifference{
			Type:        "STEP_SUCCESS",
			StepName:    stepName,
			Field:       "Success",
			Description: fmt.Sprintf("Step success status changed from %t to %t", oldStep.Success, newStep.Success),
			OldValue:    oldStep.Success,
			NewValue:    newStep.Success,
			Impact:      sc.determineSuccessChangeImpact(oldStep.Success, newStep.Success),
		})
	}

	// Compare return values using existing SDK method comparison logic
	if oldStep.ReturnValue != nil && newStep.ReturnValue != nil {
		sdkDiffs := compareSDKReturnValues(stepName, oldStep.ReturnValue, newStep.ReturnValue)
		for _, sdkDiff := range sdkDiffs {
			differences = append(differences, StoryDifference{
				Type:        fmt.Sprintf("STEP_%s", sdkDiff.Type),
				StepName:    stepName,
				Field:       sdkDiff.Field,
				Description: sdkDiff.Description,
				OldValue:    sdkDiff.OldValue,
				NewValue:    sdkDiff.NewValue,
				Impact:      sc.determineSDKDifferenceImpact(sdkDiff),
			})
		}
	} else if oldStep.ReturnValue == nil && newStep.ReturnValue != nil {
		differences = append(differences, StoryDifference{
			Type:        "STEP_RETURN_VALUE_ADDED",
			StepName:    stepName,
			Description: "Return value was added to step",
			NewValue:    newStep.ReturnValue,
			Impact:      Compatible,
		})
	} else if oldStep.ReturnValue != nil && newStep.ReturnValue == nil {
		differences = append(differences, StoryDifference{
			Type:        "STEP_RETURN_VALUE_REMOVED",
			StepName:    stepName,
			Description: "Return value was removed from step",
			OldValue:    oldStep.ReturnValue,
			Impact:      Breaking,
		})
	}

	// Compare error status
	if !reflect.DeepEqual(oldStep.Error, newStep.Error) {
		differences = append(differences, StoryDifference{
			Type:        "STEP_ERROR",
			StepName:    stepName,
			Field:       "Error",
			Description: "Step error information changed",
			OldValue:    oldStep.Error,
			NewValue:    newStep.Error,
			Impact:      Warning,
		})
	}

	return differences
}

// compareStepOrder compares the order of steps in the story
func (sc *StoryComparator) compareStepOrder(comparison *StoryComparison, oldSteps, newSteps []StepSnapshot) {
	// Create ordered step name lists
	oldOrder := make([]string, len(oldSteps))
	newOrder := make([]string, len(newSteps))

	for i, step := range oldSteps {
		oldOrder[i] = step.StepName
	}

	for i, step := range newSteps {
		newOrder[i] = step.StepName
	}

	// Find common steps and check their relative order
	commonSteps := sc.findCommonSteps(oldOrder, newOrder)
	if len(commonSteps) > 1 {
		oldPositions := sc.getStepPositions(oldOrder, commonSteps)
		newPositions := sc.getStepPositions(newOrder, commonSteps)

		if !sc.isSameOrder(oldPositions, newPositions) {
			comparison.Differences = append(comparison.Differences, StoryDifference{
				Type:        "STEP_ORDER",
				Description: "Step execution order changed",
				OldValue:    oldOrder,
				NewValue:    newOrder,
				Impact:      Warning,
			})

			comparison.CompatibilityReport.StepSequenceChanges = append(comparison.CompatibilityReport.StepSequenceChanges, StepSequenceChange{
				Type:        "ORDER_CHANGE",
				Description: "The order of step execution has changed",
				Impact:      Warning,
			})
		}
	}
}

// compareStateTransitions compares state transitions between steps
func (sc *StoryComparator) compareStateTransitions(comparison *StoryComparison, oldSteps, newSteps []StepSnapshot) {
	// Extract state transitions from steps
	oldTransitions := sc.extractStateTransitions(oldSteps)
	newTransitions := sc.extractStateTransitions(newSteps)

	// Compare transitions
	for transitionKey, newTransition := range newTransitions {
		if oldTransition, exists := oldTransitions[transitionKey]; exists {
			if !sc.compareTransitions(oldTransition, newTransition) {
				comparison.Differences = append(comparison.Differences, StoryDifference{
					Type:        "STATE_TRANSITION",
					Description: fmt.Sprintf("State transition changed: %s", transitionKey),
					OldValue:    oldTransition,
					NewValue:    newTransition,
					Impact:      Warning,
				})

				comparison.CompatibilityReport.StateTransitionChanges = append(comparison.CompatibilityReport.StateTransitionChanges, StateTransitionChange{
					Type:        "TRANSITION_CHANGE",
					Description: fmt.Sprintf("State transition %s has changed", transitionKey),
					Impact:      Warning,
					OldState:    oldTransition,
					NewState:    newTransition,
				})
			}
		} else {
			// New transition
			comparison.CompatibilityReport.StateTransitionChanges = append(comparison.CompatibilityReport.StateTransitionChanges, StateTransitionChange{
				Type:        "TRANSITION_ADDED",
				Description: fmt.Sprintf("New state transition added: %s", transitionKey),
				Impact:      Compatible,
				NewState:    newTransition,
			})
		}
	}

	// Check for removed transitions
	for transitionKey, oldTransition := range oldTransitions {
		if _, exists := newTransitions[transitionKey]; !exists {
			comparison.CompatibilityReport.StateTransitionChanges = append(comparison.CompatibilityReport.StateTransitionChanges, StateTransitionChange{
				Type:        "TRANSITION_REMOVED",
				Description: fmt.Sprintf("State transition removed: %s", transitionKey),
				Impact:      Breaking,
				OldState:    oldTransition,
			})
		}
	}
}

// Helper methods

// createBaselineComparison creates a baseline comparison for first release
func (sc *StoryComparator) createBaselineComparison(snapshot *StorySnapshot) *StoryComparison {
	return &StoryComparison{
		StoryName:      snapshot.StoryName,
		ComparisonType: "baseline",
		Timestamp:      time.Now(),
		NewSnapshot:    snapshot,
		Differences:    []StoryDifference{},
		CompatibilityReport: StoryCompatibilityReport{
			Level:   Compatible,
			Summary: "Baseline snapshot - no comparison available",
			Passed:  true,
		},
		Summary: StoryComparisonSummary{
			TotalDifferences:     0,
			OverallCompatibility: Compatible,
		},
	}
}

// determineStatusChangeImpact determines impact of story status change
func (sc *StoryComparator) determineStatusChangeImpact(oldStatus, newStatus TestStatus) CompatibilityLevel {
	if oldStatus == TestStatusPassed && newStatus != TestStatusPassed {
		return Breaking
	}
	if oldStatus != TestStatusPassed && newStatus == TestStatusPassed {
		return Compatible
	}
	return Warning
}

// determineStepCountChangeImpact determines impact of step count change
func (sc *StoryComparator) determineStepCountChangeImpact(oldCount, newCount int) CompatibilityLevel {
	if newCount > oldCount {
		return Compatible // Adding steps is generally compatible
	}
	if newCount < oldCount {
		return Breaking // Removing steps is breaking
	}
	return Compatible
}

// determineSuccessChangeImpact determines impact of step success change
func (sc *StoryComparator) determineSuccessChangeImpact(oldSuccess, newSuccess bool) CompatibilityLevel {
	if oldSuccess && !newSuccess {
		return Breaking // Success to failure is breaking
	}
	if !oldSuccess && newSuccess {
		return Compatible // Failure to success is improvement
	}
	return Compatible
}

// determineSDKDifferenceImpact maps SDK method differences to story-level impact
func (sc *StoryComparator) determineSDKDifferenceImpact(diff *SDKMethodDifference) CompatibilityLevel {
	switch diff.Type {
	case "STATUS_CODE":
		if oldCode, ok := diff.OldValue.(int); ok {
			if newCode, ok := diff.NewValue.(int); ok {
				// Success to error or vice versa is breaking
				if (oldCode >= 200 && oldCode < 300 && (newCode < 200 || newCode >= 300)) ||
					(newCode >= 200 && newCode < 300 && (oldCode < 200 || oldCode >= 300)) {
					return Breaking
				}
			}
		}
		return Warning
	case "TYPE":
		return Breaking
	case "JSON_DATA", "BODY":
		return Warning
	case "HEADERS":
		return Compatible
	default:
		return Warning
	}
}

// findCommonSteps finds steps that exist in both old and new sequences
func (sc *StoryComparator) findCommonSteps(oldOrder, newOrder []string) []string {
	oldSet := make(map[string]bool)
	for _, step := range oldOrder {
		oldSet[step] = true
	}

	var common []string
	for _, step := range newOrder {
		if oldSet[step] {
			common = append(common, step)
		}
	}

	return common
}

// getStepPositions gets positions of specified steps in the order
func (sc *StoryComparator) getStepPositions(order []string, steps []string) []int {
	positions := make([]int, len(steps))
	for i, step := range steps {
		for j, orderStep := range order {
			if step == orderStep {
				positions[i] = j
				break
			}
		}
	}
	return positions
}

// isSameOrder checks if two position arrays represent the same order
func (sc *StoryComparator) isSameOrder(oldPositions, newPositions []int) bool {
	if len(oldPositions) != len(newPositions) {
		return false
	}

	for i := 1; i < len(oldPositions); i++ {
		oldRelative := oldPositions[i] - oldPositions[i-1]
		newRelative := newPositions[i] - newPositions[i-1]
		if (oldRelative > 0) != (newRelative > 0) {
			return false
		}
	}

	return true
}

// extractStateTransitions extracts state transitions from step sequence
func (sc *StoryComparator) extractStateTransitions(steps []StepSnapshot) map[string]interface{} {
	transitions := make(map[string]interface{})

	for i := 0; i < len(steps)-1; i++ {
		currentStep := steps[i]
		nextStep := steps[i+1]

		transitionKey := fmt.Sprintf("%s->%s", currentStep.StepName, nextStep.StepName)

		// Extract relevant state information
		transition := map[string]interface{}{
			"from_step":        currentStep.StepName,
			"to_step":          nextStep.StepName,
			"from_success":     currentStep.Success,
			"to_success":       nextStep.Success,
			"from_status_code": sc.getStatusCodeFromStep(&currentStep),
			"to_status_code":   sc.getStatusCodeFromStep(&nextStep),
		}

		transitions[transitionKey] = transition
	}

	return transitions
}

// compareTransitions compares two state transitions
func (sc *StoryComparator) compareTransitions(oldTransition, newTransition interface{}) bool {
	return reflect.DeepEqual(oldTransition, newTransition)
}

// getStatusCodeFromStep extracts status code from step
func (sc *StoryComparator) getStatusCodeFromStep(step *StepSnapshot) int {
	if step.ReturnValue != nil {
		return step.ReturnValue.StatusCode
	}
	return step.StatusCode
}

// generateCompatibilityReport generates the compatibility report
func (sc *StoryComparator) generateCompatibilityReport(comparison *StoryComparison) {
	report := &comparison.CompatibilityReport

	// Determine overall compatibility level
	maxLevel := Compatible
	for _, diff := range comparison.Differences {
		if diff.Impact > maxLevel {
			maxLevel = diff.Impact
		}
	}

	report.Level = maxLevel
	report.Passed = maxLevel != Breaking

	// Generate story flow changes
	for _, diff := range comparison.Differences {
		if diff.Type == "STORY_STATUS" || diff.Type == "STORY_DURATION" || diff.Type == "STEP_COUNT" {
			report.StoryFlowChanges = append(report.StoryFlowChanges, StoryFlowChange{
				Type:        diff.Type,
				Description: diff.Description,
				Impact:      diff.Impact,
				Details:     map[string]interface{}{"old": diff.OldValue, "new": diff.NewValue},
			})
		}
	}

	// Generate summary
	if len(comparison.Differences) == 0 {
		report.Summary = "No differences detected between story snapshots"
	} else {
		report.Summary = fmt.Sprintf("Found %d differences with %s compatibility level",
			len(comparison.Differences), report.Level.String())
	}
}

// generateComparisonSummary generates the comparison summary
func (sc *StoryComparator) generateComparisonSummary(comparison *StoryComparison) {
	summary := &comparison.Summary
	summary.TotalDifferences = len(comparison.Differences)

	// Count by impact level
	for _, diff := range comparison.Differences {
		switch diff.Impact {
		case Breaking:
			summary.BreakingChanges++
		case Warning:
			summary.WarningChanges++
		case Compatible:
			summary.CompatibleChanges++
		}
	}

	// Count step changes
	for _, diff := range comparison.Differences {
		switch diff.Type {
		case "STEP_ADDED":
			summary.StepsAdded++
		case "STEP_REMOVED":
			summary.StepsRemoved++
		}
	}

	// Set overall compatibility
	summary.OverallCompatibility = comparison.CompatibilityReport.Level

	// Count steps compared
	if comparison.NewSnapshot != nil {
		summary.StepsCompared = len(comparison.NewSnapshot.Steps)
	}
}
