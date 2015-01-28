package main

import (
	"fmt"
	"reflect"
	"runtime"
	"time"
)

type Benchmark struct {
	f     func(*Benchmark)
	times int
	start time.Time
	stop  time.Time
}

type BenchmarkResult struct {
	N int           // The number of iterations.
	T time.Duration // The total time taken.

}

func (r BenchmarkResult) NsPerOp() int64 {
	if r.N <= 0 {
		return 0
	}
	return r.T.Nanoseconds() / int64(r.N)
}

func (r BenchmarkResult) MsPerOp() int64 {
	if r.N <= 0 {
		return 0
	}
	return int64((r.T / time.Millisecond)) / int64(r.N)
}

func getFunctionName(i interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
}

func (r BenchmarkResult) String() string {

	nsop := r.NsPerOp()
	msop := r.MsPerOp()
	ns := fmt.Sprintf("%10d ns/op", nsop)
	if r.N > 0 && nsop < 100 {
		// The format specifiers here make sure that
		// the ones digits line up for all three possible formats.
		if nsop < 10 {
			ns = fmt.Sprintf("%13.2f ns/op", float64(r.T.Nanoseconds())/float64(r.N))
		} else {
			ns = fmt.Sprintf("%12.1f ns/op", float64(r.T.Nanoseconds())/float64(r.N))
		}
	}

	if msop > 0 {
		ns = ns + "   " + fmt.Sprintf("%10d ms/op", msop)
	}

	return fmt.Sprintf("%8d\t%s", r.N, ns)
}

func Run(f func(*Benchmark), times int) {
	RunWithBenchmark(&Benchmark{f: f, times: times})
}

// Profiles only the function in question, and adjust the start/stop timesmaps
func (b *Benchmark) Profile(f func()) {
	b.start = time.Now()
	f()
	b.stop = time.Now()

}

func RunWithBenchmark(bench *Benchmark) {

	if bench.times == 0 {
		bench.times = 100
	}

	benchmarkChan <- bench

	newbench := <-benchmarkChan

	var duration = newbench.stop.Sub(newbench.start)

	result := BenchmarkResult{N: newbench.times, T: duration}

	fmt.Printf("%60s %s\n", getFunctionName(bench.f), result.String())
	//fmt.Println((duration.Nanoseconds() / times))

}
