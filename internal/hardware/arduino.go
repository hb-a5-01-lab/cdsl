package hardware

import (
	"fmt"
	"math"
	"os"
	"strings"
	"time"

	"github.com/kalexion/cdsl/internal/core"
	"golang.org/x/sys/unix"
)

// ========================================
// Serial Configuration
// ========================================

const (
	baud9600   = 9600
	baud19200  = 19200
	baud38400  = 38400
	baud57600  = 57600
	baud115200 = 115200

	responseTimeout = 5 * time.Second
)

// ========================================
// Baud Rate
// ========================================

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

// ========================================
// Serial Port Configuration
// ========================================

func configureSerial(fd int, baudRate uint32) error {
	termios, err := unix.IoctlGetTermios(fd, unix.TCGETS)
	if err != nil {
		return fmt.Errorf("get serial configuration: %w", err)
	}

	// Configure raw input.
	termios.Iflag &^= unix.IGNBRK |
		unix.BRKINT |
		unix.PARMRK |
		unix.ISTRIP |
		unix.INLCR |
		unix.IGNCR |
		unix.ICRNL |
		unix.IXON |
		unix.IXOFF |
		unix.IXANY

	// Disable output processing and terminal line handling.
	termios.Oflag &^= unix.OPOST

	termios.Lflag &^= unix.ECHO |
		unix.ECHONL |
		unix.ICANON |
		unix.ISIG |
		unix.IEXTEN

	// Configure 8-N-1.
	termios.Cflag &^= unix.CSIZE |
		unix.PARENB |
		unix.CSTOPB |
		unix.CBAUD

	termios.Cflag |= unix.CS8 |
		unix.CREAD |
		unix.CLOCAL |
		baudRate

	// Polling controls the timeout.
	termios.Cc[unix.VMIN] = 1
	termios.Cc[unix.VTIME] = 0

	if err := unix.IoctlSetTermios(fd, unix.TCSETS, termios); err != nil {
		return fmt.Errorf("set serial configuration: %w", err)
	}

	return nil
}

// ========================================
// Serial Read / Write
// ========================================

func writeAll(file *os.File, data []byte) error {
	for len(data) > 0 {
		n, err := file.Write(data)
		if err != nil {
			return err
		}

		if n == 0 {
			return fmt.Errorf("zero bytes written")
		}

		data = data[n:]
	}

	return nil
}

// readLine waits for one newline-terminated Arduino response.
func readLine(fd int, serial *os.File) (string, error) {
	deadline := time.Now().Add(responseTimeout)
	var response strings.Builder

	for {
		remaining := time.Until(deadline)
		if remaining <= 0 {
			return "", fmt.Errorf("timed out waiting for Arduino response")
		}

		timeoutMS := int(remaining.Milliseconds())
		if timeoutMS < 1 {
			timeoutMS = 1
		}

		fds := []unix.PollFd{
			{
				Fd:     int32(fd),
				Events: unix.POLLIN,
			},
		}

		n, err := unix.Poll(fds, timeoutMS)
		if err == unix.EINTR {
			continue
		}
		if err != nil {
			return "", fmt.Errorf("poll serial port: %w", err)
		}
		if n == 0 {
			return "", fmt.Errorf("timed out waiting for Arduino response")
		}

		if fds[0].Revents&(unix.POLLERR|unix.POLLHUP|unix.POLLNVAL) != 0 {
			return "", fmt.Errorf(
				"serial port error: events %#x",
				fds[0].Revents,
			)
		}

		if fds[0].Revents&unix.POLLIN == 0 {
			continue
		}

		var b [1]byte
		count, err := serial.Read(b[:])
		if err != nil {
			return "", fmt.Errorf("read serial port: %w", err)
		}
		if count == 0 {
			continue
		}

		if b[0] == '\n' {
			return strings.TrimRight(response.String(), "\r"), nil
		}

		response.WriteByte(b[0])
	}
}

// ========================================
// Command Transmission
// ========================================

// SendCommand sends a component ID and floating-point value to the Arduino,
// waits for its text response, and displays the result.
func SendCommand(
	config *core.Config,
	componentLabel string,
	value float32,
) error {
	if math.IsNaN(float64(value)) || math.IsInf(float64(value), 0) {
		return fmt.Errorf("invalid value: %v", value)
	}

	if value < 0 || value > 180 {
		return fmt.Errorf("value must be between 0 and 180: %v", value)
	}

	componentID, err := core.ComponentID(config, componentLabel)
	if err != nil {
		return fmt.Errorf(
			"get component ID for %q: %w",
			componentLabel,
			err,
		)
	}

	configuredBaudRate, err := core.ArduinoBaudRate(config)
	if err != nil {
		return fmt.Errorf("get Arduino baud rate: %w", err)
	}

	baudRate, err := getBaudRate(configuredBaudRate)
	if err != nil {
		return fmt.Errorf("configure baud rate: %w", err)
	}

	port, err := core.ArduinoPort(config)
	if err != nil {
		return fmt.Errorf("get Arduino port: %w", err)
	}

	// Open for both reading and writing.
	serial, err := os.OpenFile(
		port,
		os.O_RDWR|unix.O_NOCTTY,
		0,
	)
	if err != nil {
		return fmt.Errorf("open serial port %q: %w", port, err)
	}
	defer serial.Close()

	fd := int(serial.Fd())

	if err := configureSerial(fd, baudRate); err != nil {
		return fmt.Errorf("configure serial port %q: %w", port, err)
	}

	// Many Arduino boards reset when a serial connection is opened.
	fmt.Println("Waiting for Arduino to initialize...")
	time.Sleep(2 * time.Second)

	// Discard startup messages before sending the command.
	if err := unix.IoctlSetInt(fd, unix.TCFLSH, unix.TCIFLUSH); err != nil {
		return fmt.Errorf("flush serial input: %w", err)
	}

	// ----------------------------------------
	// Send plain-text command
	// ----------------------------------------

	// Format: <component_id> <float_value>\n
	command := fmt.Sprintf("%d %.6f\n", componentID, value)

	fmt.Printf("Sending: %s", command)

	if err := writeAll(serial, []byte(command)); err != nil {
		return fmt.Errorf("send command to Arduino: %w", err)
	}

	// ----------------------------------------
	// Wait for response
	// ----------------------------------------

	fmt.Println("Waiting for Arduino response...")

	response, err := readLine(fd, serial)
	if err != nil {
		return fmt.Errorf("receive Arduino response: %w", err)
	}

	fmt.Printf("Received: %s\n", response)

	// ----------------------------------------
	// Validate response
	// ----------------------------------------

	if response == "REJECTED" {
		return fmt.Errorf(
			"Arduino rejected component %d with value %.6f",
			componentID,
			value,
		)
	}

	expectedPrefix := fmt.Sprintf("EXECUTED: %d ", componentID)

	if !strings.HasPrefix(response, expectedPrefix) {
		return fmt.Errorf(
			"unexpected Arduino response: %q",
			response,
		)
	}

	fmt.Println("Command acknowledged successfully.")

	return nil
}
