package app

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"../base"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/lonord/sse"
)

type WebService struct {
	e   *echo.Echo
	cfg *ba.Config
	act *MainAction
}

func NewWebService(act *MainAction, cfg *ba.Config) *WebService {
	ec := createEcho()
	bindRouters(ec, act)
	return &WebService{
		e:   ec,
		cfg: cfg,
		act: act,
	}
}

func (s *WebService) Start() {
	addr := fmt.Sprintf("%s:%d", s.cfg.RPCHost, s.cfg.RPCPort)
	go func() {
		log.Println("server listens at http://", addr)
		if err := s.e.Start(addr); err != nil {
			log.Println("shutting down the server")
		}
	}()
}

func (s *WebService) Stop() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := s.e.Shutdown(ctx); err != nil {
		log.Fatal(err)
	}
}

func createEcho() *echo.Echo {
	ec := echo.New()
	ec.Use(middleware.Logger())
	ec.Use(middleware.Recover())
	ec.Use(middleware.CORS())
	return ec
}

func bindRouters(ec *echo.Echo, action *MainAction) {
	nss := newNetSpeedService(action)
	// get dnsmasq leases
	ec.GET("/clients", func(c echo.Context) error {
		r, err := action.GetOnlineClients()
		if err != nil {
			return err
		}
		return c.JSON(http.StatusOK, r)
	})
	// restart dnsmasq
	ec.PUT("/action/dnsmasq/restart", func(c echo.Context) error {
		err := action.RestartDnsmasq()
		if err != nil {
			return err
		}
		return c.String(http.StatusOK, "OK")
	})
	// sse net speed
	ec.Any("/netspeed", func(c echo.Context) error {
		nss.handleClient(sse.GenerateClientID(), c.Response().Writer)
		return nil
	})
}
