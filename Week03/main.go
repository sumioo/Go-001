package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"golang.org/x/sync/errgroup"
)

func hello(res http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(res, "Hello, My name is mowe")
}

func newServer() *http.Server {
	http.HandleFunc("/", hello)
	srv := http.Server{
		Addr:    "localhost:4000",
		Handler: http.DefaultServeMux,
	}
	return &srv
}

func main() {
	eg, ctx := errgroup.WithContext(context.Background())
	eg.Go(func() error {
		server := newServer()
		go func() {
			select {
			case <-ctx.Done():
				fmt.Println("http server shutdown form ctx")
				server.Shutdown(context.Background())
			}
		}()
		return server.ListenAndServe()
	})

	eg.Go(func() error {
		exitSignals := []os.Signal{os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGINT} // SIGTERM is POSIX specific
		sig := make(chan os.Signal, len(exitSignals))
		signal.Notify(sig, exitSignals...)
		select {
		case <-ctx.Done():
			fmt.Println("signal trap exit form ctx")
			return ctx.Err()
		case s := <-sig:
			fmt.Printf("receive signal %v and will quit\n", s)
			return nil
		}
	})

	eg.Go(func() error {
		fmt.Println("inject")
		time.Sleep(time.Second * 200)
		fmt.Println("inject finish")
		return errors.New("inject error")
	})

	if err := eg.Wait(); err != nil {
		fmt.Printf("first error %v", err)
	}
}
