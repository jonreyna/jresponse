package bgproute

import (
	"bytes"
	"encoding/xml"
	"os"
	"reflect"
	"testing"
)

var (
	bgpRouteXMLModel, bgpRouteJSONModel *BGPRoute
)

const (
	BGP_XML_FILE  = "show_route_protocol_bgp.xml"
	BGP_JSON_FILE = "show_route_protocol_bgp.json"
	BGP_CLI_FILE  = "show_route_protocol_bgp.cli"
)

func initBGPRouteModel() {
	bgpRouteXMLModel = &BGPRoute{
		XMLName: xml.Name{"http://xml.juniper.net/junos/12.3R6/junos-routing", "route-information"},
		RouteTable: RouteTable{
			TableName:          "inet.0",
			DestinationCount:   565525,
			TotalRouteCount:    4400004,
			ActiveRouteCount:   565520,
			HoldDownRouteCount: 0,
			HiddenRouteCount:   14,
			RT: []RT{
				{
					RTDestination: "8.8.8.0/24",
					RTEntry: []RTEntry{
						{
							ActiveTag:     "*",
							CurrentActive: "",
							LastActive:    "",
							ProtocolName:  "BGP",
							Preference:    170,
							Age: Age{
								AgeSecs: "585128",
								AgeTime: "6d 18:32:08",
							},
							Med:             0,
							LocalPreference: 130,
							LearnedFrom:     "206.126.239.251",
							AsPath:          "15169 I",
							ValidationState: "unverified",
							NH: []NH{
								{
									SelectedNextHop: new(string),
									To:              "206.126.236.21",
									Via:             "ae0.0",
								},
							},
						},
						{
							ActiveTag:    "",
							ProtocolName: "BGP",
							Preference:   170,
							Age: Age{
								AgeSecs: "585128",
								AgeTime: "6d 18:32:08",
							},
							Med:             0,
							LocalPreference: 130,
							LearnedFrom:     "206.126.239.252",
							AsPath:          "15169 I",
							ValidationState: "unverified",
							NH: []NH{
								{
									SelectedNextHop: new(string),
									To:              "206.126.236.21",
									Via:             "ae0.0",
								},
							},
						},
						{
							ActiveTag:    "",
							ProtocolName: "BGP",
							Preference:   170,
							Age: Age{
								AgeSecs: "762247",
								AgeTime: "1w1d 19:44:07",
							},
							Med:             0,
							LocalPreference: 130,
							LearnedFrom:     "76.73.165.1",
							AsPath:          "15169 I",
							ValidationState: "unverified",
							NH: []NH{
								{
									To:      "24.236.73.12",
									Via:     "ae5.0",
									LSPName: "VAASHBPO1EDGJ01>>OHIOLAHUHEDGJ01-ECMP1",
								},
								{
									SelectedNextHop: new(string),
									To:              "24.236.73.12",
									Via:             "ae5.0",
									LSPName:         "VAASHBPO1EDGJ01>>OHIOLAHUHEDGJ01-ECMP2",
								},
								{
									To:      "24.236.73.12",
									Via:     "ae5.0",
									LSPName: "VAASHBPO1EDGJ01>>OHIOLAHUHEDGJ01-ECMP3",
								},
								{
									To:      "69.73.0.136",
									Via:     "ae4.0",
									LSPName: "VAASHBPO1EDGJ01>>OHIOLAHUHEDGJ01-ECMP1",
								},
								{
									To:      "69.73.0.136",
									Via:     "ae4.0",
									LSPName: "VAASHBPO1EDGJ01>>OHIOLAHUHEDGJ01-ECMP2",
								},
								{
									To:      "69.73.0.136",
									Via:     "ae4.0",
									LSPName: "VAASHBPO1EDGJ01>>OHIOLAHUHEDGJ01-ECMP3",
								},
							},
						},
					},
				},
			},
		},
	}

	bgpRouteJSONModel = &BGPRoute{
		XMLName: xml.Name{},
		RouteTable: RouteTable{
			TableName:          "inet.0",
			DestinationCount:   565525,
			TotalRouteCount:    4400004,
			ActiveRouteCount:   565520,
			HoldDownRouteCount: 0,
			HiddenRouteCount:   14,
			RT: []RT{
				{
					RTDestination: "8.8.8.0/24",
					RTEntry: []RTEntry{
						{
							ActiveTag:     "*",
							CurrentActive: "",
							LastActive:    "",
							ProtocolName:  "BGP",
							Preference:    170,
							Age: Age{
								AgeSecs: "585128",
								AgeTime: "6d 18:32:08",
							},
							Med:             0,
							LocalPreference: 130,
							LearnedFrom:     "206.126.239.251",
							AsPath:          "15169 I",
							ValidationState: "unverified",
							NH: []NH{
								{
									SelectedNextHop: new(string),
									To:              "206.126.236.21",
									Via:             "ae0.0",
								},
							},
						},
						{
							ActiveTag:    "",
							ProtocolName: "BGP",
							Preference:   170,
							Age: Age{
								AgeSecs: "585128",
								AgeTime: "6d 18:32:08",
							},
							Med:             0,
							LocalPreference: 130,
							LearnedFrom:     "206.126.239.252",
							AsPath:          "15169 I",
							ValidationState: "unverified",
							NH: []NH{
								{
									SelectedNextHop: new(string),
									To:              "206.126.236.21",
									Via:             "ae0.0",
								},
							},
						},
						{
							ActiveTag:    "",
							ProtocolName: "BGP",
							Preference:   170,
							Age: Age{
								AgeSecs: "762247",
								AgeTime: "1w1d 19:44:07",
							},
							Med:             0,
							LocalPreference: 130,
							LearnedFrom:     "76.73.165.1",
							AsPath:          "15169 I",
							ValidationState: "unverified",
							NH: []NH{
								{
									To:      "24.236.73.12",
									Via:     "ae5.0",
									LSPName: "VAASHBPO1EDGJ01>>OHIOLAHUHEDGJ01-ECMP1",
								},
								{
									SelectedNextHop: new(string),
									To:              "24.236.73.12",
									Via:             "ae5.0",
									LSPName:         "VAASHBPO1EDGJ01>>OHIOLAHUHEDGJ01-ECMP2",
								},
								{
									To:      "24.236.73.12",
									Via:     "ae5.0",
									LSPName: "VAASHBPO1EDGJ01>>OHIOLAHUHEDGJ01-ECMP3",
								},
								{
									To:      "69.73.0.136",
									Via:     "ae4.0",
									LSPName: "VAASHBPO1EDGJ01>>OHIOLAHUHEDGJ01-ECMP1",
								},
								{
									To:      "69.73.0.136",
									Via:     "ae4.0",
									LSPName: "VAASHBPO1EDGJ01>>OHIOLAHUHEDGJ01-ECMP2",
								},
								{
									To:      "69.73.0.136",
									Via:     "ae4.0",
									LSPName: "VAASHBPO1EDGJ01>>OHIOLAHUHEDGJ01-ECMP3",
								},
							},
						},
					},
				},
			},
		},
	}
}

