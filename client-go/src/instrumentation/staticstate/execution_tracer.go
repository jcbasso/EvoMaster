package staticstate

import (
	"fmt"
	"github.com/jcbasso/EvoMaster/client-go/src/instrumentation/heuristic"
	"github.com/jcbasso/EvoMaster/client-go/src/instrumentation/shared"
	"github.com/jcbasso/EvoMaster/client-go/src/instrumentation/shared/string_specialization"
	"github.com/jcbasso/EvoMaster/client-go/src/instrumentation/shared/taint_type"
	"strings"
	"sync"
)

type ExecutionTracer struct {
	objectiveCoverage    map[string]TargetInfo // key: Descriptive ID
	objectiveCoverageMx  sync.RWMutex
	actionIndex          int64
	actionIndexMx        sync.RWMutex
	inputVariablesSet    map[string]bool
	inputVariablesSetMx  sync.RWMutex
	additionalInfoList   []*AdditionalInfo
	additionalInfoListMx sync.RWMutex

	fullObjectiveCoverage   map[string]TargetInfo // key: Descriptive ID
	fullObjectiveCoverageMx sync.RWMutex
}

var executionTracerOnce sync.Once
var executionTracerInstance *ExecutionTracer

func NewExecutionTracer() *ExecutionTracer {
	executionTracerOnce.Do(func() {
		executionTracerInstance = &ExecutionTracer{}
		executionTracerInstance.Reset()
		executionTracerInstance.actionIndexMx.Lock()
		executionTracerInstance.actionIndex = 0
		executionTracerInstance.actionIndexMx.Unlock()
	})

	return executionTracerInstance
}

func (o *ExecutionTracer) Reset() {
	o.objectiveCoverageMx.Lock()
	o.objectiveCoverage = map[string]TargetInfo{}
	o.objectiveCoverageMx.Unlock()
	o.inputVariablesSetMx.Lock()
	o.inputVariablesSet = map[string]bool{}
	o.inputVariablesSetMx.Unlock()
	o.additionalInfoListMx.Lock()
	//log.Println("[WARN] Resetting action info list")
	o.additionalInfoList = []*AdditionalInfo{NewAdditionalInfo()}
	o.additionalInfoListMx.Unlock()

	o.fullObjectiveCoverageMx.Lock()
	o.fullObjectiveCoverage = map[string]TargetInfo{}
	o.fullObjectiveCoverageMx.Unlock()
}

func (o *ExecutionTracer) SetAction(action *Action) {
	o.actionIndexMx.Lock()
	if o.actionIndex != action.Index {
		o.actionIndex = action.Index
		o.additionalInfoListMx.Lock()
		//log.Println(fmt.Sprintf("[WARN] Setting action. [action: %d]", action.Index))
		o.additionalInfoList = append(o.additionalInfoList, &AdditionalInfo{})
		o.additionalInfoListMx.Unlock()
	}
	if action.InputVariables != nil && len(action.InputVariables) > 0 {
		// TODO: Validate if it is needed to be an array in one place and a set in another one
		o.inputVariablesSetMx.Lock()
		o.inputVariablesSet = map[string]bool{}
		for _, variable := range action.InputVariables {
			o.inputVariablesSet[variable] = true
		}
		o.inputVariablesSetMx.Unlock()
	}
	o.actionIndexMx.Unlock()
}

func (o *ExecutionTracer) MarkLastExecutedStatement(lastLine string) {
	o.actionIndexMx.RLock()
	o.additionalInfoListMx.Lock()
	//log.Println(fmt.Sprintf("[WARN] MarkLastExecutedStatement [action: %d, additionalInfoList: %v, lastLine: %s]", o.actionIndex, o.additionalInfoList, lastLine))
	//for i, info := range o.additionalInfoList {
	//	log.Println(fmt.Sprintf("[WARN] additionalInfoList[%d]: %v]", i, info.LastExecutedStatementStack))
	//}
	o.additionalInfoList[o.actionIndex].PushLastExecutedStatement(lastLine)
	o.additionalInfoListMx.Unlock()
	o.actionIndexMx.RUnlock()
}

