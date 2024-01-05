package controller

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/jcbasso/EvoMaster/client-go/src/controller/api"
	"github.com/jcbasso/EvoMaster/client-go/src/controller/api/dto"
	"github.com/jcbasso/EvoMaster/client-go/src/controller/api/dto/problem"
	"net/http"
	"reflect"
	"strconv"
	"strings"
)

type EMController struct {
	sutController *SutController
	baseUrlOfSUT  string
	host          string
	port          int

	server *http.Server
}

func NewEMController(host string, port int, controller SutControllerInterface) *EMController {
	if host == "" {
		host = api.DEFAULT_CONTROLLER_HOST
	}

	if port == 0 {
		port = api.DEFAULT_CONTROLLER_PORT
	}

	sutController := NewSutController(controller)

	emController := EMController{
		sutController: sutController,
	}

	emController.createServer(host, port)

	return &emController
}

// StartTheControllerServer starts controller.
// StartTheControllerServer always returns a non-nil error. After Close,
// the returned error is ErrServerClosed.
func (e *EMController) StartTheControllerServer() error {
	return e.server.ListenAndServe()
}

// StopTheControllerServer stops SUT and EM controllers.
// StopTheControllerServer returns any error returned from closing the Server's
// underlying Listener(s).
func (e *EMController) StopTheControllerServer() error {
	if e.sutController.IsSutRunning() {
		e.sutController.StopSut()
	}

	return e.server.Close()
}

// Server creation
func (e *EMController) createServer(host string, port int) {
	r := mux.NewRouter().PathPrefix(api.BASE_PATH).Subrouter()

	r.HandleFunc(api.INFO_SUT, e.handleInfoSut).Methods(http.MethodGet)
	r.HandleFunc(api.CONTROLLER_INFO, e.handleControllerInfo).Methods(http.MethodGet)
	r.HandleFunc(api.NEW_SEARCH, e.handleNewSearch).Methods(http.MethodPost)
	r.HandleFunc(api.RUN_SUT, e.handleRunSut).Methods(http.MethodPut)
	r.HandleFunc(api.TEST_RESULTS, e.handleTestResults).Methods(http.MethodGet)
	r.HandleFunc(api.NEW_ACTION, e.handleNewAction).Methods(http.MethodPut)
	r.HandleFunc(api.ALL_TARGETS, e.handleAllTargets).Methods(http.MethodGet)
	r.HandleFunc(api.FULL_OBJECTIVE_COVERAGE, e.handleFullObjectiveCoverage).Methods(http.MethodGet)

	addr := fmt.Sprintf("%s:%d", host, port)

	e.server = &http.Server{
		Addr:    addr,
		Handler: r,
	}
}

// Server handlers
func (e *EMController) handleInfoSut(w http.ResponseWriter, r *http.Request) {
	baseUrlOfSUT := e.baseUrlOfSUT
	if baseUrlOfSUT != "" {
		baseUrlOfSUT = fmt.Sprintf("http://%s", e.baseUrlOfSUT)
	}
	sutInfo := dto.SutInfoDto{
		RestProblem: problem.RestProblemDto{},
		//GraphQLProblem:        problem.GraphQLProblemDto{},
		IsSutRunning:          e.sutController.IsSutRunning(),
		DefaultOutputFormat:   e.sutController.GetPreferredOutputFormat().String(),
		BaseUrlOfSUT:          baseUrlOfSUT,
		InfoForAuthentication: e.sutController.GetInfoForAuthentication(),
		UnitsInfoDto:          e.sutController.GetUnitsInfoDto(),
	}

	problemInfo := e.sutController.GetProblemInfo()
	if problemInfo == nil {
		e.respondError(w, "Undefined problem type in the EM Controller", http.StatusInternalServerError)
		return
	}

	switch x := problemInfo.(type) {
	case problem.RestProblemDto:
		sutInfo.RestProblem = x
		if sutInfo.RestProblem.OpenApiUrl != "" {
			sutInfo.RestProblem.OpenApiUrl = fmt.Sprintf("http://%s", sutInfo.RestProblem.OpenApiUrl)
		}
	//case problem.GraphQLProblemDto:
	//	sutInfo.GraphQLProblem = x
	default:
		e.respondError(w, "Unrecognized problem type", http.StatusInternalServerError)
		return
	}

	e.respondData(w, sutInfo, http.StatusOK)
}

func (e *EMController) handleControllerInfo(w http.ResponseWriter, r *http.Request) {
	controllerInfo := dto.ControllerInfoDto{
		// TODO will need something to identify the file to import
		FullName:            getType(e.sutController.SutControllerInterface),
		IsInstrumentationOn: e.sutController.IsInstrumentationActivated(),
	}

	e.respondData(w, controllerInfo, http.StatusOK)
}

// getType returns the type name of an interface
func getType(v interface{}) (res string) {
	t := reflect.TypeOf(v)
	for t.Kind() == reflect.Ptr {
		t = t.Elem()
		res += "*"
	}
	return res + t.Name()
}

func (e *EMController) handleNewSearch(w http.ResponseWriter, r *http.Request) {
	e.sutController.NewSearch()
	e.respondData(w, nil, http.StatusCreated)
}

