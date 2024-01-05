package instrumentation

import (
	"fmt"
	"github.com/jcbasso/EvoMaster/client-go/src/instrumentation/heuristic"
	"github.com/jcbasso/EvoMaster/client-go/src/instrumentation/staticstate/execution_tracer"
	"github.com/jcbasso/EvoMaster/client-go/src/instrumentation/staticstate/objective_recorder"
	"reflect"
)

func RegisterTargets(ids []string) {
	for _, id := range ids {
		objective_recorder.New().RegisterTarget(id)
	}
}

var tracer = execution_tracer.New()
var heuristicForBooleans = heuristic.NewHeuristicForBooleans()

func EnteringStatement(fileName string, line int, statement int) {
	tracer.EnteringStatement(fileName, line, statement)
	//log.Printf("entering statement: %s-%v-%v", fileName, line, statement)
}

func CompletedStatement(fileName string, line int, statement int) {
	tracer.CompletedStatement(fileName, line, statement)
	//log.Printf("completed statement: %s-%v-%v", fileName, line, statement)
}

// CompletionStatement used for statements like:
//   - return (with no data)
//   - continue
//   - break
//   - TODO: For now it is called always
func CompletionStatement(fileName string, line int, statement int) {
	// log.Printf("[WARN] CompletionStatement:EnteringStatement [file: %s, line: %d, statement: %d]\n", fileName, line, statement)
	EnteringStatement(fileName, line, statement)
	// log.Printf("[WARN] CompletionStatement:CompletedStatement [file: %s, line: %d, statement: %d]\n", fileName, line, statement)
	CompletedStatement(fileName, line, statement)
}

// TODO: Delete if not used? Also as how it works, it is not clear if it makes sense.
func CompletingStatement(fileName string, line int, statement int, values ...any) any {
	//log.Printf("[WARN] CompletingStatement:CompletedStatement [file: %s, line: %d, statement: %d]\n", fileName, line, statement)
	CompletedStatement(fileName, line, statement)
	return values
}

func Not(value bool) bool {
	return heuristicForBooleans.HandleNot(value)
}

func And(left func() bool, right func() bool, fileName string, line int, branchId int) bool {
	return heuristicForBooleans.EvaluateAnd(left, right, fileName, line, branchId, tracer)
}

func Or(left func() bool, right func() bool, fileName string, line int, branchId int) bool {
	return heuristicForBooleans.EvaluateOr(left, right, fileName, line, branchId, tracer)
}

func CmpUnordered(left any, op string, right any, fileName string, line int, branchId int) bool {
	lvalue := reflect.ValueOf(left)
	rvalue := reflect.ValueOf(right)

	if heuristic.BothInt(lvalue, rvalue) {
		return heuristicForBooleans.EvaluateUnorderedCmp(lvalue.Int(), op, rvalue.Int(), fileName, line, branchId, tracer)
	}

	if heuristic.BothUint(lvalue, rvalue) {
		return heuristicForBooleans.EvaluateUnorderedCmp(lvalue.Uint(), op, rvalue.Uint(), fileName, line, branchId, tracer)
	}

	if heuristic.BothFloat(lvalue, rvalue) {
		return heuristicForBooleans.EvaluateUnorderedCmp(lvalue.Float(), op, rvalue.Float(), fileName, line, branchId, tracer)
	}

	return heuristicForBooleans.EvaluateUnorderedCmp(left, op, right, fileName, line, branchId, tracer)
}

// TODO: Could combine the CMPs with this validation but not sure if it makes sense

func CmpOrdered(left any, op string, right any, fileName string, line int, branchId int) bool {
	lvalue := reflect.ValueOf(left)
	rvalue := reflect.ValueOf(right)

	if heuristic.BothInt(lvalue, rvalue) {
		return heuristic.EvaluateOrderedCmp[int64](heuristicForBooleans, lvalue.Int(), op, rvalue.Int(), fileName, line, branchId, tracer)
	}

	if heuristic.BothUint(lvalue, rvalue) {
		return heuristic.EvaluateOrderedCmp[uint64](heuristicForBooleans, lvalue.Uint(), op, rvalue.Uint(), fileName, line, branchId, tracer)
	}

	if heuristic.BothFloat(lvalue, rvalue) {
		return heuristic.EvaluateOrderedCmp[float64](heuristicForBooleans, lvalue.Float(), op, rvalue.Float(), fileName, line, branchId, tracer)
	}

	if heuristic.BothString(lvalue, rvalue) {
		return heuristic.EvaluateOrderedCmp[string](heuristicForBooleans, lvalue.String(), op, rvalue.String(), fileName, line, branchId, tracer)
	}

	panic(fmt.Sprintf("Invalid types not ordered [left:%s, right:%s]", left, right))
}
