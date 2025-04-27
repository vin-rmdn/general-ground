package middleware

import (
	"bytes"
	"fmt"
	"io"
	"log/slog"

	"github.com/labstack/echo/v4"
)

func Logger(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		request := c.Request()
		
		body, newReadCloser, err := cloneBody(request.Body)
		if err != nil {
			return fmt.Errorf("can not clone body: %w", err)
		}

		slogHTTPMetrics := []any{
			slog.String("body", string(body)),
			slog.Any("header", request.Header),
			slog.String("host", request.Host),
			slog.String("method", request.Method),
			slog.String("uri", request.RequestURI),
			slog.Any("query", request.URL.Query()),
		}

		slog.Debug("received http request", slogHTTPMetrics...)

		request.Body = newReadCloser

		if err := next(c); err != nil {
			slogHTTPMetrics = append(slogHTTPMetrics, slog.String("error", err.Error()))
			slog.Error("handler returned error", slogHTTPMetrics...)

			return err
		}

		slog.Debug("http request is successful", slogHTTPMetrics...)

		return nil
	}
}

func cloneBody(requestBody io.ReadCloser) (string, io.ReadCloser, error) {
	defer func() { _ = requestBody.Close() }()

	body, err := io.ReadAll(requestBody)
	if err != nil {
		return "", nil, fmt.Errorf("can not read request body: %w", err)
	}

	if err = requestBody.Close(); err != nil {
		return "", nil, fmt.Errorf("can not close request body: %w", err)
	}

	newReadCloser := io.NopCloser(bytes.NewReader(body))

	return string(body), newReadCloser, nil
}
