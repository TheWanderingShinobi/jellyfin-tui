package main

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/TheWanderingShinobi/jellyfin-tui/internal/config"
	"github.com/TheWanderingShinobi/jellyfin-tui/internal/jellyfin"
	"github.com/TheWanderingShinobi/jellyfin-tui/internal/ui"
	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	if err := checkDependencies(); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	cfg := config.Load()
	client := jellyfin.NewClient(cfg.ServerURL)

	m := ui.NewModel(client, cfg)
	p := tea.NewProgram(m)

	if err := p.Start(); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}

func checkDependencies() error {
	// Check if MPV is installed
	_, err := exec.LookPath("mpv")
	if err != nil {
		return fmt.Errorf("MPV is not installed or not in PATH. Please install MPV to use this application")
	}
	return nil
}
