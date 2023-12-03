package service

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/fbufler/database-monitor/internal/tester"
	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
)

type Service interface {
	Run(ctx context.Context)
}

type Config struct {
	Port             int
	InvalidationTime int
}

type Response struct {
	Results map[string]tester.Result `json:"results"`
}

type ServiceImpl struct {
	config     Config
	results    chan tester.Result
	router     *mux.Router
	resultsMap map[string]tester.Result
}

func New(config Config, results chan tester.Result, router *mux.Router) Service {
	return &ServiceImpl{
		resultsMap: make(map[string]tester.Result),
		config:     config,
		results:    results,
		router:     router,
	}
}

func (s *ServiceImpl) collectResults(ctx context.Context) {
	log.Info().Msg("Starting result collector")
	for {
		select {
		case res := <-s.results:
			s.resultsMap[res.Database] = res
		case <-ctx.Done():
			return
		}
		for _, res := range s.resultsMap {
			if time.Since(res.Timestamp) > time.Duration(s.config.InvalidationTime)*time.Second {
				delete(s.resultsMap, res.Database)
			}
		}
	}
}

func (s *ServiceImpl) getResultsHandler(w http.ResponseWriter, r *http.Request) {
	log.Debug().Msgf("Result requested from %s", r.RemoteAddr)
	json.NewEncoder(w).Encode(Response{
		Results: s.resultsMap,
	})
}

func (s *ServiceImpl) Run(ctx context.Context) {
	srv := &http.Server{
		Addr: fmt.Sprintf(":%d", s.config.Port),
	}
	go s.collectResults(ctx)
	s.router.HandleFunc("/results", s.getResultsHandler).Methods("GET")
	srv.Handler = s.router
	log.Info().Msgf("Starting service on port %d", s.config.Port)
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			if err != http.ErrServerClosed {
				log.Fatal().Msgf("service: %s", err)
			}
		}
	}()
	<-ctx.Done()
	log.Info().Msg("Shutting down service")
	srv.Shutdown(ctx)
}
