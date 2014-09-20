// Copyright 2013 The lime Authors.
// Use of this source code is governed by a 2-clause
// BSD-style license that can be found in the LICENSE file.

package keys

import (
	"encoding/json"
	"github.com/limetext/lime/backend/util"
)

type (
	// A Context definition for which a key binding
	// is to be considered.
	KeyContext struct {
		rawKeyContext
	}

	// TODO(.): HACK. This is because I want to use the default UnmarshalJSON
	// behaviour on most of the struct member, but then also do some custom
	// handling too.
	//
	// So the publicly exported KeyContext implements a custom UnmarshalJSON,
	// which then invokes the default UnMarshalJSON handling on the embedded
	// rawKeyContext, and then does it's own custom code after that.
	//
	// Is there a better way to do this?
	rawKeyContext struct {
		Key      string      //< The context's name.
		Operator util.Op     //< The operation to perform.
		Operand  interface{} //< The operand on which this operation should be performed.
		MatchAll bool        `json:"match_all"` //< Whether all selections should match the context or if it's enough for just one to match.
	}
)

func (k *KeyContext) UnmarshalJSON(d []byte) error {
	if err := json.Unmarshal(d, &k.rawKeyContext); err != nil {
		return err
	}
	if k.Operand == nil {
		k.Operand = true
	}
	return nil
}
