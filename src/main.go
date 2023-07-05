package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"net"
	"time"
)

func main() {
	err := SendMagicPacket("04:7C:16:75:80:20", "255.255.255.0:9090")
	if err != nil {
		fmt.Println("Failed to send magic packet:", err)
	}

	for {
		time.Sleep(1000)
	}
}

func SendMagicPacket(macAddr string, bcastAddr string) error {
	// Parse MAC address
	hwAddr, err := net.ParseMAC(macAddr)
	if err != nil {
		return err
	}

	// Create magic packet
	magicPacket := createMagicPacket(hwAddr)

	// Get the broadcast address
	addr, err := net.ResolveUDPAddr("udp", bcastAddr)
	if err != nil {
		return err
	}

	// Dial UDP connection
	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		return err
	}
	defer conn.Close()

	// Write magic packet to connection
	_, err = conn.Write(magicPacket)
	if err != nil {
		return err
	}

	return nil
}

func createMagicPacket(hwAddr net.HardwareAddr) []byte {
	// Create magic packet
	// 6 bytes of FF, followed by 16 repetitions of the hardware address
	var packet bytes.Buffer
	binary.Write(&packet, binary.BigEndian, []byte{255, 255, 255, 255, 255, 255})
	for i := 0; i < 16; i++ {
		binary.Write(&packet, binary.BigEndian, hwAddr)
	}

	return packet.Bytes()
}
