package parallel_sort

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"math/rand"
	"sort"
	"sync"
	"testing"
	"time"
)

type parallelTestStruct struct {
	LaneCodeList   []string
	ZoneSectorList []int64
	SortFactor     int64
}

// 并行快排
func TestParallelQuickSort(t *testing.T) {
	testQuickList := make([]*parallelTestStruct, 0)
	testMergeList := make([]*parallelTestStruct, 0)
	TestParalleList := make([]*parallelTestStruct, 0)
	totalNum := 400 * 10000
	fmt.Println("total data: ", totalNum)
	for i := 0; i < totalNum; i++ {
		v := rand.Int() % (totalNum * 10)
		p := &parallelTestStruct{
			SortFactor: int64(v),
		}
		testQuickList = append(testQuickList, p)
		testMergeList = append(testMergeList, p)
		TestParalleList = append(TestParalleList, p)
	}

	startTime := time.Now().UnixNano()
	sort.Slice(testQuickList, func(i, j int) bool {
		return testQuickList[i].SortFactor < testQuickList[j].SortFactor
	})
	fmt.Println("quick time: ", (time.Now().UnixNano()-startTime)/1000000, " ms")

	startTime = time.Now().UnixNano()
	sort.SliceStable(testMergeList, func(i, j int) bool {
		return testMergeList[i].SortFactor < testMergeList[j].SortFactor
	})
	fmt.Println("stable time: ", (time.Now().UnixNano()-startTime)/1000000, " ms")

	startTime = time.Now().UnixNano()
	ParallelQuickSort(TestParalleList, func(i, j int) bool {
		return TestParalleList[i].SortFactor < TestParalleList[j].SortFactor
	})
	fmt.Println("parallel time: ", (time.Now().UnixNano()-startTime)/1000000, " ms")

	for i := 0; i < totalNum; i++ {
		assert.Equal(t, testQuickList[i].SortFactor, TestParalleList[i].SortFactor)
		assert.Equal(t, testMergeList[i].SortFactor, TestParalleList[i].SortFactor)
	}
}

// test  int[] parallel sort
func TestParallelQuickSortInt(t *testing.T) {
	listSort := make([]int, 0)
	listMerge := make([]int, 0)
	listQuick := make([]int, 0)
	totalNum := 4000000
	fmt.Println("total data: ", totalNum)
	for i := 0; i < totalNum; i++ {
		v := rand.Int() % totalNum * 5
		listSort = append(listSort, v)
		listMerge = append(listMerge, v)
		listQuick = append(listQuick, v)
	}

	startTime := time.Now().UnixNano()
	sort.Slice(listSort, func(i, j int) bool {
		return listSort[i] < listSort[j]
	})
	fmt.Println("quick time: ", (time.Now().UnixNano()-startTime)/1000000, " ms")

	startTime = time.Now().UnixNano()
	parallelMergeSortInt(listMerge, 8)
	fmt.Println("stable time: ", (time.Now().UnixNano()-startTime)/1000000, " ms")

	startTime = time.Now().UnixNano()
	parallelQuickSortInt(listQuick, 8)
	fmt.Println("parallel time: ", (time.Now().UnixNano()-startTime)/1000000, " ms")

	for i := 0; i < totalNum; i++ {
		assert.Equal(t, listSort[i], listQuick[i])
		assert.Equal(t, listMerge[i], listQuick[i])
	}
}

func parallelMergeSortInt(a []int, level int) {
	ch := make(chan int, 1)
	defer close(ch)
	mergeSortInt(a, 0, len(a)-1, ch, 0, level)
}

func mergeSortInt(a []int, left, right int, c chan int, depth int, level int) {
	depth++
	if left < right {
		ch := make(chan int, 2)
		defer close(ch)
		mid := left + (right-left)/2
		if depth >= level {
			mergeSortInt(a, left, mid, ch, depth, level)
			mergeSortInt(a, mid+1, right, ch, depth, level)
		} else {
			go mergeSortInt(a, left, mid, ch, depth, level)
			go mergeSortInt(a, mid+1, right, ch, depth, level)
		}
		<-ch
		<-ch
		MergeInt(a, left, mid, right)
	}
	c <- 1
}

func MergeInt(a []int, left, mid, right int) {
	arr := make([]int, 0)
	i, j := left, mid+1
	for i <= mid && j <= right {
		if a[i] <= a[j] {
			arr = append(arr, a[i])
			i++
		} else {
			arr = append(arr, a[j])
			j++
		}
	}
	arr = append(arr, a[i:mid+1]...)
	arr = append(arr, a[j:right+1]...)

	for i, v := range arr {
		a[left+i] = v
	}
}

func parallelQuickSortInt(num []int, level int) {
	var innerLock sync.WaitGroup
	//如果缺少该句,由于下面第一次调用QuickSort,没有加一操作,执行完后直接lock.Done(),将导致数量减为-1而报错
	innerLock.Add(1)
	quickSortInt(num, 0, len(num)-1, 0, level, &innerLock)
	innerLock.Wait()
}

func quickSortInt(num []int, low, high int, depth int, level int, innerLock *sync.WaitGroup) {
	depth++
	defer func() {
		if depth <= level {
			innerLock.Done()
		}
	}()
	if low >= high {
		return
	}
	i, j := low, high
	key := num[low]
	for i < j {
		for j > i && num[j] >= key {
			j--
		}
		num[i] = num[j]

		for i < j && num[i] < key {
			i++
		}
		num[j] = num[i]
	}
	num[j] = key
	if depth >= level {
		quickSortInt(num, low, i-1, depth, level, innerLock)
		quickSortInt(num, i+1, high, depth, level, innerLock)
	} else {
		innerLock.Add(2)
		go quickSortInt(num, low, i-1, depth, level, innerLock)
		go quickSortInt(num, i+1, high, depth, level, innerLock)
	}
}
