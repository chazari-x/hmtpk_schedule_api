package protobuf

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"
	"time"

	schedule "github.com/chazari-x/hmtpk_schedule"
	"github.com/chazari-x/hmtpk_schedule/model"
	"github.com/chazari-x/hmtpk_schedule_api/config"
	log "github.com/sirupsen/logrus"
	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Server struct {
	cfg config.GRPC
	sch *schedule.Controller
}

func NewServer(cfg config.GRPC, sch *schedule.Controller) Server {
	return Server{cfg: cfg, sch: sch}
}

func (s Server) GetGroups(ctx context.Context, _ *Request) (*Response, error) {
	url := "https://api.vk.com/method/execute.getGroups?v=5.154&access_token=" + s.cfg.MiniAppToken
	request, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		log.Error(err)
		return nil, status.Errorf(codes.Internal, codes.Internal.String())
	}

	client := http.Client{}
	get, err := client.Do(request)
	if err != nil {
		log.Error(err)
		return nil, status.Errorf(codes.Internal, codes.Internal.String())
	}

	if get.StatusCode != 200 {
		return nil, status.Errorf(codes.Internal, codes.Internal.String())
	}

	all, err := io.ReadAll(get.Body)
	if err != nil {
		log.Error(err)
		return nil, status.Errorf(codes.Internal, codes.Internal.String())
	}

	return &Response{Message: string(all)}, nil
}

func (s Server) GetTeachers(ctx context.Context, _ *Request) (*Response, error) {
	url := "https://api.vk.com/method/execute.getTeachers?v=5.154&access_token=" + s.cfg.MiniAppToken
	request, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		log.Error(err)
		return nil, status.Errorf(codes.Internal, codes.Internal.String())
	}

	client := http.Client{}
	get, err := client.Do(request)
	if err != nil {
		log.Error(err)
		return nil, status.Errorf(codes.Internal, codes.Internal.String())
	}

	if get.StatusCode != 200 {
		return nil, status.Errorf(codes.Internal, codes.Internal.String())
	}

	all, err := io.ReadAll(get.Body)
	if err != nil {
		log.Error(err)
		return nil, status.Errorf(codes.Internal, codes.Internal.String())
	}

	return &Response{Message: string(all)}, nil
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
