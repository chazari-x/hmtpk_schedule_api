package cmd

import (
	"github.com/chazari-x/hmtpk_schedule"
	"github.com/chazari-x/hmtpk_schedule_api/domain/server"
	"github.com/chazari-x/hmtpk_schedule_api/redis"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func init() {
	cmd := &cobra.Command{
		Use:   "Server",
		Short: "Server",
		Long:  "Server",
		Run: func(cmd *cobra.Command, args []string) {
			cfg := getConfig(cmd)

			log.Trace("server starting..")
			defer log.Trace("server stopped")

			client, err := redis.Connect(&cfg.Redis)
			if err != nil {
				log.Fatal(err)
			}
			defer func() {
				_ = client.Close()
			}()

			if err := server.StartServer(cfg.Server, hmtpk_schedule.NewController(client, log.StandardLogger())); err != nil {
				log.Fatalln(err)
			}
		},
	}
	cmd.PersistentFlags().String("config", "", "dev")
	rootCmd.AddCommand(cmd)
}
