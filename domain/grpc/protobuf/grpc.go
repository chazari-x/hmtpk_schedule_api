package protobuf

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/chazari-x/hmtpk_parser/v2"
	"github.com/chazari-x/hmtpk_parser/v2/model"
	"github.com/chazari-x/hmtpk_schedule_api/config"
	log "github.com/sirupsen/logrus"
	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Server struct {
	cfg config.GRPC
	sch *hmtpk_parser.Controller
}

func NewServer(cfg config.GRPC, sch *hmtpk_parser.Controller) Server {
	return Server{cfg: cfg, sch: sch}
}

func (s Server) GetGroups(ctx context.Context, _ *Request) (*Response, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second*30)
	defer cancel()

	options, err := s.sch.GetGroupOptions(ctx)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled) {
			return nil, status.Errorf(codes.Internal, "hmtpk not working")
		}

		log.Error(err)

		return nil, status.Errorf(codes.Internal, codes.Internal.String())
	}

	marshal, err := json.Marshal(options)
	if err != nil {
		return nil, status.Errorf(codes.Internal, codes.Internal.String())
	}

	return &Response{Message: string(marshal)}, nil
}

func (s Server) GetTeachers(ctx context.Context, _ *Request) (*Response, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second*30)
	defer cancel()

	options, err := s.sch.GetTeacherOptions(ctx)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled) {
			return nil, status.Errorf(codes.Internal, "hmtpk not working")
		}

		log.Error(err)

		return nil, status.Errorf(codes.Internal, codes.Internal.String())
	}

	marshal, err := json.Marshal(options)
	if err != nil {
		return nil, status.Errorf(codes.Internal, codes.Internal.String())
	}

	return &Response{Message: string(marshal)}, nil
}

func (s Server) GetSchedule(ctx context.Context, r *ScheduleRequest) (*ScheduleResponse, error) {
	if r.Date != "" {
		_, err := time.Parse("02.01.2006", r.Date)
		if err != nil {
			return nil, status.Errorf(codes.InvalidArgument, codes.InvalidArgument.String())
		}
	} else {
		r.Date = time.Now().Format("02.01.2006")
	}

	if r.Group == "" && r.Teacher == "" {
		return nil, status.Errorf(codes.InvalidArgument, codes.InvalidArgument.String())
	}

	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	var sch []model.Schedule
	var err error

	if r.Group != "" {
		sch, err = s.sch.GetScheduleByGroup(r.Group, r.Date, ctx)
	} else {
		sch, err = s.sch.GetScheduleByTeacher(r.Teacher, r.Date, ctx)
	}

	if err != nil {
		if errors.Is(err, context.Canceled) {
			return nil, status.Errorf(codes.DeadlineExceeded, codes.DeadlineExceeded.String())
		}

		if strings.Contains(err.Error(), http.StatusText(http.StatusBadRequest)) {
			return nil, status.Errorf(codes.InvalidArgument, codes.InvalidArgument.String())
		}

		log.Error(err)
		return nil, status.Errorf(codes.Internal, codes.Internal.String())
	}

	marshal, err := json.Marshal(sch)
	if err != nil {
		log.Error(err)
		return nil, status.Errorf(codes.Internal, codes.Internal.String())
	}

	return &ScheduleResponse{Message: string(marshal)}, nil
}

func (Server) mustEmbedUnimplementedScheduleServer() {}
