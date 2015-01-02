// Copyright 2015 The lime Authors.
// Use of this source code is governed by a 2-clause
// BSD-style license that can be found in the LICENSE file.

package plist

import (
	"archive/zip"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
	"time"
)

var _ = time.Time{}

const testname = "testdata/C.plist"

func loadData(path string) (retdata string, err error) {
	var data []byte
	if strings.HasPrefix(path, "http://") {
		p2 := strings.Replace(strings.Replace(path, "http://", "http_", -1), "/", "_", -1)
		if retdata, err = loadData(p2); err != nil {
			if res, err := http.Get(path); err != nil {
				return "", err
			} else {
				defer res.Body.Close()
				if data, err = ioutil.ReadAll(res.Body); err != nil {
					return "", err
				} else if err = ioutil.WriteFile(p2, data, 0644); err != nil {
					return "", err
				} else {
					return loadData(p2)
				}
			}
		} else {
			return
		}
	} else if strings.HasSuffix(path, ".zip") {
		if zf, err := zip.OpenReader(path); err != nil {
			return "", err
		} else {
			defer zf.Close()
			f, err := zf.File[0].Open()
			if err != nil {
				return "", err
			}
			defer f.Close()
			if data, err = ioutil.ReadAll(f); err != nil {
				return "", err
			}
		}
	} else {
		data, err = ioutil.ReadFile(path)
	}
	return string(data), err
}

func TestParser(t *testing.T) {
	var p PLIST
	if data, err := loadData(testname); err != nil {
		t.Fatal(err)
	} else {
		root := p.RootNode()
		if !p.Parse(data) {

			t.Fatalf("Didn't parse correctly: %s\n", p.Error())
		} else {

			if root.Range.B != p.ParserData.Len() {
				t.Fatalf("Parsing didn't finish: %v\n%s", root, p.Error())
			}
		}
	}
}

func BenchmarkParser(b *testing.B) {
	var p PLIST
	if data, err := loadData(testname); err != nil {
		b.Fatal(err)
	} else {
		for i := 0; i < b.N; i++ {
			p.Parse(data)
		}
	}
}
