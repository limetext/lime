package loaders

import (
	sj "encoding/json"
	"errors"
	"fmt"
	"lime/backend/loaders/json"
	. "lime/backend/primitives"
)

func LoadJSON(data []byte, intf interface{}) error {
	var (
		b   = Buffer{}
		p   json.JSON
		set RegionSet
	)
	b.Insert(0, string(data))
	if !p.Parse(string(data)) {
		return fmt.Errorf("%s, %s", p.Error().String(), p.RootNode())
	} else {
		root := p.RootNode()
		for _, child := range root.Children {
			switch child.Name {
			case "BlockComment", "LineComment", "EndOfFile", "JunkComma":
				set.Add(Region{child.Range.Start, child.Range.End})
			default:
				return errors.New("Unhandled node: " + child.Name)
			}
		}
	}
	b.AddCallback(func(b *Buffer, pos, delta int) {
		set.Adjust(pos, delta)
	})
	i := 0
	for {
		l := set.Len()
		if i >= l {
			break
		}
		r := set.Get(i)
		b.Erase(r.Begin(), r.Size())
		if l2 := set.Len(); l2 != l {
			continue
		}
		i++
	}
	// TODO(q): Map any line/column errors to the actual file's line/column
	return sj.Unmarshal([]byte(b.Substr(Region{0, b.Size()})), intf)
}
