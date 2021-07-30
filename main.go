package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
)

func main() {
	ping := func(resp http.ResponseWriter, req *http.Request) {
		_, _ = fmt.Fprint(resp, "ping")
	}
	eg, ctx := errgroup.WithContext(context.Background())
	// http server
	s := http.Server{
		Addr:    "localhost:8888",
		Handler: http.HandlerFunc(ping),
	}
	eg.Go(func() error {
		return s.ListenAndServe()
	})
	eg.Go(func() error {
		select {
		case <-ctx.Done():
			return s.Shutdown(context.Background())
		}
	})
	// sig listener
	eg.Go(func() error {
		return listenSig(ctx)
	})
	log.Println("server start")
	if e := eg.Wait(); e != nil {
		log.Println("do close work, err =", e)
	}

}

func listenSig(ctx context.Context) error {
	defer func() {
		_ = recover()
	}()
	sig := make(chan os.Signal)
	signal.Notify(sig, syscall.SIGKILL, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGKILL, syscall.SIGABRT, syscall.SIGINT)
	select {
	case <-sig:
		log.Println("exit sig receive")
		return errors.New("exit sig receive")
	case <-ctx.Done():
		close(sig)
		log.Println("close listen sig")
		return nil
	}
}
