package backend

import (
	sj "encoding/json"
	"errors"
	"lime/backend/json"
)

func LoadJSON(data []byte, intf interface{}) error {
	var (
		b   = Buffer{data: string(data)}
		v   View
		p   json.JSON
		set = &v.selection
	)
	v.setBuffer(&b)
	if !p.Parse(b.data) {
		return errors.New(p.Error().String())
	} else {
		root := p.RootNode()
		for _, child := range root.Children {
			switch child.Name {
			case "BlockComment", "LineComment", "EndOfFile", "JunkComma":
				set.Regions = append(set.Regions, Region{child.Range.Start, child.Range.End})
			default:
				return errors.New("Unhandled node: " + child.Name)
			}
		}
	}
	for i := range v.selection.Regions {
		v.Erase(v.selection.Regions[i])
	}
	// TODO(q): Map any line/column errors to the actual file's line/column
	return sj.Unmarshal([]byte(b.data), intf)
}
