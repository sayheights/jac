package jac

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"testing"
	"time"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
)

func TestClient_DoAsync(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	registerAsyncHandler()

	type fields struct {
		BaseURL string
	}
	type args struct {
		ctx    context.Context
		asyncs []AsyncRequest
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   []*asyncRequestResponse
	}{
		{
			name: "AysncRequest=All success",
			fields: fields{
				BaseURL: "https://asynctest.com",
			},
			args: args{
				ctx: context.Background(),
				asyncs: []AsyncRequest{
					&asynRequest{id: "abc"},
					&asynRequest{id: "123"},
				},
			},
			want: []*asyncRequestResponse{
				{
					ID:      "abc",
					Success: true,
				},
				{
					ID:      "123",
					Success: true,
				},
			},
		},
		{
			name: "AysncRequest=All success",
			fields: fields{
				BaseURL: "https://asynctest.com",
			},
			args: args{
				ctx: context.Background(),
				asyncs: []AsyncRequest{
					&asynRequest{id: "abc"},
					&asynRequest{id: "123"},
					&asynRequest{id: "fail"},
					&asynRequest{id: "fail", state: 1, fail: true},
				},
			},
			want: []*asyncRequestResponse{
				{
					ID:      "abc",
					Success: true,
				},
				{
					ID:      "123",
					Success: true,
				},
				{
					ID:      "fail",
					Success: false,
				},
				{
					ID:      "fail",
					Success: false,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			retry := &Retry{
				Policy:    DefaultPolicy,
				Backoff:   func(int) time.Duration { return time.Millisecond * 1 },
				MaxAmount: 3,
			}
			c := &Client{
				BaseURL: tt.fields.BaseURL,
				hc:      http.DefaultClient,
				Retry:   retry,
			}

			var got chan *AsyncResponse
			got = c.DoAsync(tt.args.ctx, tt.args.asyncs...)
			var gotAsyncs []*asyncRequestResponse
			for res := range got {
				if res.Err != nil {
					gotAsyncs = append(gotAsyncs, &asyncRequestResponse{"fail", false})
					continue
				}
				var asyncRes asyncRequestResponse
				err := json.Unmarshal(res.Response.Data, &asyncRes)
				if err != nil {
					t.Errorf("%+v", res.Response)
					return
				}
				gotAsyncs = append(gotAsyncs, &asyncRes)
			}
			assert.ElementsMatch(t, tt.want, gotAsyncs)
		})
	}
}

type asynRequest struct {
	*PostRequest
	id    string
	state int
	fail  bool
}

func (a *asynRequest) Path() string {
	return fmt.Sprintf("/async/%d/%s", a.state, a.id)
}

func (a *asynRequest) Body() []byte {
	buf := &bytes.Buffer{}
	buf.WriteString(`{"id":"`)
	buf.WriteString(a.id)
	buf.WriteString(`","state":`)
	buf.WriteString(strconv.Itoa(a.state))
	buf.WriteRune('}')

	return buf.Bytes()
}

func (a *asynRequest) IsReady(res *Response) (bool, error) {
	a.state += 1
	if a.fail {
		return false, errors.New("failed")
	}

	if res.StatusCode == 200 {
		return true, nil
	}

	return false, nil
}

func (a *asynRequest) OnReady() Request {
	a.state = 2
	return a
}

func registerAsyncHandler() {
	httpmock.RegisterResponder(
		"POST",
		"https://asynctest.com/async/0/123",
		httpmock.NewStringResponder(202, `{"id":"123","state":"0"}`),
	)
	httpmock.RegisterResponder(
		"POST",
		"https://asynctest.com/async/1/123",
		httpmock.NewStringResponder(200, `{"id":"123","state":"1"}`),
	)
	httpmock.RegisterResponder(
		"POST",
		"https://asynctest.com/async/2/123",
		httpmock.NewStringResponder(200, `{"id":"123","success":true}`),
	)

	httpmock.RegisterResponder(
		"POST",
		"https://asynctest.com/async/0/abc",
		httpmock.NewStringResponder(202, `{"id":"abc","state":"0"}`),
	)
	httpmock.RegisterResponder(
		"POST",
		"https://asynctest.com/async/1/abc",
		httpmock.NewStringResponder(200, `{"id":"abc","state":"1"}`),
	)
	httpmock.RegisterResponder(
		"POST",
		"https://asynctest.com/async/2/abc",
		httpmock.NewStringResponder(200, `{"id":"abc","success":true}`),
	)
	httpmock.RegisterResponder(
		"POST",
		"https://asynctest.com/async/0/fail",
		httpmock.NewStringResponder(404, `{"id":"fail","success":false}`),
	)
	httpmock.RegisterResponder(
		"POST",
		"https://asynctest.com/async/1/fail",
		httpmock.NewStringResponder(207, `{"id":"fail","success":false}`),
	)
}

type asyncRequestResponse struct {
	ID      string
	Success bool
}
