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

// week03 homework
func main() {
	ping := func(resp http.ResponseWriter, req *http.Request) {
		_, _ = fmt.Fprint(resp, "ping")
	}
	eg, ctx := errgroup.WithContext(context.Background())
	// http server
	eg.Go(func() error {
		var e error
		c, cancel := context.WithCancel(context.Background())
		go func() {
			e = http.ListenAndServe("localhost:8888", http.HandlerFunc(ping))
			if e != nil {
				cancel()
				log.Println("cancel")
			}
		}()
		select {
		case <-ctx.Done():
			log.Println("exit sig")
		case <-c.Done():
			log.Println("server error")
			return errors.New("server error")
		}
		log.Println("error =", e)
		return e
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
