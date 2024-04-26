package staticstate

import (
	"github.com/jcbasso/EvoMaster/client-go/src/instrumentation/shared/string_specialization"
	"github.com/jcbasso/EvoMaster/client-go/src/instrumentation/shared/taint_type"
	"sync"
)

// TODO: Review this struct

type AdditionalInfo struct {
	QueryParameters map[string]bool
	Headers         map[string]bool
	/**
	 * Map from taint input name to string specializations for it
	 */
	StringSpecializations      map[string]map[string_specialization.StringSpecializationInfo]bool
	LastExecutedStatementStack []string
	NoExceptionStatement       string
	mx                         sync.RWMutex
}

func NewAdditionalInfo() *AdditionalInfo {
	return &AdditionalInfo{
		QueryParameters:            map[string]bool{},
		Headers:                    map[string]bool{},
		StringSpecializations:      map[string]map[string_specialization.StringSpecializationInfo]bool{},
		LastExecutedStatementStack: []string{},
	}
}

func (a *AdditionalInfo) PushLastExecutedStatement(lastLine string) {
	a.mx.Lock()
	// log.Println("[WARN] Pushing LastExecutedStatementStack")
	// log.Println(fmt.Sprintf("[WARN] push: %v", a.LastExecutedStatementStack))
	a.NoExceptionStatement = ""
	a.LastExecutedStatementStack = append(a.LastExecutedStatementStack, lastLine)
	a.mx.Unlock()
}

func (a *AdditionalInfo) PopLastExecutedStatement() (string, bool) {
	a.mx.Lock()
	defer a.mx.Unlock()
	// log.Println("[WARN] Popping LastExecutedStatementStack")
	// log.Println(fmt.Sprintf("[WARN] pop: %v", a.LastExecutedStatementStack))
	if len(a.LastExecutedStatementStack) <= 0 {
		// panic("Popping on empty LastExecutedStatementStack")
		return "", false
	}
	statement := a.LastExecutedStatementStack[len(a.LastExecutedStatementStack)-1]
	a.LastExecutedStatementStack = a.LastExecutedStatementStack[:len(a.LastExecutedStatementStack)-1]
	if len(a.LastExecutedStatementStack) == 0 {
		a.NoExceptionStatement = statement
	}
	return statement, true
}

func (a *AdditionalInfo) GetLastExecutedStatement() string {
	a.mx.RLock()
	defer a.mx.RUnlock()
	if len(a.LastExecutedStatementStack) != 0 {
		return a.LastExecutedStatementStack[len(a.LastExecutedStatementStack)-1]
	}

	return a.NoExceptionStatement
}

func (a *AdditionalInfo) AddSpecialization(taintInputName string, info string_specialization.StringSpecializationInfo) {
	a.mx.RLock()
	defer a.mx.RUnlock()

	if NewExecutionTracer().GetTaintType(taintInputName) == taint_type.NONE {
		panic("Invalid taint input name: " + taintInputName)
	}

	set, ok := a.StringSpecializations[taintInputName]
	if !ok {
		set = map[string_specialization.StringSpecializationInfo]bool{}
		a.StringSpecializations[taintInputName] = set
	}

	set[info] = true
}
