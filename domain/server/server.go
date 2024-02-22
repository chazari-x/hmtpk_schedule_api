package server

import (
	"crypto/tls"
	"net/http"
	"os"
	"os/user"
	"path/filepath"
	"time"

	schedule "github.com/chazari-x/hmtpk_schedule"
	"github.com/chazari-x/hmtpk_schedule_api/config"
	"github.com/chazari-x/hmtpk_schedule_api/domain/server/handler"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/acme/autocert"
)

func StartServer(cfg config.Server, sch *schedule.Controller) error {
	if cfg.HTTPSAddress == "" {
		log.Tracef("server: %s%s", cfg.Domain, cfg.HTTPAddress)

		return http.ListenAndServe(cfg.HTTPAddress, handler.Router(cfg, sch))
	}

	log.Tracef("server: %s%s and %s%s", cfg.Domain, cfg.HTTPAddress, cfg.Domain, cfg.HTTPSAddress)

	certManager := autocert.Manager{
		Prompt:     autocert.AcceptTOS,
		HostPolicy: autocert.HostWhitelist(cfg.Domain),
		Cache:      autocert.DirCache("/tmp/cache-golang-autocert"),
	}

	if u, _ := user.Current(); u != nil {
		dir := filepath.Join(os.TempDir(), "cache-golang-autocert-"+u.Username)
		if os.MkdirAll(dir, 0700) == nil {
			certManager.Cache = autocert.DirCache(dir)
		}
	}

	server := &http.Server{
		Addr:    cfg.HTTPSAddress,
		Handler: handler.Router(cfg, sch),
		TLSConfig: &tls.Config{
			GetCertificate:   certManager.GetCertificate,
			CurvePreferences: []tls.CurveID{tls.X25519, tls.CurveP256},
		},
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	go func() {
		log.Fatal(http.ListenAndServe(cfg.HTTPAddress, certManager.HTTPHandler(nil)))
	}()

	return server.ListenAndServeTLS("", "")
}
