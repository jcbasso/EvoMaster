package staticstate

import (
	"fmt"
	"github.com/jcbasso/EvoMaster/client-go/src/instrumentation/shared"
	"strings"
	"sync"
	"sync/atomic"
)

type ObjectiveRecorder struct {
	maxObjectiveCoverage   map[int64]float64
	maxObjectiveCoverageMx sync.RWMutex
	idMapping              map[string]int64
	idMappingMx            sync.RWMutex
	reversedIdMapping      map[int64]string
	reversedIdMappingMx    sync.RWMutex
	idMappingCounter       atomic.Uint64
	firstTimeEncountered   []string
	firstTimeEncounteredMx sync.RWMutex
	allTargetsSet          map[string]bool
	allTargetsSetMx        sync.RWMutex
	unitsInfo              *UnitsInfo
	unitsInfoMx            sync.RWMutex
}

var objectiveRecorderOnce sync.Once
var objectiveRecorderInstance *ObjectiveRecorder

func NewObjectiveRecorder() *ObjectiveRecorder {
	objectiveRecorderOnce.Do(func() {
		objectiveRecorderInstance = &ObjectiveRecorder{}
		objectiveRecorderInstance.Reset(true)
	})

	return objectiveRecorderInstance
}

func (o *ObjectiveRecorder) Reset(atLoadTime bool) {
	o.maxObjectiveCoverageMx.Lock()
	o.maxObjectiveCoverage = map[int64]float64{}
	o.maxObjectiveCoverageMx.Unlock()
	o.idMappingMx.Lock()
	o.idMapping = map[string]int64{}
	o.idMappingMx.Unlock()
	o.reversedIdMappingMx.Lock()
	o.reversedIdMapping = map[int64]string{}
	o.reversedIdMappingMx.Unlock()
	o.firstTimeEncounteredMx.Lock()
	o.firstTimeEncountered = []string{}
	o.firstTimeEncounteredMx.Unlock()
	if atLoadTime {
		o.allTargetsSetMx.Lock()
		o.allTargetsSet = map[string]bool{}
		o.allTargetsSetMx.Unlock()
		o.unitsInfoMx.Lock()
		o.unitsInfo = NewUnitsInfo()
		o.unitsInfoMx.Unlock()
	}
}

func (o *ObjectiveRecorder) GetNumberOfTargets() int {
	o.allTargetsSetMx.Lock()
	defer o.allTargetsSetMx.Unlock()
	return len(o.allTargetsSet)
}

func (o *ObjectiveRecorder) RegisterTarget(target string) {
	o.allTargetsSetMx.Lock()
	o.allTargetsSet[target] = true
	o.allTargetsSetMx.Unlock()
	o.unitsInfoMx.Lock()
	if strings.HasPrefix(target, shared.FILE) {
		o.unitsInfo.UnitNamesSet[shared.GetFileIdFromObjectiveName(target)] = true
	}
	if strings.HasPrefix(target, shared.LINE) {
		o.unitsInfo.LinesCount++
	}
	if strings.HasPrefix(target, shared.BRANCH) {
		o.unitsInfo.BranchCount++
	}
	if strings.HasPrefix(target, shared.STATEMENT) {
		o.unitsInfo.StatementCount++
	}
	o.unitsInfoMx.Unlock()
}

func (o *ObjectiveRecorder) GetTargetsSeenFirstTime() []string {
	o.firstTimeEncounteredMx.RLock()
	defer o.firstTimeEncounteredMx.RUnlock()
	return o.firstTimeEncountered
}

func (o *ObjectiveRecorder) ClearFirstTimeEncountered() {
	o.firstTimeEncounteredMx.Lock()
	o.firstTimeEncountered = []string{}
	o.firstTimeEncounteredMx.Unlock()
}

func (o *ObjectiveRecorder) UpdateTarget(descriptiveID string, value float64) {
	if value < 0 || value > 1 {
		panic(fmt.Sprintf("invalid Value %.2f, out of range [0,1]", value))
	}

	mappedID := o.GetMappedID(descriptiveID)
	o.maxObjectiveCoverageMx.Lock()
	if _, ok := o.maxObjectiveCoverage[mappedID]; !ok {
		o.firstTimeEncounteredMx.Lock()
		o.firstTimeEncountered = append(o.firstTimeEncountered, descriptiveID)
		o.firstTimeEncounteredMx.Unlock()
		o.maxObjectiveCoverage[mappedID] = value
	}

	oldValueAny, ok := o.maxObjectiveCoverage[mappedID]
	var oldValue float64 = 0
	if !ok {
		oldValue = oldValueAny
	}

	if value > oldValue {
		o.maxObjectiveCoverage[mappedID] = value
	}
	o.maxObjectiveCoverageMx.Unlock()
}

func (o *ObjectiveRecorder) GetMappedID(descriptiveID string) (mappedID int64) {
	mappedIDCounter := int64(o.idMappingCounter.Load())
	o.idMappingMx.RLock()
	val, ok := o.idMapping[descriptiveID]
	o.idMappingMx.RUnlock()
	if !ok {
		// New
		o.idMappingMx.Lock()
		o.idMapping[descriptiveID] = mappedIDCounter
		o.idMappingMx.Unlock()
		o.reversedIdMappingMx.Lock()
		o.reversedIdMapping[mappedIDCounter] = descriptiveID
		o.reversedIdMappingMx.Unlock()
		o.idMappingCounter.Add(1)
		mappedID = mappedIDCounter
	} else {
		// Existed
		mappedID = val
	}

	return
}

func (o *ObjectiveRecorder) GetDescriptiveID(mappedID int64) (string, error) {
	o.reversedIdMappingMx.RLock()
	value, ok := o.reversedIdMapping[mappedID]
	o.reversedIdMappingMx.RUnlock()
	if !ok {
		err := fmt.Errorf("id %v is not mapped", mappedID)
		return "", err
	}
	return value, nil
}

// TODO: Separate units to UnitsInfoRecorder

func (o *ObjectiveRecorder) GetUnitNames() map[string]bool {
	o.unitsInfoMx.RLock()
	defer o.unitsInfoMx.RUnlock()
	return o.unitsInfo.UnitNamesSet
}

func (o *ObjectiveRecorder) GetNumberOfLines() int64 {
	o.unitsInfoMx.RLock()
	defer o.unitsInfoMx.RUnlock()
	return o.unitsInfo.LinesCount
}

func (o *ObjectiveRecorder) GetNumberOfBranches() int64 {
	o.unitsInfoMx.RLock()
	defer o.unitsInfoMx.RUnlock()
	return o.unitsInfo.BranchCount
}

func (o *ObjectiveRecorder) GetAllTargets() []string {
	o.allTargetsSetMx.RLock()

	res := make([]string, len(o.allTargetsSet))
	i := 0
	for k, _ := range o.allTargetsSet {
		res[i] = k
		i++
	}

	o.allTargetsSetMx.RUnlock()
	return res
}
