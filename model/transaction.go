package model

import "time"

type Command struct {
	Cmd string `json:"command"`
}

type Response struct {
	Status    Status        `json:"status"`
	Phase     Phase         `json:"phase"`
	Remaining time.Duration `json:"remaining"`
}
