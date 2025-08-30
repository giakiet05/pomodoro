/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"net"
	"pomodoro/config"
	"pomodoro/model"
)

const CONFIG_PATH = "config.json"

var (
	work       int
	shortBreak int
	longBreak  int
)

// configCmd represents the config command
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Configure the pomodoro app",
	Run: func(cmd *cobra.Command, args []string) {
		//Check if daemon is running
		conn, err := net.Dial("unix", SERVER_ADDR)
		if err != nil {
			fmt.Println("Daemon not running?")
			return
		}
		defer conn.Close()

		cfg, err := config.LoadConfig(CONFIG_PATH)
		if err != nil {
			fmt.Println("Error loading config:", err)
			return
		}

		workChanged := cmd.Flags().Changed("work")
		shortBreakChanged := cmd.Flags().Changed("short-break")
		longBreakChanged := cmd.Flags().Changed("long-break")

		if workChanged {
			cfg.Work = work
		}
		if shortBreakChanged {
			cfg.ShortBreak = shortBreak
		}
		if longBreakChanged {
			cfg.LongBreak = longBreak
		}

		if workChanged || shortBreakChanged || longBreakChanged {
			fmt.Println("Configuration updated successfully.")
		}

		fmt.Printf("Current Configuration:\nWork: %d minutes\nShort Break: %d minutes\nLong Break: %d minutes\n", cfg.Work, cfg.ShortBreak, cfg.LongBreak)

		if err := config.SaveConfig(CONFIG_PATH, cfg); err != nil {
			fmt.Println("Error saving config:", err)
			return
		}

		//Notify the daemon to reload config
		command := model.Command{
			Cmd: "reload-config",
		}

		json.NewEncoder(conn).Encode(&command)
		var resp model.Response
		if err := json.NewDecoder(conn).Decode(&resp); err != nil {
			fmt.Println("Error receiving response from daemon!")
			return
		}

	},
}

func init() {
	rootCmd.AddCommand(configCmd)
	configCmd.Flags().IntVarP(&work, "work", "w", 25, "Work duration in minutes")
	configCmd.Flags().IntVarP(&shortBreak, "short-break", "s", 5, "Short break duration in minutes")
	configCmd.Flags().IntVarP(&longBreak, "long-break", "l", 15, "Long break duration in minutes")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// configCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// configCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
