package HaproxySocketLib

import (
	"errors"
	"fmt"
	"strings"
)

func (h *Host) SetMap(Map *MapReport, entry *MapEntry, value string) (e error) {
	var response []byte
	if response, e = h.sendCommand(fmt.Sprintf("set map #%d #%s %s", Map.Id, entry.MemoryPointer, value)); e == nil {
		if strings.TrimSpace(string(response)) != "" {
			e = errors.New("host->delAclMap: " + string(response))
			return
		}
	}
	return
}
