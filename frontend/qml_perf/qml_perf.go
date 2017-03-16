// Copyright 2014 The lime Authors.
// Use of this source code is governed by a 2-clause
// BSD-style license that can be found in the LICENSE file.

package main

import (
	//"fmt"
	"gopkg.in/qml.v1"
	"math/rand"
	"runtime"
	"runtime/debug"
	"time"
)

var benchmarkChan chan *Benchmark
var quitChan chan bool

func QmlBenchmark() error {
outer:
	for {
		//fmt.Println("QmlBenchmark")
		select {
		case bench := <-benchmarkChan:

			// make tests less fluctuated by gc calls
			runtime.GC()
			debug.FreeOSMemory()

			bench.start = time.Now()

			for i := 0; i < bench.times; i++ {
				bench.f(bench)
			}

			// if stop timer is not set by the test, stop it
			if bench.stop.Nanosecond() == 0 {
				bench.stop = time.Now()
			}

			// not enough gc
			runtime.GC()
			debug.FreeOSMemory()

			benchmarkChan <- bench

		case <-quitChan:
			//fmt.Println("quit flag received")
			break outer

		case <-time.After(time.Millisecond * 100):
			//fmt.Println("sleeping")

		}
		//fmt.Println("loop")
	}
	return nil
}

func main() {

	runtime.LockOSThread()

	benchmarkChan = make(chan *Benchmark)
	quitChan = make(chan bool)

	// do we need to do this?
	//	qml.SetupTesting()

	go RunTests()

	qml.Run(QmlBenchmark)

}

func RunTests() {

	const testRuns int = 100

	Run(QmlEngineCreationTest, testRuns)
	Run(QmlEngineCreationAndDestroyTest, testRuns)

	Run(QmlEngineLoadForm1, testRuns)
	Run(QmlEngineCreateForm1, testRuns)

	Run(QmlEngineLoadForm2, testRuns)
	Run(QmlEngineCreateForm2, testRuns)

	Run(QmlEngineLoadForm3, testRuns)
	Run(QmlEngineCreateForm3, testRuns)

	Run(QmlEngineLoadForm4, testRuns)
	Run(QmlEngineCreateForm4, testRuns)

	Run(QmlEngineCallingJsFunction, testRuns)

	Run(QmlEngineChangingObjectPropertyInJs, testRuns)
	Run(QmlEngineChangingObjectPropertyInGo, testRuns)

	Run(QmlEnginePassingGoObject, testRuns)

	Run(QmlEngineCallingGoModelFunctionFromJs, testRuns)
	Run(QmlEngineCallingGoModelFunctionFromJsWithArguments, testRuns)

	quitChan <- true
}

func QmlEngineCreationTest(*Benchmark) {

	qml.NewEngine()
	//engine.Destroy()

}

func QmlEngineCreationAndDestroyTest(*Benchmark) {

	engine := qml.NewEngine()
	engine.Destroy()

}

func QmlEngineLoadForm1(bench *Benchmark) {

	engine := qml.NewEngine()

	bench.Profile(func() {
		engine.LoadFile("testdata/form1.qml")

	})

	engine.Destroy()

}

func QmlEngineCreateForm1(bench *Benchmark) {

	engine := qml.NewEngine()

	component, _ := engine.LoadFile("testdata/form1.qml")
	bench.Profile(func() {
		component.Create(nil)

	})

	component.Destroy()
	engine.Destroy()

}

func QmlEngineLoadForm2(bench *Benchmark) {

	engine := qml.NewEngine()

	bench.Profile(func() {
		engine.LoadFile("testdata/form2.qml")
	})

	engine.Destroy()

}

func QmlEngineCreateForm2(bench *Benchmark) {

	engine := qml.NewEngine()

	component, _ := engine.LoadFile("testdata/form2.qml")
	bench.Profile(func() {
		component.Create(nil)

	})

	component.Destroy()
	engine.Destroy()

}

func QmlEngineLoadForm3(bench *Benchmark) {

	engine := qml.NewEngine()

	bench.Profile(func() {
		engine.LoadFile("testdata/form3.qml")

	})

	engine.Destroy()

}

func QmlEngineCreateForm3(bench *Benchmark) {

	engine := qml.NewEngine()

	component, _ := engine.LoadFile("testdata/form3.qml")
	bench.Profile(func() {
		component.Create(nil)

	})

	component.Destroy()
	engine.Destroy()

}

func QmlEngineLoadForm4(bench *Benchmark) {

	engine := qml.NewEngine()

	bench.Profile(func() {
		engine.LoadFile("testdata/form4.qml")

	})

	engine.Destroy()

}

func QmlEngineCreateForm4(bench *Benchmark) {

	engine := qml.NewEngine()

	component, _ := engine.LoadFile("testdata/form4.qml")
	bench.Profile(func() {
		component.Create(nil)

	})

	component.Destroy()
	engine.Destroy()

}

func QmlEngineCallingJsFunction(bench *Benchmark) {

	engine := qml.NewEngine()
	form, _ := engine.LoadFile("testdata/form3.qml")

	rc := form.CreateWindow(nil)

	bench.Profile(func() {
		rc.Call("functionShouldReturn_10")
	})

	engine.Destroy()

}

func QmlEngineChangingObjectPropertyInJs(bench *Benchmark) {

	engine := qml.NewEngine()
	form, _ := engine.LoadFile("testdata/form3.qml")

	rc := form.CreateWindow(nil)

	bench.Profile(func() {
		rc.Call("functionShouldExpandWindow")
	})

	engine.Destroy()

}

func QmlEngineChangingObjectPropertyInGo(bench *Benchmark) {

	engine := qml.NewEngine()

	form, _ := engine.LoadFile("testdata/form3.qml")

	rc := form.CreateWindow(nil)

	bench.Profile(func() {

		rnd := rand.Int()
		rc.Set("width", 101+rnd)

	})

	rc.Destroy()
	form.Destroy()
	engine.Destroy()

}

func QmlEnginePassingGoObject(bench *Benchmark) {

	engine := qml.NewEngine()
	form, _ := engine.LoadFile("testdata/form3.qml")

	rc := form.CreateWindow(nil)

	model := TestModel{InnerModel: TestModelSubModel{}}

	bench.Profile(func() {
		rc.Set("myProperty", model)
	})

	rc.Destroy()
	form.Destroy()
	engine.Destroy()

}

func QmlEngineCallingGoModelFunctionFromJs(bench *Benchmark) {

	engine := qml.NewEngine()

	model := new(TestModel)
	model.InnerModel = *new(TestModelSubModel)

	engine.Context().SetVar("myvar", model)

	form, _ := engine.LoadFile("testdata/form3.qml")

	rc := form.CreateWindow(nil)

	bench.Profile(func() {
		rc.Call("functionShouldCallmodelMyFunction")
	})

	rc.Destroy()
	form.Destroy()
	engine.Destroy()

}

func QmlEngineCallingGoModelFunctionFromJsWithArguments(bench *Benchmark) {

	engine := qml.NewEngine()

	model := new(TestModel)
	model.InnerModel = *new(TestModelSubModel)

	engine.Context().SetVar("myvar", model)

	form, _ := engine.LoadFile("testdata/form3.qml")

	rc := form.CreateWindow(nil)

	bench.Profile(func() {
		rc.Call("functionShouldCallmodelMyFunctionWithArgument")
	})

	rc.Destroy()
	form.Destroy()
	engine.Destroy()

}