func (e *EMController) handleRunSut(w http.ResponseWriter, r *http.Request) {
	sutRun := &dto.SutRunDto{}
	err := json.NewDecoder(r.Body).Decode(sutRun)
	if err != nil {
		e.respondError(w, "No provided JSON payload", http.StatusBadRequest)
		return
	}

	if !sutRun.Run {
		if sutRun.ResetState {
			e.respondError(w, "Invalid JSON: cannot reset state and stop service at same time", http.StatusBadRequest)
			return
		}

		// if on, we want to shut down the server
		if e.sutController.IsSutRunning() {
			e.sutController.StopSut()
			e.baseUrlOfSUT = ""
		}

	} else {
		// If SUT is not up and running, let's start it

		if !e.sutController.IsSutRunning() {
			e.baseUrlOfSUT = e.sutController.StartSut()
			if e.baseUrlOfSUT == "" {
				// there has been an internal failure in starting the SUT
				e.respondError(w, "Internal failure: cannot start SUT based on given configuration", http.StatusInternalServerError)
				return
			}
		} else {
			// TODO as starting should be blocking, need to check if initialized, and wait if not
		}

		//  regardless of where it was running or not, need to reset state.
		//  this is controlled by a boolean, although most likely we ll always
		//  want to do it
		if sutRun.ResetState {
			e.sutController.ResetStateOfSUT()
			e.sutController.NewTest()
		}

		// Note: here even if we start the SUT, the starting of a "New Search"
		// cannot be done here, as in this endpoint we also deal with the reset
		// of state. When we reset state for a new test run, we do not want to
		// reset all the other data regarding the whole search
	}

	e.respondData(w, nil, http.StatusNoContent)
	return
}

func (e *EMController) handleTestResults(w http.ResponseWriter, r *http.Request) {
	idsParam := r.URL.Query().Get("ids")
	ids := map[int]bool{}
	for _, id := range strings.Split(idsParam, ",") {
		if id != "" {
			intId, err := strconv.Atoi(id)
			if err != nil {
				e.respondError(w, "Ids must be integers", http.StatusBadRequest)
				return
			}
			ids[intId] = true
		}
	}

	targetInfos, err := e.sutController.GetTargetInfos(ids)
	if err != nil {
		e.respondError(w, "Internal failure: cannot get test results from SUT controller", http.StatusInternalServerError)
		return
	}

	targetInfoDtos := []dto.TargetInfoDto{}
	for _, targetInfo := range targetInfos {
		targetInfoDto := dto.TargetInfoDto{
			Id:            int(targetInfo.MappedID),
			Value:         targetInfo.Value,
			DescriptiveId: targetInfo.DescriptiveID,
			ActionIndex:   int(targetInfo.ActionIndex),
		}
		targetInfoDtos = append(targetInfoDtos, targetInfoDto)
	}

	additionalInfos := e.sutController.GetAdditionalInfoList()

	additionalInfoDtos := []dto.AdditionalInfoDto{}
	for _, additionalInfo := range additionalInfos {
		queryParameters := []string{}
		for queryParameter, _ := range additionalInfo.QueryParameters {
			queryParameters = append(queryParameters, queryParameter)
		}

		headers := []string{}
		for header, _ := range additionalInfo.Headers {
			headers = append(queryParameters, header)
		}

		additionalInfoDto := dto.AdditionalInfoDto{
			QueryParameters: queryParameters,
			Headers:         headers,
			//StringSpecializations: ,
			LastExecutedStatement: additionalInfo.GetLastExecutedStatement(),
		}

		additionalInfoDtos = append(additionalInfoDtos, additionalInfoDto)
	}

	testResultsDto := dto.TestResultsDto{
		Targets:            targetInfoDtos,
		AdditionalInfoList: additionalInfoDtos,
		ExtraHeuristics:    []dto.ExtraHeuristicsDto{},
	}

	e.respondData(w, testResultsDto, http.StatusOK)
}

func (e *EMController) handleNewAction(w http.ResponseWriter, r *http.Request) {
	req := &dto.ActionDto{}
	err := json.NewDecoder(r.Body).Decode(req)
	if err != nil {
		e.respondError(w, err.Error(), http.StatusBadRequest)
		return
	}

	e.sutController.NewAction(*req)
	e.respondData(w, nil, http.StatusNoContent)
}

func (e *EMController) handleAllTargets(w http.ResponseWriter, r *http.Request) {
	allTargetsDto := e.sutController.GetAllTargets()
	e.respondData(w, allTargetsDto, http.StatusOK)
}

func (e *EMController) handleFullObjectiveCoverage(w http.ResponseWriter, r *http.Request) {
	fullObjectiveCoverage := e.sutController.GetFullObjectiveCoverage()
	e.respondData(w, fullObjectiveCoverage, http.StatusOK)
}

func (e *EMController) respondData(w http.ResponseWriter, response interface{}, statusCode int) {
	wrapper := dto.NewWrappedResponseDtoWithNoData()
	if response != nil {
		wrapper = dto.NewWrappedResponseDtoWithData(response)
	}

	e.respond(w, wrapper, statusCode)

	return
}

func (e *EMController) respondError(w http.ResponseWriter, msg string, statusCode int) {
	wrapper := dto.NewWrappedResponseDtoWithError(msg)

	e.respond(w, wrapper, statusCode)

	return
}

// respond Writes DTO response with sent Status Code. If there is an error in the encoding it returns an 500 with the error.
func (e *EMController) respond(w http.ResponseWriter, wrapper dto.WrappedResponseDto[interface{}], statusCode int) {
	w.Header().Add("Content-Type", "application/json")
	if statusCode == http.StatusNoContent {
		w.WriteHeader(statusCode)
		return
	}

	err := e.testEncoding(wrapper)
	if err != nil {
		fmt.Printf("error encoding response: %v", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, err.Error())
		return
	}

	w.WriteHeader(statusCode)
	err = json.NewEncoder(w).Encode(wrapper)
	if err != nil {
		fmt.Printf("error final encoding response: %v", err.Error())
		fmt.Println(err.Error())
		fmt.Fprintln(w, err.Error())
		return
	}
}

func (e *EMController) testEncoding(v any) error {
	_, err := json.Marshal(v)
	return err
}
