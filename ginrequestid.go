package ginrequestid

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

var defaultHeader = "X-Request-ID"
var defaultGenerator = func(_ *gin.Context) string {
	return uuid.NewString()
}

type RequestIDGenerator func(c *gin.Context) string

type Config struct {
	Header       string             //example: X-Request-Header
	Generate     RequestIDGenerator //example: uuid.NewString
	SetGinCtx    bool               //make the Header in gin.Context available
	SetReqHeader bool               //update the gin.Context Request.Header if the header is not set
}

//DefaultCfg will set X-Request-ID with uuidV4, which will be applied gin.Context
func DefaultCfg() Config {
	return Config{
		Header:       defaultHeader,
		Generate:     defaultGenerator,
		SetGinCtx:    true,
		SetReqHeader: false,
	}
}

//RequestID adds the configured request header if it does not exist.
//If the cfg Config values are empty, the DefaultCfg config values will be applied.
func RequestID(cfg Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		//early exit
		if !cfg.SetReqHeader && !cfg.SetGinCtx {
			return
		}
		//check for existing request header
		xRequestID := c.Request.Header.Get(cfg.Header)
		if xRequestID == "" {
			xRequestID = cfg.GenerateID(c)
			updateHeader(c, &cfg, xRequestID)
		}
		updateGinCtx(c, &cfg, xRequestID)

		c.Next()
	}
}

func updateHeader(c *gin.Context, cfg *Config, xRequestID string) {
	if cfg.SetReqHeader {
		c.Request.Header.Set(cfg.Header, xRequestID)
	}
}

func updateGinCtx(c *gin.Context, cfg *Config, xRequestID string) {
	if cfg.SetGinCtx {
		c.Set(cfg.Header, xRequestID)
	}
}

func (cfg *Config) validate() *Config {
	if cfg.Header == "" {
		cfg.Header = defaultHeader
	}
	if cfg.Generate == nil {
		cfg.Generate = defaultGenerator
	}
	return cfg
}

func (cfg *Config) GenerateID(c *gin.Context) string {
	return cfg.validate().Generate(c)
}
