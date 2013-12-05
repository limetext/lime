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
)

var licenseheader = []byte(`// Copyright 2013 The lime Authors.
// Use of this source code is governed by a 2-clause
// BSD-style license that can be found in the LICENSE file.
`)

func patch(path string, fi os.FileInfo, err error) error {
	if fi.IsDir() {
		switch filepath.Base(path) {
		case "testdata", "build2", "3rdparty":
			return filepath.SkipDir
		}
	}
	if fi.IsDir() {
		return nil
	}
	switch filepath.Ext(path) {
	case ".go", ".c", ".cpp":
	default:
		return nil
	}
	cmp, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	lhn := append(licenseheader, '\n')
	changed := false
	if !bytes.Equal(licenseheader, cmp[:len(licenseheader)]) {
		cmp = append(lhn, cmp...)
		changed = true
	} else if cmp[len(licenseheader)] != '\n' {
		cmp = append(lhn, cmp[len(licenseheader):]...)
		changed = true
	}
	if changed {
		log.Println("Added license to", path)
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
	log.Println("done: ", filepath.Walk("../", patch))
}
