package main

import (
	"fmt"
	"net"
	"net/http"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

func (a *application) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func (w http.ResponseWriter, r *http.Request)  {
		defer func() {
			if err := recover(); err != nil {
				w.Header().Set("Connection", "close")
				a.serverErrorResponse(w, r, fmt.Errorf("%s", err))
			}
		}()
		next.ServeHTTP(w, r)
	})
}

func (a *application) rateLimit(next http.Handler) http.Handler {
	type client struct {
		limiter *rate.Limiter
		lastSeen time.Time
	}
	var (
		mu sync.Mutex
		clients = make(map[string]*client)
	)

	go func() {
		for {
			time.Sleep(time.Minute)

			mu.Lock()

			for ip, client := range clients {
				if time.Since(client.lastSeen) > 3*time.Second {
					delete(clients, ip)
				}
			}
			mu.Unlock()
		}
	}()


	return http.HandlerFunc(func (w http.ResponseWriter, r *http.Request) {
		if a.conf.limiter.enabled {
			ip, _, err := net.SplitHostPort(r.RemoteAddr)
			if err != nil {
				a.serverErrorResponse(w, r, err)
			}
	
			mu.Lock()
	
			if _, found := clients[ip]; !found {
				clients[ip] = &client{limiter: rate.NewLimiter(rate.Limit(a.conf.limiter.rps), a.conf.limiter.burst)}
			}
	
			clients[ip].lastSeen = time.Now()
	
			if !clients[ip].limiter.Allow() {
				mu.Unlock()
				a.rateLimitExceeded(w, r)
				return
			}
	
			mu.Unlock()
		}
		next.ServeHTTP(w, r)
	})
}