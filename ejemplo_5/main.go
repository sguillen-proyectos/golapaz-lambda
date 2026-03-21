package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	echoadapter "github.com/awslabs/aws-lambda-go-api-proxy/echo"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

var echoLambda *echoadapter.EchoLambda

func init() {
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

	e.GET("/hello", handleGet)
	e.POST("/hello", handlePost)

	echoLambda = echoadapter.New(e)
}

func handleGet(c echo.Context) error {
	name := c.QueryParam("name")
	if name == "" {
		name = "World"
	}
	return c.JSON(http.StatusOK, map[string]string{
		"message": "Hello, echo " + name + "!",
	})
}

func handlePost(c echo.Context) error {
	var payload struct {
		Name string `json:"name"`
	}
	if err := c.Bind(&payload); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "invalid JSON body",
		})
	}
	if payload.Name == "" {
		payload.Name = "World"
	}
	return c.JSON(http.StatusOK, map[string]string{
		"message": "Hello, " + payload.Name + "! (POST) echo",
	})
}

func Handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return echoLambda.ProxyWithContext(ctx, request)
}

func main() {
	runLocally := os.Getenv("RUN_LOCALLY")

	if runLocally != "" {
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
		e.GET("/hello", handleGet)
		e.POST("/hello", handlePost)
		e.Logger.Fatal(e.Start(":8080"))
		return
	}
	lambda.Start(Handler)
}
