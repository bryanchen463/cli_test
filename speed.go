package main

import (
	"encoding/json"
	"path"
	"sort"
	"sync"
	"time"
)

type SpeedInfo struct {
	Name    string
	Times   int
	Cost    int
	Success int
	Failed  int
	Err     string
}

var speedInfos sync.Map

func benchFn(fn func() error, times int, name string) {
	start := time.Now()
	success := 0
	failed := 0
	var e string
	for i := 0; i < times; i++ {
		err := fn()
		if err != nil {
			failed++
			e = err.Error()
		} else {
			success++
		}

	}
	since := time.Since(start)
	// fmt.Printf("%s run %d cost: %dms\n", name, times, since.Milliseconds())
	speedInfos.Store(name, SpeedInfo{
		Name:    name,
		Times:   times,
		Cost:    int(since.Milliseconds()),
		Success: success,
		Failed:  failed,
		Err:     e,
	})
}

func Result() string {
	res := ""
	var speedList []SpeedInfo
	speedInfos.Range(func(key, value any) bool {
		speedInfo := value.(SpeedInfo)
		speedList = append(speedList, speedInfo)
		return true
	})
	sort.Slice(speedList, func(i, j int) bool {
		return speedList[i].Cost > speedList[j].Cost
	})
	for _, l := range speedList {
		l.Name = path.Base(l.Name)
		v, _ := json.Marshal(l)
		res += string(v) + "\n"
	}
	speedInfos = sync.Map{}
	return res
}
