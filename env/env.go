package env

import (
	"sync"

	"github.com/Netflix/go-env"
	"github.com/joho/godotenv"
)

type Environment struct {
	Databases struct {
		Mongo struct {
			Host       string `env:"MONGO_HOST,required=true"`
			Port       string `env:"MONGO_PORT,default=27017"`
			Username   string `env:"MONGO_USERNAME,required=true" `
			Password   string `env:"MONGO_PASSWORD,required=true" json:"-" yaml:"-" toml:"-"`
			Database   string `env:"MONGO_DATABASE,required=true"`
			Collection string `env:"MONGO_COLLECTION,required=true"`
		}

		Redis struct {
			Host     string `env:"REDIS_HOST,required=true"`
			Port     string `env:"REDIS_PORT,default=6379"`
			Password string `env:"REDIS_PASSWORD,required=true" json:"-" yaml:"-" toml:"-"`
			User     string `env:"REDIS_USER,default=default"`
			Db       int    `env:"DB,default=0"`
		}
	}

	PORT struct {
		HTTP string `env:"PORT,default=8080"`
	}

	Nonce     string `env:"NONCE,default=gopher"`
	EmailPswd string `env:"EMAILPSWD,required=true"`

	Extras env.EnvSet
}

var once sync.Once
var environment Environment

func LoadAllEnvVars() {
	once.Do(func() {
		_init()
	})
}

func _init() {

	if err := godotenv.Load(); err != nil {
		panic(err)
	}

	es, err := env.UnmarshalFromEnviron(&environment)
	if err != nil {
		panic(err)
	}
	// Remaining environment variables.
	environment.Extras = es
}

func GetEnvVar() Environment {
	LoadAllEnvVars()
	return environment
}

func BuildURI(username, password, host string) string {
	return "mongodb+srv://" + username + ":" + password + "@" + host + "/?retryWrites=true&w=majority"
}
