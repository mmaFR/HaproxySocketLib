package HaproxySocketLib

import (
	"errors"
	"fmt"
	"strings"
)

func (h *Host) commitAclMap(target any) (e error) {
	var rsp []byte
	var rspString string
	var command string
	switch target.(type) {
	case *AclReport:
		command = fmt.Sprintf("commit acl @%d #%d", target.(*AclReport).NextVer, target.(*AclReport).Id)
	case *MapReport:
		command = fmt.Sprintf("commit map @%d #%d", target.(*MapReport).NextVer, target.(*MapReport).Id)
	default:
		e = errors.New("host->commitAclMap: unknown target type")
		return
	}
	if rsp, e = h.sendCommand(command); e == nil {
		rspString = string(rsp)
		rspString = strings.TrimSpace(rspString)
		rspString = strings.Trim(rspString, ">")
		if rspString != "" {
			e = errors.New("host->commitAclMap: " + string(rsp))
			return
		} else {
			switch target.(type) {
			case *AclReport:
				target.(*AclReport).CurrVer = target.(*AclReport).NextVer
			case *MapReport:
				target.(*MapReport).CurrVer = target.(*MapReport).NextVer
			}
		}
	}
	return
}

func (h *Host) CommitAcl(aclTarget *AclReport) error {
	return h.commitAclMap(aclTarget)
}
func (h *Host) CommitMap(mapTarget *MapReport) error {
	return h.commitAclMap(mapTarget)
}
