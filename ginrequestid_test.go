package ginrequestid

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"
)

var testGeneator RequestIDGenerator = func(c *gin.Context) string {
	return "test"
}

func init() {
	gin.SetMode(gin.ReleaseMode)
}

type idResponse struct {
	Header  string
	CtxOK   bool
	CtxText string
}

func (got idResponse) Matches(t *testing.T, want idResponse) {
	//check header
	compile, err := regexp.Compile(want.Header)
	if err != nil {
		t.Errorf("invalid Header regex: %v", err)
		return
	}
	if !compile.MatchString(got.Header) {
		t.Errorf("Header = %v, want %v", got.Header, want.Header)
	}
	//check ctx
	if got.CtxOK != want.CtxOK {
		t.Errorf("CtxOK = %v, want %v", got.CtxOK, want.CtxOK)
	}

	compile, err = regexp.Compile(want.CtxText)
	if err != nil {
		t.Errorf("invalid CtxText regex: %v", err)
		return
	}
	if !compile.MatchString(got.CtxText) {
		t.Errorf("CtxText = %v, want %v", got.Header, want.Header)
	}
}

func newTestRouter(cfg Config, k string) *gin.Engine {
	gin.SetMode("release")
	router := gin.New()
	router.Use(RequestID(cfg))
	router.GET("/", func(c *gin.Context) {
		v, ok := c.Get(k)
		if !ok {
			v = ""
		}
		c.JSON(200, &idResponse{
			Header:  c.GetHeader(k),
			CtxOK:   ok,
			CtxText: v.(string),
		})
	})
	return router
}

func TestRequestID(t *testing.T) {

	type args struct {
		cfg         Config
		expectedKey string
	}
	tests := []struct {
		name string
		args args
		want idResponse
	}{
		{
			name: "default",
			args: args{DefaultCfg(), defaultHeader},
			want: idResponse{
				Header:  "",
				CtxOK:   true,
				CtxText: "[a-f0-9]{8}(-[a-f0-9]{4}){3}-[a-f0-9]{12}", //uuid
			},
		},
		{
			name: "empty config",
			args: args{Config{}, defaultHeader},
			want: idResponse{
				Header:  "",
				CtxOK:   false,
				CtxText: "",
			},
		},
		{
			name: "update with custom header",
			args: args{Config{SetReqHeader: true, Header: "Test"}, "Test"},
			want: idResponse{
				Header:  "[a-f0-9]{8}(-[a-f0-9]{4}){3}-[a-f0-9]{12}", //uuid
				CtxOK:   false,
				CtxText: "",
			},
		},
		{
			name: "both outputs enabled",
			args: args{Config{SetReqHeader: true, SetGinCtx: true}, defaultHeader},
			want: idResponse{
				Header:  "[a-f0-9]{8}(-[a-f0-9]{4}){3}-[a-f0-9]{12}", //uuid
				CtxOK:   true,
				CtxText: "[a-f0-9]{8}(-[a-f0-9]{4}){3}-[a-f0-9]{12}", //uuid
			},
		},
		{
			name: "custom generator",
			args: args{Config{SetReqHeader: true, SetGinCtx: true, Generate: testGeneator}, defaultHeader},
			want: idResponse{
				Header:  "test", //uuid
				CtxOK:   true,
				CtxText: "test", //uuid
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := newTestRouter(tt.args.cfg, tt.args.expectedKey)
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", "/", bytes.NewReader([]byte{}))
			router.ServeHTTP(w, req)

			//test body
			body := w.Body.Bytes()
			var got idResponse
			err := json.Unmarshal(body, &got)
			if err != nil {
				t.Errorf("unmarshal %v", err)
				return
			}

			got.Matches(t, tt.want)
		})
	}
}

func benchmarkRequestID(b *testing.B, cfg Config) {
	router := newTestRouter(cfg, cfg.Header)
	w := httptest.NewRecorder()
	for i := 0; i < b.N; i++ {
		req, _ := http.NewRequest("GET", "/", bytes.NewReader([]byte{}))
		router.ServeHTTP(w, req)
		w.Flush()
	}
}

func BenchmarkRequestIDWithGinCtx(b *testing.B) {
	cfg := DefaultCfg()
	cfg.SetReqHeader = false
	cfg.SetGinCtx = true
	benchmarkRequestID(b, cfg)
}

func BenchmarkRequestIDWithReqHeader(b *testing.B) {
	cfg := DefaultCfg()
	cfg.SetReqHeader = true
	cfg.SetGinCtx = false
	benchmarkRequestID(b, cfg)
}

func BenchmarkRequestIDWithBothOutPuts(b *testing.B) {
	cfg := DefaultCfg()
	cfg.SetReqHeader = true
	cfg.SetGinCtx = true
	benchmarkRequestID(b, cfg)
}

func BenchmarkRequestIDBothDisabled(b *testing.B) {
	cfg := DefaultCfg()
	cfg.SetReqHeader = false
	cfg.SetGinCtx = false
	benchmarkRequestID(b, cfg)
}
