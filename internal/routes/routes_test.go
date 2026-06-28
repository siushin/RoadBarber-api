package routes

import (
	"io"
	"net/http/httptest"
	"strings"
	"testing"

	"roadbarber/api/internal/config"
	"roadbarber/api/pkg/utils"

	"github.com/gofiber/fiber/v2"
)

func TestHealthEndpoint(t *testing.T) {
	app := fiber.New(fiber.Config{DisableStartupMessage: true})
	cfg := &config.Config{ServerPort: "8080"}
	Setup(app, cfg, &utils.ConsoleProvider{})

	req := httptest.NewRequest("GET", "/health", nil)
	resp, err := app.Test(req, -1)
	if err != nil {
		t.Fatalf("health request failed: %v", err)
	}
	if resp.StatusCode != 200 {
		t.Fatalf("expected status 200, got %d", resp.StatusCode)
	}
	body, _ := io.ReadAll(resp.Body)
	got := string(body)
	if !strings.Contains(got, `"status":"ok"`) {
		t.Errorf("expected status:ok in body, got %s", got)
	}
}
