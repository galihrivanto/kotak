package module

import (
	"context"
	"database/sql"
	"io"
)

type Module interface {
	io.Closer
	Start(context.Context)
}

type ModuleFactory func(*sql.DB) Module

var modules = map[string]ModuleFactory{}

func RegisterModule(name string, factory ModuleFactory) {
	modules[name] = factory
}
