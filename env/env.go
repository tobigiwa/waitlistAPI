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
			Username   string `env:"MONGO_USERNAME,required=true" `
			Password   string `env:"MONGO_PASSWORD,required=true" json:"-" yaml:"-" toml:"-"`
			Database   string `env:"MONGO_DATABASE,default=waitlist"`
			Collection string `env:"MONGO_COLLECTION,default=users"`
		}
	}

	Server struct {
		Port    string `env:"PORT,default=8090"`
		Env     string `env:"ENV,default=development"`
		Version string `env:"VERSION"`
	}

	Mail struct {
		EmailAcc            string `env:"EMAILACC,required=true"`
		EmailPswd           string `env:"EMAILPSWD,required=true"`
		EmailSmtpServerHost string `env:"SMTPHOST,required=true"`
		EmailSmtpServerPort int    `env:"SMTPPORT,required=true"`
	}

	EncryptionKey string `env:"ENCRYPTIONKEY,required=true"`

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

// GetEnvVar gets a particular env. variable and should only be
// called after LoadAllEnvVars() has be called.
func GetEnvVar() Environment {
	LoadAllEnvVars()
	return environment
}

func BuildURI(username, password, host string) string {
	return "mongodb+srv://" + username + ":" + password + "@" + host + "/?retryWrites=true&w=majority"
}
