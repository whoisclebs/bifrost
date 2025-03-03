package proxy

import (
	"bifrost/proxy/internal/config"
	"bifrost/proxy/internal/handlers"
	"bifrost/proxy/pkg/discovery"
	"bytes"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"io"
	"math/rand"
	"net/http"
	"path"
	"strings"
	"time"
)

var (
	httpClient = &http.Client{Timeout: 10 * time.Second}
	rnd        = rand.New(rand.NewSource(time.Now().UnixNano()))
)

func ReverseProxyHandler(cfg config.ProxyConfig) fiber.Handler {
	return func(c *fiber.Ctx) error {
		route, prefix := findRoute(cfg, c.Path())
		if route == nil {
			return sendError(c, http.StatusNotFound, fmt.Errorf("route not found"), "Route not found")
		}
		instances, err := discovery.GetInstances(cfg, route.ServiceId)
		if err != nil || len(instances) == 0 {
			return sendError(c, http.StatusBadGateway, err, "No instances for service: "+route.ServiceId)
		}
		instance := instances[rnd.Intn(len(instances))]
		newPath := normalizePath(c.Path(), prefix)
		targetURL := fmt.Sprintf("http://%s:%d%s", instance.IPAddr, instance.Port, newPath)
		req, err := http.NewRequest(c.Method(), targetURL, bytes.NewReader(c.Body()))
		if err != nil {
			return sendError(c, http.StatusInternalServerError, err, err.Error())
		}
		copyRequestHeaders(c, req)
		resp, err := httpClient.Do(req)
		if err != nil {
			return sendError(c, http.StatusBadGateway, err, err.Error())
		}
		defer resp.Body.Close()
		copyResponseHeaders(c, resp)
		c.Status(resp.StatusCode)
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return sendError(c, http.StatusInternalServerError, err, err.Error())
		}
		return c.Send(body)
	}
}

func findRoute(cfg config.ProxyConfig, path string) (*config.SentryRouteConfig, string) {
	for _, route := range cfg.Sentry.Routes {
		prefix := strings.TrimSuffix(route.Path, "/**")
		if strings.HasPrefix(path, prefix) {
			return &route, prefix
		}
	}
	return nil, ""
}

func normalizePath(fullPath, prefix string) string {
	newPath := strings.TrimPrefix(fullPath, prefix)
	newPath = "/" + strings.TrimLeft(newPath, "/")
	return path.Clean(newPath)
}

func copyRequestHeaders(c *fiber.Ctx, req *http.Request) {
	c.Request().Header.VisitAll(func(key, value []byte) {
		req.Header.Set(string(key), string(value))
	})
}

func copyResponseHeaders(c *fiber.Ctx, resp *http.Response) {
	for key, values := range resp.Header {
		for _, v := range values {
			c.Set(key, v)
		}
	}
}

func sendError(c *fiber.Ctx, status int, err error, message string) error {
	defErr := handlers.NewDefaultErrorMessage(c.Path(), status, message, err)
	return c.Status(status).JSON(defErr)
}
