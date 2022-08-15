package config

import (
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"log"
)

var (
	// Configuration
	/**
	   I don't love this global approach, but since the client is an envoy filter I don't believe I have access to modify the Check
	   endpoint or use the grpc metadata.  I prefer to pass explicitly or via the context at the least.  Maybe if someone knows
	   better than me, they can provide a better solution. Even if we use flags, we would have to make those values global to access
	   them in other packages, since we still can't pass them
	*/

	Configuration *Config
)

type Config struct {
	ConcurrentStreams         uint32 `default:"10"`
	//ShowReflection            bool   `default:"false"`
	ShowReflection            bool   `default:"true"`
	//PrivateKeyProductCoreAuth string `required:"true"`
	PrivateKeyProductCoreAuth string `required:"false"`
	ExpirationTimeFromNow     int    `default:"100"` //in hours - will append these many hours after the current time
}

func ConfigSetup() {
	// we do not care if there is no .env file.
	_ = godotenv.Overload(".env")

	var cfg Config
	if err := envconfig.Process("", &cfg); err != nil {
		log.Fatal(err)
	}

	// there was some circular dependency, so we just use the default logger in this case
	// logrus.ConfigureLogger(&cfg)
	printConfigValues(&cfg)

	Configuration = &cfg
}

func printConfigValues(cfg *Config) {
	log.Printf("EXPIRATION_TIME_FROM_NOW %d\n", cfg.ExpirationTimeFromNow)
	log.Printf("CONCURRENT_STREAMS %v\n", cfg.ConcurrentStreams)
	log.Printf("SHOW_REFLECTION %t\n", cfg.ShowReflection)
}
