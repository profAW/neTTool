package infrastructure

import (
	"fmt"
	"log"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/pcap"
)

// SavedPacketsAdapter get Data from saved pcap file
type SavedPacketsAdapter struct {
	FileAndFolder string
}

var (
	pcapFile string
	handle   *pcap.Handle
	err      error
)

var packets = make(map[int]gopacket.Packet)

func (e SavedPacketsAdapter) Read() map[int]gopacket.Packet {
	// Open file instead of device
	pcapFile = e.FileAndFolder
	handle, err = pcap.OpenOffline(pcapFile)
	if err != nil {
		log.Fatal(err)
		fmt.Println("ERROR LOAD PCAP")
	}
	defer handle.Close()

	// Loop through packets in file
	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())

	i := 0
	for packet := range packetSource.Packets() {
		packets[i] = packet
		i = i + 1
	}

	//fmt.Println(packets)
	return packets
}

//  --------------------------------------------------------

var (
	device      string        //   = "/Device/NPF_{3EFE15FE-3499-4237-92D6-7AB5B4B9FD9C}" //"eth0"
	snapshotLen int32         //= 1024
	promiscuous bool          //= true
	timeout     time.Duration //= 30 * time.Second
)

// LivePacketsAdapter get Data from live network stream of defined device
type LivePacketsAdapter struct {
	Device      string
	SnapshotLen int32
	Promiscuous bool
	Timeout     time.Duration
}

func (e LivePacketsAdapter) Read() map[int]gopacket.Packet {

	device = e.Device
	snapshotLen = e.SnapshotLen
	promiscuous = e.Promiscuous
	timeout = e.Timeout * time.Second
	handle, err = pcap.OpenLive(device, snapshotLen, promiscuous, timeout)
	if err != nil {
		log.Fatal(err)
	}
	defer handle.Close()

	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
	i := 0
	for packet := range packetSource.Packets() {
		packets[i] = packet
		i = i + 1
	}

	//fmt.Println(packets)
	return packets
}
