package app

import (
	"context"
	"fmt"
	"log"
	"time"

	"../base"
	"github.com/labstack/echo"
)

type WebService struct {
	e   *echo.Echo
	cfg *ba.Config
	act *MainAction
}

func NewWebService(act *MainAction, cfg *ba.Config) *WebService {
	ec := createEcho()
	return &WebService{
		e:   ec,
		cfg: cfg,
		act: act,
	}
}

func (s *WebService) Start(port int, hostname string) {
	addr := fmt.Sprintf(hostname, ":", port)
	go func() {
		log.Println("server listens at http://", hostname, ":", port)
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
	// TODO init
	return ec
}
