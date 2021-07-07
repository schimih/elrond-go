package core

import (
	"bytes"
	"regexp"
	"runtime"
	"runtime/pprof"
	"sort"
	"strings"
	"time"
)

const (
	newRoutine    = "new"
	oldRoutine    = "old"
	closedRoutine = "closed"
)

func NewGoroutinesDumper() *GoRoutinesDumper {
	latestData := make(map[string]map[string]*GoRoutineData)
	latestData[newRoutine] = make(map[string]*GoRoutineData)
	latestData[oldRoutine] = make(map[string]*GoRoutineData)
	latestData[closedRoutine] = make(map[string]*GoRoutineData)
	return &GoRoutinesDumper{
		latestData,
	}
}

type GoRoutinesDumper struct {
	LatestData map[string]map[string]*GoRoutineData
}

type GoRoutineData struct {
	ID              string
	FirstOccurrence time.Time
	AllData         string
}

// DumpGoRoutinesToLog will print the currently running go routines in the log
func DumpGoRoutinesToLog(goRoutinesNumberStart int) {
	buffer := new(bytes.Buffer)
	err := pprof.Lookup("goroutine").WriteTo(buffer, 2)
	if err != nil {
		log.Error("could not dump goroutines", "error", err)
	}
	log.Debug("go routines number",
		"start", goRoutinesNumberStart,
		"end", runtime.NumGoroutine())

	log.Debug(buffer.String())
}

// DumpGoRoutinesToLogWithTypes will print the currently running go routines in the log
func (grd *GoRoutinesDumper) DumpGoRoutinesToLogWithTypes() {
	buffer := new(bytes.Buffer)
	err := pprof.Lookup("goroutine").WriteTo(buffer, 2)
	if err != nil {
		log.Error("could not dump goroutines", "error", err)
	}
	log.Debug("DumpGoRoutinesToLogWithTypes", "end", runtime.NumGoroutine())

	allGoRoutinesString := buffer.String()
	splits := strings.Split(allGoRoutinesString, "\n\n")

	newOldGoRoutine := make(map[string]*GoRoutineData)

	currentTime := time.Now()
	for _, st := range splits {
		gId := getGoroutineId(st)
		val, ok := grd.LatestData[oldRoutine][gId]
		if !ok {
			grd.LatestData[newRoutine][gId] = &GoRoutineData{
				ID:              gId,
				FirstOccurrence: currentTime,
				AllData:         st,
			}
			continue
		}
		newOldGoRoutine[val.ID] = val
	}

	grd.LatestData[oldRoutine] = newOldGoRoutine

	goRoutines := make([]*GoRoutineData, 0)

	for _, val := range newOldGoRoutine {
		goRoutines = append(goRoutines, val)
	}

	sort.Slice(goRoutines, func(i, j int) bool {
		if goRoutines[i].FirstOccurrence.Equal(goRoutines[j].FirstOccurrence) {
			return strings.Compare(goRoutines[i].ID, goRoutines[j].ID) < 0
		}
		return goRoutines[i].FirstOccurrence.Before(goRoutines[j].FirstOccurrence)
	})

	for _, gr := range goRoutines {
		runningTime := currentTime.Sub(gr.FirstOccurrence).Seconds()
		if runningTime > 3600 {
			log.Debug("\nremaining routine more than an hour", "ID", gr.ID, "running seconds", runningTime)
		}
		log.Debug("\nremaining routine", "ID", gr.ID, "running seconds", runningTime)
	}

	for k, val := range grd.LatestData[newRoutine] {
		log.Debug("\nnew routine", "ID", val.ID, "\nData", val.AllData+"\n")
		grd.LatestData[oldRoutine][k] = val
	}

	grd.LatestData[newRoutine] = make(map[string]*GoRoutineData)
}

func getGoroutineId(goroutineString string) string {

	rg := regexp.MustCompile(`goroutine \d+`)

	submatchall := rg.FindAllString(goroutineString, -1)
	for _, element := range submatchall {
		replaced := strings.ReplaceAll(element, "goroutine ", "")
		return replaced
	}
	return ""
}
