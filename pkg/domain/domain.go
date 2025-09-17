package domain

import (
	"context"

	"github.com/rs/zerolog"
	"github.com/tcncloud/sati-go/pkg/interfaces"
)

// Domain is the main domain object for the application.
// It implements 5 services:
// - ConfigWatcher, it instantiate a config watcher that notify the application when config file is present and/or updated.
type Domain struct {
	log           *zerolog.Logger
	configWatcher interfaces.ConfigWatcher
}

func NewDomain(log *zerolog.Logger) *Domain {
	return &Domain{
		log: log,
	}
}

// SetConfigWatcher sets the configuration watcher for the domain.
func (d *Domain) SetConfigWatcher(watcher interfaces.ConfigWatcher) {
	d.configWatcher = watcher
}

func (d *Domain) StartConfigWatcher(ctx context.Context) error {
	if d.configWatcher == nil {
		return nil // No watcher configured
	}
	return d.configWatcher.Start(ctx)
}

func (d *Domain) StartGateClient() {

}

func (d *Domain) StartPollEvents() {

}

func (d *Domain) StartStreamJobs() {

}
