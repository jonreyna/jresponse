package ping

import (
	"bytes"
	"encoding/xml"
	"os"
	"reflect"
	"testing"
)

var (
	pingXMLModel, pingJSONModel *Ping
)

const (
	PING_XML_FILE  = "ping_8.8.8.8.xml"
	PING_JSON_FILE = "ping_8.8.8.8.json"
	PING_CLI_FILE = "ping_8.8.8.8.cli"
)

func initPingModel() {

	pingXMLModel = &Ping{
		XMLName:     xml.Name{"http://xml.juniper.net/junos/12.3R7/junos-probe-tests", "ping-results"},
		TargetHost:  "8.8.8.8",
		TargetIP:    "8.8.8.8",
		PacketSize:  1200,
				ProbeResult: []ProbeResult{
					{
						DateDetermined: 1447350764,
						ProbeIndex:     1,
						ProbeSuccess:   new(string),						
						SequenceNumber: 0,
						IPAddress:      "8.8.8.8",
						TimeToLive:     62,
						ResponseSize:   1208,
						RTT:            690,
					},
					{
						DateDetermined: 1447350765,
						ProbeIndex:     2,						
						ProbeSuccess:   new(string),						
						SequenceNumber: 1,
						IPAddress:      "8.8.8.8",
						TimeToLive:     62,
						ResponseSize:   1208,
						RTT:            644,
					},
					{
						DateDetermined: 1447350766,
						ProbeIndex:     3,						
						ProbeSuccess:   new(string),						
						SequenceNumber: 2,
						IPAddress:      "8.8.8.8",
						TimeToLive:     62,
						ResponseSize:   1208,
						RTT:            681,
					},
					{
						DateDetermined: 1447350767,
						ProbeIndex:     4,						
						ProbeSuccess:   new(string),						
						SequenceNumber: 3,
						IPAddress:      "8.8.8.8",
						TimeToLive:     62,
						ResponseSize:   1208,
						RTT:            645,
					},
					{
						DateDetermined: 1447350768,
						ProbeIndex:     5,						
						ProbeSuccess:   new(string),						
						SequenceNumber: 4,
						IPAddress:      "8.8.8.8",
						TimeToLive:     62,
						ResponseSize:   1208,
						RTT:            686,
					},								
				},
                ProbeResultsSummary: ProbeResultsSummary{
						ProbesSent:            5,						
						ResponsesReceived:     5,						
						PacketLoss:            0,						
						RTTMinimum:            644,						
						RTTMaximum:            690,						
						RTTAverage:            669,						
						RTTStdDev:            20,             
                },

	}

	pingJSONModel = &Ping{
		XMLName:     xml.Name{},
		TargetHost:  "8.8.8.8",
		TargetIP:    "8.8.8.8",
		PacketSize:  40,
				ProbeResult: []ProbeResult{
					{
						DateDetermined: 1439961690,
						ProbeIndex:     1,
						ProbeSuccess:   new(string),						
						SequenceNumber: 0,
						IPAddress:      "8.8.8.8",
						TimeToLive:     62,
						ResponseSize:   1208,
						RTT:            690,
					},
					{
						DateDetermined: 1439961690,
						ProbeIndex:     2,						
						ProbeSuccess:   new(string),						
						SequenceNumber: 1,
						IPAddress:      "8.8.8.8",
						TimeToLive:     62,
						ResponseSize:   1208,
						RTT:            644,
					},
					{
						DateDetermined: 1439961690,
						ProbeIndex:     3,						
						ProbeSuccess:   new(string),						
						SequenceNumber: 2,
						IPAddress:      "8.8.8.8",
						TimeToLive:     62,
						ResponseSize:   1208,
						RTT:            681,
					},
					{
						DateDetermined: 1447350767,
						ProbeIndex:     4,						
						ProbeSuccess:   new(string),						
						SequenceNumber: 3,
						IPAddress:      "8.8.8.8",
						TimeToLive:     62,
						ResponseSize:   1208,
						RTT:            645,
					},
					{
						DateDetermined: 1447350768,
						ProbeIndex:     5,						
						ProbeSuccess:   new(string),						
						SequenceNumber: 4,
						IPAddress:      "8.8.8.8",
						TimeToLive:     62,
						ResponseSize:   1208,
						RTT:            686,
					},									
				},
                ProbeResultsSummary: ProbeResultsSummary{
						ProbesSent:            5,						
						ResponsesReceived:     5,						
						PacketLoss:            0,						
						RTTMinimum:            644,						
						RTTMaximum:            690,						
						RTTAverage:            669,						
						RTTStdDev:             20,                  
                },

	}
}

