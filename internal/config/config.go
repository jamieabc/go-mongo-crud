package config

import "github.com/spf13/viper"

type ServerInfo struct {
	IP       string
	Port     int
	Database string
}

type Config struct {
	Server ServerInfo
}

// Parse - parse config
func Parse(configFile string) (Config, error) {
	viper.SetConfigType("yaml")
	viper.SetConfigFile(configFile)
	viper.AddConfigPath(".")

	err := viper.ReadInConfig()
	if nil != err {
		return Config{}, err
	}

	return Config{
		ServerInfo{
			IP:       viper.GetString("server.ip"),
			Port:     viper.GetInt("server.port"),
			Database: viper.GetString("server.database"),
		},
	}, nil
}
