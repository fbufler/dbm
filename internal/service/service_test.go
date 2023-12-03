package service

import (
	"context"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/fbufler/database-monitor/internal/tester"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	type args struct {
		config  Config
		results chan tester.Result
		router  *mux.Router
	}
	tests := []struct {
		name string
		args args
		want Service
	}{
		{
			name: "New",
			args: args{
				config: Config{
					Port:             8080,
					InvalidationTime: 60,
				},
				results: make(chan tester.Result),
				router:  mux.NewRouter(),
			},
			want: &ServiceImpl{
				config:     Config{Port: 8080, InvalidationTime: 60},
				results:    make(chan tester.Result),
				router:     mux.NewRouter(),
				resultsMap: make(map[string]tester.Result),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := New(tt.args.config, tt.args.results, tt.args.router)
			if got == nil {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestServiceImplCollectResults(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	results := make(chan tester.Result)
	s := New(Config{Port: 8080, InvalidationTime: 1}, results, mux.NewRouter())
	go s.(*ServiceImpl).collectResults(ctx)
	results <- tester.Result{
		Database:       "test",
		Connectable:    true,
		Writable:       true,
		Readable:       true,
		ConnectionTime: 0,
		Timestamp:      time.Now(),
		WriteTime:      0,
		ReadTime:       0,
	}
	result, ok := s.(*ServiceImpl).resultsMap["test"]
	assert.Equal(t, true, ok)
	assert.Equal(t, "test", result.Database)
	assert.Equal(t, true, result.Connectable)
	assert.Equal(t, true, result.Writable)
	assert.Equal(t, true, result.Readable)
	assert.Equal(t, time.Duration(0), result.ConnectionTime)
	assert.Equal(t, time.Duration(0), result.WriteTime)
	assert.Equal(t, time.Duration(0), result.ReadTime)
}

func TestServiceImplCollectResultsInvalidation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	results := make(chan tester.Result)
	s := New(Config{Port: 8080, InvalidationTime: 1}, results, mux.NewRouter())
	go s.(*ServiceImpl).collectResults(ctx)
	results <- tester.Result{
		Database:       "test",
		Connectable:    true,
		Writable:       true,
		Readable:       true,
		ConnectionTime: 0,
		Timestamp:      time.Now().Add(-2 * time.Second),
		WriteTime:      0,
		ReadTime:       0,
	}
	_, ok := s.(*ServiceImpl).resultsMap["test"]
	assert.Equal(t, false, ok)
}

func TestRunResults(t *testing.T) {
	router := mux.NewRouter()
	results := make(chan tester.Result)
	s := New(Config{Port: 8080, InvalidationTime: 1}, results, router)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go s.Run(ctx)
	now := time.Now()
	results <- tester.Result{
		Database:       "test",
		Connectable:    true,
		Writable:       true,
		Readable:       true,
		ConnectionTime: 0,
		Timestamp:      now,
		WriteTime:      0,
		ReadTime:       0,
	}
	time.Sleep(50 * time.Millisecond)
	server := httptest.NewServer(router)
	request := httptest.NewRequest("GET", server.URL+"/results", nil)
	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)
	assert.Equal(t, 200, response.Code)
	assert.Equal(t, "{\"results\":{\"test\":{\"database\":\"test\",\"connectable\":true,\"connection_time\":0,\"writable\":true,\"write_time\":0,\"readable\":true,\"read_time\":0,\"timestamp\":\""+now.Format(time.RFC3339Nano)+"\"}}}\n", response.Body.String())
}
