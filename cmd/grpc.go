package cmd

import (
	"github.com/chazari-x/hmtpk_schedule"
	"github.com/chazari-x/hmtpk_schedule_api/domain/grpc"
	"github.com/chazari-x/hmtpk_schedule_api/redis"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func init() {
	cmd := &cobra.Command{
		Use:   "grpc",
		Short: "grpc",
		Long:  "grpc",
		Run: func(cmd *cobra.Command, args []string) {
			cfg := getConfig(cmd)

			log.Trace("grpc starting..")
			defer log.Trace("grpc stopped")

			client, err := redis.Connect(&cfg.Redis)
			if err != nil {
				log.Fatal(err)
			}
			defer func() {
				_ = client.Close()
			}()

			schedule := hmtpk_schedule.NewController(client, log.StandardLogger())

			if err := grpc.Start(cfg.Server.GRPC, schedule); err != nil {
				log.Error(err)
			}
		},
	}
	cmd.PersistentFlags().String("config", "", "dev")
	rootCmd.AddCommand(cmd)
}
