package HaproxySocketLib

import (
	"errors"
	"fmt"
	"strings"
)

func (h *Host) AddMap(Map *MapReport, key string, value string) (e error) {
	var rsp []byte
	var rspString string
	var command string = fmt.Sprintf("add map #%d %s %s", Map.Id, key, value)
	if Map.CurrVer != Map.NextVer {
		command = fmt.Sprintf("add map @%d #%d %s %s", Map.NextVer, Map.Id, key, value)
	}
	if rsp, e = h.sendCommand(command); e == nil {
		rspString = string(rsp)
		rspString = strings.TrimSpace(rspString)
		rspString = strings.Trim(rspString, ">")
		if rspString != "" {
			e = errors.New("host->AddMap: " + string(rsp))
			return
		}
	}
	return
}
