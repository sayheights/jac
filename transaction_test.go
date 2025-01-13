package jac

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/kr/pretty"
	"github.com/stretchr/testify/assert"
)

func Test_txnState_isDone(t *testing.T) {
	tests := []struct {
		name string
		tr   txnState
		want bool
	}{
		{
			tr:   txnRetryable,
			want: false,
		},
		{
			tr:   txnUnrecoverable,
			want: true,
		},
		{
			tr:   txnInitial,
			want: false,
		},
		{
			tr:   txnExhausted,
			want: true,
		},
		{
			tr:   txnSuccessful,
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				if got := tt.tr.isDone(); got != tt.want {
					t.Errorf("txnState.isDone() = %v, want %v", got, tt.want)
				}
			},
		)
	}
}

func Test_transaction_init(t *testing.T) {
	tests := []struct {
		name    string
		req     *http.Request
		wantErr bool
	}{
		{
			req: &http.Request{
				Body: &errReader{},
			},
			wantErr: true,
		},
		{
			req: &http.Request{
				Body: io.NopCloser(strings.NewReader(`{"test": "testvalue", "amount": 10}"`)),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				tr := &transaction{
					req: tt.req,
				}
				if err := tr.init(); (err != nil) != tt.wantErr {
					t.Errorf("transaction.init() error = %v, wantErr %v", err, tt.wantErr)
				}
			},
		)
	}
}

func Test_transaction_resolveState(t *testing.T) {
	type fields struct {
		res   *http.Response
		count int
		err   error
		ret   *Retry
	}
	tests := []struct {
		name   string
		fields fields
		want   txnState
	}{
		{
			name: "retries exhausted",
			fields: fields{
				count: 6,
				err:   errors.New("404"),
				ret: &Retry{
					Policy:    nil,
					Backoff:   nil,
					MaxAmount: 5,
				},
			},
			want: txnExhausted,
		},
		{
			name: "error and retry",
			fields: fields{
				count: 4,
				err:   errors.New("404"),
				ret: &Retry{
					Policy:    nil,
					Backoff:   nil,
					MaxAmount: 5,
				},
			},
			want: txnRetryable,
		},
		{
			name: "success",
			fields: fields{
				count: 4,
				err:   nil,
				ret: &Retry{
					Policy:    nil,
					Backoff:   nil,
					MaxAmount: 1,
				},
				res: &http.Response{
					StatusCode: 201,
					Body:       io.NopCloser(strings.NewReader(`{"test": "testvalue", "amount": 10}"`)),
				},
			},
			want: txnSuccessful,
		},
		{
			name: "error response",
			fields: fields{
				count: 4,
				err:   nil,
				ret: &Retry{
					Policy:    RetryOn(defRetryCodes),
					Backoff:   nil,
					MaxAmount: 5,
				},
				res: &http.Response{
					StatusCode: 502,
					Body:       io.NopCloser(strings.NewReader(`{"test": "testvalue", "amount": 10}"`)),
				},
			},
			want: txnRetryable,
		},
		{
			name: "error response",
			fields: fields{
				count: 4,
				err:   nil,
				ret: &Retry{
					Policy:    RetryOn(defRetryCodes),
					Backoff:   nil,
					MaxAmount: 1,
				},
				res: &http.Response{
					StatusCode: 401,
					Body:       io.NopCloser(strings.NewReader(`{"test": "testvalue", "amount": 10}"`)),
				},
			},
			want: txnUnrecoverable,
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				tr := &transaction{
					res:   tt.fields.res,
					count: tt.fields.count,
					err:   tt.fields.err,
					ret:   tt.fields.ret,
				}
				if got := tr.resolveState(); got != tt.want {
					t.Errorf("transaction.resolveState() = %v, want %v", got, tt.want)
				}
			},
		)
	}
}

type attemptDoer struct{}

func (t *attemptDoer) Do(req *http.Request) (*http.Response, error) {
	return &http.Response{
		StatusCode: 201,
		Body:       io.NopCloser(strings.NewReader("success")),
	}, nil
}

