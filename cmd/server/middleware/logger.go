package middleware

import (
	"bytes"
	"fmt"
	"io"
	"log/slog"
	"net/http"
)

func Logger(next HandlerWithError) HandlerWithError {
	return func(w http.ResponseWriter, r *http.Request) error {
		body, newReadCloser, err := cloneBody(r.Body)
		if err != nil {
			return fmt.Errorf("can not clone body: %w", err)
		}
		
		slogHTTPMetrics := []any{
			slog.String("body", string(body)),
			slog.Any("header", r.Header),
			slog.String("host", r.Host),
			slog.String("method", r.Method),
			slog.String("uri", r.RequestURI),
			slog.Any("query", r.URL.Query()),
		}

		slog.Debug("received http request", slogHTTPMetrics...)

		r.Body = newReadCloser

		if err := next(w, r); err != nil {
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
