package main

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
	"os/signal"
	"pomodoro/config"
	"pomodoro/model"
	"pomodoro/sound"
	"syscall"
)

func main() {
	Start() //Start the deamon
}

func Start() {
	sound.InitSpeaker() //Initialize sound speaker
	socketPath := "/tmp/pomodoro.sock"
	os.Remove(socketPath) //Remove old socket if exists

	ln, err := net.Listen("unix", socketPath) // Listen to the socket

	if err != nil {
		panic(err)
	}

	defer ln.Close()

	//catch kill signal to remove socket
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-sigs
		fmt.Println("\nShutting down daemon...")
		ln.Close()
		os.Remove(socketPath)
		os.Exit(0)
	}()

	fmt.Println("Pomodoro daemon started at", socketPath)

	//Get config and create pomodoro instance
	cfg, err := config.LoadConfig("config.json")
	p := model.NewPomodoro(cfg) // Create pomodoro instance

	for {
		conn, err := ln.Accept()
		if err != nil {
			continue
		}

		go handleConn(conn, p)
	}
}

func handleConn(conn net.Conn, p *model.Pomodoro) {
	defer conn.Close()

	var cmd model.Command
	if err := json.NewDecoder(conn).Decode(&cmd); err != nil {
		return
	}

	switch cmd.Cmd {
	case "start":
		if p.Status == model.StatusRunning {
			json.NewEncoder(conn).Encode(model.Response{Status: model.StatusAlreadyRunning})
			return
		}
		p.Start()
		json.NewEncoder(conn).Encode(p.GetStatus())
	case "pause":
		if p.Status != model.StatusRunning {
			json.NewEncoder(conn).Encode(model.Response{Status: model.StatusAlreadyStopped})
			return
		}
		p.Pause()
		json.NewEncoder(conn).Encode(p.GetStatus())
	case "reset":
		p.Reset()
		json.NewEncoder(conn).Encode(p.GetStatus())
	case "status":
		json.NewEncoder(conn).Encode(p.GetStatus())
	case "reload-config":
		newCfg, err := config.LoadConfig("config.json")
		if err != nil {
			json.NewEncoder(conn).Encode(model.Response{Status: model.StatusError})
			return
		}
		p.ReloadConfig(newCfg)
		json.NewEncoder(conn).Encode(model.Response{Status: model.StatusConfigReloaded})

	default:
		json.NewEncoder(conn).Encode(model.Response{Status: model.StatusUnknownCommand})
	}
}
