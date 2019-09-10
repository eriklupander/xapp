package worker

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/callistaenterprise/xapp/internal/app/filehandler"
	"github.com/callistaenterprise/xapp/internal/app/imageloader"
	"github.com/callistaenterprise/xapp/internal/app/imageprocessor"
	"github.com/callistaenterprise/xapp/internal/app/model"
	"github.com/callistaenterprise/xapp/internal/app/persistence"
	"github.com/dghubble/go-twitter/twitter"
	"github.com/sirupsen/logrus"
	"image"
	"image/png"
	"strings"
	"time"
)

type TweetWorker struct {
	imgProcessor imageprocessor.ImageProcessor
	db           persistence.Database
	filehandler  filehandler.FileHandler
	imageLoader  imageloader.ImageLoader
	tweetChan    chan *twitter.Tweet
}

func NewTweetWorker(imgProcessor imageprocessor.ImageProcessor, db persistence.Database, filehandler filehandler.FileHandler, imageLoader imageloader.ImageLoader, tweetChan chan *twitter.Tweet) *TweetWorker {
	return &TweetWorker{imgProcessor: imgProcessor, db: db, filehandler: filehandler, imageLoader: imageLoader, tweetChan: tweetChan}
}

func (tw *TweetWorker) Start() {

	// start 8 workers
	for i := 0; i < 8; i++ {
		go tw.process()
	}
}

func (tw *TweetWorker) process() {
	logrus.Debug("Worker running...")

	for twt := range tw.tweetChan {
		tw.processTweet(twt)
	}
}

func (tw *TweetWorker) filter(tweet *twitter.Tweet) (model.Tweet, error) {
	// Perform some extra filtering.
	if tweet.PossiblySensitive || tweet.InReplyToStatusID != 0 || tweet.Text[0:2] == "RT" {
		return model.Tweet{}, errors.New("filtered out")
	}

	// Check if tweet has an image
	if tweet.Entities != nil {
		if len(tweet.Entities.Media) > 0 {

			// ConsumeStream each media as a message
			for _, media := range tweet.Entities.Media {
				if media.Type == "photo" {

					created, _ := tweet.CreatedAtTime()
					return model.Tweet{
						Author:    tweet.User.ScreenName,
						Text:      tweet.Text,
						CreatedAt: created,
						URL:       media.MediaURL,
					}, nil
				}
			}
		}
	}
	return model.Tweet{}, errors.New("filtered out")
}

func (tw *TweetWorker) processTweet(twt *twitter.Tweet) {
	start := time.Now()

	tweet, err := tw.filter(twt)
	if err != nil {
		return
	}

	// Check if already fetched
	if tw.db.ExistsByURL(tweet.URL) {
		logrus.Infof("image '%v' already fetched", tweet.URL)
		return
	}

	// Fetch image
	data, err := tw.imageLoader.Load(tweet.URL)
	if err != nil {
		logrus.WithError(err).Errorf("unable to load image from url %v", tweet.URL)
		return
	}

	var imgData image.Image
	var decodeErr error
	suffix := "jpeg"
	if strings.HasSuffix(tweet.URL, ".png") {
		suffix = "png"
		imgData, decodeErr = png.Decode(bytes.NewReader(data))
	} else {
		imgData, suffix, decodeErr = image.Decode(bytes.NewReader(data))
	}

	if decodeErr != nil {
		logrus.WithError(decodeErr).Errorf("unable to decode image. URL: %v", tweet.URL)
		return
	}

	// Apply filter
	out := &bytes.Buffer{}
	err = tw.imgProcessor.Hipsterize(imgData, out)
	if err != nil {
		logrus.WithError(err).Fatalf("unable to apply Hipsterize filter")
		return
	}

	// ConsumeStream image to disk...
	err = tw.filehandler.Write(fmt.Sprintf("%v_%v.%v", tweet.Author, time.Now().Unix(), suffix), out.Bytes())
	if err != nil {
		logrus.WithError(err).Fatalf("unable to store converted image on disk")
		return
	}

	// Store tweet in DB
	err = tw.db.Persist(tweet)
	if err != nil {
		logrus.WithError(err).Fatalf("unable to store tweet data in DB")
		return
	}

	logrus.WithField("duration", time.Now().Sub(start)).
		WithField("suffix", suffix).
		Info("successfully processed an image!")
}
