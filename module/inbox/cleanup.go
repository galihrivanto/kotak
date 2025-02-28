package inbox

import (
	"context"
	"time"

	"github.com/galihrivanto/kotak/config"
	"github.com/galihrivanto/kotak/db"
	"github.com/galihrivanto/kotak/log"
	"github.com/galihrivanto/kotak/module"
)

type Cleanup struct {
	ctx    context.Context
	cancel context.CancelFunc
	db     *db.DB
	cfg    *config.Config
}

func (c *Cleanup) Start(ctx context.Context) error {
	c.ctx, c.cancel = context.WithCancel(ctx)

	interval := c.cfg.Inbox.CleanupInterval
	if interval <= 0 {
		interval = 5 * time.Minute
	}

	age := c.cfg.Inbox.MaxAge
	if age <= 0 {
		age = 24 * time.Hour
	}

	go func() {
		for {
			select {
			case <-c.ctx.Done():
				return
			case <-time.After(interval):
				log.Info("Cleaning up inbox with age %vs", age)
				if err := c.db.Cleanup(int(age.Hours())); err != nil {
					log.Error("Failed to cleanup inbox: %v", err)
				}
			}
		}
	}()

	return nil
}

func (c *Cleanup) Close() error {
	c.cancel()
	return nil
}

func NewCleanup(cfg *config.Config, db *db.DB) *Cleanup {
	return &Cleanup{cfg: cfg, db: db}
}

func init() {
	module.RegisterModule("inbox_cleanup", func(config *config.Config, db *db.DB) module.Module {
		return NewCleanup(config, db)
	})
}
