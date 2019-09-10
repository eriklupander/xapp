package persistence

import "github.com/callistaenterprise/xapp/internal/app/model"

type Database interface {
	Persist(tweet model.Tweet) error
	ExistsByURL(url string) bool
	Close()
	Ping() error
	Migrate()
	// Get() interface{}
}
