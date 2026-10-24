package conduit

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	cloudtasks "cloud.google.com/go/cloudtasks/apiv2"
	"cloud.google.com/go/cloudtasks/apiv2/cloudtaskspb"
	"github.com/maxence-charriere/go-app/v11/pkg/app"
	"github.com/maxence-charriere/go-app/v11/pkg/errors"
	"github.com/maxence-charriere/go-app/v11/pkg/logs"
	"github.com/maxence-charriere/sunforge/pkg/stash"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// Task represents a unit of work scheduled for execution at a specific time.
type Task struct {
	// Path is the target path to invoke when the task executes.
	Path string

	// Method is the HTTP method used when invoking Path. Defaults to POST.
	Method string

	// Scope is used to group tasks for deduplication.
	Scope string

	// Time is the scheduled execution time of the task. If zero, the task
	// executes immediately.
	Time time.Time

	// Body is the optional payload sent with the request when the task
	// executes.
	Body any
}

// TaskScheduler schedules tasks for future execution.
type TaskScheduler interface {
	// Post schedules a task for execution.
	Post(context.Context, Task)

	// Close releases any resources held by the scheduler.
	Close() error
}

// SecondTask returns a UTC time bucket identifier with second precision.
func SecondTask() string {
	now := time.Now().UTC()

	return fmt.Sprintf("%v-%v-%v_%v-%v-%v",
		now.Year(),
		now.Month(),
		now.Day(),
		now.Hour(),
		now.Minute(),
		now.Second(),
	)
}

// MinuteTask returns a UTC time bucket identifier with minute precision.
func MinuteTask() string {
	now := time.Now().UTC()

	return fmt.Sprintf("%v-%v-%v_%v-%v",
		now.Year(),
		now.Month(),
		now.Day(),
		now.Hour(),
		now.Minute(),
	)
}

// SemiDailyTask returns a UTC time bucket identifier split by AM/PM.
func SemiDailyTask() string {
	now := time.Now().UTC()
	var part string
	if now.Hour() < 12 {
		part = "am"
	} else {
		part = "pm"
	}

	return fmt.Sprintf("%v-%v-%v-%s",
		now.Year(),
		now.Month(),
		now.Day(),
		part,
	)
}

// ThriceDailyTask returns a UTC time bucket identifier split into night,
// morning, and afternoon.
func ThriceDailyTask() string {
	now := time.Now().UTC()
	var part string
	if hour := now.UTC().Hour(); hour < 8 {
		part = "night"
	} else if hour < 16 {
		part = "morning"
	} else {
		part = "afternoon"
	}

	return fmt.Sprintf("%v-%v-%v-%s",
		now.Year(),
		now.Month(),
		now.Day(),
		part,
	)
}

// DailyTask returns a UTC time bucket identifier with day precision.
func DailyTask() string {
	now := time.Now().UTC()

	return fmt.Sprintf("%v-%v-%v",
		now.Year(),
		now.Month(),
		now.Day(),
	)
}

// WeeklyTask returns a UTC time bucket identifier with week-of-month precision.
func WeeklyTask() string {
	now := time.Now().UTC()

	week := 1
	switch d := now.Day(); {
	case d > 28:
		week = 5

	case d > 21:
		week = 4

	case d > 14:
		week = 3

	case d > 7:
		week = 2
	}

	return fmt.Sprintf("%v-%v-%vw",
		now.Year(),
		now.Month(),
		week,
	)
}

// LocalTaskScheduler schedules tasks in-memory and executes them locally via a
// worker pool.
type LocalTaskScheduler struct {
	// Cache is used to deduplicate tasks.
	Cache stash.Cache

	// Transport is used for task HTTP requests.
	Transport http.RoundTripper

	// WorkerCount is the number of worker goroutines.
	WorkerCount int

	// Endpoint is the base endpoint used to execute tasks.
	Endpoint string

	// WorkerHTTPHeader is added to requests executed by workers.
	WorkerHTTPHeader string

	// BypassCacheURL lists URL substrings that skip deduplication.
	BypassCacheURL []string

	// Encode serializes task bodies before they are attached to HTTP requests.
	// If nil, json.Marshal is used.
	Encode func(any) ([]byte, error)

	once          sync.Once
	mu            sync.Mutex
	queue         *queue[Task]
	workers       []chan Task
	nextWorkerIdx int
}

