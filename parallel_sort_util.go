package parallel_sort

import (
	"reflect"
	"runtime"
	"sync"
)

func ParallelQuickSort(slice interface{}, less func(i, j int) bool) {
	totalParallel := runtime.NumCPU() / 2
	ParallelQuickSortWithOccurrence(slice, less, totalParallel)
}

func ParallelQuickSortWithOccurrence(slice interface{}, less func(i, j int) bool, totalParallel int) {
	rv := reflectValueOf(slice)
	swap := reflectSwapper(slice)
	length := rv.Len()

	var innerLock sync.WaitGroup
	innerLock.Add(1)

	quickSortFunc(lessSwap{less, swap}, 0, length, maxDepth(length), 0, totalParallel, &innerLock, true)
	//quickSortByTwoGroupsRecursion(lessSwap{less, swap}, 0, length, maxDepth(length), 0, totalParallel, &innerLock, true)

	innerLock.Wait()
}

var reflectValueOf = reflect.ValueOf
var reflectSwapper = reflect.Swapper
