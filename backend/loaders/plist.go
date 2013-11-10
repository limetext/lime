package loaders

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/quarnster/parser"
	"lime/backend/loaders/plist"
	"strings"
)

func plistconv(buf *bytes.Buffer, node *parser.Node) error {
	switch node.Name {
	case "Key":
		buf.WriteString("\"" + node.Data() + "\": ")
	case "String":
		n := node.Data()
		n = strings.Replace(n, "\\", "\\\\", -1)
		n = strings.Replace(n, "\"", "\\\"", -1)
		n = strings.Replace(n, "\n", "\\n", -1)
		n = strings.Replace(n, "\t", "\\t", -1)
		n = strings.Replace(n, "&gt;", ">", -1)
		n = strings.Replace(n, "&lt;", "<", -1)
		buf.WriteString("\"" + n + "\"")
	case "Dictionary":
		buf.WriteString("{\n\t")
		for i, child := range node.Children {
			if i != 0 && i&1 == 0 {
				buf.WriteString(",\n\t")
			}
			if err := plistconv(buf, child); err != nil {
				return err
			}
		}

		buf.WriteString("}\n")
	case "Array":
		buf.WriteString("[\n\t")
		for i, child := range node.Children {
			if i != 0 {
				buf.WriteString(",\n\t")
			}

			if err := plistconv(buf, child); err != nil {
				return err
			}
		}

		buf.WriteString("]\n\t")
	case "EndOfFile":
	default:
		return errors.New(fmt.Sprintf("Unhandled node: %s", node.Name))
	}
	return nil
}

func LoadPlist(data []byte, intf interface{}) error {
	var (
		p plist.PLIST
	)
	if !p.Parse(strings.Replace(string(data), "\r", "", -1)) {
		return p.Error()
	} else {
		var (
			root = p.RootNode()
			buf  bytes.Buffer
		)
		for _, child := range root.Children {
			if err := plistconv(&buf, child); err != nil {
				return err
			}
		}
		return LoadJSON(buf.Bytes(), intf)
	}
}
