package loaders

import (
	sj "encoding/json"
	"errors"
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
		return errors.New(p.Error().String())
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
	b.AddCallback(set.Adjust)
	for i := 0; i < set.Len(); i++ {
		r := set.Get(i)
		b.Erase(r.Begin(), r.Size())
	}
	// TODO(q): Map any line/column errors to the actual file's line/column
	return sj.Unmarshal([]byte(b.Substr(Region{0, b.Size()})), intf)
}
