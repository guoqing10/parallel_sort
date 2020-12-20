package parallel_sort

import "sync"

// Auto-generated variant of sort.go:insertionSort
func insertionSortFunc(data lessSwap, a, b int) {
	for i := a + 1; i < b; i++ {
		for j := i; j > a && data.Less(j, j-1); j-- {
			data.Swap(j, j-1)
		}
	}
}

// Auto-generated variant of sort.go:siftDown
func siftDownFunc(data lessSwap, lo, hi, first int) {
	root := lo
	for {
		child := 2*root + 1
		if child >= hi {
			break
		}
		if child+1 < hi && data.Less(first+child, first+child+1) {
			child++
		}
		if !data.Less(first+root, first+child) {
			return
		}
		data.Swap(first+root, first+child)
		root = child
	}
}

// Auto-generated variant of sort.go:heapSort
func heapSortFunc(data lessSwap, a, b int) {
	first := a
	lo := 0
	hi := b - a
	for i := (hi - 1) / 2; i >= 0; i-- {
		siftDownFunc(data, i, hi, first)
	}
	for i := hi - 1; i >= 0; i-- {
		data.Swap(first, first+i)
		siftDownFunc(data, lo, i, first)
	}
}

// Auto-generated variant of sort.go:medianOfThree
func medianOfThreeFunc(data lessSwap, m1, m0, m2 int) {
	if data.Less(m1, m0) {
		data.Swap(m1, m0)
	}
	if data.Less(m2, m1) {
		data.Swap(m2, m1)
		if data.Less(m1, m0) {
			data.Swap(m1, m0)
		}
	}
}

// Auto-generated variant of sort.go:doPivot
func doPivotFunc(data lessSwap, lo, hi int) (midlo, midhi int) {
	m := int(uint(lo+hi) >> 1)
	if hi-lo > 40 {
		s := (hi - lo) / 8
		medianOfThreeFunc(data, lo, lo+s, lo+2*s)
		medianOfThreeFunc(data, m, m-s, m+s)
		medianOfThreeFunc(data, hi-1, hi-1-s, hi-1-2*s)
	}
	medianOfThreeFunc(data, lo, m, hi-1)
	pivot := lo
	a, c := lo+1, hi-1
	for ; a < c && data.Less(a, pivot); a++ {
	}
	b := a
	for {
		for ; b < c && !data.Less(pivot, b); b++ {
		}
		for ; b < c && data.Less(pivot, c-1); c-- {
		}
		if b >= c {
			break
		}
		data.Swap(b, c-1)
		b++
		c--
	}
	protect := hi-c < 5
	if !protect && hi-c < (hi-lo)/4 {
		dups := 0
		if !data.Less(pivot, hi-1) {
			data.Swap(c, hi-1)
			c++
			dups++
		}
		if !data.Less(b-1, pivot) {
			b--
			dups++
		}
		if !data.Less(m, pivot) {
			data.Swap(m, b-1)
			b--
			dups++
		}
		protect = dups > 1
	}
	if protect {
		for {
			for ; a < b && !data.Less(b-1, pivot); b-- {
			}
			for ; a < b && data.Less(a, pivot); a++ {
			}
			if a >= b {
				break
			}
			data.Swap(a, b-1)
			a++
			b--
		}
	}
	data.Swap(pivot, b-1)
	return b - 1, c
}

// Auto-generated variant of sort.go:quickSort
func quickSortFunc(data lessSwap, a, b, maxDepth int, curParallel int, totalParallel int, innerLock *sync.WaitGroup, needDone bool) {
	curParallel++
	defer func() {
		if needDone {
			innerLock.Done()
		}
	}()
	for b-a > 12 {
		if maxDepth == 0 {
			heapSortFunc(data, a, b)
			return
		}
		maxDepth--
		mlo, mhi := doPivotFunc(data, a, b)
		if curParallel <= totalParallel {
			innerLock.Add(1)
			if mlo-a < b-mhi {
				go quickSortFunc(data, a, mlo, maxDepth, curParallel, totalParallel, innerLock, true)
				a = mhi
			} else {
				go quickSortFunc(data, mhi, b, maxDepth, curParallel, totalParallel, innerLock, true)
				b = mlo
			}
		} else {
			if mlo-a < b-mhi {
				quickSortFunc(data, a, mlo, maxDepth, curParallel, totalParallel, innerLock, false)
				a = mhi
			} else {
				quickSortFunc(data, mhi, b, maxDepth, curParallel, totalParallel, innerLock, false)
				b = mlo
			}
		}
	}
	if b-a > 1 {
		for i := a + 6; i < b; i++ {
			if data.Less(i, i-6) {
				data.Swap(i, i-6)
			}
		}
		insertionSortFunc(data, a, b)
	}
}

// Auto-generated variant of sort.go:quickSort
func quickSortByTwoGroupsRecursion(data lessSwap, a, b, maxDepth int, curParallel int, totalParallel int, innerLock *sync.WaitGroup, needNone bool) {
	curParallel++
	defer func() {
		if needNone {
			innerLock.Done()
		}
	}()

	if maxDepth == 0 {
		heapSortFunc(data, a, b)
		return
	}
	maxDepth--
	mlo, mhi := doPivotFunc(data, a, b)
	if curParallel <= totalParallel {
		innerLock.Add(2)
		go quickSortByTwoGroupsRecursion(data, a, mlo, maxDepth, curParallel, totalParallel, innerLock, true)
		go quickSortByTwoGroupsRecursion(data, mhi, b, maxDepth, curParallel, totalParallel, innerLock, true)
	} else {
		if mlo-a < b-mhi {
			quickSortByTwoGroupsRecursion(data, a, mlo, maxDepth, curParallel, totalParallel, innerLock, false)
			a = mhi
		} else {
			quickSortByTwoGroupsRecursion(data, mhi, b, maxDepth, curParallel, totalParallel, innerLock, false)
			b = mlo
		}
	}
}

// maxDepth returns a threshold at which quicksort should switch
// to heapsort. It returns 2*ceil(lg(n+1)).
func maxDepth(n int) int {
	var depth int
	for i := n; i > 0; i >>= 1 {
		depth++
	}
	return depth * 2
}

// lessSwap is a pair of Less and Swap function for use with the
// auto-generated func-optimized variant of sort.go in
// zfuncversion.go.
type lessSwap struct {
	Less func(i, j int) bool
	Swap func(i, j int)
}
