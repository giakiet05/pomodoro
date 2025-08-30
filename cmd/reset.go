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

// resetCmd represents the reset command
var resetCmd = &cobra.Command{
	Use:   "reset",
	Short: "Reset the timer",
	Long:  "Reset the timer",
	Run: func(cmd *cobra.Command, args []string) {
		conn, err := net.Dial("unix", SERVER_ADDR)
		if err != nil {
			fmt.Println("Daemon not running?")
			return
		}

		defer conn.Close()

		command := model.Command{
			Cmd: "reset",
		}

		json.NewEncoder(conn).Encode(&command)

		var resp model.Response
		if err := json.NewDecoder(conn).Decode(&resp); err != nil {
			fmt.Println("Error receiving response from daemon!")
			return
		}

		fmt.Println("Pomodoro has been reset.")

	},
}

func init() {
	rootCmd.AddCommand(resetCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// resetCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// resetCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
