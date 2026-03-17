// Package sound provides audio playback for pomo sessions.
package sound

import (
	"context"
	"log"
	"os/exec"
	"sync"
)

// Player handles audio playback with proper cleanup.
type Player struct {
	cancel context.CancelFunc
	wg     *sync.WaitGroup
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
// Uses tight loop with minimal delay for near-gapless playback.
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

// PlayCommandLoop starts looping with a full command (e.g., mpv --loop URL).
func (p *Player) PlayCommandLoop(cmd []string) {
	if len(cmd) < 1 {
		return
	}

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
				c := exec.CommandContext(ctx, cmd[0], cmd[1:]...)
				_ = c.Run()
			}
		}
	}()
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
