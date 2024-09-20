package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/chazari-x/hmtpk_parser/v2"
	"github.com/chazari-x/hmtpk_schedule_api/config"
	"github.com/go-chi/chi/v5"
	log "github.com/sirupsen/logrus"
	"golang.org/x/net/context"
)

type Handler struct {
	cfg config.HTTP
	sch *hmtpk_parser.Controller
}

func Router(cfg config.HTTP, sch *hmtpk_parser.Controller) *chi.Mux {
	h := &Handler{
		cfg: cfg,
		sch: sch,
	}

	router := chi.NewRouter()
	router.Get("/*", router.NotFoundHandler())
	router.Get("/groups", h.groups)
	router.Get("/teachers", h.teachers)
	router.Get("/schedule", h.schedule)
	router.Get("/announces", h.announces)

	router.Get("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "domain/http/images/favicon.ico")
	})

	return router
}

type Error struct {
	Error string `json:"error"`
}

func write(w http.ResponseWriter, statusCode int, data interface{}) {
	if statusCode != http.StatusOK {
		w.WriteHeader(statusCode)
	}

	marshal, err := json.Marshal(data)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	_, _ = w.Write(marshal)
}

func (h *Handler) teachers(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), time.Second*30)
	defer cancel()

	options, err := h.sch.GetTeacherOptions(ctx)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled) {
			write(w, http.StatusInternalServerError, Error{"hmtpk not working"})
			return
		}

		log.Error(err)

		write(w, http.StatusInternalServerError, Error{http.StatusText(http.StatusInternalServerError)})
		return
	}

	write(w, http.StatusOK, options)
}

func (h *Handler) groups(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), time.Second*30)
	defer cancel()

	options, err := h.sch.GetGroupOptions(ctx)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled) {
			write(w, http.StatusInternalServerError, Error{"hmtpk not working"})
			return
		}

		log.Error(err)

		write(w, http.StatusInternalServerError, Error{http.StatusText(http.StatusInternalServerError)})
		return
	}

	write(w, http.StatusOK, options)
}

func (h *Handler) schedule(w http.ResponseWriter, r *http.Request) {
	log.Trace(r.URL.String())

	date := r.URL.Query().Get("date")
	if date != "" {
		if _, err := time.Parse("02.01.2006", date); err != nil {
			write(w, http.StatusBadRequest, Error{http.StatusText(http.StatusBadRequest)})
			return
		}
	} else {
		date = time.Now().Format("02.01.2006")
	}

	group := r.URL.Query().Get("group")
	if group != "" {
		ctx, cancel := context.WithTimeout(r.Context(), time.Second*30)
		defer cancel()

		scheduleByGroup, err := h.sch.GetScheduleByGroup(group, date, ctx)
		if err != nil {
			if errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled) {
				write(w, http.StatusInternalServerError, Error{"hmtpk not working"})
				return
			} else if strings.Contains(err.Error(), http.StatusText(http.StatusBadRequest)) {
				write(w, http.StatusBadRequest, Error{http.StatusText(http.StatusBadRequest)})
				return
			}

			log.Error(err)

			write(w, http.StatusInternalServerError, Error{http.StatusText(http.StatusInternalServerError)})
			return
		}

		write(w, http.StatusOK, scheduleByGroup)
		return
	}

	teacher := r.URL.Query().Get("teacher")
	if teacher != "" {
		ctx, cancel := context.WithTimeout(r.Context(), time.Second*30)
		defer cancel()

		scheduleByTeacher, err := h.sch.GetScheduleByTeacher(teacher, date, ctx)
		if err != nil {
			if errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled) {
				write(w, http.StatusInternalServerError, Error{"hmtpk not working"})
				return
			} else if strings.Contains(err.Error(), http.StatusText(http.StatusBadRequest)) {
				write(w, http.StatusBadRequest, Error{http.StatusText(http.StatusBadRequest)})
				return
			}

			log.Error(err)

			write(w, http.StatusInternalServerError, Error{http.StatusText(http.StatusInternalServerError)})
			return
		}

		write(w, http.StatusOK, scheduleByTeacher)
		return
	}

	write(w, http.StatusBadRequest, Error{http.StatusText(http.StatusBadRequest)})
}

func (h *Handler) announces(w http.ResponseWriter, r *http.Request) {
	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil {
		write(w, http.StatusBadRequest, Error{http.StatusText(http.StatusBadRequest)})
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), time.Second*30)
	defer cancel()

	announces, err := h.sch.GetAnnounces(ctx, page)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled) {
			write(w, http.StatusInternalServerError, Error{"hmtpk not working"})
			return
		}

		log.Error(err)

		write(w, http.StatusInternalServerError, Error{http.StatusText(http.StatusInternalServerError)})
		return
	}

	write(w, http.StatusOK, announces)
}
