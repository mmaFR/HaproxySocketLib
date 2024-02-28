package HaproxySocketLib

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
)

var prepareAclMapRegexValidation *regexp.Regexp = regexp.MustCompile(`^New\sversion\screated:\s([0-9]+)`)

func (h *Host) prepareAclMap(target any) (id uint64, e error) {
	var response []byte
	var g [][]byte
	var command string
	switch target.(type) {
	case *AclReport:
		command = fmt.Sprintf("prepare acl #%d", target.(*AclReport).Id)
	case *MapReport:
		command = fmt.Sprintf("prepare map #%d", target.(*MapReport).Id)
	default:
		e = errors.New("host->prepareAclMap: unknown target type")
		return
	}
	if response, e = h.sendCommand(command); e == nil {
		if !prepareAclMapRegexValidation.Match(response) {
			e = errors.New("host->prepareAclMap: " + string(response))
			return
		} else {
			g = prepareAclMapRegexValidation.FindSubmatch(response)
			id, _ = strconv.ParseUint(string(g[1]), 10, 64)
		}
	}
	switch target.(type) {
	case *AclReport:
		target.(*AclReport).NextVer = id
	case *MapReport:
		target.(*MapReport).NextVer = id
	}
	return
}

func (h *Host) PrepareAcl(aclTarget *AclReport) (uint64, error) {
	return h.prepareAclMap(aclTarget)
}
func (h *Host) PrepareMap(mapTarget *MapReport) (uint64, error) {
	return h.prepareAclMap(mapTarget)
}