func (o *ExecutionTracer) CompletedLastExecutedStatement(lastLine string) {
	o.actionIndexMx.RLock()
	o.additionalInfoListMx.Lock()
	//log.Println(fmt.Sprintf("[WARN] CompletedLastExecutedStatement [action: %d, additionalInfoList: %v, lastLine: %s]", o.actionIndex, o.additionalInfoList, lastLine))
	//for i, info := range o.additionalInfoList {
	//	log.Println(fmt.Sprintf("[WARN] additionalInfoList[%d]: %v]", i, info.LastExecutedStatementStack))
	//}
	stmt, ok := o.additionalInfoList[o.actionIndex].PopLastExecutedStatement()
	o.additionalInfoListMx.Unlock()
	o.actionIndexMx.RUnlock()
	if ok && stmt != lastLine {
		/*
		   actually we cannot have such check. We might end in such situation:

		   X calls F in non-instrumented framework, which then call Y (both X and Y being of SUT).
		   If Y crashes with a catch in F, then X will wrongly pop for Y.

		   TODO could have such check with a parameter, to have only in the tests
		*/
		//return Error(`Expected to pop ${lastLine} instead of ${stmt}`);
	}
}

func (o *ExecutionTracer) GetInternalReferenceToObjectiveCoverage() map[string]TargetInfo {
	o.objectiveCoverageMx.RLock()
	res := map[string]TargetInfo{}
	for key, value := range o.objectiveCoverage {
		res[key] = TargetInfo{
			MappedID:      value.MappedID,
			DescriptiveID: value.DescriptiveID,
			Value:         value.Value,
			ActionIndex:   value.ActionIndex,
		}
	}
	o.objectiveCoverageMx.RUnlock()
	return res
}

func (o *ExecutionTracer) GetFullInternalReferenceToObjectiveCoverage() map[string]TargetInfo {
	o.fullObjectiveCoverageMx.RLock()
	res := map[string]TargetInfo{}
	for key, value := range o.fullObjectiveCoverage {
		res[key] = TargetInfo{
			MappedID:      value.MappedID,
			DescriptiveID: value.DescriptiveID,
			Value:         value.Value,
			ActionIndex:   value.ActionIndex,
		}
	}
	o.fullObjectiveCoverageMx.RUnlock()
	return res
}

func (o *ExecutionTracer) GetNumberOfObjectives(prefix string) int64 {
	counter := int64(0)
	o.objectiveCoverageMx.RLock()
	for key, _ := range o.objectiveCoverage {
		if strings.HasPrefix(key, prefix) {
			counter++
		}
	}
	o.objectiveCoverageMx.RUnlock()

	return counter
}

func (o *ExecutionTracer) GetNumberOfNonCoveredObjectives(prefix string) int {
	return len(o.GetNonCoveredObjectives(prefix))
}

func (o *ExecutionTracer) GetNonCoveredObjectives(prefix string) map[string]bool {
	res := map[string]bool{}
	for descriptiveID, targetInfo := range o.GetInternalReferenceToObjectiveCoverage() {
		if strings.HasPrefix(descriptiveID, prefix) && targetInfo.Value < 1 {
			res[descriptiveID] = true
		}
	}

	return res
}

func (o *ExecutionTracer) GetValue(descriptiveID string) float64 {
	o.objectiveCoverageMx.RLock()
	val, ok := o.objectiveCoverage[descriptiveID]
	if !ok {
		fmt.Printf("Value %s not found on objctiveCoverage\n", descriptiveID)
		return 0
	}
	o.objectiveCoverageMx.RUnlock()
	return val.Value
}

func (o *ExecutionTracer) UpdateObjective(descriptiveID string, value float64) {
	if value < 0 || value > 1 {
		panic(fmt.Sprintf("invalid Value %.2f, out of range [0,1]", value))
	}
	targetInfo := TargetInfo{
		DescriptiveID: descriptiveID,
		Value:         value,
		ActionIndex:   o.actionIndex,
	}
	o.objectiveCoverageMx.Lock()
	o.fullObjectiveCoverageMx.Lock()
	previous, ok := o.objectiveCoverage[descriptiveID]
	if ok && value > previous.Value {
		o.objectiveCoverage[descriptiveID] = targetInfo
		o.fullObjectiveCoverage[descriptiveID] = targetInfo
	}
	if !ok {
		o.objectiveCoverage[descriptiveID] = targetInfo
		o.fullObjectiveCoverage[descriptiveID] = targetInfo
	}
	o.fullObjectiveCoverageMx.Unlock()
	o.objectiveCoverageMx.Unlock()
	NewObjectiveRecorder().UpdateTarget(descriptiveID, value)
}

