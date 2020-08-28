package rpc

import "google.golang.org/grpc"

func New() *Registers {
	return &Registers{Registers: make([]Register, 0)}
}

type (
	Register interface {
		Register(grpcServer *grpc.Server)
	}

	Registers struct {
		Registers []Register
	}
)

func (r *Registers) Register(grpcServer *grpc.Server) {

	for _, register := range r.Registers {
		register.Register(grpcServer)
	}

}

func (r *Registers) Add(register Register) {
	r.Registers = append(r.Registers, register)
}
