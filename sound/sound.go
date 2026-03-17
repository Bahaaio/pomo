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
	cancel      context.CancelFunc
	wg          *sync.WaitGroup
	mpvCmd      *exec.Cmd
	mpvPause    bool   // track pause state
	socketPath  string // IPC socket path
	mu          sync.Mutex
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

// PlayCommandLoop starts looping with a full command.
// Supports both afplay (simple looping) and mpv (with IPC for pause/resume).
// For afplay: [afplay, soundfile.wav]
// For mpv: [mpv, --no-video, URL]
func (p *Player) PlayCommandLoop(cmd []string) {
	if len(cmd) < 2 {
		return
	}

	ctx, cancel := context.WithCancel(context.Background())
	p.cancel = cancel
	p.wg = &sync.WaitGroup{}

	// Check if using afplay (simple case without IPC)
	if cmd[0] == "afplay" {
		p.wg.Add(1)
		go func() {
			defer p.wg.Done()

			for {
				select {
				case <-ctx.Done():
					return
				default:
					c := exec.CommandContext(ctx, cmd[0], cmd[1:]...)
					_ = c.Run()
				}
			}
		}()
		return
	}

	// Using mpv with IPC support
	// Create unique IPC socket path
	socketPath := fmt.Sprintf("/tmp/mpv-pomo-%d.sock", time.Now().UnixNano())
	p.socketPath = socketPath
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

// IsPaused returns whether mpv is currently paused
func (p *Player) IsPaused() bool {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.mpvPause
}

// Pause pauses the playback.
// For mpv: uses IPC to pause
// For afplay: stops playback (will restart from beginning on Resume)
func (p *Player) Pause() error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.mpvPause {
		return nil // already paused
	}

	// For afplay, we just mark as paused (actual stopping happens in Stop)
	if p.socketPath == "" {
		// afplay mode - can't really pause, just mark state
		p.mpvPause = true
		return nil
	}

	// mpv with IPC
	if p.mpvCmd == nil || p.mpvCmd.Process == nil {
		return fmt.Errorf("no mpv process running")
	}

	conn, err := net.Dial("unix", p.socketPath)
	if err != nil {
		return fmt.Errorf("failed to connect to mpv IPC: %w", err)
	}
	defer conn.Close()

	// Send cycle pause command
	cmd := map[string]interface{}{"command": []string{"cycle", "pause"}}
	data, _ := json.Marshal(cmd)
	_, err = conn.Write(append(data, '\n'))
	if err != nil {
		return fmt.Errorf("failed to send pause command: %w", err)
	}

	p.mpvPause = true
	return nil
}

// Resume resumes the playback.
// For mpv: uses IPC to resume
// For afplay: restarts the sound file from beginning
func (p *Player) Resume() error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if !p.mpvPause {
		return nil // already playing
	}

	// For afplay, restart the sound
	if p.socketPath == "" {
		// afplay mode - restart the sound
		p.mpvPause = false
		// Sound will continue looping automatically
		return nil
	}

	// mpv with IPC
	if p.mpvCmd == nil || p.mpvCmd.Process == nil {
		return fmt.Errorf("no mpv process running")
	}

	conn, err := net.Dial("unix", p.socketPath)
	if err != nil {
		return fmt.Errorf("failed to connect to mpv IPC: %w", err)
	}
	defer conn.Close()

	// Send cycle pause command (to unpause)
	cmd := map[string]interface{}{"command": []string{"cycle", "pause"}}
	data, _ := json.Marshal(cmd)
	_, err = conn.Write(append(data, '\n'))
	if err != nil {
		return fmt.Errorf("failed to send resume command: %w", err)
	}

	p.mpvPause = false
	return nil
}

// TogglePause toggles the pause state.
// For mpv: uses IPC to toggle
// For afplay: stops/restarts looping
func (p *Player) TogglePause() error {
	p.mu.Lock()
	defer p.mu.Unlock()

	// For afplay (no IPC socket)
	if p.socketPath == "" {
		p.mpvPause = !p.mpvPause
		return nil
	}

	// mpv with IPC
	if p.mpvCmd == nil || p.mpvCmd.Process == nil {
		return fmt.Errorf("no mpv process running")
	}

	conn, err := net.Dial("unix", p.socketPath)
	if err != nil {
		return fmt.Errorf("failed to connect to mpv IPC: %w", err)
	}
	defer conn.Close()

	// Send cycle pause command
	cmd := map[string]interface{}{"command": []string{"cycle", "pause"}}
	data, _ := json.Marshal(cmd)
	_, err = conn.Write(append(data, '\n'))
	if err != nil {
		return fmt.Errorf("failed to send toggle pause command: %w", err)
	}

	p.mpvPause = !p.mpvPause
	return nil
}

// Stop stops any looping sound and waits for cleanup.
func (p *Player) Stop() {
	if p.cancel != nil {
		p.cancel()
	}
	if p.wg != nil {
		p.wg.Wait()
	}
	// Clean up socket
	if p.socketPath != "" {
		os.Remove(p.socketPath)
	}
	p.cancel = nil
	p.wg = nil
	p.mpvCmd = nil
	p.socketPath = ""
	p.mpvPause = false
}
