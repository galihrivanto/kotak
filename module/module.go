package module

import (
	"context"
	"io"

	"github.com/galihrivanto/kotak/config"
	"github.com/galihrivanto/kotak/db"
)

type Module interface {
	io.Closer
	Start(context.Context) error
}

type ModuleFactory func(*config.Config, *db.DB) Module

var (
	modules        = map[string]ModuleFactory{}
	moduleInstance = map[string]Module{}
)

func RegisterModule(name string, factory ModuleFactory) {
	modules[name] = factory
}

func Start(ctx context.Context, cfg *config.Config, db *db.DB) error {
	for name, factory := range modules {
		module := factory(cfg, db)
		if err := module.Start(ctx); err != nil {
			return err
		}
		moduleInstance[name] = module
	}
	return nil
}

func Stop() error {
	for _, module := range moduleInstance {
		module.Close()
	}
	return nil
}
