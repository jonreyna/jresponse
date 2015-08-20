package main

import (
	"bytes"
	"log"
	"os"

	"github.com/JReyLBC/jresponse/command/traceroute"
	"github.com/Juniper/go-netconf/netconf"
)

func main() {

	rawMethod := netconf.RawMethod("<traceroute><host>8.8.8.8</host></traceroute>")
	s := new(netconf.Session)
	var err error

	if s, err = netconf.DialSSH("192.168.1.1", netconf.SSHConfigPassword("username", "password")); err != nil {
		log.Fatal(err)
	} else {
		defer s.Close()
	}

	tr := new(traceroute.TraceRoute)
	if ncReply, err := s.Exec(rawMethod); err != nil {
		log.Fatal(err)
	} else if _, err := tr.ReadXMLFrom(bytes.NewBuffer([]byte(ncReply.Data))); err != nil {
		log.Fatal(err)
	} else {
		tr.WriteCLITo(os.Stdout)
	}
}
