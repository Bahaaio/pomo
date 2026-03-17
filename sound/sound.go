// Package sound provides audio playback for pomo sessions.
package sound

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
	"os/exec"
	"sync"
	"time"
)

// Player handles audio playback with proper cleanup.
type Player struct {
	cancel   context.CancelFunc
	wg       *sync.WaitGroup
	mpvCmd   *exec.Cmd
	mpvPause bool // track pause state
	mu       sync.Mutex
}

// NewPlayer creates a new sound player.
func NewPlayer() *Player {
	return &Player{}
}

// PlayOnce plays a sound file once (fire and forget).
func PlayOnce(soundFile string) {
	go func() {
		cmd := exec.Command("afplay", soundFile)
		if err := cmd.Run(); err != nil {
			log.Printf("failed to play sound %q: %v", soundFile, err)
		}
	}()
}

// PlayCommandOnce plays a command once (fire and forget).
func PlayCommandOnce(cmd []string) {
	go func() {
		c := exec.Command(cmd[0], cmd[1:]...)
		if err := c.Run(); err != nil {
			log.Printf("failed to play command %q: %v", cmd, err)
		}
	}()
}

// PlayLoop starts looping a sound file until Stop is called.
func (p *Player) PlayLoop(soundFile string) {
	ctx, cancel := context.WithCancel(context.Background())
	p.cancel = cancel
	p.wg = &sync.WaitGroup{}

	p.wg.Add(1)
	go func() {
		defer p.wg.Done()

		for {
			select {
			case <-ctx.Done():
				return
			default:
				cmd := exec.CommandContext(ctx, "afplay", soundFile)
				_ = cmd.Run()
			}
		}
	}()
}

// PlayCommandLoop starts looping with a full command using mpv with IPC support.
// The command should be in format: [mpv, --no-video, URL]
func (p *Player) PlayCommandLoop(cmd []string) {
	if len(cmd) < 1 {
		return
	}

	ctx, cancel := context.WithCancel(context.Background())
	p.cancel = cancel
	p.wg = &sync.WaitGroup{}

	// Create unique IPC socket path
	socketPath := fmt.Sprintf("/tmp/mpv-pomo-%d.sock", time.Now().UnixNano())
	os.Remove(socketPath)

	// Build mpv command with IPC socket
	mpvArgs := append([]string{"--no-video", "--loop", "--input-ipc-server=" + socketPath}, cmd[1:]...)
	p.mpvCmd = exec.CommandContext(ctx, cmd[0], mpvArgs...)

	p.wg.Add(1)
	go func() {
		defer p.wg.Done()
		_ = p.mpvCmd.Run()
		// Clean up socket
		os.Remove(socketPath)
	}()
}

// Pause pauses the mpv playback
func (p *Player) Pause() error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.mpvCmd == nil || p.mpvCmd.Process == nil {
		return fmt.Errorf("no mpv process running")
	}

	// Find the IPC socket
	socketPath := p.getSocketPath()
	if socketPath == "" {
		return fmt.Errorf("no IPC socket found")
	}

	conn, err := net.Dial("unix", socketPath)
	if err != nil {
		return fmt.Errorf("failed to connect to mpv IPC: %w", err)
	}
	defer conn.Close()

	// Send pause command
	cmd := map[string]interface{}{"command": []string{"cycle", "pause"}}
	data, _ := json.Marshal(cmd)
	_, err = conn.Write(append(data, '\n'))
	if err != nil {
		return fmt.Errorf("failed to send pause command: %w", err)
	}

	p.mpvPause = true
	return nil
}

// Resume resumes the mpv playback
func (p *Player) Resume() error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if !p.mpvPause {
		return nil // already playing
	}

	socketPath := p.getSocketPath()
	if socketPath == "" {
		return fmt.Errorf("no IPC socket found")
	}

	conn, err := net.Dial("unix", socketPath)
	if err != nil {
		return fmt.Errorf("failed to connect to mpv IPC: %w", err)
	}
	defer conn.Close()

	// Send pause command (cycle to unpause)
	cmd := map[string]interface{}{"command": []string{"cycle", "pause"}}
	data, _ := json.Marshal(cmd)
	_, err = conn.Write(append(data, '\n'))
	if err != nil {
		return fmt.Errorf("failed to send resume command: %w", err)
	}

	p.mpvPause = false
	return nil
}

// getSocketPath tries to find the mpv IPC socket
func (p *Player) getSocketPath() string {
	// Simple approach: try the most recent socket
	entries, _ := os.ReadDir("/tmp")
	var latest string
	var latestTime int64

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		if len(entry.Name()) > 12 && entry.Name()[:12] == "mpv-pomo-" {
			info, _ := entry.Info()
			if info != nil && info.ModTime().Unix() > latestTime {
				latest = "/tmp/" + entry.Name()
				latestTime = info.ModTime().Unix()
			}
		}
	}

	return latest
}

// Stop stops any looping sound and waits for cleanup.
func (p *Player) Stop() {
	if p.cancel != nil {
		p.cancel()
	}
	if p.wg != nil {
		p.wg.Wait()
	}
	p.cancel = nil
	p.wg = nil
}
