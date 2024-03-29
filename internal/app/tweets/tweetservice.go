package tweets

import (
	"github.com/callistaenterprise/xapp/cmd"
	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
	"github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Streamer struct {
	client        *twitter.Client
	tweetsChannel chan *twitter.Tweet
}

func NewStreamer(cfg cmd.TwitterConfig, tweetsChannel chan *twitter.Tweet) *Streamer {
	logrus.Info("Starting TwitterStreamer...")
	if cfg.ConsumerKey == "" || cfg.ConsumerSecret == "" || cfg.AccessToken == "" || cfg.AccessSecret == "" {
		logrus.Fatal("one or more twitter API tokens were nil or empty")
	}
	config := oauth1.NewConfig(cfg.ConsumerKey, cfg.ConsumerSecret)
	token := oauth1.NewToken(cfg.AccessToken, cfg.AccessSecret)
	httpClient := config.Client(oauth1.NoContext, token)

	// Twitter Client
	client := twitter.NewClient(httpClient)
	return &Streamer{client: client, tweetsChannel: tweetsChannel}
}

func (s *Streamer) ConsumeStream(keyword string) {

	var counter = 0
	countChan := make(chan int)
	// Convenience Demux demultiplexed stream messages
	demux := twitter.NewSwitchDemux()
	demux.Tweet = func(tweet *twitter.Tweet) {
		// Pass incoming tweets to our own tweetschannel
		counter++
		countChan <- counter
		s.tweetsChannel <- tweet
	}

	go func() {

		printChan := time.Tick(time.Second * 5)
		last := 0
		for {
			select {
			case _ = <-printChan:
				logrus.Infof("%v tweets processed", last)
			case last = <-countChan:
			}
		}
	}()

	logrus.Infof("Starting Stream for keyword %v...", keyword)

	// FILTER
	filterParams := &twitter.StreamFilterParams{
		Track:         []string{keyword, "exclude:retweets", "exclude:replies"},
		Language:      []string{"en"},
		StallWarnings: twitter.Bool(true),
	}
	stream, err := s.client.Streams.Filter(filterParams)
	if err != nil {
		logrus.WithError(err).Fatal("error applying filter to stream")
	}

	// Receive messages until stopped or stream quits
	go demux.HandleChan(stream.Messages)

	// Block here!
	s.awaitSigterm(stream)
}

func (s *Streamer) awaitSigterm(stream *twitter.Stream) {
	termChan := make(chan os.Signal)
	signal.Notify(termChan, syscall.SIGINT, syscall.SIGTERM)

	<-termChan
	stream.Stop()
	logrus.Info("Stopped stream")
}
