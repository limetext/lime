// Copyright 2013 The lime Authors.
// Use of this source code is governed by a 2-clause
// BSD-style license that can be found in the LICENSE file.

package commands

import (
    . "github.com/limetext/lime/backend"
    . "github.com/limetext/text"
)

type (
    // Transpose: Swap the characters on either side of the cursor,
    // then move the cursor forward one character.
    TransposeCommand struct {
        DefaultCommand
    }
)

func (c *TransposeCommand) Run(v *View, e *Edit) error {
    /*
    	Correct behavior of Transpose:
    		- Swap the characters on either side of the cursor(s), then move
    		  forward one character. If a region is selected, do nothing.
    */

    if v.Sel().HasNonEmpty() {
        return nil
    }

    rs := v.Sel().Regions()
    for i := range rs {
        r := rs[i]
        if r.A == 0 || r.A >= v.Buffer().Size() {
            continue
        }
        s := Region{r.A - 1, r.A + 1}
        rns := v.Buffer().SubstrR(s)
        rnd := make([]rune, len(rns))
        rnd[0] = rns[1]
        rnd[1] = rns[0]
        v.Replace(e, s, string(rnd))
    }

    return nil
}

func init() {
    register([]Command{
        &TransposeCommand{},
    })
}
