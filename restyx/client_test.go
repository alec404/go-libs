package restyx

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"
)

func TestClientGet(t *testing.T) {
	c := NewClient()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		if r.Method != http.MethodGet {
			t.Fatalf("expected GET, got %s", r.Method)
		}

		w.WriteHeader(http.StatusOK)
		sec := r.Header.Get("x-sec")
		var resp = make(map[string]interface{})
		resp["status"] = 200
		resp["message"] = "ok"
		resp["data"] = map[string]interface{}{
			"a":   r.URL.Query().Get("a"),
			"sec": sec,
		}
		_ = json.NewEncoder(w).Encode(resp)
	}))
	defer srv.Close()

	// 发送GET请求获取内容
	var result map[string]interface{}

	resp, err := c.R().SetHeader("x-sec", "sec-123").SetQueryParam("a", "1").SetResult(&result).Get(srv.URL)
	if err != nil {
		t.Fatalf("get request failed: %v", err)
	}

	if resp.IsError() {
		t.Fatalf("unexpected error response: %s", resp.Status())
	}

	if resp.StatusCode() != http.StatusOK {
		t.Fatalf("unexpected status: %d", resp.StatusCode())
	}

	status, ok := result["status"].(float64)
	if !ok {
		t.Fatalf("status missing or not number: %v", result["status"])
	}
	if status != 200 {
		t.Fatalf("unexpected status: %v", status)
	}

	msg, ok := result["message"].(string)
	if !ok || msg != "ok" {
		t.Fatalf("unexpected message: %#v", result["message"])
	}

	data, ok := result["data"].(map[string]interface{})
	if !ok {
		t.Fatalf("unexpected data: %#v", result["data"])
	}
	if data["a"].(string) != "1" || data["sec"].(string) != "sec-123" {
		t.Fatalf("unexpected data: %+v", data)
	}
}

type ResponseBody struct {
	Status  int                    `json:"status"`
	Message string                 `json:"message"`
	Data    map[string]interface{} `json:"data"`
}

func TestClientGetWithRetry(t *testing.T) {
	c := NewClient()
	c.SetRetryCount(1)

	var callCount int32

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&callCount, 1)

		w.Header().Set("Content-Type", "application/json")

		if r.Method != http.MethodGet {
			t.Fatalf("expected GET, got %s", r.Method)
		}

		// 第一次：返回 500（触发 retry）
		if atomic.LoadInt32(&callCount) == 1 {
			w.WriteHeader(http.StatusInternalServerError)
			_ = json.NewEncoder(w).Encode(ResponseBody{
				Status:  500,
				Message: "fail once",
			})
			return
		}

		w.WriteHeader(http.StatusOK)
		sec := r.Header.Get("x-sec")
		_ = json.NewEncoder(w).Encode(ResponseBody{
			Status:  200,
			Message: "ok",
			Data: map[string]interface{}{
				"a":   r.URL.Query().Get("a"),
				"sec": sec,
			},
		})
	}))
	defer srv.Close()

	// 发送GET请求获取内容
	var result ResponseBody

	resp, err := c.R().SetHeader("x-sec", "sec-123").SetQueryParam("a", "1").SetResult(&result).Get(srv.URL)
	if err != nil {
		t.Fatalf("get request failed: %v", err)
	}

	if resp.IsError() {
		t.Fatalf("unexpected error response: %s", resp.Status())
	}
	if got := atomic.LoadInt32(&callCount); got != 2 {
		t.Fatalf("expected 2 calls (1 retry), got %d", got)
	}
	if resp.StatusCode() != http.StatusOK {
		t.Fatalf("unexpected status: %d", resp.StatusCode())
	}
	if result.Status != 200 || result.Message != "ok" {
		t.Fatalf("unexpected result: %+v", result)
	}
}

func withMinTimeout(ctx context.Context, d time.Duration) (context.Context, context.CancelFunc) {
	if d <= 0 {
		return ctx, func() {}
	}

	want := time.Now().Add(d)
	if dl, ok := ctx.Deadline(); ok && dl.Before(want) {
		// 外层更短，就尊重外层，不再创建新 ctx
		return ctx, func() {}
	}
	return context.WithDeadline(ctx, want)
}

func TestClientGetSetContextWithTimeout(t *testing.T) {
	c := NewClient()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("expected GET, got %s", r.Method)
		}
		time.Sleep(200 * time.Millisecond)

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	}))
	defer srv.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 自定义超时时间，取短的时间
	reqCtx, cancel := withMinTimeout(ctx, 100*time.Millisecond)
	defer cancel()

	_, err := c.R().SetContext(reqCtx).Get(srv.URL)

	if !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("get request failed: %v", err)
	}

}
