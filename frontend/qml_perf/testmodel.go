package main

// test model to pass to qml and back
type TestModelSubModel struct {
	Text string
}

type TestModel struct {
	IntField    int
	DoubleField float32
	StringField string
	InnerModel  TestModelSubModel
}

func (model *TestModel) MyFunction() {

}

func (model *TestModel) MyFunctionWithInt(i int) int {

	return 102
}