func (s *LocalTaskScheduler) init() {
	if s.Cache == nil {
		s.Cache = &stash.ExpiringCache{ItemTTL: time.Hour}
	}

	if s.WorkerCount < 1 {
		s.WorkerCount = 1
	}

	if s.Transport == nil {
		s.Transport = http.DefaultTransport
	}

	if s.Encode == nil {
		s.Encode = json.Marshal
	}

	workers := make([]chan Task, s.WorkerCount)
	for i := 0; i < s.WorkerCount; i++ {
		workers[i] = make(chan Task)
	}
	s.workers = workers

	s.queue = newQueue[Task](300000)
}

// Post schedules a task for local execution.
func (s *LocalTaskScheduler) Post(ctx context.Context, t Task) {
	s.once.Do(s.init)

	taskID, err := s.taskID(t)
	if err != nil {
		app.Log(errors.New("encoding local task body failed").
			WithTag("path", t.Path).
			Wrap(err))
		return
	}

	if !s.shouldExecute(ctx, taskID, t.Path) {
		return
	}

	if delay := time.Until(t.Time); delay > 0 {
		go func() {
			time.Sleep(delay)
			s.queue.Enqueue(t)
		}()
		return
	}
	s.queue.Enqueue(t)
}

// Close stops the scheduler and releases resources.
func (s *LocalTaskScheduler) Close() error {
	return nil
}

func (s *LocalTaskScheduler) Start(ctx context.Context) {
	s.once.Do(s.init)

	for i, w := range s.workers {
		go s.startWorker(ctx, i+1, w)
	}

	for {
		select {
		case <-ctx.Done():
			return

		default:
			task, ok := s.queue.Dequeue()
			if !ok {
				time.Sleep(time.Millisecond * 10)
				continue
			}

			s.nextWorker() <- task
		}
	}
}

func (s *LocalTaskScheduler) startWorker(ctx context.Context, id int, worker chan Task) {
	for {
		select {
		case task := <-worker:
			taskID, _ := s.taskID(task)
			app.Log(logs.New("executing task").
				WithTag("worker-id", id).
				WithTag("id", taskID))
			if err := s.execTask(ctx, id, task); err != nil {
				app.Log(errors.New("executing task failed").
					WithTag("worker-id", id).
					WithTag("path", task.Path).
					Wrap(err))
			}

		case <-ctx.Done():
			return
		}
	}
}

func (s *LocalTaskScheduler) nextWorker() chan Task {
	w := s.workers[s.nextWorkerIdx]
	s.nextWorkerIdx = (s.nextWorkerIdx + 1) % len(s.workers)
	return w
}

func (s *LocalTaskScheduler) execTask(ctx context.Context, workerID int, task Task) error {
	url := strings.TrimSuffix(s.Endpoint, "/") + task.Path

	method := task.Method
	if method == "" {
		method = http.MethodPost
	}

	var body io.Reader
	if method == http.MethodPost && task.Body != nil {
		jsonBody, err := s.Encode(task.Body)
		if err != nil {
			return errors.New("encoding task body failed").Wrap(err)
		}
		body = bytes.NewBuffer(jsonBody)
	}

	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return errors.New("creating task request failed").
			WithTag("url", url).
			Wrap(err)
	}
	if s.WorkerHTTPHeader != "" {
		req.Header.Set(s.WorkerHTTPHeader, strconv.Itoa(workerID))
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	res, err := s.Transport.RoundTrip(req)
	if err != nil {
		return errors.New("task request failed").Wrap(err)
	}
	defer res.Body.Close()

	return nil
}

func (s *LocalTaskScheduler) taskID(task Task) (string, error) {
	method := task.Method
	if method == "" {
		method = http.MethodPost
	}

	taskID := fmt.Sprintf("%s-%s", task.Scope, task.Path)
	if method != http.MethodPost || task.Body == nil {
		return taskID, nil
	}

	jsonBody, err := s.Encode(task.Body)
	if err != nil {
		return "", err
	}
	return taskID + "-" + hash(jsonBody), nil
}

