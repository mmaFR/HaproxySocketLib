package HaproxySocketLib

import (
	"errors"
	"fmt"
	"strings"
)

func (h *Host) delAclMap(target any, entry any) (e error) {
	var response []byte
	var command string
	switch target.(type) {
	case *AclReport:
		command = fmt.Sprintf("del acl #%d #%s", target.(*AclReport).Id, entry.(*AclEntry).MemoryPointer)
	case *MapReport:
		command = fmt.Sprintf("del map #%d #%s", target.(*MapReport).Id, entry.(*MapEntry).MemoryPointer)
	default:
		e = errors.New("host->delAclMap: unknown target type")
	}
	if response, e = h.sendCommand(command); e == nil {
		if strings.TrimSpace(string(response)) != "" {
			e = errors.New("host->delAclMap: " + string(response))
			return
		}
	}
	return
}

func (h *Host) DelAcl(acl *AclReport, entry *AclEntry) error {
	return h.delAclMap(acl, entry)
}

func (h *Host) DelMap(Map *MapReport, entry *MapEntry) error {
	return h.delAclMap(Map, entry)
}
