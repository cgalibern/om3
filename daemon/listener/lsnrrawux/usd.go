package lsnrrawux

import (
	"context"
	"errors"
	"net"
	"os"
	"time"

	"opensvc.com/opensvc/daemon/listener/routehttp"
	"opensvc.com/opensvc/daemon/listener/routeraw"
)

func (t *T) stop() error {
	if t.listener == nil {
		t.log.Info().Msg("listener already stopped")
		return nil
	}
	if err := (*t.listener).Close(); err != nil {
		t.log.Error().Err(err).Msg("close failed")
		return err
	}
	t.log.Info().Msg("listener stopped")
	return nil
}

func (t *T) start(ctx context.Context) error {
	if err := os.RemoveAll(t.addr); err != nil {
		t.log.Error().Err(err).Msg("RemoveAll")
		return err
	}
	mux := routeraw.New(routehttp.New(ctx, false), t.log, 5*time.Second)
	listener, err := net.Listen("unix", t.addr)
	if err != nil {
		t.log.Error().Err(err).Msg("listen failed")
		time.Sleep(time.Second)
		if listener, err = net.Listen("unix", t.addr); err != nil {
			return err
		}
	}
	c := make(chan bool)
	go func() {
		c <- true
		for {
			conn, err := listener.Accept()
			if err != nil {
				if errors.Is(err, net.ErrClosed) {
					break
				} else {
					t.log.Error().Err(err).Msg("Accept")
					continue
				}
			}
			go mux.Serve(conn)
		}
	}()
	t.listener = &listener
	<-c
	t.log.Info().Msg("listener started " + t.addr)
	return nil
}
