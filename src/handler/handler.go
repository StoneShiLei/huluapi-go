package handler

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"huluapi/src/model"
	"net"

	"github.com/gin-gonic/gin"
)

func OpenComputerHandler(c *gin.Context) {

	var request model.OpenComputerRequest
	if err := c.BindJSON(&request); err != nil {
		c.JSON(400, model.Response{Message: err.Error()})
		return
	}

	if request.MacAddr == "" || request.BcastAddr == "" {
		c.JSON(400, model.Response{
			Message: "mac or bcast is required",
		})
		return
	}

	err := sendMagicPacket(request.MacAddr, request.BcastAddr)

	if err != nil {
		fmt.Println("Failed to send magic packet:", err)
	}

	c.JSON(200, model.Response{
		Message: "success",
	})
}

func CloseComputerHandler(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}

func sendMagicPacket(macAddr string, bcastAddr string) error {
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
