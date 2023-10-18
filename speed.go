package main

import (
	"encoding/json"
	"path"
	"sort"
	"sync"
	"time"
)

type Cost struct {
	P99   int64 `json:"p99"`
	P90   int64 `json:"p90"`
	P50   int64 `json:"p50"`
	Total int64 `json:"total"`
}

type SpeedInfo struct {
	Name    string
	Times   int
	Cost    Cost
	Success int
	Failed  int
	Err     string
}

func ConvertCost(costs []int64, total int64) Cost {
	sort.Slice(costs, func(i, j int) bool {
		return costs[i] < costs[j]
	})
	p99index := int(float64(len(costs)) * 0.99)
	p90index := int(float64(len(costs)) * 0.90)
	p50index := int(float64(len(costs)) * 0.50)
	return Cost{
		P99:   costs[p99index],
		P90:   costs[p90index],
		P50:   costs[p50index],
		Total: total,
	}
}

var speedInfos sync.Map

func benchFn(fn func() error, times int, name string) {
	start := time.Now()
	success := 0
	failed := 0
	var e string
	costs := make([]int64, 0, times)
	for i := 0; i < times; i++ {
		perStart := time.Now()
		err := fn()
		if err != nil {
			failed++
			e = err.Error()
		} else {
			success++
		}
		costs = append(costs, int64(time.Since(perStart).Milliseconds()))
	}
	since := time.Since(start)
	// fmt.Printf("%s run %d cost: %dms\n", name, times, since.Milliseconds())

	speedInfos.Store(name, SpeedInfo{
		Name:    name,
		Times:   times,
		Cost:    ConvertCost(costs, since.Milliseconds()),
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
		return speedList[i].Cost.Total > speedList[j].Cost.Total
	})
	for _, l := range speedList {
		l.Name = path.Base(l.Name)
		v, _ := json.Marshal(l)
		res += string(v) + "\n"
	}
	speedInfos = sync.Map{}
	return res
}
