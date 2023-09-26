package health_check

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"git.edenfarm.id/cuxs/cuxs"
	"git.edenfarm.id/cuxs/dbredis"
	"git.edenfarm.id/cuxs/env"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Handler struct{}

// URLMapping : function to map url with it's handler and add authorization validation
func (h *Handler) URLMapping(r *echo.Group) {
	r.GET("", h.read)
}

type responseData struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Status  string `json:"status"`
}

func (h Handler) read(c echo.Context) (err error) {
	ctx := c.(*cuxs.Context)

	expPeriod := time.Duration(60) * time.Second
	var response *responseData
	if isLimit(ctx.RealIP(), expPeriod) {
		_, err = ConnectToMongoDB()
		if err != nil {
			response = &responseData{
				Code:    27017,
				Message: err.Error(),
				Status:  "fail",
			}
			ctx.ResponseData = response
			return
		}
		response = &responseData{
			Code:    200,
			Message: "OK",
			Status:  "success",
		}
		ctx.ResponseData = response
	} else {
		response = &responseData{
			Code:    429,
			Message: "Rate limit exceeded",
			Status:  "fail",
		}
		ctx.ResponseData = response
	}

	return ctx.Serve(err)

}

func ConnectToMongoDB() (*mongo.Client, error) {

	conf := loadConfig()
	clientOptions := options.Client().ApplyURI(conf.MongoDBHostPlain)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %w", err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		client.Disconnect(ctx)
		return nil, fmt.Errorf("failed to ping MongoDB: %w", err)
	}

	fmt.Println("Connected to MongoDB!")
	return client, nil
}

type config struct {
	MongoDBHostPlain string
	MongoDBHost      string
	MongoDBUsername  string
	MongoDBPassword  string
	MongoDBName      string
}

// loadConfig set config value from environment variable.
// If not exists, it will have a default values.
func loadConfig() *config {
	c := new(config)
	//mongodb://<username>:<password>@localhost:27017
	// ApplyURI("mongodb+srv://<username>:<password>@cluster0-zzart.mongodb.net/test?retryWrites=true&w=majority")
	c.MongoDBHostPlain = env.GetString("MONGO_DB_HOST", "mongodb://<username>:<password>@localhost:27017")
	c.MongoDBUsername = env.GetString("MONGO_DB_USERNAME", "root")
	c.MongoDBPassword = env.GetString("MONGO_DB_PASSWORD", "secret")
	c.MongoDBName = env.GetString("MONGO_DB_NAME", "eden_v2")
	c.MongoDBHost = c.MongoDBHostPlain
	c.MongoDBHost = strings.ReplaceAll(c.MongoDBHost, "<username>", c.MongoDBUsername)
	c.MongoDBHost = strings.ReplaceAll(c.MongoDBHost, "<password>", c.MongoDBPassword)
	return c
}

func isLimit(key string, period time.Duration) bool {
	var c int64
	var err error
	if dbredis.Redis.Client != nil {
		// Check if the key exists in Redis
		exists := dbredis.Redis.CheckExistKey(key)

		// If the key doesn't exist, set the initial count and expiration
		if !exists {
			err = dbredis.Redis.SetCache(key, 1, period)
			if err != nil {
				log.Fatal("Error setting rate limit key:", err)
			}
			return true
		}

		// Get the current count for the key
		err = dbredis.Redis.GetCache(key, &c)
		if err != nil {
			log.Fatal("Error retrieving rate limit count:", err)
		}
		// Increment the count and update the key in Redis
		if c != 10 {
			err = dbredis.Redis.SetCache(key, c+1, period)
			if err != nil {
				log.Fatal("Error updating rate limit count:", err)
			}
			return true
		}

		// Rate limit has been reached
		return false
	}

	// Fallback behavior if Redis client is not available
	return true
}
