package main

import (
	"context"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
)

var ctx = context.Background()

func getEnv(key, def string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		value = def
	}
	return value
}

func getEnvInt(key string, def int) int {
	var intvalue int = def
	value, exists := os.LookupEnv(key)
	if exists {
		intvalue, err := strconv.Atoi(value)
		if err != nil {
			intvalue = def
		}
		return intvalue
	}
	return intvalue
}

func getSquareKey(c *gin.Context) {
	key := c.Param("key")

	rdb := redis.NewClient(&redis.Options{
		Addr:     getEnv("REDIS_ADDR", "10.0.1.101:6415"),
		Password: getEnv("REDIS_PASSWD", ""), // no password set
		DB:       getEnvInt("REDIS_DB", 0),   // use default DB
	})

	val, err := rdb.Get(key).Result()
	if err == redis.Nil {
		c.IndentedJSON(http.StatusNotFound, key+" is not available")
		return
	} else if err != nil {
		// panic(err)
		c.IndentedJSON(http.StatusInternalServerError, "redis connection error")
		return
	} else {
		var intval int = 0
		intval, err := strconv.Atoi(val)
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, key+" value is not integer")
			return
		}
		newkey := key + "2"
		err2 := rdb.Set(newkey, strconv.Itoa(intval*intval), 0).Err()
		if err2 != nil {
			c.IndentedJSON(http.StatusNotModified, "redis was not set")
			return
		}
		c.IndentedJSON(http.StatusOK, newkey)
		return
	}

}

func main() {
	router := gin.Default()
	router.GET("/go/rest/square/:key", getSquareKey)
	router.Run("localhost:9301")
}

// $resp=Invoke-WebRequest -Uri "http://127.0.0.1:9301/go/rest/square/A" -Method 'GET'
// $resp.Content
