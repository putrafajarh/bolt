package controllers

import (
	"io"
	"net/http"

	"github.com/bytedance/sonic"
	"github.com/gofiber/fiber/v2"
	"github.com/putrafajarh/bolt/internal/middlewares"
	"github.com/putrafajarh/bolt/pkg/httpclient"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type PingResponse struct {
	Message string `json:"message"`
}

// @Summary      Ping
// @Description  Responds with a "pong" message to test server availability
// @Tags         Health
// @Produce      json
// @Success      200  {object} controllers.PingResponse
// @Router       /v1/ping [get]
func HandlePing(c *fiber.Ctx) error {
	logger, _ := middlewares.GetLogger(c)
	parentSpan := trace.SpanFromContext(c.UserContext())
	isSuccess := c.QueryBool("success")
	tracer := middlewares.GetTracer(c)

	ctx := c.UserContext()
	if parentSpan.IsRecording() {
		childCtx, span := tracer.Start(c.UserContext(), "HandlePing")
		ctx = childCtx
		span.SetAttributes(attribute.String("request_id", middlewares.GetRequestID(c)))
		defer span.End()
	}

	client := httpclient.NewHeimdall(httpclient.HeimdallConfig{
		Headers: http.Header{
			"User-Agent":   []string{"bolt"},
			"X-Request-ID": []string{middlewares.GetRequestID(c)},
		},
		OtelEnabled: parentSpan.IsRecording(),
	})
	response, err := client.Get(ctx, "https://jsonplaceholder.typicode.com/posts")
	if err != nil {
		return err
	}
	defer response.Body.Close()

	spanClient := trace.SpanFromContext(response.Request.Context())
	spanClient.SetAttributes(
		attribute.Bool("is_success", isSuccess),
	)

	tx, err := middlewares.GetTx(c)
	if err != nil {
		return err
	}

	if !isSuccess {
		return fiber.NewError(fiber.StatusInternalServerError, "ping failed")
	}

	if tx = tx.Exec("SELECT 1"); tx.Error != nil {
		return tx.Error
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		logger.Error().Err(err).Msg("Error reading response body")
		return err
	}

	var jsonData any
	err = sonic.Unmarshal(body, &jsonData)
	if err != nil {
		logger.Error().Err(err).Msg("Error unmarshaling JSON")
		return err
	}

	return c.Status(200).JSON(jsonData)
}
