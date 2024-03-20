package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	schedule "github.com/chazari-x/hmtpk_schedule"
	"github.com/chazari-x/hmtpk_schedule/model"
	"github.com/chazari-x/hmtpk_schedule_api/config"
	"github.com/go-chi/chi/v5"
	log "github.com/sirupsen/logrus"
	"golang.org/x/net/context"
)

type Handler struct {
	cfg config.HTTP
	sch *schedule.Controller
}

func Router(cfg config.HTTP, sch *schedule.Controller) *chi.Mux {
	h := &Handler{
		cfg: cfg,
		sch: sch,
	}

	router := chi.NewRouter()
	router.Get("/*", router.NotFoundHandler())
	router.Get("/groups", h.groups)
	router.Get("/teachers", h.teachers)
	router.Get("/schedule", h.get)
	router.Get("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "domain/http/images/favicon.ico")
	})

	return router
}

type Error struct {
	Error string `json:"error"`
}

func (h *Handler) teachers(w http.ResponseWriter, _ *http.Request) {
	get, err := http.Get("https://api.vk.com/method/execute.getTeachers?v=5.154&access_token=" + h.cfg.MiniAppToken)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		if marshal, err := json.Marshal(Error{http.StatusText(http.StatusInternalServerError)}); err == nil {
			_, _ = w.Write(marshal)
		}
		return
	}

	if get.StatusCode != 200 {
		w.WriteHeader(http.StatusInternalServerError)
		if marshal, err := json.Marshal(Error{fmt.Sprintf("vk api response status code: %s", http.StatusText(get.StatusCode))}); err == nil {
			_, _ = w.Write(marshal)
		}

		return
	}

	all, err := io.ReadAll(get.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		if marshal, err := json.Marshal(Error{http.StatusText(http.StatusInternalServerError)}); err == nil {
			_, _ = w.Write(marshal)
		}
		return
	}

	if _, err = w.Write(all); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		if marshal, err := json.Marshal(Error{http.StatusText(http.StatusInternalServerError)}); err == nil {
			_, _ = w.Write(marshal)
		}
	}
}

func (h *Handler) groups(w http.ResponseWriter, _ *http.Request) {
	get, err := http.Get("https://api.vk.com/method/execute.getGroups?v=5.154&access_token=" + h.cfg.MiniAppToken)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		if marshal, err := json.Marshal(Error{http.StatusText(http.StatusInternalServerError)}); err == nil {
			_, _ = w.Write(marshal)
		}
		return
	}

	if get.StatusCode != 200 {
		w.WriteHeader(http.StatusInternalServerError)
		if marshal, err := json.Marshal(Error{fmt.Sprintf("vk api response status code: %s", http.StatusText(get.StatusCode))}); err == nil {
			_, _ = w.Write(marshal)
		}

		return
	}

	all, err := io.ReadAll(get.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		if marshal, err := json.Marshal(Error{http.StatusText(http.StatusInternalServerError)}); err == nil {
			_, _ = w.Write(marshal)
		}
		return
	}

	if _, err = w.Write(all); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		if marshal, err := json.Marshal(Error{http.StatusText(http.StatusInternalServerError)}); err == nil {
			_, _ = w.Write(marshal)
		}
	}
}

func (h *Handler) get(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	date := r.URL.Query().Get("date")
	if date != "" {
		if _, err := time.Parse("02.01.2006", date); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			if marshal, err := json.Marshal(Error{http.StatusText(http.StatusBadRequest)}); err == nil {
				_, _ = w.Write(marshal)
			}
			return
		}
	} else {
		date = time.Now().Format("02.01.2006")
	}

	group := r.URL.Query().Get("group")
	teacher := r.URL.Query().Get("teacher")

	if group == "" && teacher == "" {
		w.WriteHeader(http.StatusBadRequest)
		if marshal, err := json.Marshal(Error{http.StatusText(http.StatusBadRequest)}); err == nil {
			_, _ = w.Write(marshal)
		}

		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), time.Second*10)
	defer cancel()

	var sch []model.Schedule
	var err error

	if group != "" {
		sch, err = h.sch.GetScheduleByGroup(group, date, ctx)
	} else {
		sch, err = h.sch.GetScheduleByTeacher(teacher, date, ctx)
	}

	if err != nil {
		if errors.Is(err, context.Canceled) {
			w.WriteHeader(http.StatusRequestTimeout)
			if marshal, err := json.Marshal(Error{http.StatusText(http.StatusRequestTimeout)}); err == nil {
				_, _ = w.Write(marshal)
			}
			return
		}

		if strings.Contains(err.Error(), http.StatusText(http.StatusBadRequest)) {
			w.WriteHeader(http.StatusBadRequest)
			if marshal, err := json.Marshal(Error{http.StatusText(http.StatusBadRequest)}); err == nil {
				_, _ = w.Write(marshal)
			}
			return
		}

		log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		if marshal, err := json.Marshal(Error{http.StatusText(http.StatusInternalServerError)}); err == nil {
			_, _ = w.Write(marshal)
		}

		return
	}

	marshal, err := json.Marshal(sch)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		if marshal, err := json.Marshal(Error{http.StatusText(http.StatusInternalServerError)}); err == nil {
			_, _ = w.Write(marshal)
		}

		return
	}

	if _, err = w.Write(marshal); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		if marshal, err := json.Marshal(Error{http.StatusText(http.StatusInternalServerError)}); err == nil {
			_, _ = w.Write(marshal)
		}
	}
}
