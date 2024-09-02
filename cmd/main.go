package main

import (
	"database/sql"
	"fmt"
	"go-expert-desafio-clean-architecture/config"
	"go-expert-desafio-clean-architecture/internal/event"
	"go-expert-desafio-clean-architecture/internal/infra/database"
	"go-expert-desafio-clean-architecture/internal/infra/graphql"
	"go-expert-desafio-clean-architecture/internal/infra/web"
	"go-expert-desafio-clean-architecture/internal/infra/web/server"
	"go-expert-desafio-clean-architecture/internal/usecase/create_order"
	"go-expert-desafio-clean-architecture/internal/usecase/list_orders"
	"net"
	"net/http"

	graphqlhandler "github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/streadway/amqp"
	"go-expert-desafio-clean-architecture/internal/event/handler"
	"go-expert-desafio-clean-architecture/internal/infra/grpc/pb"
	"go-expert-desafio-clean-architecture/internal/infra/grpc/service"
	"go-expert-desafio-clean-architecture/pkg/events"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	// mysql
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	configs, err := config.LoadConfig(".")
	if err != nil {
		panic(err)
	}

	db, err := sql.Open(configs.DBDriver, fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", configs.DBUser, configs.DBPassword, configs.DBHost, configs.DBPort, configs.DBName))
	if err != nil {
		panic(err)
	}
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS orders
		(
			id          VARCHAR(36)   NOT NULL PRIMARY KEY,
			price       NUMERIC(6, 2) NOT NULL,
			tax         NUMERIC(6, 2) NOT NULL,
			final_price NUMERIC(6, 2) NOT NULL
		);
    `)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	rabbitMQChannel := getRabbitMQChannel()

	eventDispatcher := events.NewEventDispatcher()
	eventDispatcher.Register("OrderCreated", &handler.OrderCreatedHandler{RabbitMQChannel: rabbitMQChannel})

	orderRepository := database.NewOrderRepository(db)
	createOrderUseCase := create_order.New(orderRepository, event.NewOrderCreated(), eventDispatcher)
	listOrdersUseCase := list_orders.New(orderRepository)

	webserver := server.NewWebServer(configs.WebServerPort)
	webOrderHandler := web.NewWebOrderHandler(eventDispatcher, orderRepository, event.NewOrderCreated())
	webserver.AddHandler("POST", "/order", webOrderHandler.Create)
	webserver.AddHandler("GET", "/order", webOrderHandler.List)
	fmt.Println("Starting web server on port", configs.WebServerPort)
	go webserver.Start()

	grpcServer := grpc.NewServer()
	createOrderService := service.NewOrderService(createOrderUseCase, listOrdersUseCase)
	pb.RegisterOrderServiceServer(grpcServer, createOrderService)
	reflection.Register(grpcServer)

	fmt.Println("Starting gRPC server on port", configs.GRPCServerPort)
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", configs.GRPCServerPort))
	if err != nil {
		panic(err)
	}
	go grpcServer.Serve(lis)

	srv := graphqlhandler.NewDefaultServer(graphql.NewExecutableSchema(graphql.Config{
		Resolvers: &graphql.Resolver{
			CreateOrderUseCase: createOrderUseCase,
			ListOrdersUseCase:  listOrdersUseCase,
		},
	}))
	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	fmt.Println("Starting GraphQL server on port", configs.GraphQLServerPort)
	http.ListenAndServe(":"+configs.GraphQLServerPort, nil)
}

func getRabbitMQChannel() *amqp.Channel {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		panic(err)
	}
	ch, err := conn.Channel()
	if err != nil {
		panic(err)
	}
	return ch
}
