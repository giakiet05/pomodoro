/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"net"
	"pomodoro/model"
)

// pauseCmd represents the pause command
var pauseCmd = &cobra.Command{
	Use:   "pause",
	Short: "Pause the timer",
	Long:  "Pause the timer",
	Run: func(cmd *cobra.Command, args []string) {
		conn, err := net.Dial("unix", SERVER_ADDR)
		if err != nil {
			fmt.Println("Daemon not running?")
			return
		}

		defer conn.Close()
		command := model.Command{
			Cmd: "pause",
		}

		json.NewEncoder(conn).Encode(&command)

		var resp model.Response
		if err := json.NewDecoder(conn).Decode(&resp); err != nil {
			fmt.Println("Error receiving response from daemon!")
			return
		}

		if resp.Status == model.StatusAlreadyStopped {
			fmt.Println("Pomorodo is already stopped!")
			return
		}

		fmt.Printf("Status: %s, Phase: %s, Remaining: %v\n", resp.Status, resp.Phase, resp.Remaining)

	},
}

func init() {
	rootCmd.AddCommand(pauseCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// pauseCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// pauseCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
