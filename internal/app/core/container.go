package core

import (
	"go-modular-monolith/internal/domain/product"
	"go-modular-monolith/internal/domain/user"

	repoMongo "go-modular-monolith/internal/modules/product/repository/mongo"
	repoSQL "go-modular-monolith/internal/modules/product/repository/sql"

	serviceUnimplemented "go-modular-monolith/internal/modules/product/service/noop"
	serviceV1 "go-modular-monolith/internal/modules/product/service/v1"

	handlerUnimplemented "go-modular-monolith/internal/modules/product/handler/noop"
	handlerV1 "go-modular-monolith/internal/modules/product/handler/v1"

	handlerV1User "go-modular-monolith/internal/modules/user/handler/v1"
	repoSQLUser "go-modular-monolith/internal/modules/user/repository/sql"
	serviceV1User "go-modular-monolith/internal/modules/user/service/v1"

	"github.com/jmoiron/sqlx"
	"go.mongodb.org/mongo-driver/mongo"
)

type Container struct {
	ProductRepository product.ProductRepository
	ProductService    product.ProductService
	ProductHandler    product.ProductHandler
	UserRepository    user.UserRepository
	UserService       user.UserService
	UserHandler       user.UserHandler
}

func NewContainer(
	featureFlag FeatureFlag,
	db *sqlx.DB,
	mongo *mongo.Client,
) *Container {
	var (
		productRepository product.ProductRepository
		productService    product.ProductService
		productHandler    product.ProductHandler
		userRepository    user.UserRepository
		userService       user.UserService
		userHandler       user.UserHandler
	)

	// repo
	switch featureFlag.Repository.Product {
	case "mongo":
		productRepository = repoMongo.NewMongoRepository(mongo, "appdb")
	case "postgres":
		productRepository = repoSQL.NewSQLRepository(db)
	default:
		// productRepository = repoUnimplemented.NewUnimplementedRepository()
	}

	// service
	switch featureFlag.Service.Product {
	case "v1":
		productService = serviceV1.NewServiceV1(productRepository)
	default:
		productService = serviceUnimplemented.NewUnimplementedService()
	}

	// handler
	switch featureFlag.Handler.Product {
	case "v1":
		productHandler = handlerV1.NewHandler(productService)
	default:
		productHandler = handlerUnimplemented.NewUnimplementedHandler()
	}

	// user repo
	switch featureFlag.Repository.User {
	case "postgres":
		userRepository = repoSQLUser.NewSQLRepository(db)
	default:
	}

	// user service
	switch featureFlag.Service.User {
	case "v1":
		userService = serviceV1User.NewServiceV1(userRepository)
	default:
	}

	// user handler
	switch featureFlag.Handler.User {
	case "v1":
		userHandler = handlerV1User.NewHandler(userService)
	default:
	}

	return &Container{
		ProductService: productService,
		ProductHandler: productHandler,
		UserRepository: userRepository,
		UserService:    userService,
		UserHandler:    userHandler,
	}
}
