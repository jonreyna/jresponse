package traceroute

import (
	"bytes"
	"encoding/xml"
	"os"
	"reflect"
	"testing"
)

var (
	traceRouteXMLModel, traceRouteJSONModel *TraceRoute
)

const (
	traceRouteXMLFile  = "traceroute_8.8.8.8.xml"
	traceRouteJSONFile = "traceroute_8.8.8.8.json"
	traceRouteCLIFile  = "traceroute_8.8.8.8.cli"
)

func initTraceRouteModel() {

	traceRouteXMLModel = &TraceRoute{
		XMLName:     xml.Name{"http://xml.juniper.net/junos/12.1X46/junos-probe-tests", "traceroute-results"},
		TargetHost:  "8.8.8.8",
		TargetIP:    "8.8.8.8",
		MaxHopIndex: 30,
		PacketSize:  40,
		Hops: []Hop{
			{
				TTLValue:     1,
				LastIPAddr:   "10.226.0.1",
				LastHostName: "10.226.0.1",
				ProbeResult: []ProbeResult{
					{
						DateDetermined: 1439961690,
						ProbeIndex:     1,
						IPAddress:      "10.226.0.1",
						HostName:       "10.226.0.1",
						RTT:            13876,
						ProbeSuccess:   new(string),
					},
					{
						DateDetermined: 1439961690,
						ProbeIndex:     2,
						IPAddress:      "10.226.0.1",
						HostName:       "10.226.0.1",
						RTT:            11752,
						ProbeSuccess:   new(string),
					},
					{
						DateDetermined: 1439961690,
						ProbeIndex:     3,
						IPAddress:      "10.226.0.1",
						HostName:       "10.226.0.1",
						RTT:            10973,
						ProbeSuccess:   new(string),
					},
				},
			},
		},
	}

	traceRouteJSONModel = &TraceRoute{
		XMLName:     xml.Name{},
		TargetHost:  "8.8.8.8",
		TargetIP:    "8.8.8.8",
		MaxHopIndex: 30,
		PacketSize:  40,
		Hops: []Hop{
			{
				TTLValue:     1,
				LastIPAddr:   "10.226.0.1",
				LastHostName: "10.226.0.1",
				ProbeResult: []ProbeResult{
					{
						DateDetermined: 1439961690,
						ProbeIndex:     1,
						IPAddress:      "10.226.0.1",
						HostName:       "10.226.0.1",
						RTT:            26792,
						ProbeSuccess:   new(string),
					},
					{
						DateDetermined: 1439961690,
						ProbeIndex:     2,
						IPAddress:      "10.226.0.1",
						HostName:       "10.226.0.1",
						RTT:            14184,
						ProbeSuccess:   new(string),
					},
					{
						DateDetermined: 1439961690,
						ProbeIndex:     3,
						IPAddress:      "10.226.0.1",
						HostName:       "10.226.0.1",
						RTT:            15821,
						ProbeSuccess:   new(string),
					},
				},
			},
		},
	}
}

func TestMain(m *testing.M) {
	initTraceRouteModel()
	os.Exit(m.Run())
}

func TestReadXMLFrom(t *testing.T) {

	tr := new(TraceRoute)

	if file, err := os.Open(traceRouteXMLFile); err != nil {
		t.Error(err)
	} else if _, err := tr.ReadXMLFrom(file); err != nil {
		t.Error(err)
	} else if !reflect.DeepEqual(tr, traceRouteXMLModel) {
		t.Log(traceRouteXMLModel)
		t.Log(tr)
		t.Error("unmarshalled XML does not match trace route model")
	}
}

func TestReadJSONFrom(t *testing.T) {

	tr := new(TraceRoute)

	if file, err := os.Open(traceRouteJSONFile); err != nil {
		t.Error(err)
	} else if _, err := tr.ReadJSONFrom(file); err != nil {
		t.Error(err)
	} else if !reflect.DeepEqual(tr, traceRouteJSONModel) {
		t.Error("unmarshalled JSON does not match trace route model")
	}
}

func TestWriteCLITo(t *testing.T) {

	modelBuf := bytes.Buffer{}

	if err := traceRouteXMLModel.WriteCLITo(&modelBuf); err != nil {
		t.Error(err)
	}

	fileBuf := bytes.Buffer{}
	if file, err := os.Open(traceRouteCLIFile); err != nil {
		t.Error(err)
	} else if _, err := fileBuf.ReadFrom(file); err != nil {
		t.Error(err)
	}

	mbuf := bytes.TrimSpace(modelBuf.Bytes())
	fbuf := bytes.TrimSpace(fileBuf.Bytes())
	if !bytes.Equal(mbuf, fbuf) {
		f1, f2 := "model.cli", "current.cli"
		if file1, err := os.Create(f1); err != nil {
			t.Error(err)
		} else if file2, err := os.Create(f2); err != nil {
			t.Error(err)
		} else {
			t.Logf("writing output files for diff: %s, %s", f1, f2)
			file1.Write(mbuf)
			file2.Write(fbuf)
			t.Error("model for CLI parse does not match")
		}
	}
}

func TestNumberFormat(t *testing.T) {
	msSlice := []uint{4635, 2639, 9478, 12850, 10227, 10796, 10612, 13286, 10605}
	correctOutput := []string{"4.635", "2.639", "9.478", "12.85", "10.227", "10.796", "10.612", "13.286", "10.605"}

	for index, num := range msSlice {
		if strNum := formatAsMs(num); strNum != correctOutput[index] {
			t.Errorf("Incorrect millisecond conversion: %s should be %s", strNum, correctOutput[index])
		}
	}
}

func TestWriteXMLTo(t *testing.T) {
	modelBuf := bytes.Buffer{}
	if _, err := traceRouteXMLModel.WriteXMLTo(&modelBuf); err != nil {
		t.Error(err)
	}

	fileBuf := bytes.Buffer{}
	if file, err := os.Open(traceRouteXMLFile); err != nil {
		t.Error(err)
	} else if _, err := fileBuf.ReadFrom(file); err != nil {
		t.Error(err)
	}

	tr := new(TraceRoute)
	tr.ReadXMLFrom(&fileBuf)

	fileBuf.Reset()
	tr.WriteXMLTo(&fileBuf)

	if !bytes.Equal(modelBuf.Bytes(), fileBuf.Bytes()) {
		t.Log(tr)
		t.Log(traceRouteXMLModel)
		t.Error("XML bytes not equal")
	}
}

func TestWriteJSONTo(t *testing.T) {
	modelBuf := bytes.Buffer{}
	if _, err := traceRouteJSONModel.WriteJSONTo(&modelBuf); err != nil {
		t.Error(err)
	}

	fileBuf := bytes.Buffer{}
	if file, err := os.Open(traceRouteJSONFile); err != nil {
		t.Error(err)
	} else if _, err := fileBuf.ReadFrom(file); err != nil {
		t.Error(err)
	}

	tr := new(TraceRoute)
	tr.ReadJSONFrom(&fileBuf)
	fileBuf.Reset()
	tr.WriteJSONTo(&fileBuf)

	if !bytes.Equal(modelBuf.Bytes(), fileBuf.Bytes()) {
		t.Log(tr)
		t.Log(traceRouteJSONModel)
		t.Error("JSON bytes not equal")
	}
}
