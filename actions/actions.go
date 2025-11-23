// Package actions provides functionality to run post actions after a task is completed.
package actions

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"sync"
	"time"

	"github.com/Bahaaio/pomo/config"
	"github.com/gen2brain/beeep"
)

// RunPostActions sends task notification and runs post commands using goroutines
//
// returns a wait group to wait for their completion
func RunPostActions(task *config.Task) *sync.WaitGroup {
	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		sendNotification(task.Notification)
	}()

	go func() {
		defer wg.Done()
		runPostCommands(task.Then)
	}()

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
func runPostCommands(cmds [][]string) {
	log.Println("running post commands")

	for _, cmd := range cmds {
		c := exec.Command(cmd[0], cmd[1:]...)

		if err := c.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "failed to run command '%q': %v\n", cmd, err)
		}

		// wait some time before running the next command
		time.Sleep(50 * time.Millisecond)
	}
}
