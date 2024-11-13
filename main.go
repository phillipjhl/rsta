package main

import (
	"time"
	"bytes"
	"encoding/binary"
	"fmt"
	"net"
	"os"
)

type PacketHeader struct {
	PacketFormat            uint16
	GameMajorVersion        uint8
	GameMinorVersion        uint8
	PacketVersion           uint8
	PacketId                uint8
	SessionUID              uint64
	SessionTime             float32
	FrameIdentifier         uint32
	PlayerCarIndex          uint8
	SecondaryPlayerCarIndex uint8
}

type Packet struct {
	Header PacketHeader
	Data   interface{}
}

type CarMotionData struct {
	WorldPositionX     float32 // World space X position
	WorldPositionY     float32 // World space Y position
	WorldPositionZ     float32 // World space Z position
	WorldVelocityX     float32 // Velocity in world space X
	WorldVelocityY     float32 // Velocity in world space Y
	WorldVelocityZ     float32 // Velocity in world space Z
	WorldForwardDirX   int16   // World space forward X direction (normalised)
	WorldForwardDirY   int16   // World space forward Y direction (normalised)
	WorldForwardDirZ   int16   // World space forward Z direction (normalised)
	WorldRightDirX     int16   // World space right X direction (normalised)
	WorldRightDirY     int16   // World space right Y direction (normalised)
	WorldRightDirZ     int16   // World space right Z direction (normalised)
	GForceLateral      float32 // Lateral G-Force component
	GForceLongitudinal float32 // Longitudinal G-Force component
	GForceVertical     float32 // Vertical G-Force component
	Yaw                float32 // Yaw angle in radians
	Pitch              float32 // Pitch angle in radians
	Roll               float32 // Roll angle in radians
}

type PacketMotionData struct {
	Header        PacketHeader      // Header
	CarMotionData [22]CarMotionData // Data for all cars on track
	// Extra player car ONLY data
	SuspensionPosition     [4]float32 // All wheel arrays have order: RL, RR, FL, FR
	SuspensionVelocity     [4]float32 // RL, RR, FL, FR
	SuspensionAcceleration [4]float32 // RL, RR, FL, FR
	WheelSpeed             [4]float32 // Speed of each wheel
	WheelSlip              [4]float32 // Slip ratio for each wheel
	LocalVelocityX         float32    // Velocity in local space
	LocalVelocityY         float32    // Velocity in local space
	LocalVelocityZ         float32    // Velocity in local space
	AngularVelocityX       float32    // Angular velocity x-component
	AngularVelocityY       float32    // Angular velocity y-component
	AngularVelocityZ       float32    // Angular velocity z-component
	AngularAccelerationX   float32    // Angular acceleration x-component
	AngularAccelerationY   float32    // Angular acceleration y-component
	AngularAccelerationZ   float32    // Angular acceleration z-component
	FrontWheelsAngle       float32    // Current front wheels angle in radians
}

func main() {
	// write raw data to log file
	timestamp := time.Now().UnixMilli()
	filename := fmt.Sprintf("%d_udp_telemetry_raw.log", timestamp)
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("Could not open file: %v", err)
	}
	defer file.Close()

	// start listening for UDP packets
	addr := net.UDPAddr{
		Port: 20777,
		IP:   net.ParseIP("0.0.0.0"),
	}
	ln, err := net.ListenUDP("udp", &addr)

	if err != nil {
		fmt.Println("Error listening:", err)
		panic(err)
	}

	defer ln.Close()

	fmt.Println("Listening 0.0.0.0:20777")

	// Allocate buffer for each packet
	buffer := make([]byte, 2048)
	for {
		rlen, remote, err := ln.ReadFromUDP(buffer)

		if err != nil {
			fmt.Printf("Error reading from UDP: %f", err)
			continue
		}
		fmt.Printf("Data In. len: %v, addr: %v\n", rlen, remote)

		// Read only the length of the buffer
		data := buffer[:rlen]
		// write raw buffer to log file
		if _, err := file.Write(data); err != nil {
			fmt.Println("Failed to write to log file: ", err)
		}

		var header PacketHeader
		reader := bytes.NewReader(data)
		err = binary.Read(reader, binary.LittleEndian, &header)
		if err != nil {
			fmt.Printf("Could not read packet: %v \n", err)
		}

		fmt.Printf("Header: %v \n", header)
		// fmt.Printf("WorldPostionX: %v \n", packet.CarMotionData.WorldPositionX)
		// fmt.Printf("WorldPostionY: %v \n", packet.CarMotionData.WorldPositionY)
	}

}
