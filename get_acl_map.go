package HaproxySocketLib

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

var getAclRegexes []*regexp.Regexp = []*regexp.Regexp{
	regexp.MustCompile(`^type=([a-z]+),\scase=([a-z]+),\smatch=([a-z]+),\sidx=([a-z]+),\spattern="(.+?)"`),
	regexp.MustCompile(`^type=([a-z]+),\scase=([a-z]+),\smatch=([a-z]+)`),
}

var getMapRegexes []*regexp.Regexp = []*regexp.Regexp{
	regexp.MustCompile(`^type=([a-z]+),\scase=([a-z]+),\sfound=([a-z]+),\sidx=([a-z]+),\skey="(.+?)",\svalue="(.+?)",\stype="(.+?)"`),
	regexp.MustCompile(`^type=([a-z]+),\scase=([a-z]+),\sfound=([a-z]+)`),
}

type GetAclResult struct {
	MatchingType string
	Case         string
	Match        bool
	Idx          string
	Pattern      string
}

type GetMapResult struct {
	MatchingType string
	Case         string
	Found        bool
	Idx          string
	Key          string
	Value        string
	ReturnType   string
}

func (h *Host) getAclMap(target any, value string) (res any, e error) {
	var response []byte
	var command string
	var g [][]byte
	switch target.(type) {
	case *AclReport:
		command = fmt.Sprintf("get acl #%d %s", target.(*AclReport).Id, value)
	case *MapReport:
		command = fmt.Sprintf("get map #%d", target.(*MapReport).Id, value)
	default:
		e = errors.New("host->getAclMap: unknown target type")
		return
	}
	if response, e = h.sendCommand(command); e == nil {
		if strings.TrimSpace(string(response)) != "" {
			e = errors.New("host->getAclMap: " + string(response))
			return
		}
		switch target.(type) {
		case *AclReport:
			for i, r := range getAclRegexes {
				g = r.FindSubmatch(response)
				if r.Match(response) {
					switch i {
					case 0:
						res = &GetAclResult{
							MatchingType: string(g[1]),
							Case:         string(g[2]),
							Match:        string(g[3]) == "yes",
							Idx:          string(g[4]),
							Pattern:      string(g[5]),
						}
					case 1:
						res = &GetAclResult{
							MatchingType: string(g[1]),
							Case:         string(g[2]),
							Match:        string(g[3]) == "yes",
						}
					}
				}
			}
		case *MapReport:
			for i, r := range getMapRegexes {
				g = r.FindSubmatch(response)
				if r.Match(response) {
					switch i {
					case 0:
						res = &GetMapResult{
							MatchingType: string(g[1]),
							Case:         string(g[2]),
							Found:        string(g[3]) == "yes",
							Idx:          string(g[4]),
							Key:          string(g[5]),
							Value:        string(g[6]),
							ReturnType:   string(g[7]),
						}
					case 1:
						res = &GetMapResult{
							MatchingType: string(g[1]),
							Case:         string(g[2]),
							Found:        string(g[3]) == "yes",
						}
					}
				}
			}
		}
	}
	return
}

func (h *Host) GetAcl(acl *AclReport, value string) (r *GetAclResult, e error) {
	var res any
	res, e = h.getAclMap(acl, value)
	r = res.(*GetAclResult)
	return
}
func (h *Host) GetMap(Map *MapReport, key string) (r *GetMapResult, e error) {
	var res any
	res, e = h.getAclMap(Map, key)
	r = res.(*GetMapResult)
	return
}
