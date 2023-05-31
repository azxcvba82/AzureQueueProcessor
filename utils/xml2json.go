package utils

import (
	"strings"

	xj "github.com/basgys/goxml2json"
)

func XML2JSON(xmlString string) string {
	xml := strings.NewReader(xmlString)
	json, err := xj.Convert(xml)
	if nil != err {
		return ""
	}
	return json.String()
}
