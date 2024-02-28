package HaproxySocketLib

import (
	"errors"
	"fmt"
	"strings"
)

func (h *Host) AddAcl(acl *AclReport, value string) (e error) {
	var rsp []byte
	var rspString string
	var command string = fmt.Sprintf("add acl #%d %s", acl.Id, value)
	if acl.CurrVer != acl.NextVer {
		command = fmt.Sprintf("add acl @%d #%d %s", acl.NextVer, acl.Id, value)
	}
	if rsp, e = h.sendCommand(command); e == nil {
		rspString = string(rsp)
		rspString = strings.TrimSpace(rspString)
		rspString = strings.Trim(rspString, ">")
		if rspString != "" {
			e = errors.New("host->AddAcl: " + string(rsp))
			return
		}
	}
	return
}
