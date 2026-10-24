package conduit

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/maxence-charriere/go-app/v11/pkg/app"
	"github.com/stretchr/testify/require"
)

func TestMain(m *testing.M) {
	app.DefaultLogger = func(string, ...any) {}
	os.Exit(m.Run())
}

func TestLocalJobsExecutedOnMultipeWorkers(t *testing.T) {
	var mu sync.Mutex
	var wg sync.WaitGroup

	jobCounts := map[string]int{}

	workerCounts := map[string]int{
		"1": 0,
		"2": 0,
	}

	s := httptest.NewServer(nil)
	defer s.Close()

	tasks := LocalTaskScheduler{
		WorkerCount:      len(workerCounts),
		Endpoint:         s.URL,
		Transport:        http.DefaultTransport,
		WorkerHTTPHeader: "Test-Worker",
		BypassCacheURL:   []string{"/task/3"},
	}
	ctx, stop := startLocalTaskScheduler(t, &tasks)
	defer stop()

	s.Config.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Log(r.URL)

		mu.Lock()
		defer mu.Unlock()

		jobCounts[r.URL.Path]++
		workerCounts[r.Header.Get(tasks.WorkerHTTPHeader)]++
		wg.Done()
	})

	wg.Add(4)
	tasks.Post(ctx, Task{Path: "/task/1", Scope: DailyTask()})
	tasks.Post(ctx, Task{Path: "/task/1", Scope: DailyTask()})

	task2Time := time.Now().Add(time.Millisecond * 10)
	tasks.Post(ctx, Task{
		Path:  "/task/2",
		Scope: DailyTask(),
		Time:  task2Time,
	})
	tasks.Post(ctx, Task{Path: "/task/2", Scope: DailyTask()})
	tasks.Post(ctx, Task{Path: "/task/3", Scope: DailyTask()})
	tasks.Post(ctx, Task{Path: "/task/3", Scope: DailyTask()})

	wg.Wait()

	require.Equal(t, 1, jobCounts["/task/1"])
	require.Equal(t, 1, jobCounts["/task/2"])
	require.Equal(t, 2, jobCounts["/task/3"])

	for k, v := range workerCounts {
		require.Equal(t, 2, v, k)
	}
}

func TestTaskScopes(t *testing.T) {
	t.Log(MinuteTask())
	t.Log(SemiDailyTask())
	t.Log(ThriceDailyTask())
	t.Log(DailyTask())
	t.Log(WeeklyTask())
}

func TestLocalTasksDedupUsesBodyHash(t *testing.T) {
	type payload struct {
		Value string `json:"value"`
	}

	var mu sync.Mutex
	var wg sync.WaitGroup
	bodies := map[string]int{}

	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		require.NoError(t, err)

		mu.Lock()
		bodies[string(body)]++
		mu.Unlock()
		wg.Done()
	}))
	defer s.Close()

	tasks := LocalTaskScheduler{
		Endpoint:  s.URL,
		Transport: http.DefaultTransport,
	}
	ctx, stop := startLocalTaskScheduler(t, &tasks)
	defer stop()

	wg.Add(2)
	tasks.Post(ctx, Task{Path: "/task/body", Scope: DailyTask(), Body: payload{Value: "a"}})
	tasks.Post(ctx, Task{Path: "/task/body", Scope: DailyTask(), Body: payload{Value: "a"}})
	tasks.Post(ctx, Task{Path: "/task/body", Scope: DailyTask(), Body: payload{Value: "b"}})
	wg.Wait()

	require.Equal(t, 1, bodies[`{"value":"a"}`])
	require.Equal(t, 1, bodies[`{"value":"b"}`])
}

func TestLocalTasksKeepFirstPostedSchedule(t *testing.T) {
	var count atomic.Int32

	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		count.Add(1)
	}))
	defer s.Close()

	tasks := LocalTaskScheduler{
		Endpoint:  s.URL,
		Transport: http.DefaultTransport,
	}
	ctx, stop := startLocalTaskScheduler(t, &tasks)
	defer stop()

	tasks.Post(ctx, Task{
		Path:  "/task/scheduled",
		Scope: DailyTask(),
		Time:  time.Now().Add(80 * time.Millisecond),
	})
	tasks.Post(ctx, Task{Path: "/task/scheduled", Scope: DailyTask()})

	time.Sleep(30 * time.Millisecond)
	require.Equal(t, int32(0), count.Load())

	require.Eventually(t, func() bool {
		return count.Load() == 1
	}, time.Second, 10*time.Millisecond)
}

func startLocalTaskScheduler(t *testing.T, tasks *LocalTaskScheduler) (context.Context, func()) {
	t.Helper()

	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan struct{})

	go func() {
		defer close(done)
		tasks.Start(ctx)
	}()

	return ctx, func() {
		t.Helper()

		cancel()
		select {
		case <-done:
		case <-time.After(time.Second):
			t.Fatal("local task scheduler did not stop")
		}
	}
}
