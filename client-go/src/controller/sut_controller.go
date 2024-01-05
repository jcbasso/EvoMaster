package controller

import (
	"github.com/jcbasso/EvoMaster/client-go/src/controller/api/dto"
	"github.com/jcbasso/EvoMaster/client-go/src/controller/api/dto/problem"
	"github.com/jcbasso/EvoMaster/client-go/src/instrumentation/shared"
	"github.com/jcbasso/EvoMaster/client-go/src/instrumentation/staticstate"
	"github.com/jcbasso/EvoMaster/client-go/src/instrumentation/staticstate/execution_tracer"
	"github.com/jcbasso/EvoMaster/client-go/src/instrumentation/staticstate/objective_recorder"
	"strings"
)

type SutControllerInterface interface {
	SutHandler

	// GetInfoForAuthentication returns a list of valid authentication credentials, or nil if none is necessary
	GetInfoForAuthentication() []dto.AuthenticationDto

	// GetPreferredOutputFormat the format in which the test cases should be generated
	GetPreferredOutputFormat() dto.OutputFormat

	// IsSutRunning Check if the system under test (SUT) is running and fully initialized
	// returns true if the SUT is running
	IsSutRunning() bool

	// GetProblemInfo Depending on which kind of SUT we are dealing with (eg, REST, GraphQL or SPA frontend),
	// there is different info that must be provided
	// returns an instance of object with all the needed data for the specific addressed problem
	GetProblemInfo() problem.ProblemInfo
}

type SutController struct {
	SutControllerInterface
}

func NewSutController(sutControllerInterface SutControllerInterface) *SutController {
	return &SutController{
		SutControllerInterface: sutControllerInterface,
	}
}

// GetAdditionalInfoList returns additional info for each action in the test.
// The list is ordered based on the action index.
func (s *SutController) GetAdditionalInfoList() []*staticstate.AdditionalInfo {
	return execution_tracer.New().ExposeAdditionalInfoList()
}

func (s *SutController) GetUnitsInfoDto() dto.UnitsInfoDto {
	unitNames := []string{}
	for key, _ := range objective_recorder.New().GetUnitNames() {
		unitNames = append(unitNames, key)
	}

	return dto.UnitsInfoDto{
		UnitNames:        unitNames,
		NumberOfLines:    int(objective_recorder.New().GetNumberOfLines()),
		NumberOfBranches: int(objective_recorder.New().GetNumberOfBranches()),
	}
}

// IsInstrumentationActivated Check if instrumentation is on.
// returns true if the instrumentation is on
func (s *SutController) IsInstrumentationActivated() bool {
	return objective_recorder.New().GetNumberOfTargets() > 0
}

// NewSearch Re-initialize all internal data to enable a completely new search phase
// which should be independent of previous ones
func (s *SutController) NewSearch() {
	execution_tracer.New().Reset()
	objective_recorder.New().Reset(false)
}

// NewTest Re-initialize some internal data needed before running a new test
func (s *SutController) NewTest() {
	execution_tracer.New().Reset()

	/*
	   Note: it should be fine but, if for any reason EM did not do
	   a GET on the targets, then all those newly encountered targets
	   would be lost, as EM will have no way to ask for them later, unless
	   we explicitly say to return ALL targets
	*/
	objective_recorder.New().ClearFirstTimeEncountered()
}

func (s *SutController) GetTargetInfos(ids map[int]bool) ([]staticstate.TargetInfo, error) {

	list := []staticstate.TargetInfo{}

	objectives := execution_tracer.New().GetInternalReferenceToObjectiveCoverage()

	for id, _ := range ids {
		descriptiveId, err := objective_recorder.New().GetDescriptiveID(int64(id))
		if err != nil {
			return nil, err
		}

		info, ok := objectives[descriptiveId]
		if !ok {
			info = staticstate.NewNotReachedTargetInfo(int64(id))
		} else {
			info.MappedID = int64(id)
			info.DescriptiveID = ""
		}

		list = append(list, info)
	}

	// If new targets were found, we add them even if not requested by EM
	for _, s := range objective_recorder.New().GetTargetsSeenFirstTime() {
		info := objectives[s]
		info.MappedID = objective_recorder.New().GetMappedID(s)

		list = append(list, info)
	}

	return list, nil
}

// NewAction As some heuristics are based on which action (eg HTTP call, or click of button)
// in the test sequence is executed, and their order, we need to keep track of which
// action does cover what.
// @param dto is the DTO with the information about the action (eg its index in the test)
func (s *SutController) NewAction(dto dto.ActionDto) {
	execution_tracer.New().SetAction(
		staticstate.NewAction(
			int64(dto.Index),
			dto.InputVariables,
		),
	)
}

func (s *SutController) GetAllTargets() dto.AllTargetsDto {
	allTargets := objective_recorder.New().GetAllTargets()
	files := []string{}
	lines := []string{}
	branches := []string{}
	statements := []string{}

	for _, target := range allTargets {
		target := target
		if strings.HasPrefix(target, shared.FILE) {
			files = append(files, target)
		}
		if strings.HasPrefix(target, shared.LINE) {
			lines = append(lines, target)
		}
		if strings.HasPrefix(target, shared.BRANCH) {
			branches = append(branches, target)
		}
		if strings.HasPrefix(target, shared.STATEMENT) {
			statements = append(statements, target)
		}
	}

	return dto.AllTargetsDto{
		Files:      files,
		Lines:      lines,
		Branches:   branches,
		Statements: statements,
	}
}

func (s *SutController) GetFullObjectiveCoverage() dto.FullObjectiveCoverageDto {
	objectiveCoverages := execution_tracer.New().GetFullInternalReferenceToObjectiveCoverage()
	files := []dto.CoverageDto{}
	lines := []dto.CoverageDto{}
	branches := []dto.CoverageDto{}
	statements := []dto.CoverageDto{}

	for _, targetInfo := range objectiveCoverages {
		targetInfo := targetInfo
		if strings.HasPrefix(targetInfo.DescriptiveID, shared.FILE) {
			files = append(files, dto.CoverageDto{
				DescriptiveID: targetInfo.DescriptiveID,
				MappedID:      targetInfo.MappedID,
				Value:         targetInfo.Value,
				ActionIndex:   targetInfo.ActionIndex,
			})
		}
		if strings.HasPrefix(targetInfo.DescriptiveID, shared.LINE) {
			lines = append(lines, dto.CoverageDto{
				DescriptiveID: targetInfo.DescriptiveID,
				MappedID:      targetInfo.MappedID,
				Value:         targetInfo.Value,
				ActionIndex:   targetInfo.ActionIndex,
			})
		}
		if strings.HasPrefix(targetInfo.DescriptiveID, shared.BRANCH) {
			branches = append(branches, dto.CoverageDto{
				DescriptiveID: targetInfo.DescriptiveID,
				MappedID:      targetInfo.MappedID,
				Value:         targetInfo.Value,
				ActionIndex:   targetInfo.ActionIndex,
			})
		}
		if strings.HasPrefix(targetInfo.DescriptiveID, shared.STATEMENT) {
			statements = append(statements, dto.CoverageDto{
				DescriptiveID: targetInfo.DescriptiveID,
				MappedID:      targetInfo.MappedID,
				Value:         targetInfo.Value,
				ActionIndex:   targetInfo.ActionIndex,
			})
		}
	}

	return dto.FullObjectiveCoverageDto{
		Files:      files,
		Lines:      lines,
		Branches:   branches,
		Statements: statements,
	}
}
