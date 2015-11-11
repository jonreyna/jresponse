// Package response encapsulates structures used to marshal
// JSON responses to client requests.
package traceroute

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	tmpl "text/template"

	log "github.com/Sirupsen/logrus"
	"strings"
)

const traceRouteTmplStr = "traceroute to {{ .TargetHost }} " +
	"({{ .TargetIP }}), {{ .MaxHopIndex }} hops max, {{ .PacketSize }} byte packets\n" +

	"{{range $i, $element := .Hops }} " +
	"{{.TTLValue}}  {{ .TrimmedLastHostName }} ({{ .LastIPAddr }})  " +
	"{{ range $_, $probeResult := .ProbeResult }}{{ formatAsMs $probeResult.RTT }} ms  " +
	"{{end}}\n{{end}}"

func formatAsMs(ms uint) string {
	return fmt.Sprintf("%g", float32(ms)/1000)
}

func add(a, b int) int {
	return a + b
}

var (
	traceRouteTempl *tmpl.Template
)

func init() {
	fmtFuncMap := tmpl.FuncMap{"add": add, "formatAsMs": formatAsMs}

	var err error
	if traceRouteTempl, err = tmpl.New("traceRouteTmpl").
		Funcs(fmtFuncMap).
		Parse(traceRouteTmplStr); err != nil {
		log.Fatalln(err)
	}
}

type ICMPCode struct {
	IntegerCodeValue uint   `xml:"integer-code-value,omitempty"    json:"integer-code-value,omitempty"`
	ICMPTimxceed     string `xml:"icmp-timxceed-intrans,omitempty" json:"icmp-timxceed-intrans,omitempty"`
	ICMPUnreachPort  string `xml:"icmp-unreach-port,omitempty"     json:"icmp-unreach-port,omitempty"`
}

type ICMPType struct {
	IntegerTypeValue uint   `xml:"integer-type-value,attr,omitempty" json:"integer-type-value,omitempty"`
	ICMPTimxceed     string `xml:"icmp-timxceed,omitempty"           json:"icmp-timxceed,omitempty"`
	ICMPUnreach      string `xml:"icmp-unreach,omitempty"            json:"icpm-unreach,omitempty"`
}

type ProbeResult struct {
	DateDetermined uint    `xml:"date-determined,attr,omitempty" json:"date-determined,omitempty"`
	ProbeIndex     uint    `xml:"probe-index,omitempty"          json:"probe-index,omitempty"`
	IPAddress      string  `xml:"ip-address,omitempty"           json:"ip-address,omitempty"`
	HostName       string  `xml:"host-name,omitempty"            json:"host-name,omitempty"`
	ProbeSuccess   *string `xml:"probe-success,omitempty"        json:"probe-success,omitempty"`
	ProbeFailure   *string `xml:"probe-failure,omitempty"        json:"probe-failure,omitempty"`
	ProbeReached   string  `xml:"probe-reached,omitempty"        json:"probe-reached,omitempty"`
	RTT            uint    `xml:"rtt,omitempty"                  json:"rtt,omitempty"`
}

type Hop struct {
	TTLValue     uint          `xml:"ttl-value,omitempty"       json:"ttl-value,omitempty"`
	LastIPAddr   string        `xml:"last-ip-address,omitempty" json:"last-ip-address,omitempty"`
	LastHostName string        `xml:"last-host-name,omitempty"  json:"last-host-name,omitempty"`
	ProbeResult  []ProbeResult `xml:"probe-result,omitempty"    json:"probe-result,omitempty"`
}

func (h *Hop) TrimmedLastHostName() string {
	return strings.TrimSpace(h.LastHostName)
}

// Represents the trace route XML structure, and is used to convert it
// from XML to JSON.
type TraceRoute struct {
	XMLName           xml.Name   `xml:"traceroute-results,omitempty" json:"-"`
	TargetHost        string     `xml:"target-host,omitempty"        json:"target-host,omitempty"`
	TargetIP          string     `xml:"target-ip,omitempty"          json:"target-ip,omitempty"`
	MaxHopIndex       uint       `xml:"max-hop-index,omitempty"      json:"max-hop-index,omitempty"`
	PacketSize        uint       `xml:"packet-size,omitempty"        json:"packet-size,omitempty"`
	Hops              []Hop      `xml:"hop,omitempty"                json:"hop,omitempty"`
	Errors            []RPCError `xml:"rpc-error,omitempty"          json:"rpc-error,omitempty"`
	TraceRouteFailure string     `xml:"traceroute-failure,omitempty" json:"traceroute-failure,omitempty"`
	OriginHost        string `json:"originhost,omitempty"` 
	OriginIP          string `json:"originip,omitempty"`		
}

type RPCError struct {
	Type     string `xml:"error-type"     json:"error-type"`
	Tag      string `xml:"error-tag"      json:"error-tag"`
	Severity string `xml:"error-severity" json:"error-severity"`
	Path     string `xml:"error-path"     json:"error-path"`
	Message  string `xml:"error-message"  json:"error-message"`
	Info     string `xml:",innerxml"      json:",string`
}

func (traceRoute *TraceRoute) WriteXMLTo(w io.Writer) (n int64, err error) {
	if s, err := xml.Marshal(traceRoute); err != nil {
		return 0, err
	} else {
		return bytes.NewBuffer(s).WriteTo(w)
	}
}

func (traceRoute *TraceRoute) WriteJSONTo(w io.Writer) (n int64, err error) {
	if s, err := json.Marshal(traceRoute); err != nil {
		return 0, nil
	} else {
		return bytes.NewBuffer(s).WriteTo(w)
	}
}

func (traceRoute *TraceRoute) WriteCLITo(w io.Writer) error {
	return traceRouteTempl.Execute(w, traceRoute)
}

func (traceRoute *TraceRoute) ReadXMLFrom(r io.Reader) (n int64, err error) {

	buf := bytes.Buffer{}

	if n, err := buf.ReadFrom(r); err != nil {
		return n, err
	}

	pNoNewlines := bytes.Replace(buf.Bytes(), []byte("\n"), []byte(""), -1)

	if err := xml.Unmarshal(pNoNewlines, traceRoute); err != nil {
		return 0, err
	} else {
		return n, nil
	}
}

func (traceRoute *TraceRoute) ReadJSONFrom(r io.Reader) (n int64, err error) {

	buf := bytes.Buffer{}

	if n, err = buf.ReadFrom(r); err != nil {
		return n, err
	} else if err = json.Unmarshal(buf.Bytes(), traceRoute); err != nil {
		return n, err
	} else {
		return n, nil
	}
}
