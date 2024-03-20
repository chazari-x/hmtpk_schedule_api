package cmd

import (
	"github.com/chazari-x/hmtpk_schedule"
	"github.com/chazari-x/hmtpk_schedule_api/domain/http"
	"github.com/chazari-x/hmtpk_schedule_api/redis"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func init() {
	cmd := &cobra.Command{
		Use:   "http",
		Short: "http",
		Long:  "http",
		Run: func(cmd *cobra.Command, args []string) {
			cfg := getConfig(cmd)

			log.Trace("http starting..")
			defer log.Trace("http stopped")

			client, err := redis.Connect(&cfg.Redis)
			if err != nil {
				log.Fatal(err)
			}
			defer func() {
				_ = client.Close()
			}()

			schedule := hmtpk_schedule.NewController(client, log.StandardLogger())

			if err := http.Start(cfg.Server.HTTP, schedule); err != nil {
				log.Error(err)
			}
		},
	}
	cmd.PersistentFlags().String("config", "", "dev")
	rootCmd.AddCommand(cmd)
}
