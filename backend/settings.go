package backend

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
		onChangeCallbacks map[string]OnChangeCallback
		data              settingsMap
		Parent            SettingsInterface
	}
)

func (s *HasSettings) Settings() *Settings {
	if s.settings.data == nil {
		s.settings = NewSettings()
	}
	return &s.settings
}

func NewSettings() Settings {
	return Settings{make(map[string]OnChangeCallback), make(settingsMap), nil}
}

func (s *Settings) AddOnChange(key string, cb OnChangeCallback) {
	s.onChangeCallbacks[key] = cb
}

func (s *Settings) ClearOnChange(key string) {
	s.onChangeCallbacks[key] = nil
}

func (s *Settings) Get(name string, def ...interface{}) interface{} {
	if v, ok := s.data[name]; ok {
		return v
	} else if s.Parent != nil {
		return s.Parent.Settings().Get(name, def...)
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
