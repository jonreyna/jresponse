// Package response encapsulates structures used to marshal
// JSON responses to client requests.
package ping

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	tmpl "text/template"
	log "github.com/Sirupsen/logrus"
)

const pingTemplateString = "PING {{ .TargetHost }} ({{ .TargetIP }}): {{ .PacketSize }} data bytes\n" +
	"{{range $i, $element := .ProbeResult }}" +
	"{{.ResponseSize}} bytes from {{ .IPAddress }}: icmp_seq={{ .SequenceNumber }} ttl={{ .TimeToLive }} time=" +
	"{{ formatAsMs .RTT }} ms\n{{end}}"

func formatAsMs(ms uint) string {
	return fmt.Sprintf("%.3f",float32(ms)/1000)
}

func add(a, b int) int {
	return a + b
}

var (
	pingTemplate *tmpl.Template
)

func init() {
	fmtFuncMap := tmpl.FuncMap{"add": add, "formatAsMs": formatAsMs}

	var err error
	if pingTemplate, err = tmpl.New("pingTemplate").
		Funcs(fmtFuncMap).
		Parse(pingTemplateString); err != nil {
		log.Fatalln(err)
	}
}


type ProbeResult struct {
	DateDetermined uint    `xml:"date-determined,attr,omitempty" json:"date-determined,omitempty"`
	ProbeIndex     uint    `xml:"probe-index,omitempty"          json:"probe-index,omitempty"`
	ProbeSuccess   *string `xml:"probe-success,omitempty"        json:"probe-success,omitempty"`
	ProbeFailure   *string `xml:"probe-failure,omitempty"        json:"probe-failure,omitempty"`
	SequenceNumber uint    `xml:"sequence-number,omitempty"      json:"sequence-number,omitempty"`
	IPAddress      string  `xml:"ip-address,omitempty"           json:"ip-address,omitempty"`
    TimeToLive     uint    `xml:"time-to-live,omitempty"         json:"time-to-live,omitempty"`	
    ResponseSize   uint    `xml:"response-size,omitempty"        json:"response-size,omitempty"`	
	ProbeReached   string  `xml:"probe-reached,omitempty"        json:"probe-reached,omitempty"`
	RTT            uint    `xml:"rtt,omitempty"                  json:"rtt,omitempty"`
}

type ProbeResultsSummary struct {
	ProbesSent     uint        `xml:"probes-sent,omitempty"     json:"probes-sent,omitempty"`
	ResponsesReceived uint     `xml:"responses-received,omitempty"  json:"responses-received,omitempty"`
	PacketLoss     uint        `xml:"packet-loss,omitempty"     json:"packet-loss,omitempty"`
	RTTMinimum     uint        `xml:"rtt-minimum,omitempty"     json:"rtt-minimum,omitempty"`
	RTTMaximum     uint        `xml:"rtt-maximum,omitempty"     json:"rtt-maximum,omitempty"`
	RTTAverage     uint        `xml:"rtt-average,omitempty"     json:"rtt-average,omitempty"`
	RTTStdDev      uint        `xml:"rtt-stddev,omitempty"      json:"rtt-stddev,omitempty"`		
}

// Represents the trace route XML structure, and is used to convert it
// from XML to JSON.
type Ping struct {
	XMLName           xml.Name   `xml:"ping-results,omitempty" json:"-"`
	TargetHost        string     `xml:"target-host,omitempty"        json:"target-host,omitempty"`
	TargetIP          string     `xml:"target-ip,omitempty"          json:"target-ip,omitempty"`
	PacketSize        uint       `xml:"packet-size,omitempty"        json:"packet-size,omitempty"`
	ProbeResult  []ProbeResult   `xml:"probe-result,omitempty"       json:"probe-result,omitempty"`	
	ProbeResultsSummary ProbeResultsSummary `xml:"probe-results-summary,omitempty" json:"probe-results-summary,omitempty"`	
	Errors            []RPCError `xml:"rpc-error,omitempty"          json:"rpc-error,omitempty"`	
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

func (ping *Ping) WriteXMLTo(w io.Writer) (n int64, err error) {
	if s, err := xml.Marshal(ping); err != nil {
		return 0, err
	} else {
		return bytes.NewBuffer(s).WriteTo(w)
	}
}

func (ping *Ping) WriteJSONTo(w io.Writer) (n int64, err error) {
	if s, err := json.Marshal(ping); err != nil {
		return 0, nil
	} else {
		return bytes.NewBuffer(s).WriteTo(w)
	}
}

func (ping *Ping) WriteCLITo(w io.Writer) error {
	return pingTemplate.Execute(w, ping)
}

func (ping *Ping) ReadXMLFrom(r io.Reader) (n int64, err error) {

	buf := bytes.Buffer{}

	if n, err := buf.ReadFrom(r); err != nil {
		return n, err
	}

	pNoNewlines := bytes.Replace(buf.Bytes(), []byte("\n"), []byte(""), -1)

	if err := xml.Unmarshal(pNoNewlines, ping); err != nil {
		return 0, err
	} else {
		return n, nil
	}
}

func (ping *Ping) ReadJSONFrom(r io.Reader) (n int64, err error) {

	buf := bytes.Buffer{}

	if n, err = buf.ReadFrom(r); err != nil {
		return n, err
	} else if err = json.Unmarshal(buf.Bytes(), ping); err != nil {
		return n, err
	} else {
		return n, nil
	}
}
