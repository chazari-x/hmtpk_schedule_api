package cmd

import (
	"fmt"
	"os"
	"runtime"

	"github.com/chazari-x/hmtpk_schedule_api/config"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var configFile = "etc/"

type Config struct {
	Redis  config.Redis  `yaml:"redis"`
	Server config.Server `yaml:"server"`
}

func getConfig(cmd *cobra.Command) *Config {
	log.SetFormatter(&log.TextFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
		FullTimestamp:   true,
		CallerPrettyfier: func(frame *runtime.Frame) (function string, file string) {
			return "", fmt.Sprintf(" %s:%d", frame.File, frame.Line)
		},
	})

	log.SetReportCaller(true)

	log.SetLevel(log.TraceLevel)

	var cfg Config

	file, err := cmd.Flags().GetString("config")
	if err != nil {
		log.Fatalf("get flag err: %s", err)
	} else if file != "" {
		file += "."
	}

	configFile += fmt.Sprintf("config.%syaml", file)

	f, err := os.Open(configFile)
	if err != nil {
		log.Fatalf("open config file \"%s\": %s", configFile, err)
	}

	if err = yaml.NewDecoder(f).Decode(&cfg); err != nil {
		log.Fatalf("decode config file: %s", err)
	}

	return &cfg
}
