package livekit

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

type ServerManager struct {
	cmd *exec.Cmd
}

func NewServerManager() *ServerManager {
	return &ServerManager{}
}

func (m *ServerManager) Start(apiKey, apiSecret, port string) error {
	exePath, err := os.Executable()
	if err != nil {
		return err
	}
	// Assume livekit-server.exe is in the same directory as mystia-server
	lkPath := filepath.Join(filepath.Dir(exePath), "livekit-server.exe")
	if _, err := os.Stat(lkPath); os.IsNotExist(err) {
		// Fallback to current working directory
		lkPath = "./livekit-server.exe"
	}

	log.Printf("Starting LiveKit server: %s", lkPath)

	m.cmd = exec.Command(lkPath, "--dev", "--bind", "0.0.0.0", "--port", port, "--node-ip", "127.0.0.1")
	log.Printf("Executing LiveKit command: %v", m.cmd.Args)

	// Set environment variables for API keys
	// Note: --dev mode in LiveKit usually uses default keys, but we want to be explicit if possible.
	// However, --dev mode is very convenient for local setup.
	// To use custom keys, we might need a config file or LIVEKIT_KEYS env var.
	m.cmd.Env = append(os.Environ(),
		fmt.Sprintf("LIVEKIT_KEYS=%s: %s", apiKey, apiSecret),
	)

	m.cmd.Stdout = os.Stdout
	m.cmd.Stderr = os.Stderr

	return m.cmd.Start()
}

func (m *ServerManager) Stop() error {
	if m.cmd != nil && m.cmd.Process != nil {
		log.Printf("Stopping LiveKit server (PID: %d)...", m.cmd.Process.Pid)
		return m.cmd.Process.Kill()
	}
	return nil
}
