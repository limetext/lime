package textmate

import (
	"code.google.com/p/log4go"
	"encoding/json"
	"fmt"
	"image/color"
	"io/ioutil"
	"lime/backend/loaders"
	"lime/backend/render"
	"strconv"
	"strings"
)

type (
	Color color.RGBA

	// TODO(q): personally I don't care about the font style attributes
	Settings map[string]Color

	ScopeSetting struct {
		Name     string
		Scope    string
		Settings Settings
	}
	Theme struct {
		GutterSettings Settings
		Name           string
		Settings       []ScopeSetting
		UUID           UUID
	}
)

func LoadTheme(filename string) (*Theme, error) {
	var scheme Theme
	if d, err := ioutil.ReadFile(filename); err != nil {
		return nil, fmt.Errorf("Unable to load colorscheme definition: %s", err)
	} else if err := loaders.LoadPlist(d, &scheme); err != nil {
		return nil, fmt.Errorf("Unable to load colorscheme definition: %s", err)
	}

	return &scheme, nil
}

func (t Theme) String() (ret string) {
	ret = fmt.Sprintf("%s - %s\n", t.Name, t.UUID)
	for i := range t.Settings {
		ret += fmt.Sprintf("\t%s\n", t.Settings[i])
	}
	return
}

func (c Color) String() string {
	return fmt.Sprintf("0x%02X%02X%02X%02X", c.A, c.R, c.G, c.B)
}

func (c *Color) UnmarshalJSON(data []byte) error {
	i64, err := strconv.ParseInt(string(data[2:len(data)-1]), 16, 64)
	if err != nil {
		log4go.Warn("Couldn't properly load color from %s: %s", string(data), err)
	}
	c.A = uint8((i64 >> 24) & 0xff)
	c.R = uint8((i64 >> 16) & 0xff)
	c.G = uint8((i64 >> 8) & 0xff)
	c.B = uint8((i64 >> 0) & 0xff)
	return nil
}

func (s *Settings) UnmarshalJSON(data []byte) error {
	*s = make(Settings)
	tmp := make(map[string]json.RawMessage)
	if err := json.Unmarshal(data, &tmp); err != nil {
		return err
	}
	for k, v := range tmp {
		if strings.HasPrefix(k, "font") {
			continue
		} else {
			var c Color
			if err := json.Unmarshal(v, &c); err != nil {
				return err
			}
			(*s)[k] = c
		}
	}
	return nil
}

func (t *Theme) ClosestMatchingSetting(scope string) *ScopeSetting {
	na := scope
	for len(na) > 0 {
		sn := na
		i := strings.LastIndex(sn, " ")
		if i != -1 {
			sn = sn[i+1:]
		}

		for j := range t.Settings {
			if t.Settings[j].Scope == sn {
				return &t.Settings[j]
			}
		}
		if i2 := strings.LastIndex(na, "."); i2 == -1 {
			break
		} else if i > i2 {
			na = na[:i]
		} else {
			na = strings.TrimSpace(na[:i2])
		}
	}
	return &t.Settings[0]
}

func (t *Theme) Spice(vr *render.ViewRegions) (ret render.Flavour) {
	s := t.ClosestMatchingSetting(vr.Scope)
	ret.Foreground = render.Colour(s.Settings["foreground"])
	if bg, ok := s.Settings["background"]; ok {
		ret.Background = render.Colour(bg)
	}
	return
}
