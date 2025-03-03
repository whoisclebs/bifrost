package entrypoint

import (
	"bifrost/proxy/internal/config"
	"bifrost/proxy/internal/handlers"
	"bifrost/proxy/internal/interceptors"
	"bifrost/proxy/internal/proxy"
	"encoding/json"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"log"
)

func StartServer() {
	config.Init()
	cfg := config.Get()
	StartServerConfig(cfg)
}

func StartServerConfig(cfg config.ProxyConfig) {
	app := fiber.New(fiber.Config{
		JSONEncoder:  json.Marshal,
		JSONDecoder:  json.Unmarshal,
		ErrorHandler: handlers.DefaultErrorHandler,
	})
	icpts := interceptors.CreateInterceptors()
	for _, handler := range icpts {
		app.Use(handler)
	}
	app.All("/*", proxy.ReverseProxyHandler(cfg))
	addr := fmt.Sprintf(":%d", cfg.App.Port)
	if err := app.Listen(addr); err != nil {
		log.Fatal(err)
	}
}
