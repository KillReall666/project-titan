package handlePacket

import (
	"net/http"

	"titan/internal/observability/metrics"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

func SetupRouters(r *gin.Engine) {
	r.POST("/packet", HandlePacket)
}

func HandlePacket(c *gin.Context) {
	ctx, span := otel.Tracer("gateway").Start(c.Request.Context(), "handle_packet")
	defer span.End()

	span.SetAttributes(
		attribute.String("ip", c.ClientIP()),
	)

	// Прокидываем ctx дальше (когда появятся HTTP/Kafka/DB вызовы)
	c.Request = c.Request.WithContext(ctx)

	metric.Requests.Inc()
	c.Status(http.StatusOK)
}
