package sublime

import (
	"github.com/qur/gopy/lib"
	"log"
)

func init() {
	methods := []py.Method{}
	py.Initialize()
	m, err := py.InitModule("sublime", methods)
	if err != nil {
		log.Fatal(err)
	}

	c, err := _regionClass.Create()
	if err != nil {
		log.Fatal(err)
	}
	if err := m.AddObject("Region", c); err != nil {
		log.Fatal(err)
	}
}
