package bgproute

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"io"
	tmpl "text/template"

	log "github.com/Sirupsen/logrus"
)

var (
	bgpRouteTmpl *tmpl.Template
)

func init() {
	cntxLog := log.WithFields(log.Fields{
		"func": "response.init()",
	})

	fmtFuncMap := tmpl.FuncMap{"isNextHop": isNextHop}

	var err error
	if bgpRouteTmpl, err = tmpl.
		New("bgpRouteTmpl").
		Funcs(fmtFuncMap).
		Parse(bgpRouteTmplStr); err != nil {
		cntxLog.Fatalln(err)
	}
}

const bgpRouteTmplStr = "{{.RouteTable.TableName}}: {{.RouteTable.DestinationCount}} destinations, " +
	"{{.RouteTable.TotalRouteCount}} routes ({{.RouteTable.ActiveRouteCount}}, " +
	"{{.RouteTable.HoldDownRouteCount}} holddown, {{.RouteTable.HiddenRouteCount}} hidden)\n" +
	"@ = Routing Use Only, # = Forwarding Use Only\n" +
	"+ = Active Route, - = Last Active, * = Both\n\n" +

	"{{range $_, $rt := .RouteTable.RT}}" +
	"{{$rt.RTDestination}}" +

	"{{range $i, $rtEntry := $rt.RTEntry}}" +
	"{{if eq $i 0}}      {{else}}                {{end}}{{$rtEntry.ActiveTag}}" +
	"[{{$rtEntry.ProtocolName}}/{{$rtEntry.Preference}}] {{$rtEntry.Age.AgeTime}}, " +
	"MED {{$rtEntry.Med}}, localpref {{$rtEntry.LocalPreference}}, " +
	"from {{$rtEntry.LearnedFrom}}\n" +
	"                  AS path: {{$rtEntry.AsPath}}, " +
	"validation-state: {{$rtEntry.ValidationState}}\n" +

	"{{range $_, $nh := $rtEntry.NH}}" +
	"                {{isNextHop $nh.SelectedNextHop}} to {{$nh.To}} via {{$nh.Via}}" +
	"{{if eq $nh.LSPName \"\"}}{{else}}, label-switched-path {{$nh.LSPName}}{{end}}\n" +

	"{{end}}{{end}}{{end}}"

func isNextHop(nextHopInd *string) string {
	switch nextHopInd {
	case nil:
		return " "
	default:
		return ">"
	}
}

type NH struct {
	// <selected-next-hop> is either present as an empty tag, or not present
	// so we need a pointer to distinguish its presence (not nil), or lack thereof (nil)
	SelectedNextHop *string `xml:"selected-next-hop"  json:"selected-next-hop"`
	To              string  `xml:"to,omitempty"       json:"to,omitempty"`
	Via             string  `xml:"via,omitempty"      json:"via,omitempty"`
	LSPName         string  `xml:"lsp-name,omitempty" json:"lsp-name,omitempty"`
}

type Age struct {
	AgeSecs string `xml:"seconds,attr" json:"age-seconds,omitempty"`
	AgeTime string `xml:",chardata"    json:"age,omitempty"`
}

type RTEntry struct {
	ActiveTag       string `xml:"active-tag,omitempty"       json:"active-date,omitempty"`
	CurrentActive   string `xml:"current-active,omitempty"   json:"current-active,omitempty"`
	LastActive      string `xml:"last-active,omitempty"      json:"last-active,omitempty"`
	ProtocolName    string `xml:"protocol-name,omitempty"    json:"protocol-name,omitempty"`
	Preference      int    `xml:"preference,omitempty"       json:"preference,omitempty"`
	Age             Age    `xml:"age,omitempty"              json:"age,omitempty"`
	Med             int    `xml:"med,omitempty"              json:"med,omitempty"`
	LocalPreference int    `xml:"local-preference,omitempty" json:"local-preference,omitempty"`
	LearnedFrom     string `xml:"learned-from,omitempty"     json:"learned-from,omitempty"`
	AsPath          string `xml:"as-path,omitempty"          json:"as-path,omitempty"`
	ValidationState string `xml:"validation-state,omitempty" json:"validation-state,omitempty"`
	NH              []NH   `xml:"nh,omitempty"               json:"nh,omitempty"`
}

type RT struct {
	RTDestination string    `xml:"rt-destination,omitempty" json:"rt-destination,omitempty"`
	RTEntry       []RTEntry `xml:"rt-entry,omitempty"       json:"rt-entry,omitempty"`
}

type RouteTable struct {
	TableName          string `xml:"table-name,omitempty"           json:"table-name,omitempty"`
	DestinationCount   int    `xml:"destination-count,omitempty"    json:"destination-count,omitempty"`
	TotalRouteCount    int    `xml:"total-route-count,omitempty"    json:"total-route-count,omitempty"`
	ActiveRouteCount   int    `xml:"active-route-count,omitempty"   json:"active-route-count,omitempty"`
	HoldDownRouteCount int    `xml:"holddown-route-count"           json:"holddown-route-count"`
	HiddenRouteCount   int    `xml:"hidden-route-count,omitempty"   json:"hidden-route-count,omitempty"`
	RT                 []RT   `xml:"rt,omitempty"                   json:"rt,omitempty"`
}

// Represents the BGP route XML structure, and is used to convert it
// from XML to JSON.
type BGPRoute struct {
	XMLName    xml.Name   `xml:"route-information"     json:"-"`
	RouteTable RouteTable `xml:"route-table,omitempty" json:"route-table,omitempty"`
	OriginHost        string `json:"originhost,omitempty"`
	OriginIP          string `json:"originip,omitempty"`	
}

func (bgpRoute *BGPRoute) WriteXMLTo(w io.Writer) (n int64, err error) {
	if s, err := xml.Marshal(bgpRoute); err != nil {
		return 0, err
	} else {
		buf := bytes.NewBuffer(s)
		return buf.WriteTo(w)
	}
}

func (bgpRoute *BGPRoute) WriteJSONTo(w io.Writer) (n int64, err error) {
	if s, err := json.Marshal(bgpRoute); err != nil {
		return 0, nil
	} else {
		buf := bytes.NewBuffer(s)
		return buf.WriteTo(w)
	}
}

func (bgpRoute *BGPRoute) WriteCLITo(w io.Writer) error {
	return bgpRouteTmpl.Execute(w, bgpRoute)
}

func (bgpRoute *BGPRoute) ReadXMLFrom(r io.Reader) (n int64, err error) {

	buf := bytes.Buffer{}

	if n, err := buf.ReadFrom(r); err != nil {
		return n, err
	}

	pNoNewlines := bytes.Replace(buf.Bytes(), []byte("\n"), []byte(""), -1)

	if err := xml.Unmarshal(pNoNewlines, bgpRoute); err != nil {
		return n, err
	} else {
		return n, nil
	}
}

func (bgpRoute *BGPRoute) ReadJSONFrom(r io.Reader) (n int64, err error) {

	buf := bytes.Buffer{}

	if n, err = buf.ReadFrom(r); err != nil {
		return n, err
	} else if err = json.Unmarshal(buf.Bytes(), bgpRoute); err != nil {
		return n, err
	} else {
		return n, nil
	}
}

// TODO:
//func (bgpRoute *BGPRoute) ReadCLI(p []byte) (n int, err error) {
//}
