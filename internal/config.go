package core

import (
	"github.com/fsnotify/fsnotify"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"time"
)

const (
	EnvLocal string = "local"
	EnvDev   string = "dev"
	EnvProd  string = "prod"

	Version string = "1.0.1"

	LogLevelPanic   zerolog.Level = 5
	LogLevelFatal   zerolog.Level = 4
	LogLevelError   zerolog.Level = 3
	LogLevelWarning zerolog.Level = 2
	LogLevelInfo    zerolog.Level = 1
	LogLevelDebug   zerolog.Level = 0
	LogLevelTrace   zerolog.Level = -1

	ApplicationLogTag string = "application"
	DatabaseLogTag    string = "database"
	WebsocketLogTag   string = "websocket"
)

var v *viper.Viper

type app struct {
	Env         string        `mapstructure:"APP_ENV"`
	LogLevel    zerolog.Level `mapstructure:"APP_LOG_LEVEL"`
	Timezone    *time.Location
	BookmakerID uuid.UUID
}

type db struct {
	HostMaster string `mapstructure:"DB_HOST_MASTER"`
	HostSlave  string `mapstructure:"DB_HOST_SLAVE"`

	PortMaster int `mapstructure:"DB_PORT_MASTER"`
	PortSlave  int `mapstructure:"DB_PORT_SLAVE"`

	NameMaster string `mapstructure:"DB_NAME_MASTER"`
	NameSlave  string `mapstructure:"DB_NAME_SLAVE"`

	User     string `mapstructure:"DB_USER"`
	Password string `mapstructure:"DB_PASS"`
}

type gorm struct {
	TablePrefix                              string        `mapstructure:"GORM_TABLE_PREFIX"`
	SingularTableNames                       bool          `mapstructure:"GORM_SINGULAR_TABLE_NAMES"`
	MaxConnIdleTime                          time.Duration `mapstructure:"GORM_MAX_CONN_IDLE_TIME"`
	MaxConnLifeTime                          time.Duration `mapstructure:"GORM_MAX_CONN_LIFE_TIME"`
	MaxIdleConnections                       int           `mapstructure:"GORM_MAX_IDLE_CONNS"`
	MaxOpenConnections                       int           `mapstructure:"GORM_MAX_OPEN_CONNS"`
	SkipDefaultTransaction                   bool          `mapstructure:"GORM_SKIP_DEFAULT_TRANSACTION"`
	FullSaveAssociations                     bool          `mapstructure:"GORM_FULL_SAVE_ASSOCIATIONS"`
	DryRun                                   bool          `mapstructure:"GORM_DRY_RUN"`
	PrepareStatement                         bool          `mapstructure:"GORM_PREPARE_STATEMENT"`
	DisableAutomaticPing                     bool          `mapstructure:"GORM_DISABLE_AUTOMATIC_PING"`
	DisableForeignKeyConstraintWhenMigrating bool          `mapstructure:"GORM_DISABLE_FOREIGN_KEY_CONSTRAINT_WHEN_MIGRATING"`
	DisableNestedTransaction                 bool          `mapstructure:"GORM_DISABLE_NESTED_TRANSACTION"`
	AllowGlobalUpdate                        bool          `mapstructure:"GORM_ALLOW_GLOBAL_UPDATE"`
	QueryFields                              bool          `mapstructure:"GORM_QUERY_FIELDS"`
	CreateBatchSize                          int           `mapstructure:"GORM_CREATE_BATCH_SIZE"`
}

type websocket struct {
	Protocol  string `mapstructure:"WEBSOCKET_PROTOCOL"`
	Host      string `mapstructure:"WEBSOCKET_HOST"`
	Path      string `mapstructure:"WEBSOCKET_PATH"`
	Debug     bool   `mapstructure:"WEBSOCKET_DEBUG"`
	AuthKey   string `mapstructure:"WEBSOCKET_AUTH_KEY"`
	Bookmaker string `mapstructure:"WEBSOCKET_BOOKMAKER"`
}