func TestMain(m *testing.M) {
	initPingModel()
	os.Exit(m.Run())
}

func TestReadXMLFrom(t *testing.T) {

	ping := new(Ping)

	if file, err := os.Open(PING_XML_FILE); err != nil {
		t.Error(err)
	} else if _, err := ping.ReadXMLFrom(file); err != nil {
		t.Error(err)
	} else if !reflect.DeepEqual(ping, pingXMLModel) {
		t.Log(pingXMLModel)
		t.Log(ping)
		t.Error("unmarshalled XML does not match the test ping json model")
	}
}

func TestReadJSONFrom(t *testing.T) {

	ping := new(Ping)

	if file, err := os.Open(PING_JSON_FILE); err != nil {
		t.Error(err)
	} else if _, err := ping.ReadJSONFrom(file); err != nil {
		t.Error(err)
	} else if !reflect.DeepEqual(ping, pingJSONModel) {
    	t.Log(pingJSONModel)
    	t.Log(ping)
		t.Error("unmarshalled JSON does not match the test ping json model")
	}
}

func TestWriteCLITo(t *testing.T) {

	modelBuf := bytes.Buffer{}

	if err := pingXMLModel.WriteCLITo(&modelBuf); err != nil {
		t.Error(err)
	}

	fileBuf := bytes.Buffer{}
	if file, err := os.Open(PING_CLI_FILE); err != nil {
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
	correctOutput := []string{"4.635", "2.639", "9.478", "12.850", "10.227", "10.796", "10.612", "13.286", "10.605"}

	for index, num := range msSlice {
		if strNum := formatAsMs(num); strNum != correctOutput[index] {
			t.Errorf("Incorrect millisecond conversion: %s should be %s", strNum, correctOutput[index])
		}
	}
}

func TestWriteXMLTo(t *testing.T) {
	modelBuf := bytes.Buffer{}
	if _, err := pingXMLModel.WriteXMLTo(&modelBuf); err != nil {
		t.Error(err)
	}

	fileBuf := bytes.Buffer{}
	if file, err := os.Open(PING_XML_FILE); err != nil {
		t.Error(err)
	} else if _, err := fileBuf.ReadFrom(file); err != nil {
		t.Error(err)
	}

	ping := new(Ping)
	ping.ReadXMLFrom(&fileBuf)

	fileBuf.Reset()
	ping.WriteXMLTo(&fileBuf)

	if !bytes.Equal(modelBuf.Bytes(), fileBuf.Bytes()) {
		t.Log(ping)
		t.Log(pingXMLModel)
		t.Error("XML bytes not equal")
	}
}

func TestWriteJSONTo(t *testing.T) {
	modelBuf := bytes.Buffer{}
	if _, err := pingJSONModel.WriteJSONTo(&modelBuf); err != nil {
		t.Error(err)
	}

	fileBuf := bytes.Buffer{}
	if file, err := os.Open(PING_JSON_FILE); err != nil {
		t.Error(err)
	} else if _, err := fileBuf.ReadFrom(file); err != nil {
		t.Error(err)
	}

	ping := new(Ping)
	ping.ReadJSONFrom(&fileBuf)
	fileBuf.Reset()
	ping.WriteJSONTo(&fileBuf)

	if !bytes.Equal(modelBuf.Bytes(), fileBuf.Bytes()) {
		t.Log(ping)
		t.Log(pingJSONModel)
		t.Error("JSON bytes not equal")
	}
}
