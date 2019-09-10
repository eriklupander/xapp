package cmd

type Config struct {
	Environment string `arg:"env:ENVIRONMENT"`
	ImageFolder string `arg:"env:IMAGE_FOLDER"`
	HTTPServerConfig
	TwitterConfig
	DGraphConfig
	PostgresConfig
}

type HTTPServerConfig struct {
	BindAddress string `arg:"env:HTTP_BIND_ADDRESS"`
}

type DGraphConfig struct {
	DGraphDSN string `arg:"env:DGRAPH_DSN"`
}

type PostgresConfig struct {
	PostgresDSN string `arg:"env:POSTGRES_DSN"`
}

type TwitterConfig struct {
	ConsumerKey    string `arg:"env:TWITTER_CONSUMER_KEY"`
	ConsumerSecret string `arg:"env:TWITTER_CONSUMER_SECRET"`
	AccessToken    string `arg:"env:TWITTER_ACCESS_TOKEN"`
	AccessSecret   string `arg:"env:TWITTER_ACCESS_SECRET"`
}

func DefaultConfiguration() *Config {
	return &Config{
		Environment: "test",
		ImageFolder: "./images",
		HTTPServerConfig: HTTPServerConfig{
			BindAddress: ":9090",
		},
		TwitterConfig: TwitterConfig{
			ConsumerKey:    "",
			ConsumerSecret: "",
			AccessToken:    "",
			AccessSecret:   "",
		},
		DGraphConfig: DGraphConfig{
			DGraphDSN: "127.0.0.1:9080",
		},
		PostgresConfig: PostgresConfig{
			PostgresDSN: "postgres://xapp:xapp123@127.0.0.1:5432/xapp?sslmode=disable",
		},
	}
}
