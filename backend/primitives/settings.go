package primitives

import (
	"fmt"
	"sync"
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
		HasId
		lock              sync.Mutex
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
	s.lock.Lock()
	defer s.lock.Unlock()
	return s.parent
}

func (s *Settings) SetParent(p SettingsInterface) {
	s.lock.Lock()
	defer s.lock.Unlock()
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
	s.lock.Lock()
	defer s.lock.Unlock()
	s.onChangeCallbacks[key] = cb
}

func (s *Settings) ClearOnChange(key string) {
	s.lock.Lock()
	defer s.lock.Unlock()
	delete(s.onChangeCallbacks, key)
}

func (s *Settings) Get(name string, def ...interface{}) interface{} {
	s.lock.Lock()
	defer s.lock.Unlock()
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
	s.lock.Lock()
	s.data[name] = val
	s.lock.Unlock()
	s.onChange()
}

func (s *Settings) Has(name string) bool {
	s.lock.Lock()
	defer s.lock.Unlock()
	_, ok := s.data[name]
	return ok
}

func (s *Settings) onChange() {
	for _, v := range s.onChangeCallbacks {
		v()
	}
}

func (s *Settings) Erase(name string) {
	s.lock.Lock()
	defer s.lock.Unlock()
	delete(s.data, name)
}

func (s *Settings) merge(other settingsMap) {
	s.lock.Lock()
	for k, v := range other {
		s.data[k] = v
	}
	s.lock.Unlock()
	s.onChange()
}
