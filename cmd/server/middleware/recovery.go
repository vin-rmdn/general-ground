package middleware

import (
	"fmt"
	"sync"

	"github.com/labstack/echo/v4"
)

func Recovery(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		var wg sync.WaitGroup
		var err error

		catchPanic := func() {
			defer wg.Done()

			if panicValue := recover(); panicValue != nil {
				if errPanic, ok := panicValue.(error); ok {
					err = fmt.Errorf("recovered from panic: %w", errPanic)

					return
				}

				err = fmt.Errorf("recovered from panic: %v", panicValue)
			}
		}

		safeguardHandlerFromPanic := func() {
			defer catchPanic()

			err = next(c)
		}

		wg.Add(1)
		go safeguardHandlerFromPanic()

		wg.Wait()

		return err
	}
}
