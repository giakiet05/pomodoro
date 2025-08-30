package model

import (
	"context"
	"fmt"
	"pomodoro/config"
	"pomodoro/sound"
	"time"
)

type Pomodoro struct {
	Status    Status             // Current status
	Phase     Phase              // Current phase
	WorkCount int                // Number of completed work sessions
	Duration  time.Duration      // Total duration of the current phase
	Remaining time.Duration      // Time left in the current phase
	ticker    *time.Ticker       // Tick every second for countdown
	ctx       context.Context    // Control goroutine lifecycle
	cancel    context.CancelFunc // Cancels ctx to stop goroutines
	Config    *config.Config     // Configuration settings
}

func NewPomodoro(cfg *config.Config) *Pomodoro {
	return &Pomodoro{
		Config:    cfg,
		Phase:     PhaseWork,
		WorkCount: 0,
		Duration:  time.Duration(cfg.Work) * time.Second,
		Remaining: time.Duration(cfg.Work) * time.Second,
		Status:    StatusStopped,
	}
}

func (p *Pomodoro) Start() {
	if p.Status == StatusRunning {
		return
	}
	p.ctx, p.cancel = context.WithCancel(context.Background())
	p.ticker = time.NewTicker(time.Second)
	p.Status = StatusRunning

	go func() {
		for {
			select {
			case <-p.ticker.C:
				p.Remaining -= time.Second
				fmt.Printf("\rRemaining: %v", p.Remaining)
				if p.Remaining <= 0 {
					fmt.Println("\nTime's up")
					go sound.PlaySound("time_up.mp3")
					p.Status = StatusPhaseDone
					p.ticker.Stop()
					p.cancel()
					p.nextPhase()
					return
				}
			case <-p.ctx.Done():
				return
			}
		}
	}()

	if p.Duration == p.Remaining && p.WorkCount == 0 {
		fmt.Println("\nStarting...")
	} else if p.Duration != p.Remaining {
		fmt.Println("\nResuming...")
	}
	fmt.Println(p.Phase.String()) //Print current phase
}

func (p *Pomodoro) Pause() {
	if p.Status != StatusRunning {
		return
	}
	p.Status = StatusStopped
	p.ticker.Stop()
	p.cancel()
	fmt.Println("\nPaused")
}

func (p *Pomodoro) Reset() {
	if p.ticker != nil {
		p.ticker.Stop()
	}

	p.Phase = PhaseWork
	p.WorkCount = 0
	p.Duration = time.Duration(p.Config.Work) * time.Second
	p.Remaining = p.Duration
	p.Status = StatusStopped
	p.cancel()
	fmt.Println("\nResetted")

}

func (p *Pomodoro) nextPhase() {
	switch p.Phase {
	case PhaseWork:
		p.WorkCount++
		if p.WorkCount == 4 {
			p.Phase = PhaseLongBreak
			p.Duration = time.Duration(p.Config.LongBreak) * time.Second
		} else {
			p.Phase = PhaseShortBreak
			p.Duration = time.Duration(p.Config.ShortBreak) * time.Second
		}
	case PhaseShortBreak, PhaseLongBreak:
		p.Phase = PhaseWork
		p.Duration = time.Duration(p.Config.Work) * time.Second
	}
	p.Remaining = p.Duration
}

func (p *Pomodoro) GetStatus() Response {
	return Response{
		Status:    p.Status,
		Phase:     p.Phase,
		Remaining: p.Remaining,
	}
}

func (p *Pomodoro) ReloadConfig(cfg *config.Config) {
	p.Config = cfg

	// Cập nhật lại duration cho phase hiện tại
	switch p.Phase {
	case PhaseWork:
		p.Duration = time.Duration(cfg.Work) * time.Second
	case PhaseShortBreak:
		p.Duration = time.Duration(cfg.ShortBreak) * time.Second
	case PhaseLongBreak:
		p.Duration = time.Duration(cfg.LongBreak) * time.Second
	}

	// Reset remaining cho chắc (tùy mày muốn reset hay giữ nguyên progress)
	p.Remaining = p.Duration
}
