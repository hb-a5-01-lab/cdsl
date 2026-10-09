// Package logging provides file-based logging for CDSL.
package logging

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/kalexion/cdsl/internal/customTypes"
)

var logFile *os.File

// Log format: [timestamp][instruction][SUCCESS/FAILURE][message]
type Log struct {
	Timestamp         time.Time
	Success           bool
	InstructionString string
	Message           string
}

// Init creates the required directory and log file if absent.
// Call once when the main program starts.
func Init() error {
	if logFile != nil {
		return nil
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("get home directory: %w", err)
	}

	logDir := filepath.Join(homeDir, ".local", "share", "cdsl")
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return fmt.Errorf("create log directory: %w", err)
	}

	logPath := filepath.Join(logDir, "cdsl.log")
	file, err := os.OpenFile(
		logPath,
		os.O_CREATE|os.O_APPEND|os.O_WRONLY,
		0644,
	)
	if err != nil {
		return fmt.Errorf("open log file: %w", err)
	}

	logFile = file
	return nil
}

// Close closes the log file. Call when the program exits.
func Close() error {
	if logFile == nil {
		return nil
	}

	err := logFile.Close()
	logFile = nil
	return err
}

// FormatEntry formats a log entry into a single line.
// Success = true; Failure = false.
func FormatEntry(entry Log) string {
	status := "FAILURE"
	if entry.Success {
		status = "SUCCESS"
	}

	return fmt.Sprintf(
		"[%s][%s][%s][%s]",
		entry.Timestamp.Format("2006-01-02T15:04:05.000-07:00"),
		entry.InstructionString,
		status,
		entry.Message,
	)
}

// LogIt writes a timestamped log entry to the initialized log file.
// success: true = success; false = failure.
func LogIt(success bool, message string, instruction customTypes.Instruction) error {
	if logFile == nil {
		return fmt.Errorf("logging not initialized: call Init first")
	}

	entry := Log{
		Timestamp:         time.Now(),
		Success:           success,
		InstructionString: strings.Join(instruction, " "),
		Message:           message,
	}

	if _, err := fmt.Fprintln(logFile, FormatEntry(entry)); err != nil {
		return fmt.Errorf("write log entry: %w", err)
	}

	return nil
}
