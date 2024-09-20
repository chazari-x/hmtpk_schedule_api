package http

import (
	"crypto/tls"
	"net/http"
	"os"
	"os/user"
	"path/filepath"
	"time"

	"github.com/chazari-x/hmtpk_parser/v2"
	"github.com/chazari-x/hmtpk_schedule_api/config"
	"github.com/chazari-x/hmtpk_schedule_api/domain/http/handler"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/acme/autocert"
)

func Start(cfg config.HTTP, sch *hmtpk_parser.Controller) error {
	if cfg.HTTPSAddress == "" {
		log.Tracef("http server: %s%s", cfg.Domain, cfg.HTTPAddress)

		return http.ListenAndServe(cfg.HTTPAddress, handler.Router(cfg, sch))
	}

	certManager := autocert.Manager{
		Prompt: autocert.AcceptTOS,
		Cache:  autocert.DirCache("/tmp/cache-golang-autocert"),
	}

	if cfg.Domain != "localhost" {
		certManager.HostPolicy = autocert.HostWhitelist(cfg.Domain)
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
		ReadTimeout:  1 * time.Minute,
		WriteTimeout: 1 * time.Minute,
	}

	log.Tracef("http server: %s%s and %s%s", cfg.Domain, cfg.HTTPAddress, cfg.Domain, cfg.HTTPSAddress)

	go func() {
		log.Fatal(http.ListenAndServe(cfg.HTTPAddress, certManager.HTTPHandler(handler.Router(cfg, sch))))
	}()

	return server.ListenAndServeTLS("", "")
}
