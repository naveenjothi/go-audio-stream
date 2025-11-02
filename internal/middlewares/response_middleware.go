package middlewares

import (
	"bytes"
	"encoding/json"

	"github.com/labstack/echo/v4"
)

type ResponseStruct struct {
	Data interface{} `json:"data"`
}

type CustomResponseWriter struct {
	echo.Response
	body *bytes.Buffer
}

func (w *CustomResponseWriter) Write(b []byte) (int, error) {
	return w.body.Write(b)
}

func CustomResponseMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// Create a custom response writer to capture the response
		originalWriter := c.Response().Writer
		buffer := &bytes.Buffer{}
		customWriter := &CustomResponseWriter{
			Response: *c.Response(),
			body:     buffer,
		}
		c.Response().Writer = customWriter

		// Execute the next handler
		err := next(c)

		// Restore the original writer
		c.Response().Writer = originalWriter

		if err != nil {
			return err
		}

		// Only transform successful responses (2xx status codes)
		if c.Response().Status >= 200 && c.Response().Status < 300 {
			// Get the original response body
			originalBody := buffer.Bytes()

			// If there's content to wrap
			if len(originalBody) > 0 {
				var originalData interface{}

				// Try to parse the original response as JSON
				if json.Unmarshal(originalBody, &originalData) == nil {
					// Wrap the data in ResponseStruct
					wrappedResponse := ResponseStruct{
						Data: originalData,
					}

					// Send the wrapped response
					return c.JSON(c.Response().Status, wrappedResponse)
				}
			}
		}

		// For non-2xx responses or if JSON parsing fails, return the original response
		c.Response().Writer.Write(buffer.Bytes())
		return nil
	}
}
