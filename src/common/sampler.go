package common

import (
	"math/rand"
	"sort"
	"time"
)

func getSamples(start int, end int, count int) []int {
	rand.New(rand.NewSource(time.Now().UnixNano()))
	if end < start || (end-start) < count {
		return nil
	}
	nums := make([]int, 0)
	for len(nums) < count {
		num := rand.Intn(end-start) + start
		exist := false
		for _, v := range nums {
			if v == num {
				exist = true
				break
			}
		}
		if !exist {
			nums = append(nums, num)
		}
	}
	sort.Ints(nums)
	return nums
}

func GetSample[T any](sampleCount int, data []T) []T {
	if sampleCount < 1 {
		return data
	}
	totalItems := len(data)
	if sampleCount >= totalItems {
		return data
	}
	output := make([]T, sampleCount)
	for i, index := range getSamples(0, totalItems, sampleCount) {
		output[i] = data[index]
	}
	return output
}
