package HaproxySocketLib

import (
	"bufio"
	"bytes"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

var aclRegex *regexp.Regexp = regexp.MustCompile(`^(?P<id>[0-9]+)\s\(\)\sacl\s'(?P<fetch>.+?)'\sfile\s'(?P<config_file>.+)'\sline\s(?P<line>[0-9]+)\.\scurr_ver=(?P<curr_ver>[0-9]+)\snext_ver=(?P<next_ver>[0-9]+)\sentry_cnt=(?P<entry_cnt>[0-9]+)`)
var aclFileRegex *regexp.Regexp = regexp.MustCompile(`^(?P<id>[0-9]+)\s\((?P<acl_file>.+?)\)\spattern\sloaded\sfrom\sfile\s'.+?'\sused\sby\sacl\sat\sfile\s'(?P<config_file>.+?)'\sline\s(?P<line>[0-9]+)\.\scurr_ver=(?P<curr_ver>[0-9]+)\snext_ver=(?P<next_ver>[0-9]+)\sentry_cnt=(?P<entry_cnt>[0-9]+)`)

type AclEntry struct {
	MemoryPointer string
	Value         string
}

func newAclEntry(b []byte) *AclEntry {
	var x []string = strings.Split(string(b), " ")
	return &AclEntry{
		MemoryPointer: x[0],
		Value:         x[1],
	}
}

type AclContent []*AclEntry

type AclReport struct {
	Id         uint64
	AclFile    string
	Fetch      string
	ConfigFile string
	Line       uint64
	CurrVer    uint64
	NextVer    uint64
	EntryCnt   uint64
}

func NewAclReport(b []byte) *AclReport {
	var a *AclReport = new(AclReport)
	var g [][]byte
	g = aclRegex.FindSubmatch(b)
	a.Id, _ = strconv.ParseUint(string(g[1]), 10, 64)
	a.Fetch = string(g[2])
	a.ConfigFile = string(g[3])
	a.Line, _ = strconv.ParseUint(string(g[4]), 10, 64)
	a.CurrVer, _ = strconv.ParseUint(string(g[5]), 10, 64)
	a.NextVer, _ = strconv.ParseUint(string(g[6]), 10, 64)
	a.EntryCnt, _ = strconv.ParseUint(string(g[7]), 10, 64)
	return a
}
func NewAclFileReport(b []byte) *AclReport {
	var a *AclReport = new(AclReport)
	var g [][]byte
	g = aclFileRegex.FindSubmatch(b)
	a.Id, _ = strconv.ParseUint(string(g[1]), 10, 64)
	a.AclFile = string(g[2])
	a.ConfigFile = string(g[3])
	a.Line, _ = strconv.ParseUint(string(g[4]), 10, 64)
	a.CurrVer, _ = strconv.ParseUint(string(g[5]), 10, 64)
	a.NextVer, _ = strconv.ParseUint(string(g[6]), 10, 64)
	a.EntryCnt, _ = strconv.ParseUint(string(g[7]), 10, 64)
	return a
}

type AclList []*AclReport

func (h *Host) ShowAcl() (al AclList, e error) {
	var response []byte
	if response, e = h.sendCommand("show acl"); e == nil {
		al = make(AclList, 0)
		var scanner *bufio.Scanner = bufio.NewScanner(bytes.NewReader(response))
		for scanner.Scan() {
			switch {
			case aclRegex.Match(scanner.Bytes()):
				al = append(al, NewAclReport(scanner.Bytes()))
			case aclFileRegex.Match(scanner.Bytes()):
				al = append(al, NewAclFileReport(scanner.Bytes()))
			}
		}
	}
	return
}

func (h *Host) ShowAclContent(ar *AclReport, version ...uint64) (ac AclContent, e error) {
	ac = make(AclContent, ar.EntryCnt)
	var response []byte
	var command string = fmt.Sprintf("show acl #%d", ar.Id)
	if len(version) > 0 {
		command = fmt.Sprintf("show acl @%d #%d", version[0], ar.Id)
	}
	if response, e = h.sendCommand(command); e == nil {
		var scanner *bufio.Scanner = bufio.NewScanner(bytes.NewReader(response))
		for i := 0; scanner.Scan(); i++ {
			if strings.HasPrefix(scanner.Text(), "0x") {
				ac[i] = newAclEntry(scanner.Bytes())
			}
		}
	}
	return
}
