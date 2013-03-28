package backend

type HasSettings struct {
	settings Settings
}

func (s *HasSettings) Settings() *Settings {
	if s.settings.data == nil {
		s.settings = NewSettings()
	}
	return &s.settings
}

type OnChangeCallback func()
type settingsMap map[string]interface{}
type Settings struct {
	onChangeCallbacks map[string]OnChangeCallback
	data              settingsMap
}

func NewSettings() Settings {
	return Settings{make(map[string]OnChangeCallback), make(settingsMap)}
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
	}
	return def
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
