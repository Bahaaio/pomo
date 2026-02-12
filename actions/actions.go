// Package actions provides functionality to run post actions after a task is completed.
package actions

import (
	"context"
	"log"
	"os/exec"
	"sync"
	"time"

	"github.com/Bahaaio/pomo/config"
	"github.com/gen2brain/beeep"
)

var CommandTimeout = 5 * time.Second

// RunPostActions sends task notification and runs post commands using goroutines
//
// returns a wait group to wait for their completion
func RunPostActions(ctx context.Context, task *config.Task) *sync.WaitGroup {
	var wg sync.WaitGroup

	wg.Go(func() {
		sendNotification(task.Notification)
	})

	wg.Go(func() {
		runPostCommands(ctx, task.Then)
	})

	return &wg
}

// sends a notification using the beeep package
func sendNotification(notification config.Notification) {
	if !notification.Enabled {
		log.Println("notification disabled")
		return
	}

	log.Println("sending notification")

	// use the embedded icon
	var icon any = config.Icon

	// if the user has specified an icon
	// use that instead
	if len(notification.Icon) > 0 {
		icon = notification.Icon
	}

	var err error
	if notification.Urgent {
		err = beeep.Alert(notification.Title, notification.Message, icon)
	} else {
		err = beeep.Notify(notification.Title, notification.Message, icon)
	}

	if err != nil {
		log.Println("failed to send notification:", err)
	}
}

// runs the post commands specified in the task
func runPostCommands(ctx context.Context, cmds [][]string) {
	log.Println("running post commands")

	for _, cmd := range cmds {
		c := exec.CommandContext(ctx, cmd[0], cmd[1:]...)

		if err := c.Run(); err != nil {
			// TODO: show error message
			log.Printf("failed to run command '%q': %v\n", cmd, err)
		}
	}
}
