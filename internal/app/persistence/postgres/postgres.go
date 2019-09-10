package postgres

import (
	"github.com/callistaenterprise/xapp/internal/app/model"
	"github.com/callistaenterprise/xapp/internal/app/persistence"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type postgres struct {
	db *gorm.DB
}

func New(config string) persistence.Database {
	db, err := gorm.Open("postgres", config)
	if err != nil {
		logrus.WithError(err).Fatal("Cannot connect to database")
	}
	return &postgres{
		db,
	}
}

func (p *postgres) ExistsByURL(url string) bool {
	var count = 0
	err := p.db.Model(&model.Tweet{}).Where("url = ?", url).Count(&count).Error
	if err != nil {
		logrus.WithError(err).Fatal("error counting occurance of image by its URL")
	}
	return count > 0
}

func (p *postgres) Migrate() {
	err := p.db.AutoMigrate(&model.Tweet{}).Error
	if err != nil {
		logrus.WithError(err).Fatalf("unable to migrate postgres DB schema")
	}
}

func (p *postgres) Persist(tweet model.Tweet) error {
	tx := p.db.Create(&tweet)
	if tx.Error != nil {
		return errors.Wrap(tx.Error, "error occurred persisting tweet")
	}
	return nil
}

func (p *postgres) Ping() error {
	return p.db.Exec("SELECT 1").Error
}

func (p *postgres) Close() {
	p.db.Close()
}
