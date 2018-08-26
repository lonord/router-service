package app

import (
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/lonord/sse"
)

const readInterval = time.Second * 3

type netSpeedService struct {
	sse    *sse.Service
	action *MainAction
	reader *WrappedNetSpeedReader
	lck    sync.Mutex
}

func newNetSpeedService(action *MainAction) *netSpeedService {
	return &netSpeedService{
		sse: sse.NewServiceWithOption(sse.Option{
			Headers: map[string]string{"X-Accel-Buffering": "no"},
		}),
		action: action,
	}
}

func (s *netSpeedService) handleClient(clientID interface{}, w http.ResponseWriter) error {
	c, err := s.sse.HandleClient(clientID, w)
	if err != nil {
		return err
	}
	go s.runNetSpeedReader()
	<-c
	return nil
}

func (s *netSpeedService) runNetSpeedReader() {
	s.lck.Lock()
	if s.reader == nil {
		rd, err := s.action.CreateNetSpeedReader()
		if err != nil {
			log.Println(err)
		} else {
			s.reader = rd
		}
		s.lck.Unlock()
		timer := time.NewTimer(readInterval)
		for {
			select {
			case <-timer.C:
				if s.sse.GetClientCount() == 0 {
					s.lck.Lock()
					s.reader = nil
					s.lck.Unlock()
					return
				}
				result, _ := s.reader.Read()
				s.sse.Broadcast(sse.Event{
					Data: result,
				})
				timer.Reset(readInterval)
			}
		}
	} else {
		s.lck.Unlock()
	}
}
