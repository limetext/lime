// Copyright 2013 The lime Authors.
// Use of this source code is governed by a 2-clause
// BSD-style license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"go/format"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

var year = strconv.FormatInt(int64(time.Now().Year()), 10)
var licenseheader = []byte(`// Copyright ` + year + ` The lime Authors.
// Use of this source code is governed by a 2-clause
// BSD-style license that can be found in the LICENSE file.
`)

func patch(path string, fi os.FileInfo, err error) error {
	if fi.IsDir() {
		switch filepath.Base(path) {
		case "testdata", "build2", "3rdparty":
			return filepath.SkipDir
		}
		return nil
	}

	switch filepath.Ext(path) {
	case ".go", ".c", ".cpp":
	default:
		return nil
	}

	changed := false
	cmp, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	lhn := append(licenseheader, '\n')
	if !bytes.Equal([]byte("// Copyright"), cmp[:12]) {
		cmp = append(lhn, cmp...)
		log.Println("Added license to", path)
		changed = true
	}

	fmt, err := format.Source(cmp)
	if err != nil {
		return err
	}

	if !bytes.Equal(cmp, fmt) {
		log.Println("Formatted", path)
		changed = true
	}

	if changed {
		return ioutil.WriteFile(path, fmt, fi.Mode().Perm())
	}

	return nil
}

func main() {
	log.Println("done: ", filepath.Walk("./", patch))
}