func (o *ExecutionTracer) EnteringStatement(fileName string, line int, statement int) {
	fileID := shared.FileObjectiveName(fileName)
	lineID := shared.LineObjectiveName(fileName, line)
	statementID := shared.StatementObjectiveName(fileName, line, statement)
	o.UpdateObjective(fileID, 1)
	o.UpdateObjective(lineID, 1)
	o.UpdateObjective(statementID, 0.5)
	o.MarkLastExecutedStatement(fmt.Sprintf("%s_%v_%v", fileName, line, statement))
}

func (o *ExecutionTracer) CompletedStatement(fileName string, line int, statement int) {
	statementID := shared.StatementObjectiveName(fileName, line, statement)
	o.UpdateObjective(statementID, 1)
	o.CompletedLastExecutedStatement(fmt.Sprintf("%s_%v_%v", fileName, line, statement))
	heuristic.NewHeuristicForBooleans().ClearLastEvaluation()
}

func (o *ExecutionTracer) UpdateBranch(fileName string, line int, branch int, truthness *heuristic.Truthness) {
	thenBranch := shared.BranchObjectiveName(fileName, line, branch, true)
	elseBranch := shared.BranchObjectiveName(fileName, line, branch, false)
	o.UpdateObjective(thenBranch, truthness.OfTrue)
	o.UpdateObjective(elseBranch, truthness.OfFalse)
}

func (o *ExecutionTracer) ExposeAdditionalInfoList() []*AdditionalInfo {
	o.additionalInfoListMx.RLock()
	defer o.additionalInfoListMx.RUnlock()
	return append([]*AdditionalInfo(nil), o.additionalInfoList...)
}

// IsTaintInput Check if the given input represented a tainted value from the test cases.
// This could be based on static info of the input (eg, according to a precise
// name convention given by TaintInputName), or dynamic info given directly by
// the test itself (eg, the test at action can register a list of values to check
// for)
func (o *ExecutionTracer) IsTaintInput(input string) bool {
	return taint_type.IsTaintInput(input) || o.inputVariablesSet[input]
}

func (o *ExecutionTracer) HandleTaintForStringEquals(left string, right string, ignoreCase bool) {
	taintedLeft := o.IsTaintInput(left)
	taintedRight := o.IsTaintInput(right)

	taintType := taint_type.FULL_MATCH

	if taintedLeft && taintedRight {
		if left == right ||
			ignoreCase && (strings.ToLower(left) == strings.ToLower(right)) {
			//tainted, but compared to itself. so shouldn't matter
			return
		}

		/*
		   We consider binding only for base versions of taint, ie we ignore
		   the special strings provided by the Core, as it would lead to nasty
		   side-effects
		*/
		if !taint_type.IsTaintInput(left) || !taint_type.IsTaintInput(right) {
			return
		}

		//TODO could have EQUAL_IGNORE_CASE
		id := left + "___" + right
		o.AddStringSpecialization(left, string_specialization.NewStringSpecializationInfo(string_specialization.EQUAL, id, taintType))
		o.AddStringSpecialization(right, string_specialization.NewStringSpecializationInfo(string_specialization.EQUAL, id, taintType))
		return
	}

	specialization := string_specialization.CONSTANT
	if ignoreCase {
		specialization = string_specialization.CONSTANT_IGNORE_CASE
	}

	if taintedLeft || taintedRight {
		if taintedLeft {
			o.AddStringSpecialization(left, string_specialization.NewStringSpecializationInfo(specialization, right, taintType))
		} else {
			o.AddStringSpecialization(right, string_specialization.NewStringSpecializationInfo(specialization, left, taintType))
		}
	}
}

func (o *ExecutionTracer) AddStringSpecialization(taintInputName string, info string_specialization.StringSpecializationInfo) {
	o.additionalInfoListMx.Lock()
	defer o.additionalInfoListMx.Unlock()
	o.additionalInfoList[o.actionIndex].AddSpecialization(taintInputName, info)
}

func (o *ExecutionTracer) GetTaintType(input string) taint_type.TaintType {
	if input == "" {
		return taint_type.NONE
	}

	if o.IsTaintInput(input) {
		return taint_type.FULL_MATCH
	}

	if taint_type.IncludesTaintInput(input) {
		return taint_type.PARTIAL_MATCH
	}
	for s := range o.inputVariablesSet {
		if strings.Contains(input, s) {
			return taint_type.PARTIAL_MATCH
		}
	}

	return taint_type.NONE
}
