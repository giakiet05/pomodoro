package model

type Status int

const (
	StatusRunning Status = iota
	StatusStopped
	StatusAlreadyRunning
	StatusAlreadyStopped
	StatusUnknownCommand
	StatusError
	StatusConfigReloaded
	StatusPhaseDone
)

func (s Status) String() string {
	switch s {
	case StatusRunning:
		return "Running"
	case StatusStopped:
		return "Stopped"
	case StatusAlreadyRunning:
		return "Already Running"
	case StatusAlreadyStopped:
		return "Already Stopped"
	case StatusUnknownCommand:
		return "Unknown Command"
	case StatusError:
		return "Error"
	case StatusConfigReloaded:
		return "Config Reloaded"
	case StatusPhaseDone:
		return "Phase Done"
	default:
		return "Unknown GetStatus"
	}
}
