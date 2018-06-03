package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/labstack/echo"
)

type WebService struct {
	e *echo.Echo
}

func NewWebService() *WebService {
	ec := echo.New()
	// TODO init
	return &WebService{
		e: ec,
	}
}

func (s *WebService) start(port int, hostname string) {
	addr := fmt.Sprintf(hostname, ":", port)
	go func() {
		if err := s.e.Start(addr); err != nil {
			log.Println("shutting down the server")
		}
	}()
}

func (s *WebService) stop() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := s.e.Shutdown(ctx); err != nil {
		log.Fatal(err)
	}
}
