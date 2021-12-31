package logger

import (
	"context"
	"net/http"
	"time"
)

type ctxKey string

func (c ctxKey) String() string {
	return string(c)
}

const (
	newHTTPWriterCtxKey ctxKey = "NewHTTPWriter"
	dataMapCtxKey       ctxKey = "DataMap"
)

// keys in dataMap context
const (
	errorWithStackKey = "ErrStack"
	httpRequestKey    = "HTTPRequest"
	errorID           = "ErrID"
)

type grpcMessage struct {
	ServiceName string      `json:"ServiceName,omitempty"`
	ReqID       string      `json:"ReqID,omitempty"`
	Method      string      `json:"Method,omitempty"`
	Request     interface{} `json:"Request,omitempty"`
	Status      string      `json:"Status,omitempty"`
	Duration    string      `json:"Duration,omitempty"`
	Response    interface{} `json:"Response,omitempty"`
	Error       string      `json:"Error,omitempty"`
	ErrID       string      `json:"ErrID,omitempty"`
}

// HTTPWriter ...
type HTTPWriter struct {
	W http.ResponseWriter
	a []byte
}

func (w *HTTPWriter) Write(b []byte) (int, error) {
	w.a = b
	return w.W.Write(b)
}

// Header ...
func (w *HTTPWriter) Header() http.Header {
	return w.W.Header()
}

// WriteHeader ...
func (w *HTTPWriter) WriteHeader(statusCode int) {
	w.W.WriteHeader(statusCode)
}

// WrapperServer ...
type WrapperServer struct {
	TwirpServer
}

// TwirpServer ...
type TwirpServer interface {
	http.Handler
	ServiceDescriptor() ([]byte, int)
	ProtocGenTwirpVersion() string
	PathPrefix() string
}

func (s *WrapperServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var nw = new(HTTPWriter)
	nw.W = w
	r = r.WithContext(context.WithValue(r.Context(), newHTTPWriterCtxKey, nw))

	dataMap := make(map[string]interface{})
	r = r.WithContext(context.WithValue(r.Context(), dataMapCtxKey, dataMap))

	s.TwirpServer.ServeHTTP(nw, r)
}

var reqStartTimestampKey = new(int)

func markReqStart(ctx context.Context) context.Context {
	return context.WithValue(ctx, reqStartTimestampKey, time.Now())
}

func getReqStart(ctx context.Context) (time.Time, bool) {
	t, ok := ctx.Value(reqStartTimestampKey).(time.Time)
	return t, ok
}

func getReq(ctx context.Context) interface{} {
	dataMap, _ := getDataMap(ctx)
	return dataMap[httpRequestKey]
}

func getResp(ctx context.Context) (*HTTPWriter, bool) {
	t, ok := ctx.Value(newHTTPWriterCtxKey).(*HTTPWriter)
	return t, ok
}

func getDataMap(ctx context.Context) (map[string]interface{}, bool) {
	t, ok := ctx.Value(dataMapCtxKey).(map[string]interface{})
	return t, ok
}
