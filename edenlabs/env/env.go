package env

import (
	"time"

	"github.com/spf13/viper"
)

// Provider the config provider
type Provider interface {
	ConfigFileUsed() string
	Get(key string) interface{}
	GetBool(key string) bool
	GetDuration(key string) time.Duration
	GetFloat64(key string) float64
	GetInt(key string) int
	GetInt64(key string) int64
	GetSizeInBytes(key string) uint
	GetString(key string) string
	GetStringMap(key string) map[string]interface{}
	GetStringMapString(key string) map[string]string
	GetStringMapStringSlice(key string) map[string][]string
	GetStringSlice(key string) []string
	GetTime(key string) time.Time
	InConfig(key string) bool
	IsSet(key string) bool
}

func readConfig(fileName string) (v *viper.Viper, err error) {
	v = viper.New()
	v.AddConfigPath(".")
	v.AddConfigPath("./config")
	v.AddConfigPath("./env")
	v.SetConfigName(fileName)
	viper.SetConfigType("toml")

	err = v.ReadInConfig()
	if err != nil {
		return nil, err
	}

	return
}

// Env return provider so that you can read config anywhere
func Env(fileName string) (provider Provider, err error) {
	return readConfig(fileName)
}