type config struct {
	// App configuration
	App app `mapstructure:",squash"`

	// Db configuration
	Db db `mapstructure:",squash"`

	// Gorm configuration
	Gorm gorm `mapstructure:",squash"`

	// Websocket configuration
	Websocket websocket `mapstructure:",squash"`
}

var Config config

func init() {
	v = viper.New()

	setDefaults()

	// Automatically refresh environment variables
	v.AutomaticEnv()

	if err := v.UnmarshalExact(&Config); nil != err {
		panic(err)
	}

	if EnvLocal == Config.App.Env {
		v.SetConfigName(".env")
		v.SetConfigType("dotenv")
		v.AddConfigPath(".")

		// Read configuration
		if err := v.ReadInConfig(); nil != err {
			panic(err)
		}

		if err := v.UnmarshalExact(&Config); nil != err {
			panic(err)
		}

		v.WatchConfig()
		v.OnConfigChange(func(in fsnotify.Event) {
			if err := v.UnmarshalExact(&Config); nil != err {
				panic(err)
			}
		})
	}

	setTimezone()
	configureLogger()

	log.Debug().Msg("config => loaded")
}

func setTimezone() {
	utc, _ := time.LoadLocation("UTC")
	Config.App.Timezone = utc
}

func configureLogger() {
	zerolog.SetGlobalLevel(Config.App.LogLevel)
}

func setDefaults() {
	// Set default App configuration
	v.SetDefault("APP_ENV", EnvLocal)
	v.SetDefault("APP_LOG_LEVEL", LogLevelTrace)

	// Set default database configuration
	v.SetDefault("DB_HOST_MASTER", "localhost")
	v.SetDefault("DB_HOST_SLAVE", "localhost")
	v.SetDefault("DB_PORT_MASTER", 5432)
	v.SetDefault("DB_PORT_SLAVE", 5432)
	v.SetDefault("DB_NAME_MASTER", "console")
	v.SetDefault("DB_NAME_SLAVE", "console")

	v.SetDefault("DB_USER", "console")
	v.SetDefault("DB_PASS", "console")

	// Set default gorm configuration
	v.SetDefault("GORM_TABLE_PREFIX", "")
	v.SetDefault("GORM_SINGULAR_TABLE_NAMES", false)
	v.SetDefault("GORM_MAX_CONN_IDLE_TIME", "10s")
	v.SetDefault("GORM_MAX_CONN_LIFE_TIME", "30s")
	v.SetDefault("GORM_MAX_IDLE_CONNS", 2)
	v.SetDefault("GORM_MAX_OPEN_CONNS", 10)
	v.SetDefault("GORM_SKIP_DEFAULT_TRANSACTION", false)
	v.SetDefault("GORM_FULL_SAVE_ASSOCIATIONS", false)
	v.SetDefault("GORM_DRY_RUN", false)
	v.SetDefault("GORM_PREPARE_STATEMENT", false)
	v.SetDefault("GORM_DISABLE_AUTOMATIC_PING", false)
	v.SetDefault("GORM_DISABLE_FOREIGN_KEY_CONSTRAINT_WHEN_MIGRATING", false)
	v.SetDefault("GORM_DISABLE_NESTED_TRANSACTION", false)
	v.SetDefault("GORM_ALLOW_GLOBAL_UPDATE", false)
	v.SetDefault("GORM_QUERY_FIELDS", false)
	v.SetDefault("GORM_CREATE_BATCH_SIZE", 0)

	v.SetDefault("WEBSOCKET_PROTOCOL", "ws")
	v.SetDefault("WEBSOCKET_HOST", "api.oddscp.com:8001")
	v.SetDefault("WEBSOCKET_PATH", "/")
	v.SetDefault("WEBSOCKET_DEBUG", false)
	v.SetDefault("WEBSOCKET_AUTH_KEY", "")
	v.SetDefault("WEBSOCKET_BOOKMAKER", "")
}