func (s *LocalTaskScheduler) shouldExecute(ctx context.Context, taskID, path string) bool {
	for _, bypass := range s.BypassCacheURL {
		if strings.Contains(path, bypass) {
			return true
		}
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	_, cached := s.Cache.Get(ctx, taskID)
	if cached {
		return false
	}

	s.Cache.Set(ctx, taskID, 1)
	return true
}

// CloudTaskScheduler schedules tasks using Google Cloud Tasks.
type CloudTaskScheduler struct {
	// Endpoint is the base URL used to build the full request URL for a task.
	Endpoint string

	// ProjectID is the Google Cloud project identifier.
	ProjectID string

	// LocationID is the Google Cloud Tasks location (region) identifier.
	LocationID string

	// DedupID is a stable identifier used as part of the task name for
	// deduplication.
	DedupID string

	// Authorization is the value of the Authorization HTTP header to include
	// in requests made by the task. If empty, no Authorization header is sent.
	Authorization string

	// QueueID returns the queue identifier to use for the given task.
	QueueID func(Task) string

	// Encode serializes task bodies before they are attached to HTTP tasks.
	// If nil, json.Marshal is used.
	Encode func(any) ([]byte, error)

	once   sync.Once
	client *cloudtasks.Client
}

// Post creates a Cloud Tasks task for the provided Task.
func (s *CloudTaskScheduler) Post(ctx context.Context, task Task) {
	s.once.Do(func() {
		s.init(ctx)
	})

	queueID := s.QueueID(task)

	var scheduledTime *timestamppb.Timestamp
	if !task.Time.IsZero() {
		scheduledTime = timestamppb.New(task.Time)
	}

	var method cloudtaskspb.HttpMethod
	switch task.Method {
	case http.MethodGet:
		method = cloudtaskspb.HttpMethod_GET

	case http.MethodPut:
		method = cloudtaskspb.HttpMethod_PUT

	case http.MethodDelete:
		method = cloudtaskspb.HttpMethod_DELETE

	default:
		method = cloudtaskspb.HttpMethod_POST
	}

	var taskBody []byte

	taskHeaders := make(map[string]string)
	if s.Authorization != "" {
		taskHeaders["Authorization"] = s.Authorization
	}

	taskName := fmt.Sprintf("projects/%s/locations/%s/queues/%s/tasks/%s_%s_%s",
		s.ProjectID,
		s.LocationID,
		queueID,
		s.DedupID,
		task.Scope,
		hash([]byte(task.Path)),
	)
	if method == cloudtaskspb.HttpMethod_POST && task.Body != nil {
		jsonBody, err := s.Encode(task.Body)
		if err != nil {
			app.Log(errors.New("encoding cloud task body failed").
				WithTag("path", task.Path).
				Wrap(err))
			return
		}
		taskName += "_" + hash(jsonBody)
		taskHeaders["Content-Type"] = "application/json"
		taskBody = jsonBody
	}

	req := &cloudtaskspb.CreateTaskRequest{
		Parent: fmt.Sprintf("projects/%s/locations/%s/queues/%s",
			s.ProjectID,
			s.LocationID,
			queueID,
		),
		Task: &cloudtaskspb.Task{
			Name:         taskName,
			ScheduleTime: scheduledTime,
			MessageType: &cloudtaskspb.Task_HttpRequest{
				HttpRequest: &cloudtaskspb.HttpRequest{
					HttpMethod: method,
					Url:        strings.TrimSuffix(s.Endpoint, "/") + task.Path,
					Headers:    taskHeaders,
					Body:       taskBody,
				},
			},
		},
	}

	_, err := s.client.CreateTask(ctx, req)
	if status.Code(err) == codes.AlreadyExists {
		return
	}

	if err != nil {
		app.Log(errors.New("creating cloud task failed").
			WithTag("path", task.Path).
			Wrap(err))
	}
}

// Close stops the scheduler and releases resources.
func (s *CloudTaskScheduler) Close() error {
	if s.client != nil {
		return s.client.Close()
	}
	return nil
}

func (s *CloudTaskScheduler) init(ctx context.Context) {
	client, err := cloudtasks.NewClient(ctx)
	if err != nil {
		panic(errors.New("creating cloud task client failed").Wrap(err))
	}
	s.client = client

	if s.Encode == nil {
		s.Encode = json.Marshal
	}

	s.DedupID = strings.ReplaceAll(s.DedupID, ".", "-")
}

func hash(v []byte) string {
	sum := sha256.Sum256(v)
	return hex.EncodeToString(sum[:])
}