func TestMain(m *testing.M) {
	initBGPRouteModel()
	os.Exit(m.Run())
}

func TestReadXMLFrom(t *testing.T) {

	b := new(BGPRoute)

	if file, err := os.Open(BGP_XML_FILE); err != nil {
		t.Error(err)
	} else if _, err := b.ReadXMLFrom(file); err != nil {
		t.Error(err)
	} else if !reflect.DeepEqual(b, bgpRouteXMLModel) {
		t.Log(bgpRouteXMLModel)
		t.Log(b)
		t.Error("unmarshalled XML does not match BGP route model")
	}
}

func TestReadJSONFrom(t *testing.T) {

	b := new(BGPRoute)

	if file, err := os.Open(BGP_JSON_FILE); err != nil {
		t.Error(err)
	} else if _, err := b.ReadJSONFrom(file); err != nil {
		t.Error(err)
	} else if !reflect.DeepEqual(b, bgpRouteJSONModel) {
		t.Log(bgpRouteJSONModel)
		t.Log(b)
		t.Error("unmarshalled JSON does not match BGP route model")
	}
}

func TestWriteCLITo(t *testing.T) {

	modelBuf := bytes.Buffer{}

	if err := bgpRouteXMLModel.WriteCLITo(&modelBuf); err != nil {
		t.Error(err)
	}

	fileBuf := bytes.Buffer{}
	if file, err := os.Open(BGP_CLI_FILE); err != nil {
		t.Error(err)
	} else if _, err := fileBuf.ReadFrom(file); err != nil {
		t.Error(err)
	}

	if !bytes.Equal(modelBuf.Bytes(), fileBuf.Bytes()) {
		f1, f2 := "model.cli", "current.cli"
		if file1, err := os.Create(f1); err != nil {
			t.Error(err)
		} else if file2, err := os.Create(f2); err != nil {
			t.Error(err)
		} else {
			t.Logf("writing output files for diff: %s, %s", f1, f2)
			file1.Write(modelBuf.Bytes())
			file2.Write(fileBuf.Bytes())
			t.Error("model for CLI parse does not match")
		}
	}
}

func TestIsNextHop(t *testing.T) {
	testStr := ""
	if str := isNextHop(nil); str != " " {
		t.Error("did not return a space when passed a nil pointer")
	} else if str := isNextHop(&testStr); str != ">" {
		t.Error(`did not return '>' when passing empty string`)
	}
}

func TestWriteXMLTo(t *testing.T) {

	modelBuf := bytes.Buffer{}
	if _, err := bgpRouteXMLModel.WriteXMLTo(&modelBuf); err != nil {
		t.Error(err)
	}

	fileBuf := bytes.Buffer{}
	if file, err := os.Open(BGP_XML_FILE); err != nil {
		t.Error(err)
	} else if _, err := fileBuf.ReadFrom(file); err != nil {
		t.Error(err)
	}

	b := new(BGPRoute)
	b.ReadXMLFrom(&fileBuf)

	fileBuf.Reset()
	b.WriteXMLTo(&fileBuf)

	if !bytes.Equal(modelBuf.Bytes(), fileBuf.Bytes()) {
		t.Log(b)
		t.Log(bgpRouteXMLModel)
		t.Error("XML bytes not equal")
	}
}

func TestWriteJSONTo(t *testing.T) {

	modelBuf := bytes.Buffer{}
	if _, err := bgpRouteJSONModel.WriteJSONTo(&modelBuf); err != nil {
		t.Error(err)
	}

	fileBuf := bytes.Buffer{}
	if file, err := os.Open(BGP_JSON_FILE); err != nil {
		t.Error(err)
	} else if _, err := fileBuf.ReadFrom(file); err != nil {
		t.Error(err)
	}

	b := new(BGPRoute)
	b.ReadJSONFrom(&fileBuf)
	fileBuf.Reset()
	b.WriteJSONTo(&fileBuf)

	if !bytes.Equal(modelBuf.Bytes(), fileBuf.Bytes()) {
		t.Log(b)
		t.Log(bgpRouteJSONModel)
		t.Error("JSON bytes not equal")
	}
}
