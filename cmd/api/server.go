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
)

func (a *application) serve() error {

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", a.conf.port),
		Handler: a.routes(),
		// ErrorLog: log.New(logger, "", 0),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	shutDownError := make(chan error)

	go func() {
		quit := make(chan os.Signal, 1)

		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

		s := <-quit

		a.logger.PrintInfo("shutting down server", map[string]string{
			"signal": s.String(),
		})

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		err := srv.Shutdown(ctx)
		if err != nil {
			shutDownError <- err
		}

		// Log a message to say that we're waiting for any background goroutines to
		// complete their tasks.
		a.logger.PrintInfo("completing background tasks", map[string]string{
			"addr": srv.Addr,
		})
		// Call Wait() to block until our WaitGroup counter is zero --- essentially
		// blocking until the background goroutines have finished. Then we return nil on
		// the shutdownError channel, to indicate that the shutdown completed without
		// any issues.
		a.wg.Wait()
		shutDownError <- nil
	}()

	a.logger.PrintInfo("starting server", map[string]string{
		"addr": srv.Addr,
		"env":  a.conf.env,
	})

	err := srv.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	err = <-shutDownError

	if err != nil {
		return err
	}

	a.logger.PrintInfo("stopped server", map[string]string{
		"addr": srv.Addr,
	})
	return nil
}
