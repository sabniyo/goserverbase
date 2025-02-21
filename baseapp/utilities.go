package baseapp

import (
	"bytes"
	"context"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/sabariramc/goserverbase/log"
	"github.com/sabariramc/goserverbase/utils"
)

func (b *BaseApp) GetHttpCorrelationParams(r *http.Request) *log.CorrelationParam {
	correlationId := r.Header.Get("x-correlation-id")
	if correlationId == "" {
		return log.GetDefaultCorrelationParams(b.c.ServiceName)
	}
	return &log.CorrelationParam{
		CorrelationId: correlationId,
		ScenarioId:    r.Header.Get("x-scenario-id"),
		ScenarioName:  r.Header.Get("x-scenario-name"),
		SessionId:     r.Header.Get("x-session-id"),
	}
}

func (b *BaseApp) PrintRequest(ctx context.Context, r *http.Request) {
	h := r.Header
	popList := make(map[string][]string)
	for _, key := range b.lConfig.AuthHeaderKeyList {
		val := h.Values(key)
		if len(val) != 0 {
			popList[key] = val
			h.Set(key, "---redacted---")
		}
	}
	b.log.Info(ctx, "Request", map[string]interface{}{
		"Method":        r.Method,
		"Header":        h,
		"URL":           r.URL,
		"Proto":         r.Proto,
		"ContentLength": r.ContentLength,
		"Host":          r.Host,
		"RemoteAddr":    r.RemoteAddr,
		"RequestURI":    r.RequestURI,
	})
	if r.ContentLength > 0 {
		body := r.Body
		defer body.Close()
		blobBody, _ := ioutil.ReadAll(body)
		data := make(map[string]any)
		json.Unmarshal(blobBody, &data)
		r.Body = io.NopCloser(bytes.NewReader(blobBody))
		b.log.Debug(ctx, "Request Body", data)
	}
	for key, value := range popList {
		h.Del(key)
		for _, v := range value {
			h.Add(key, v)
		}
	}
}

func (b *BaseApp) GetCorrelationContext(ctx context.Context, c *log.CorrelationParam) context.Context {
	ctx = context.WithValue(ctx, log.ContextKeyCorrelation, c)
	return ctx
}

func (b *BaseApp) GetPort() string {
	return fmt.Sprintf("%v:%v", b.c.Host, b.c.Port)
}

func (b *BaseApp) SetHandlerError(ctx context.Context, err error) {
	iSetter := ctx.Value(ContextKeyError)
	if iSetter == nil {
		return
	}
	setter, ok := iSetter.(func(error))
	if !ok {
		panic(fmt.Errorf("context error handler corrupted, error to handle: %w", err))
	}
	setter(err)
}

func GetBytes(key interface{}) ([]byte, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(key)
	if err != nil {
		return nil, fmt.Errorf("baseapp.GetBytes: %w", err)
	}
	return buf.Bytes(), nil
}

type Filter struct {
	PageNo int64  `json:"pageNo" schema:"pageNo"`
	Limit  int64  `json:"limit" schema:"limit"`
	SortBy string `json:"sortBy" schema:"sortBy"`
	Asc    *bool  `json:"asc" schema:"asc"`
}

func SetDefaultPagination(filter interface{}, defaultSortBy string) error {
	var defaultFilter Filter
	err := utils.StrictJsonTransformer(filter, &defaultFilter)
	if err != nil {
		return fmt.Errorf("baseapp.SetDefault : %w", err)
	}
	if defaultFilter.PageNo <= 0 {
		defaultFilter.PageNo = 1
	}
	if defaultFilter.Limit <= 0 {
		defaultFilter.Limit = 10
	}
	if defaultFilter.SortBy == "" {
		defaultFilter.SortBy = defaultSortBy
	}
	if defaultFilter.Asc == nil {
		v := true
		defaultFilter.Asc = &v
	}
	err = utils.StrictJsonTransformer(&defaultFilter, filter)
	if err != nil {
		return fmt.Errorf("app.SetDefault : %w", err)
	}
	return nil
}

func WriteJsonWithStatusCode(w http.ResponseWriter, statusCode int, responseBody any) {
	var err error
	blob, ok := responseBody.([]byte)
	if !ok {
		blob, err = json.Marshal(responseBody)
		if err != nil {
			panic(fmt.Errorf("response marshal error: %w", err))
		}
	}
	w.Header().Add(HttpHeaderContentType, HttpContentTypeJSON)
	w.WriteHeader(statusCode)
	w.Write(blob)
}

func WriteJson(w http.ResponseWriter, responseBody any) {
	WriteJsonWithStatusCode(w, http.StatusOK, responseBody)
}
