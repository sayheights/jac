package jac

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"time"

	"github.com/darrae/jac/internal/jachttp"
	"github.com/google/uuid"
)

var errRetriesExhausted = errors.New("Retry attempts exhausted")

type txnState int

const (
	txnInitial txnState = 1 << iota
	txnSuccessful
	txnExhausted
	txnUnrecoverable
	txnRetryable
	txnResponseReady
)

const (
	txnDone = txnSuccessful | txnExhausted | txnUnrecoverable
	txnFail = txnExhausted | txnUnrecoverable
)

func (t txnState) isDone() bool {
	return t&txnDone != 0
}

func (t txnState) isFail() bool {
	return t&txnFail != 0
}

type doer interface {
	Do(*http.Request) (*http.Response, error)
}

type transaction struct {
	id           string
	req          *http.Request
	res          *http.Response
	body         *bytes.Reader
	count        int
	wait         time.Duration
	err          error
	start        time.Time
	end          time.Time
	ret          *Retry
	state        txnState
	response     *Response
	hc           doer
	isSuccessful func(*http.Response) bool
}

func (t *transaction) init() error {
	if t.req.Body != nil {
		reqData, err := io.ReadAll(t.req.Body)
		if err != nil {
			return err
		}
		t.body = bytes.NewReader(reqData)
		t.req.Body = io.NopCloser(t.body)
	}
	t.id = uuid.NewString()
	t.start = time.Now()
	t.state = txnInitial
	if t.hc == nil {
		t.hc = jachttp.Client
	}

	return nil
}

func (t *transaction) attempt() txnState {
	select {
	case <-t.req.Context().Done():
		t.err = t.req.Context().Err()
		return txnUnrecoverable
	case <-time.After(t.wait):
	}

	t.count++
	t.res, t.err = t.hc.Do(t.req)

	t.state = t.resolveState()

	return t.state
}

func (t *transaction) resolveState() txnState {
	if t.err != nil {
		return t.noResponse()
	}
	return t.receivedResponse()
}

func (t *transaction) noResponse() txnState {
	if t.ret.MaxAmount <= t.count {
		t.err = errRetriesExhausted
		return txnExhausted
	}
	return txnRetryable
}

func (t *transaction) receivedResponse() txnState {
	if checkIfSuccessful(t.isSuccessful, t.res) {
		return txnSuccessful
	}
	t.err = errors.New("Unexpected status code")
	if t.res != nil && t.res.Body != nil {
		if data, err := io.ReadAll(t.res.Body); err == nil {
			t.err = errors.New("Unexpected status code, error body: " + string(data))
		}
	}
	if t.ret.Policy(t.res) {
		if t.ret.MaxAmount <= t.count {
			t.err = errRetriesExhausted
			return txnExhausted
		}
		return txnRetryable
	}

	return txnUnrecoverable
}

func (t *transaction) reset() {
	_, err := t.body.Seek(0, 0)
	if err != nil {
		panic(err)
	}

	if t.res != nil && t.state != txnSuccessful {
		t.res.Body.Close()
	}

	t.wait = t.ret.Backoff(t.count)
}

func (t *transaction) result() (*Response, error) {
	if t.state != txnSuccessful {
		if t.res == nil {
			return nil, t.err
		}
		data, err := io.ReadAll(t.res.Body)
		defer t.res.Body.Close()
		if err != nil {
			return nil, err
		}

		t.response = &Response{
			Data: data,
		}
		return t.response, t.err
	}

	data, err := io.ReadAll(t.res.Body)
	defer t.res.Body.Close()
	if err != nil {
		return nil, err
	}

	t.response = &Response{
		URL:          t.req.URL,
		RequestURI:   t.res.Request.URL.RequestURI(),
		Data:         data,
		Header:       t.res.Header,
		Duration:     t.end.Sub(t.start),
		AttemptCount: t.count,
		StatusCode:   t.res.StatusCode,
	}

	t.state = txnResponseReady

	return t.response, nil
}

func checkIfSuccessful(fn func(*http.Response) bool, res *http.Response) bool {
	if fn == nil {
		fn = defaultIsSuccessful
	}
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return false
	}
	defer res.Body.Close()

	buf := bytes.NewReader(body)
	res.Body = io.NopCloser(buf)
	isSuccessful := fn(res)
	_, err = buf.Seek(0, 0)
	if err != nil {
		res.Body = io.NopCloser(buf)
	}

	return isSuccessful
}
