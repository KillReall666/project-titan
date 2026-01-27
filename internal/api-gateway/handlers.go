package api_gateway

import (
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"net/http"
)

func SetupRouters(r *gin.Engine) {
	r.POST("/packet", HandlePacket)
}

func HandlePacket(c *gin.Context) {
	ctx, span := trace.Tracer("gateway").Start(c.Request.Context(), "handle_packet")
	defer span.End()
	span.SetAttributes(attribute.String("ip", c.ClientIP()))

	common.Requests.Inc
	c.Status(http.StatusOK)
}
