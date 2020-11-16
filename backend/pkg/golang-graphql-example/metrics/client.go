package metrics

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Client Client metrics interface
//go:generate mockgen -destination=./mocks/mock_Client.go -package=mocks github.com/oxyno-zeta/golang-graphql-example/pkg/golang-graphql-example/metrics Client
type Client interface {
	// Instrument web server
	Instrument(serverName string) gin.HandlerFunc
	// Get prometheus handler for http expose
	GetPrometheusHTTPHandler() http.Handler
	// Get GORM middleware
	GormMiddleware(labels map[string]string, refresh int) gorm.Plugin
}

// NewMetricsClient will generate a new Client.
func NewMetricsClient() Client {
	ctx := &prometheusMetrics{}
	ctx.register()

	return ctx
}
