/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"encoding/json"
	"fmt"
	"net"
	"pomodoro/model"

	"github.com/spf13/cobra"
)

// statusCmd represents the status command
var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "See pomodoro timer status",
	Run: func(cmd *cobra.Command, args []string) {
		watch, _ := cmd.Flags().GetBool("watch")
		conn, err := net.Dial("unix", SERVER_ADDR)
		if err != nil {
			fmt.Println("Daemon not running?")
			return
		}

		defer conn.Close()

		command := model.Command{
			Cmd: "status",
		}

		json.NewEncoder(conn).Encode(&command)
		var resp model.Response
		if err := json.NewDecoder(conn).Decode(&resp); err != nil {
			fmt.Println("Error receiving response from daemon!")
			return
		}

		if watch {
			watchStatus()
			return
		}

		fmt.Printf("Status: %s, Phase: %s, Remaining: %v\n", resp.Status, resp.Phase, resp.Remaining)

	},
}

func init() {
	rootCmd.AddCommand(statusCmd)
	statusCmd.Flags().BoolP("watch", "w", false, "Watch the status until interrupted")
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// statusCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// statusCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
