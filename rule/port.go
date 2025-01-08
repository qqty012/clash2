package rules

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	C "github.com/Dreamacro/clash/constant"
)

type Port struct {
	adapter  string
	port     string
	isSource bool
	compare  string
}

func (p *Port) RuleType() C.RuleType {
	if p.isSource {
		return C.SrcPort
	}
	return C.DstPort
}

func _comparePort(p *Port, port string) bool {
	sport, serr := strconv.ParseUint(port, 10, 16)
	if serr != nil {
		return false
	}
	rport, rerr := strconv.ParseUint(p.port, 10, 16)
	if rerr != nil {
		return false
	}
	switch p.compare {
	case ">", "GT":
		return sport > rport
	case ">=", "GTE":
		return sport >= rport
	case "<", "LT":
		return sport < rport
	case "<=", "LTE":
		return sport <= rport
	case "!=", "NOT", "<>":
		return sport != rport
	default:
		return sport == rport
	}
}

func (p *Port) Match(metadata *C.Metadata) bool {
	if p.isSource {
		return _comparePort(p, metadata.SrcPort)
	}
	return _comparePort(p, metadata.DstPort)
}

func (p *Port) Adapter() string {
	return p.adapter
}

func (p *Port) Payload() string {
	return fmt.Sprintf("%s %s", p.compare, p.port)
}

func (p *Port) ShouldResolveIP() bool {
	return false
}

func (p *Port) ShouldFindProcess() bool {
	return false
}

func NewPort(port string, adapter string, isSource bool, params []string) (*Port, error) {

	re, err := regexp.Compile("[0-9]+")
	if err != nil {
		return nil, err
	}
	_port := re.FindString(port)
	compare := ""

	if len(params) > 0 {
		compare = strings.Trim(params[0], " ")
	}
	compare = strings.Replace(port, _port, "", 1)
	if compare == "" {
		compare = "="
	}

	_, err = strconv.ParseUint(_port, 10, 16)
	if err != nil {
		return nil, errPayload
	}
	return &Port{
		adapter:  adapter,
		port:     _port,
		isSource: isSource,
		compare:  strings.Trim(compare, " "),
	}, nil
}
