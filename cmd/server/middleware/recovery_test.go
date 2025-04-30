package middleware_test

import (
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/vin-rmdn/general-ground/cmd/server/middleware"
)

func TestRecovery(t *testing.T) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/", http.NoBody)

	t.Run("when handler panics, recovery to return wrapped error", func(t *testing.T) {
		t.Run("when panic is not an error, recovery to return new error", func(t *testing.T) {
			e := echo.New()
			ctx := e.NewContext(r, w)

			panicHandler := func(c echo.Context) error {
				panic("unexpected panic")
			}

			wrappedHandler := middleware.Recovery(panicHandler)
			err := wrappedHandler(ctx)
			if err == nil {
				t.Fatal("expected error, got nil")
			}

			if !strings.Contains(err.Error(), "recovered from panic: unexpected panic") {
				t.Errorf("expected error to be 'recovered from panic: unexpected panic', instead got %s", err.Error())
			}
		})

		t.Run("when panic is an error, recovery to return wrapped error", func(t *testing.T) {
			e := echo.New()
			ctx := e.NewContext(r, w)

			expectedError := errors.New("expected error")
			panicHandler := func(c echo.Context) error {
				panic(expectedError)
			}

			wrappedHandler := middleware.Recovery(panicHandler)
			err := wrappedHandler(ctx)
			if err == nil {
				t.Fatal("expected error, got nil")
			}

			if !strings.Contains(err.Error(), "recovered from panic: expected error") {
				t.Errorf("expected error to be 'recovered from panic: expected error', instead got %s", err.Error())
			}

			if !errors.Is(err, expectedError) {
				t.Errorf("error '%s' does not wrap expectedError", err.Error())
			}
		})
	})

	t.Run("when handler returns no error, recovery to return nil", func(t *testing.T) {
		router := echo.New()
		router.GET("/", func(c echo.Context) error {
			return c.String(http.StatusAccepted, "no problem")
		})

		router.Use(middleware.Recovery)

		router.ServeHTTP(w, r)

		if w.Code != http.StatusAccepted {
			t.Errorf("expected status code is accepted, got %d instead", w.Code)
		}

		responseBody, err := io.ReadAll(w.Body)
		if err != nil {
			t.Fatalf("failed to read response body: %s", err.Error())
		}

		if string(responseBody) != "no problem" {
			t.Errorf("expected response string to be 'no problem', got %s instead", string(responseBody))
		}
	})
}
