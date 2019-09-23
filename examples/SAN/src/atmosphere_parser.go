package src

import (
	"encoding/json"
	"math/rand"
	"sort"
)

type (
	ValueFail struct {
		Value float64 `json:"value"`
		Fail  float64 `json:"fail"`
	}
	AtmosphereDependency struct {
		arr   []float64
		table map[float64]float64
	}
)

func ParseAtmosphereInputDependency(fileName string) map[string][]*ValueFail {

	dataFail := make(map[string][]*ValueFail)
	bytes := ParseFileAndUnmarshal(fileName)

	err := json.Unmarshal(bytes, &dataFail)

	if err != nil {
		panic(err)
	}

	return dataFail
}

func GetAtmosphereDependencyRoot(parameter string, data map[string][]*ValueFail) *AtmosphereDependency {
	dep, ok := data[parameter]
	arr := make([]float64, 0)
	table := make(map[float64]float64)

	if !ok {
		panic("No such atmospheric dependency")
	}

	for _, vf := range dep {
		table[vf.Value] = vf.Fail
		arr = append(arr, vf.Value)
	}

	sort.Float64s(arr)

	return &AtmosphereDependency{
		table: table,
		arr:   arr,
	}
}

func binarySearch(arr []float64, target float64) float64 {
	low := 0
	high := len(arr) - 1

	if target < arr[0] {
		return arr[0]
	}

	if target > arr[len(arr)-1] {
		return arr[len(arr)-1]
	}

	for low <= high {
		middle := (low + high) / 2

		if arr[middle] < target {
			low = middle + 1
		} else if arr[middle] > target {
			high = middle - 1
		} else {
			return target
		}
	}

	if arr[low]-target < target-arr[high] {
		return arr[low]
	} else {
		return arr[high]
	}
}

func MonteCarlo(failure float64) bool {
	rValue := rand.Float64()
	if failure < rValue {
		return true
	}
	return false
}

func Reader() {

}
