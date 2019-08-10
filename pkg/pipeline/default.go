package pipeline

import (
	"strings"

	"github.com/EdmundMartin/trackpix/pkg/server"
)

// LangFromHeader extracts a language from browser headers
func LangFromHeader(data *server.TrackingData, dataObj *DataRecord) *DataRecord {
	val, ok := data.Headers["Accept-Language"]
	if ok && len(val) > 1 {
		dataObj.Language = strings.ToLower(val[0:2])
	}
	return dataObj
}

// MapValues passes values on across the two data types
func MapValues(data *server.TrackingData, dataObj *DataRecord) *DataRecord {
	dataObj.IPAddr = data.Remote
	dataObj.Datetime = data.Datetime
	val, _ := data.Headers["User-Agent"]
	dataObj.UserAgent = val
	return dataObj
}
