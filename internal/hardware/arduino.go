// Package hardware deals with all hardware communication
package hardware

import (
	"encoding/binary"
	"fmt"
	"math"
	"os"

	"github.com/znwng/cdsl/internal/core"
	"golang.org/x/sys/unix"
)

// ========================================
// Serial Communication
// ========================================

const (
	startByte = 0xAA

	crc16InitialValue = 0xFFFF
	crc16Polynomial   = 0x8408

	bitsPerByte = 8

	baud9600   = 9600
	baud19200  = 19200
	baud38400  = 38400
	baud57600  = 57600
	baud115200 = 115200
)

const packetSize = 1 + 2 + 1 + 4 + 2

// getBaudRate converts a configured baud rate into its corresponding
// Unix terminal baud-rate constant.
func getBaudRate(baudRate int) (uint32, error) {
	switch baudRate {
	case baud9600:
		return unix.B9600, nil
	case baud19200:
		return unix.B19200, nil
	case baud38400:
		return unix.B38400, nil
	case baud57600:
		return unix.B57600, nil
	case baud115200:
		return unix.B115200, nil
	default:
		return 0, fmt.Errorf("unsupported baud rate: %d", baudRate)
	}
}

// crc16 calculates the CRC-16 checksum for the supplied packet data.
func crc16(data []byte) uint16 {
	crc := uint16(crc16InitialValue)

	for _, b := range data {
		crc ^= uint16(b)

		for range bitsPerByte {
			if crc&1 != 0 {
				crc = (crc >> 1) ^ crc16Polynomial
			} else {
				crc >>= 1
			}
		}
	}

	return crc
}

// ========================================
// Command Transmission
// ========================================

// SendCommand sends a component command to the Arduino over the configured
// serial port.
//
// The command contains the component ID, target value, packet ID, and a
// CRC-16 checksum for packet integrity.
func SendCommand(config *core.Config, componentLabel string, value float32) error {
	componentID, err := core.ComponentID(config, componentLabel)
	if err != nil {
		return fmt.Errorf("get component ID for %q: %w", componentLabel, err)
	}

	configBaudRate, err := core.ArduinoBaudRate(config)
	if err != nil {
		return fmt.Errorf("get arduino baud rate: %w", err)
	}

	baudRate, err := getBaudRate(configBaudRate)
	if err != nil {
		return fmt.Errorf("configure baud rate: %w", err)
	}

	configPort, err := core.ArduinoPort(config)
	if err != nil {
		return fmt.Errorf("get arduino port: %w", err)
	}

	serial, err := os.OpenFile(
		configPort,
		os.O_WRONLY|unix.O_NOCTTY,
		0,
	)
	if err != nil {
		return fmt.Errorf("open serial port %q: %w", configPort, err)
	}

	defer func() {
		if closeErr := serial.Close(); closeErr != nil {
			// There is no useful way to return this error from a
			// deferred function without changing the function's
			// named return values.
			fmt.Printf("warning: close serial port %q: %v\n", configPort, closeErr)
		}
	}()

	fd := int(serial.Fd())

	termios, err := unix.IoctlGetTermios(fd, unix.TCGETS)
	if err != nil {
		return fmt.Errorf(
			"get serial configuration for %q: %w",
			configPort,
			err,
		)
	}

	termios.Cflag = (termios.Cflag &^ unix.CBAUD) | baudRate

	termios.Cflag &^= unix.PARENB
	termios.Cflag &^= unix.CSTOPB
	termios.Cflag &^= unix.CSIZE
	termios.Cflag |= unix.CS8

	termios.Cflag |= unix.CREAD | unix.CLOCAL

	termios.Lflag &^= unix.ICANON | unix.ECHO | unix.ECHOE | unix.ISIG
	termios.Iflag &^= unix.IXON | unix.IXOFF | unix.IXANY
	termios.Oflag &^= unix.OPOST

	if err := unix.IoctlSetTermios(fd, unix.TCSETS, termios); err != nil {
		return fmt.Errorf(
			"configure serial port %q: %w",
			configPort,
			err,
		)
	}

	packet := make([]byte, 0, packetSize)

	packet = append(packet, startByte)

	packetID := staticPacketID

	var packetIDBytes [2]byte
	binary.LittleEndian.PutUint16(packetIDBytes[:], packetID)
	packet = append(packet, packetIDBytes[:]...)

	packet = append(packet, byte(componentID))

	var valueBytes [4]byte
	binary.LittleEndian.PutUint32(
		valueBytes[:],
		math.Float32bits(value),
	)
	packet = append(packet, valueBytes[:]...)

	crc := crc16(packet)

	var crcBytes [2]byte
	binary.LittleEndian.PutUint16(crcBytes[:], crc)
	packet = append(packet, crcBytes[:]...)

	totalWritten := 0

	for totalWritten < len(packet) {
		n, err := serial.Write(packet[totalWritten:])
		if err != nil {
			return fmt.Errorf(
				"write command to serial port %q: %w",
				configPort,
				err,
			)
		}

		if n == 0 {
			return fmt.Errorf(
				"write command to serial port %q: zero bytes written",
				configPort,
			)
		}

		totalWritten += n
	}

	staticPacketID++

	return nil
}

// ========================================
// Packet State
// ========================================

var staticPacketID uint16
