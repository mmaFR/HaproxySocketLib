package HaproxySocketLib

import (
	"errors"
	"fmt"
	"strings"
)

func (h *Host) clearAclMap(target any) (e error) {
	var response []byte
	var command string
	switch target.(type) {
	case *AclReport:
		if target.(*AclReport).CurrVer != target.(*AclReport).NextVer {
			command = fmt.Sprintf("clear acl @%d #%d", target.(*AclReport).NextVer, target.(*AclReport).Id)
		} else {
			command = fmt.Sprintf("clear acl #%d", target.(*AclReport).Id)
		}
	case *MapReport:
		if target.(*MapReport).CurrVer != target.(*MapReport).NextVer {
			command = fmt.Sprintf("clear map @%d #%d", target.(*MapReport).NextVer, target.(*MapReport).Id)
		} else {
			command = fmt.Sprintf("clear map #%d", target.(*MapReport).Id)
		}
	default:
		e = errors.New("host->clearAclMap: unknown target type")
		return
	}
	if response, e = h.sendCommand(command); e == nil {
		if strings.TrimSpace(string(response)) != "" {
			e = errors.New("host->clearAclMap: " + string(response))
			return
		}
	}
	return
}

func (h *Host) ClearAcl(acl *AclReport) error {
	return h.clearAclMap(acl)
}

func (h *Host) ClearMap(Map *MapReport) error {
	return h.clearAclMap(Map)
}
