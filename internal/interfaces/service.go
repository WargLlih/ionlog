package interfaces

type IService interface {
	Status() ServiceStatus
	Start() error
	Stop()
}

type ServiceStatus int

const (
	Stopped ServiceStatus = iota
	Running
)
