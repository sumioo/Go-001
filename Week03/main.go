package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"golang.org/x/sync/errgroup"
)

func hello(res http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(res, "Hello, My name is mowe")
}

func myserver(ctx context.Context) error {
	http.HandleFunc("/", hello)
	srv := http.Server{
		Addr:    "localhost:4000",
		Handler: http.DefaultServeMux,
	}
	go func() {
		select {
		case <-ctx.Done():
			fmt.Println("other service quit and myserver will quit too")
			srv.Shutdown(context.Background())
		}
	}()
	return srv.ListenAndServe()
}

func trapSignal(ctx context.Context) error {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)
	select {
	case <-ctx.Done():
		fmt.Println("other service quit and trapSignal will quit too")
		return nil
	case s := <-c:
		fmt.Printf("receive signal %v and will quit\n", s)
		return fmt.Errorf("receive terminate signal")
	}
}

func fakeService(ctx context.Context, duration time.Duration) error {
	select {
	case <-ctx.Done():
		fmt.Println("other service quit and fakeService will quit too")
	case <-time.After(duration):
		fmt.Println("timeout and quit")
		return fmt.Errorf("fakeService timeout")
	}
	return nil
}

func main() {
	eg, ctx := errgroup.WithContext(context.Background())
	eg.Go(func() error {
		return trapSignal(ctx)
	})

	eg.Go(func() error {
		return myserver(ctx)
	})

	eg.Go(func() error {
		return fakeService(ctx, time.Duration(5)*time.Second)
	})

	eg.Go(func() error {
		return fakeService(ctx, time.Duration(100)*time.Second)
	})

	if err := eg.Wait(); err != nil {
		fmt.Printf("first error %v", err)
	}
}
