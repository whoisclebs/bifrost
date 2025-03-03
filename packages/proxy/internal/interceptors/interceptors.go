package interceptors

import (
	"bifrost/proxy/internal/config"
	"github.com/gofiber/fiber/v2"
)

var registry = map[string]func() fiber.Handler{}

func CreateInterceptors() []fiber.Handler {
	var list []fiber.Handler
	for _, interceptorName := range config.Get().App.Interceptors {
		if creator, ok := registry[interceptorName]; ok {
			list = append(list, creator())
		}
	}
	return list
}
