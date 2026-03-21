package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	echoadapter "github.com/awslabs/aws-lambda-go-api-proxy/echo"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

var echoLambda *echoadapter.EchoLambda

type TelegramUpdate struct {
	Message *TelegramMessage `json:"message"`
}

type TelegramMessage struct {
	Chat TelegramChat `json:"chat"`
	Text string       `json:"text"`
}

type TelegramChat struct {
	ID int64 `json:"id"`
}

func init() {
	e := newEcho()
	echoLambda = echoadapter.New(e)
}

func newEcho() *echo.Echo {
	e := echo.New()
	e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogStatus: true,
		LogURI:    true,
		LogMethod: true,
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			fmt.Printf("method=%s uri=%s status=%d\n", v.Method, v.URI, v.Status)
			return nil
		},
	}))
	e.Use(middleware.Recover())
	e.POST("/webhook", handleWebhook)
	return e
}

func handleWebhook(c echo.Context) error {
	body, err := io.ReadAll(c.Request().Body)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "failed to read body"})
	}
	fmt.Printf("Incoming payload: %s\n", string(body))

	var update TelegramUpdate
	if err := json.Unmarshal(body, &update); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request"})
	}

	if update.Message == nil {
		return c.JSON(http.StatusOK, map[string]string{"status": "ignored"})
	}

	chatID := update.Message.Chat.ID
	fmt.Printf("Received message from chat %d: %s\n", chatID, update.Message.Text)

	if err := sendMessage(chatID, "Hello World!"); err != nil {
		fmt.Printf("Error sending message: %v\n", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to send message"})
	}

	return c.JSON(http.StatusOK, map[string]string{"status": "ok"})
}

func sendMessage(chatID int64, text string) error {
	token := os.Getenv("TELEGRAM_BOT_TOKEN")
	if token == "" {
		return fmt.Errorf("TELEGRAM_BOT_TOKEN is not set")
	}

	apiURL := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", token)
	resp, err := http.PostForm(apiURL, url.Values{
		"chat_id": {fmt.Sprintf("%d", chatID)},
		"text":    {text},
	})
	if err != nil {
		return fmt.Errorf("error calling Telegram API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var body map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&body)
		return fmt.Errorf("Telegram API error: %v", body)
	}

	return nil
}

func Handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return echoLambda.ProxyWithContext(ctx, request)
}

func main() {
	runLocally := os.Getenv("RUN_LOCALLY")

	if runLocally != "" {
		e := newEcho()
		e.Logger.Fatal(e.Start(":8080"))
		return
	}
	lambda.Start(Handler)
}
