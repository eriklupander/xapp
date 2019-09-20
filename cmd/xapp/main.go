package main

import (
	"github.com/alexflint/go-arg"
	"github.com/callistaenterprise/xapp/internal/app/filehandler/local"
	"github.com/callistaenterprise/xapp/internal/app/httprouter"
	"github.com/callistaenterprise/xapp/internal/app/imageprocessor"
	"github.com/callistaenterprise/xapp/internal/app/persistence/postgres"
	"github.com/callistaenterprise/xapp/internal/app/tweets"
	"github.com/callistaenterprise/xapp/internal/app/worker"
	"github.com/dghubble/go-twitter/twitter"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"os"
	"os/signal"
	"syscall"
	"time"
)
import "github.com/callistaenterprise/xapp/cmd"
import "github.com/sirupsen/logrus"

func main() {
	logrus.Info("starting xapp")

	// Apply config from envvars, args etc.
	cfg := cmd.DefaultConfiguration()
	arg.MustParse(cfg)

	// Create Postgres storage and run auto-migration
	db := postgres.New(cfg.PostgresDSN)
	db.Migrate()

	// Create image processor
	imageProcessor := imageprocessor.NewGiftImageProcessor(nil)

	// Create buffered channel to pass tweets from twitter lib to workers
	tweetChan := make(chan *twitter.Tweet, 1)

	// Start tweet stream service, getting cats for us!
	tweetService := tweets.NewStreamer(cfg.TwitterConfig, tweetChan)
	go tweetService.ConsumeStream("cat")

	// Start workers
	workers := worker.NewTweetWorker(imageProcessor, db, local.NewDiskStorage(cfg.ImageFolder), tweetChan)
	workers.Start()

	// Start HTTP server
	srv := httprouter.NewServer(cfg, db)
	srv.SetupRoutes()
	go srv.Start()

	// wait until sigterm.
	awaitSigterm()
}

func awaitSigterm() {
	termChan := make(chan os.Signal)
	signal.Notify(termChan, syscall.SIGINT, syscall.SIGTERM)
	<-termChan
	logrus.Info("Shutdown initiated, waiting for workers to finish...")
	time.Sleep(time.Second * 5)
	logrus.Info("All done, shutting down!")
}
