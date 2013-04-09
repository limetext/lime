package backend

import (
	"fmt"
	"lime/backend/primitives"
)

type (
	HasSettings struct {
		settings Settings
	}
	SettingsInterface interface {
		Settings() *Settings
	}
	OnChangeCallback func()
	settingsMap      map[string]interface{}
	Settings         struct {
		primitives.HasId
		onChangeCallbacks map[string]OnChangeCallback
		data              settingsMap
		parent            SettingsInterface
	}
)

func (s *HasSettings) Settings() *Settings {
	if s.settings.data == nil {
		s.settings = NewSettings()
	}
	return &s.settings
}

func NewSettings() Settings {
	return Settings{onChangeCallbacks: make(map[string]OnChangeCallback), data: make(settingsMap), parent: nil}
}

func (s *Settings) Parent() SettingsInterface {
	return s.parent
}

func (s *Settings) SetParent(p SettingsInterface) {
	if s.parent != nil {
		old := s.parent.Settings()
		old.ClearOnChange(fmt.Sprintf("lime.child.%d", s.Id()))
	}
	s.parent = p

	if s.parent != nil {
		ns := s.parent.Settings()
		ns.AddOnChange(fmt.Sprintf("lime.child.%d", s.Id()), s.onChange)
	}
}

func (s *Settings) AddOnChange(key string, cb OnChangeCallback) {
	s.onChangeCallbacks[key] = cb
}

func (s *Settings) ClearOnChange(key string) {
	delete(s.onChangeCallbacks, key)
}

func (s *Settings) Get(name string, def ...interface{}) interface{} {
	if v, ok := s.data[name]; ok {
		return v
	} else if s.parent != nil {
		return s.parent.Settings().Get(name, def...)
	} else if len(def) > 0 {
		return def[0]
	}
	return nil
}

func (s *Settings) Set(name string, val interface{}) {
	s.data[name] = val
	s.onChange()
}

func (s *Settings) onChange() {
	for _, v := range s.onChangeCallbacks {
		v()
	}
}

func (s *Settings) Erase(name string) {
	s.data[name] = nil
}

func (s *Settings) merge(other settingsMap) {
	for k, v := range other {
		s.data[k] = v
	}
	s.onChange()
}
