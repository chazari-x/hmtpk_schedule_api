package grpc

import (
	"net"

	schedule "github.com/chazari-x/hmtpk_schedule"
	"github.com/chazari-x/hmtpk_schedule_api/config"
	"github.com/chazari-x/hmtpk_schedule_api/domain/grpc/protobuf"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

func Start(cfg config.GRPC, sch *schedule.Controller) (err error) {
	lis, err := net.Listen("tcp", cfg.Address)
	if err != nil {
		return
	}

	s := grpc.NewServer()

	protobuf.RegisterScheduleServer(s, protobuf.NewServer(cfg, sch))

	log.Tracef("grpc server: %s%s", cfg.Domain, cfg.Address)

	return s.Serve(lis)
}
