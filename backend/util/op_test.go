// Copyright 2013 The lime Authors.
// Use of this source code is governed by a 2-clause
// BSD-style license that can be found in the LICENSE file.

package util

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestUnmarshalJSONError(t *testing.T) {
	var o Op
	err := o.UnmarshalJSON(nil)
	if err != nil {
		fmt.Print("Error on unmarshaling")
	}
}

func TestUnmarshalJSONDefault(t *testing.T) {
	var o Op
	d, err := json.Marshal("")
	if err != nil {
		fmt.Println(err)
		return
	}
	err = o.UnmarshalJSON(d)
	if err != nil {
		fmt.Print("Error on unmarshaling")
	}
}

func TestUnmarshalJSON_Not_equal(t *testing.T) {
	var o Op
	d, err := json.Marshal("not_equal")
	if err != nil {
		fmt.Println(err)
		return
	}
	err = o.UnmarshalJSON(d)
	if err != nil {
		fmt.Print("Error on unmarshaling")
	}
}

func TestUnmarshalJSON_Regex_match(t *testing.T) {
	var o Op
	d, err := json.Marshal("regex_match")
	if err != nil {
		fmt.Println(err)
		return
	}
	err = o.UnmarshalJSON(d)
	if err != nil {
		fmt.Print("Error on unmarshaling")
	}
}
func TestUnmarshalJSON_Not_regex_match(t *testing.T) {
	var o Op
	d, err := json.Marshal("not_regex_match")
	if err != nil {
		fmt.Println(err)
		return
	}
	err = o.UnmarshalJSON(d)
	if err != nil {
		fmt.Print("Error on unmarshaling")
	}
}
func TestUnmarshalJSON_Regex_contains(t *testing.T) {
	var o Op
	d, err := json.Marshal("regex_contains")
	if err != nil {
		fmt.Println(err)
		return
	}
	err = o.UnmarshalJSON(d)
	if err != nil {
		fmt.Print("Error on unmarshaling")
	}
}
func TestUnmarshalJSON_Not_regex_contains(t *testing.T) {
	var o Op
	d, err := json.Marshal("not_regex_contains")
	if err != nil {
		fmt.Println(err)
		return
	}
	err = o.UnmarshalJSON(d)
	if err != nil {
		fmt.Print("Error on unmarshaling")
	}
}
func TestUnmarshalJSON_Notequal(t *testing.T) {
	var o Op
	d, err := json.Marshal("not_equal")
	if err != nil {
		fmt.Println(err)
		return
	}
	err = o.UnmarshalJSON(d)
	if err != nil {
		fmt.Print("Error on unmarshaling")
	}
}