func Test_transaction_attempt(t *testing.T) {
	a := attemptDoer{}
	type fields struct {
		id  string
		req *http.Request
		hc  doer
	}
	tests := []struct {
		name   string
		fields fields
		want   txnState
	}{
		{
			name: "",
			fields: fields{
				id:  "",
				req: &http.Request{},
				hc:  &a,
			},
			want: txnSuccessful,
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				tr := &transaction{
					id:           tt.fields.id,
					req:          tt.fields.req,
					hc:           tt.fields.hc,
					isSuccessful: defaultIsSuccessful,
				}
				if got := tr.attempt(); got != tt.want {
					t.Errorf("transaction.attempt() = %v, want %v", got, tt.want)
				}
			},
		)
	}
}

func Test_transaction_reset(t *testing.T) {
	type fields struct {
		res   *http.Response
		body  *bytes.Reader
		count int
		wait  time.Duration
		ret   *Retry
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{
			name: "",
			fields: fields{
				count: 2,
				res:   nil,
				body:  &bytes.Reader{},
				wait:  0,
				ret: &Retry{
					Policy:    nil,
					Backoff:   LinearBackoff(time.Second * 5),
					MaxAmount: 0,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				tr := &transaction{
					res:   tt.fields.res,
					body:  tt.fields.body,
					count: tt.fields.count,
					wait:  tt.fields.wait,
					ret:   tt.fields.ret,
				}
				tr.reset()
				// If any further test is required change the structure to other go test structure
				// Because in this assert there is an assumption to only having one test case
				assert.Equal(t, time.Second*10, tr.wait)
			},
		)
	}
}

func Test_transaction_result(t *testing.T) {
	data, _ := io.ReadAll(io.NopCloser(strings.NewReader(`{"test": "testvalue", "amount": 10}"`)))
	type fields struct {
		res   *http.Response
		count int
		err   error
		start time.Time
		end   time.Time
		state txnState
		req   *http.Request
	}
	tests := []struct {
		name    string
		fields  fields
		want    *Response
		wantErr bool
	}{
		{
			name: "",
			fields: fields{
				count: 0,
				err:   errors.New("test error"),
				start: time.Time{},
				end:   time.Time{},
				state: txnDone,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "",
			fields: fields{
				start: time.Time{},
				end:   time.Time{},
				res: &http.Response{
					Body: io.NopCloser(strings.NewReader(`{"test": "testvalue", "amount": 10}"`)),
					Request: &http.Request{
						URL: &url.URL{
							Path: "test",
						},
					},
				},
				req: &http.Request{
					URL: &url.URL{
						Path: "test",
					},
				},
				state: txnSuccessful,
			},
			want: &Response{
				URL: &url.URL{
					Path: "test",
				},
				RequestURI: "test",
				Data:       data,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				tr := &transaction{
					req:   tt.fields.req,
					res:   tt.fields.res,
					count: tt.fields.count,
					err:   tt.fields.err,
					start: tt.fields.start,
					end:   tt.fields.end,
					state: tt.fields.state,
				}
				got, err := tr.result()
				if (err != nil) != tt.wantErr {
					t.Errorf("transaction.result() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				if !reflect.DeepEqual(got, tt.want) {
					fmt.Printf("%# v", pretty.Formatter(got.Data))
					// fmt.Printf("%# v", pretty.Formatter(got.Data))
					t.Errorf("transaction.result() = %v, want %v", got, tt.want)
				}
			},
		)
	}
}

type retryDoer struct {
	failAmount int
	failCode   int

	successCode int
	attempt     int
}

func (t *retryDoer) Do(req *http.Request) (*http.Response, error) {
	res := &http.Response{StatusCode: t.failCode}
	if t.attempt >= t.failAmount {
		res.Body = req.Body
		res.Request = req
		res.StatusCode = t.successCode
	}
	t.attempt++

	return res, nil
}

type errReader struct{}

func (e *errReader) Close() error {
	return nil
}

func (e *errReader) Read(p []byte) (n int, err error) {
	return 0, errors.New("jac: forced error")
}

func Test_checkIfSuccessful(t *testing.T) {
	isSuccessFunc := func(res *http.Response) bool {
		body, err := io.ReadAll(res.Body)
		if err != nil {
			return false
		}

		return string(body) != "40100"
	}

	response := &http.Response{
		StatusCode: 200,
		Body:       io.NopCloser(strings.NewReader("40100")),
	}

	got := checkIfSuccessful(isSuccessFunc, response)
	assert.Equal(t, false, got)

	body, err := io.ReadAll(response.Body)
	assert.NoError(t, err)
	assert.Equal(t, "40100", string(body))
}
