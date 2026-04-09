package config

type LogLevel string

const (
	LogDebug LogLevel = "debug"
	LogInfo  LogLevel = "info"
)

type Config struct {
	Workers  int
	Limit    int
	LogLevel LogLevel
}

func Default() Config {
	return Config{
		Workers:  4,
		Limit:    200,
		LogLevel: LogInfo,
	}
}
