package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"net"
	"os"
	"os/signal"
	"pomodoro/model"
	"strings"
	"syscall"
	"time"
)

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start or continue the timer",
	Long:  "Start the timer or continue from where it paused",
	Run: func(cmd *cobra.Command, args []string) {
		conn, err := net.Dial("unix", SERVER_ADDR)
		watch, _ := cmd.Flags().GetBool("watch")

		if err != nil {
			fmt.Println("Daemon not running?")
			return
		}

		defer conn.Close()

		command := model.Command{
			Cmd: "start",
		}

		json.NewEncoder(conn).Encode(&command)

		var resp model.Response
		if err := json.NewDecoder(conn).Decode(&resp); err != nil {
			fmt.Println("Error receiving response from daemon!")
			return
		}

		if resp.Status == model.StatusAlreadyRunning {
			//Even if already running, show status if watch is enabled
			if !watch {
				fmt.Println("Pomorodo is already running!")
				return
			}
		}

		if watch {
			watchStatus()
			return
		}

		fmt.Printf("Status: %s, Phase: %s, Remaining: %v\n", resp.Status, resp.Phase, resp.Remaining)

	},
}

func askContinue() bool {
	var input string
	fmt.Print("\nDo you want to start the next phase? (Y/n): ")
	fmt.Scanln(&input)

	input = strings.ToLower(strings.TrimSpace(input))
	if input == "" || input == "y" {
		return true // y or enter defaults to yes
	}
	return false
}

func watchStatus() {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)

	// Cleanup function to restore cursor
	defer func() {
		fmt.Print("\033[?25h") // Show cursor
		fmt.Println()
	}()

	fmt.Println("Press Ctrl+C to stop watching...")
	for {
		select {
		case <-ticker.C:
			conn, err := net.Dial("unix", SERVER_ADDR)
			if err != nil {
				fmt.Println("Daemon connection lost!")
				return
			}

			command := model.Command{
				Cmd: "status",
			}

			json.NewEncoder(conn).Encode(&command)

			var resp model.Response
			if err := json.NewDecoder(conn).Decode(&resp); err != nil {
				fmt.Println("Error receiving response from daemon!")
				return
			}
			conn.Close()

			fmt.Printf("\rStatus: %-15s Phase: %-12s Remaining: %v  ", resp.Status, resp.Phase, resp.Remaining)

			//handle phase completion
			if resp.Status == model.StatusPhaseDone {
				//If user agrees, start next phase automatically
				if askContinue() == true {
					conn, err := net.Dial("unix", SERVER_ADDR)
					if err != nil {
						fmt.Println("Daemon connection lost!")
						return
					}

					command := model.Command{
						Cmd: "start",
					}

					json.NewEncoder(conn).Encode(&command)

					var resp model.Response
					if err := json.NewDecoder(conn).Decode(&resp); err != nil {
						fmt.Println("Error receiving response from daemon!")
						return
					}
					conn.Close()
				} else {
					fmt.Println("Pomodoro stopped. Use 'start' to begin next phase.\n")
					return
				}
			}

		case <-sigCh:
			fmt.Println("\nStopped watching")
			return
		}

	}
}

func init() {
	rootCmd.AddCommand(startCmd)
	startCmd.Flags().BoolP("watch", "w", false, "Watch the timer in real-time")

}
