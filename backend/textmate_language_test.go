package backend

import (
	"errors"
	"io/ioutil"
	"lime/backend/textmate"
	"testing"
)

type tmp map[string]*textmate.Language

func (t tmp) GetLanguage(id string) (*textmate.Language, error) {
	if v, ok := t[id]; !ok {
		return nil, errors.New("Can't handle id " + id)
	} else {
		return v, nil
	}
}

func TestTmLanguage(t *testing.T) {
	t2 := make(tmp)

	files := []string{
		"../3rdparty/bundles/property-list.tmbundle/Syntaxes/Property List (XML).tmLanguage",
		"../3rdparty/bundles/xml.tmbundle/Syntaxes/XML.plist",
	}
	d0 := ""
	for i, fn := range files {
		if d, err := ioutil.ReadFile(fn); err != nil {
			t.Logf("Couldn't load file %s: %s", fn, err)
		} else {
			if i == 0 {
				d0 = string(d)
			}
			var l textmate.Language
			if err := LoadPlist(d, &l); err != nil {
				t.Fatal(err)
			} else {
				t2[l.ScopeName] = &l
			}
		}
	}
	d0 = `
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
	<key>fileTypes</key>
	<array>
		<string>xml</string>
		<string>xsd</string>
		<string>tld</string>
		<string>jsp</string>
		<string>pt</string>
		<string>cpt</string>
		<string>dtml</string>
		<string>rss</string>
		<string>opml</string>
	</array>
	<key>keyEquivalent</key>
	<string>^~X</string>
	<key>name</key>
	<string>XML</string>
</dict>
</plist>
`
	textmate.Provider = t2
	lp := textmate.LanguageParser{Language: t2["text.xml.plist"]}

	t.Logf("parse result: %v", lp.Parse(d0))
}
