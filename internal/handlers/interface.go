package handlers

import "context"

type Server interface {
	Setup(context.Context) error
	Start(context.Context) error
	Stop(context.Context) error
}
