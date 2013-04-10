package backend

import (
	"code.google.com/p/log4go"
	"strings"
)

type (
	ViewEventCallback func(v *View)
	ViewEvent         []ViewEventCallback

	QueryContextReturn   int
	QueryContextCallback func(v *View, key string, operator Op, operand interface{}, match_all bool) QueryContextReturn
	QueryContextEvent    []QueryContextCallback
)

const (
	True QueryContextReturn = iota
	False
	Unknown
)

func (ve *ViewEvent) Add(cb ViewEventCallback) {
	*ve = append(*ve, cb)
}

func (ve ViewEvent) Call(v *View) {
	for i := range ve {
		ve[i](v)
	}
}

func (qe *QueryContextEvent) Add(cb QueryContextCallback) {
	*qe = append(*qe, cb)
}

func (qe QueryContextEvent) Call(v *View, key string, operator Op, operand interface{}, match_all bool) QueryContextReturn {
	for i := range qe {
		r := qe[i](v, key, operator, operand, match_all)
		if r != Unknown {
			return r
		}
	}
	log4go.Fine("Unknown context: %s", key)
	return Unknown
}

var (
	OnNew               ViewEvent
	OnLoad              ViewEvent
	OnActivated         ViewEvent
	OnDeactivated       ViewEvent
	OnClose             ViewEvent
	OnPreSave           ViewEvent
	OnPostSave          ViewEvent
	OnModified          ViewEvent
	OnSelectionModified ViewEvent
	OnQueryContext      QueryContextEvent
)

func init() {
	OnQueryContext.Add(func(v *View, key string, operator Op, operand interface{}, match_all bool) QueryContextReturn {
		if strings.HasPrefix(key, "setting.") && operator == OpEqual {
			c, ok := v.Settings().Get(key[8:]).(bool)
			if c && ok {
				return True
			}
			return False
		}
		return Unknown
	})
}
