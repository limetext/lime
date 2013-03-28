package backend

type Editor struct {
	HasSettings
}

type Window struct {
	HasSettings
	views []*View
}
