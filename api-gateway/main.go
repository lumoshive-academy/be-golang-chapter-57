package main

import (
	"api-gateway/middleware"
	"context"
	"log"
	"net/http"

	pbProduct "path/to/product/proto"
	pbUser "path/to/user/proto"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
    router := gin.Default()

    // Connect to gRPC services
    productConn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
    if err != nil {
        log.Fatalf("Failed to connect to Product Service: %v", err)
    }
    defer productConn.Close()
    productClient := pbProduct.NewProductServiceClient(productConn)

    userConn, err := grpc.Dial("localhost:50052", grpc.WithTransportCredentials(insecure.NewCredentials()))
    if err != nil {
        log.Fatalf("Failed to connect to User Service: %v", err)
    }
    defer userConn.Close()
    userClient := pbUser.NewUserServiceClient(userConn)

    // User Routes - No Authentication Required
    router.POST("/user/register", func(c *gin.Context) {
        // Handle Register
        var req pbUser.RegisterRequest
        if err := c.ShouldBindJSON(&req); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
            return
        }

        res, err := userClient.Register(context.Background(), &req)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }

        c.JSON(http.StatusOK, gin.H{"message": res.Message})
    })

    router.POST("/user/login", func(c *gin.Context) {
        // Handle Login
        var req pbUser.LoginRequest
        if err := c.ShouldBindJSON(&req); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
            return
        }

        res, err := userClient.Login(context.Background(), &req)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }

        c.JSON(http.StatusOK, gin.H{"token": res.Token})
    })

    // Product Routes - Authentication Required
    productRoutes := router.Group("/product")
    productRoutes.Use(middleware.AuthMiddleware(userClient))
    {
        productRoutes.POST("/", func(c *gin.Context) {
            var req pbProduct.CreateProductRequest
            if err := c.ShouldBindJSON(&req); err != nil {
                c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
                return
            }

            res, err := productClient.CreateProduct(context.Background(), &req)
            if err != nil {
                c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
                return
            }

            c.JSON(http.StatusOK, gin.H{"product": res.Product})
        })

        productRoutes.GET("/:id", func(c *gin.Context) {
            id := c.Param("id")
            req := &pbProduct.GetProductRequest{Id: uint32(id)}

            res, err := productClient.GetProduct(context.Background(), req)
            if err != nil {
                c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
                return
            }

            c.JSON(http.StatusOK, gin.H{"product": res.Product})
        })

        productRoutes.GET("/", func(c *gin.Context) {
            req := &pbProduct.ListProductsRequest{}
            res, err := productClient.ListProducts(context.Background(), req)
            if err != nil {
                c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
                return
            }

            c.JSON(http.StatusOK, gin.H{"products": res.Products})
        })
    }

    router.Run(":8080")
}
