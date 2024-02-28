package HaproxySocketLib

import (
	"bufio"
	"bytes"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

var mapRegexes []*regexp.Regexp = []*regexp.Regexp{
	regexp.MustCompile(`^([0-9]+)\s\((.+?)\)`),
	regexp.MustCompile(`(\sby\smap\sat\sfile\s'(.+?)'\sline\s([0-9]+))+`),
	regexp.MustCompile(`curr_ver=([0-9]+)\snext_ver=([0-9]+)\sentry_cnt=([0-9]+)$`),
}
var mapRegexValidation *regexp.Regexp = regexp.MustCompile(`^([0-9]+)\s\((.+?)\)\spattern\sloaded`)

type MapEntry struct {
	MemoryPointer string
	Key           string
	Value         string
}

func newMapEntry(b []byte) *MapEntry {
	var x []string = strings.Split(string(b), " ")
	return &MapEntry{
		MemoryPointer: x[0],
		Key:           x[1],
		Value:         x[2],
	}
}

type MapContent []*MapEntry

type PositionInConfig struct {
	File string
	Line uint64
}

type MapReport struct {
	Id       uint64
	MapFile  string
	UsedHere []*PositionInConfig
	CurrVer  uint64
	NextVer  uint64
	EntryCnt uint64
}

func NewMapReport(b []byte) *MapReport {
	var m *MapReport = new(MapReport)
	m.UsedHere = make([]*PositionInConfig, 0)
	var g2 [][]byte
	var g3 [][][]byte
	var line uint64
	for i, r := range mapRegexes {
		switch i {
		case 0:
			g2 = r.FindSubmatch(b)
			m.Id, _ = strconv.ParseUint(string(g2[1]), 10, 64)
			m.MapFile = string(g2[2])
		case 1:
			g3 = r.FindAllSubmatch(b, -1)
			for _, g := range g3 {
				line, _ = strconv.ParseUint(string(g[3]), 10, 64)
				m.UsedHere = append(m.UsedHere, &PositionInConfig{File: string(g[2]), Line: line})
			}
		case 2:
			g2 = r.FindSubmatch(b)
			m.CurrVer, _ = strconv.ParseUint(string(g2[1]), 10, 64)
			m.NextVer, _ = strconv.ParseUint(string(g2[2]), 10, 64)
			m.EntryCnt, _ = strconv.ParseUint(string(g2[3]), 10, 64)
		}
	}
	return m
}

type MapList []*MapReport

func (h *Host) ShowMap() (ml MapList, e error) {
	var response []byte
	if response, e = h.sendCommand("show map"); e == nil {
		ml = make(MapList, 0)
		var scanner *bufio.Scanner = bufio.NewScanner(bytes.NewReader(response))
		for scanner.Scan() {
			if mapRegexValidation.Match(scanner.Bytes()) {
				ml = append(ml, NewMapReport(scanner.Bytes()))
			}
		}
	}
	return
}

func (h *Host) ShowMapContent(mr *MapReport) (mc MapContent, e error) {
	mc = make(MapContent, mr.EntryCnt)
	var response []byte
	if response, e = h.sendCommand(fmt.Sprintf("show map #%d", mr.Id)); e == nil {
		var scanner *bufio.Scanner = bufio.NewScanner(bytes.NewReader(response))
		for i := 0; scanner.Scan(); i++ {
			if strings.HasPrefix(scanner.Text(), "0x") {
				mc[i] = newMapEntry(scanner.Bytes())
			}
		}
	}
	return
}
