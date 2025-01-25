package main

import (
	"common"
	apperrors "common/errors"
	"context"
	"encoding/json"
	"fmt"
	"github.com/apex/log"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"net/http"
	"order/application"
	"order/application/usecase"
	"order/domain"
	"order/infrastructure"
	"os"
	"time"
)

var orderApplication *application.OrderApplication

func init() {

	// Load environment variables for MongoDB configuration
	mongoUrl := os.Getenv("MONGO_URL")
	mongoDatabaseName := os.Getenv("MONGO_DB_NAME")
	mongoDatabaseUsername := os.Getenv("DB_USERNAME")
	mongoDatabasePassword := os.Getenv("DB_PASSWORD")
	//secretARN := os.Getenv("MONGO_SECRET_ARN")

	fullMongoURI := fmt.Sprintf("mongodb://%s:%s@%s/sample-database?tls=true&replicaSet=rs0&readpreference=secondaryPreferred", mongoDatabaseUsername, mongoDatabasePassword, mongoUrl)

	log.Infof("Connecto to mongo: %s", fullMongoURI)
	mongoConnectionTimeout := common.GetEnvDuration("MONGO_CONNECTION_TIMEOUT", 10*time.Second)

	// Set MongoDB connection timeout
	mongoCtx, cancel := context.WithTimeout(context.Background(), mongoConnectionTimeout)
	defer cancel()

	// Connect to MongoDB
	mongoClient, err := mongo.Connect(mongoCtx, options.Client().ApplyURI(fullMongoURI).SetTLSConfig(nil).SetRetryWrites(false))
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	//Ping MongoDB to check if the connection was successful
	if err := mongoClient.Ping(mongoCtx, nil); err != nil {
		log.Fatalf("Failed to ping MongoDB: %v", err)
	}

	// Log successful connection
	log.Infof("Successfully connected to MongoDB at %s", mongoUrl)

	orderRepository := infrastructure.NewOrderRepository(mongoClient, mongoDatabaseName)
	getOrderQueryHandler := usecase.NewGetOrderQueryHandler(orderRepository)
	getOrderAllOrdersQueryHandler := usecase.NewGetAllOrdersQueryHandler(orderRepository)
	createOrderCommandHandler := usecase.NewCreateOrderCommandHandler(orderRepository)

	orderApplication = application.NewOrderApplication(
		getOrderQueryHandler,
		getOrderAllOrdersQueryHandler,
		createOrderCommandHandler,
	)
}

type MongoCredentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func CreateOrderHandler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	fmt.Println(fmt.Sprintf("CreateOrderHandler: %s", request.HTTPMethod))
	orderId, isSpecificOrder := request.PathParameters["orderId"]

	switch request.HTTPMethod {
	case "POST":
		// Handle creating a new order
		return createOrder(ctx, request)
	case "GET":
		if isSpecificOrder {
			// Return specific order by ID
			return getOrder(ctx, orderId, request)
		} else {
			// Return all orders
			return getAllOrders(ctx, request)
		}
	default:
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusMethodNotAllowed,
			Body:       "Unsupported HTTP method",
		}, nil
	}
}

// Retrieve an order (GET /orders/{orderID})
func getOrder(ctx context.Context, orderID string, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	orderResult, err := orderApplication.GetOrderQueryHandler.Execute(ctx, usecase.GetOrderQuery{Id: orderID})

	if err != nil {
		return common.SerializeError(err)
	}
	return common.SerializeResponse(http.StatusOK, orderResult)
}

// Retrieve an order (GET /orders/{orderID})
func getAllOrders(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	pageFilter := common.ParsePageFilter(request.QueryStringParameters)

	result, err := orderApplication.GetAllOrdersQueryHandler.Execute(ctx, usecase.GetAllOrdersQuery{
		Filter: &domain.OrderFilter{},
		Page:   pageFilter,
	})
	if err != nil {
		log.WithError(err).Warn("Request failed")
		return common.SerializeError(err)
	}
	return common.SerializeResponse(http.StatusOK, result)
}

func createOrder(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var createOrderCommand usecase.CreateOrderCommand
	err := json.Unmarshal([]byte(request.Body), &createOrderCommand)
	if err != nil {
		return common.SerializeError(apperrors.InvalidRequest("Failed to parse request", err))
	}

	orderResult, err := orderApplication.CreateOrderCommandHandler.Execute(ctx, createOrderCommand)

	if err != nil {
		log.WithError(err).Warn("Request failed")
		return common.SerializeError(err)
	}

	return common.SerializeResponse(http.StatusCreated, orderResult)
}

func main() {
	lambda.Start(CreateOrderHandler)
}
