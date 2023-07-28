package handler

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"huluapi/src/model"
	"log"
	"net"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/ssh"
)

func Test(c *gin.Context) {
	var json map[string]interface{}

	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	testBody := json["testBody"]

	header := c.Request.Header.Get("TestHeader")

	c.JSON(200, model.Response{
		Message: "响应值：TestHeader[" + header + "];TestBody[" + testBody.(string) + "]",
	})
}

// 开电脑 "04:7C:16:75:80:20", "255.255.255.255:9"
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

// 关电脑 "192.168.66.2" "Stone" "84022499"
func CloseComputerHandler(c *gin.Context) {

	var request model.CloseComputerRequest
	if err := c.BindJSON(&request); err != nil {
		c.JSON(400, model.Response{Message: err.Error()})
		return
	}

	if request.IpAddr == "" || request.Username == "" || request.Password == "" {
		c.JSON(400, model.Response{
			Message: "ipaddr or username or password is required",
		})
		return
	}

	shutdownComputer(request.Username, request.Password, request.IpAddr)

	c.JSON(200, model.Response{
		Message: "success",
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

func shutdownComputer(userName string, password string, ipAddr string) {
	// 设置SSH客户端配置
	config := &ssh.ClientConfig{
		User: userName,
		Auth: []ssh.AuthMethod{
			ssh.Password(password),
		},
		// 必须验证服务器，否则可能会受到中间人攻击。
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	// 连接到SSH服务器
	client, err := ssh.Dial("tcp", ipAddr+":22", config)
	if err != nil {
		log.Fatal("Failed to dial: ", err)
	}
	defer client.Close()

	// 创建一个会话
	session, err := client.NewSession()
	if err != nil {
		log.Fatal("Failed to create session: ", err)
	}
	defer session.Close()

	// 运行远程命令来关闭计算机
	if err := session.Run("shutdown /s /t 0"); err != nil {
		log.Fatal("Failed to run: " + err.Error())
	}
}
