package controllers

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/Aanu1995/url-shortner/api/database"
	"github.com/Aanu1995/url-shortner/api/helpers"
	"github.com/Aanu1995/url-shortner/api/models"
	"github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v9"
	"github.com/google/uuid"
)



func ResolveURL(ctx *gin.Context){
	url := ctx.Param("url")

	if url == ""{
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid url"})
		return
	}

	defaultClient := database.CreateClient(0)
	defer defaultClient.Close()


	result, err := defaultClient.Get(context.Background(), url).Result()
	if err == redis.Nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "link does not exist"})
		return
	} else if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.Redirect(http.StatusPermanentRedirect, result)
}



func ShortenURL(ctx *gin.Context){
	var requestBody models.Request

	if err := ctx.BindJSON(&requestBody); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	quota, err := strconv.Atoi(os.Getenv("APP_QUOTA"))
	if err != nil{
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	client, limit, err := setQuota(ctx.RemoteIP(), quota)
	if err != nil{
		if limit != 0 {
			ctx.JSON(http.StatusForbidden, gin.H{"error": err.Error(), "rate_limit": limit})
		} else {
			ctx.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		}
		return
	}
	defer client.Close()

	// checks if input is an actual url
	if !govalidator.IsURL(requestBody.URL){
		ctx.JSON(http.StatusServiceUnavailable, gin.H{"error": "invalid URL"})
		return
	}

	// check domain error
	if !helpers.RemoveDomainError(requestBody.URL) {
		ctx.JSON(http.StatusServiceUnavailable, gin.H{"error": "Oops! you can't hack our server"})
		return
	}

	// enforce https, ssl
	requestBody.URL = helpers.EnforceHttp(requestBody.URL)
	uid := requestBody.CustomShort

	if uid == ""{
		uid = uuid.NewString()[:6]
	}

	defaultClient := database.CreateClient(0)
	defer defaultClient.Close()

	result, _ := defaultClient.Get(context.Background(), uid).Result()
	if result != "" {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "URL custom short already in use"})
		return
	}

	if requestBody.Expiry == 0 {
		requestBody.Expiry = 24
	}
	if err := defaultClient.Set(context.Background(), uid, requestBody.URL, requestBody.Expiry * 3600 * time.Second).Err(); err != nil{
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := models.Response{
		URL: requestBody.URL,
		CustomShort: "",
		Expiry: requestBody.Expiry,
		XRateRemaining: quota,
		XRateLimitReset: 30 * 60, // 30 minutes
	}

	client.Decr(context.Background(), ctx.RemoteIP())
	// get the rate remaining
	rateRemaining, _ := client.Get(context.Background(), ctx.RemoteIP()).Result()
	response.XRateRemaining, _ = strconv.Atoi(rateRemaining)

	// get rate limit reset
	ttl, _ := client.TTL(context.Background(), ctx.RemoteIP()).Result()
	response.XRateLimitReset = ttl / (time.Nanosecond * time.Minute)

	// set custom short
	response.CustomShort = fmt.Sprint(os.Getenv("DOMAIN"), "/", uid)

	ctx.JSON(http.StatusCreated, response)
}


func setQuota(ip string, quota int) (client *redis.Client, limit time.Duration, err error){
	client = database.CreateClient(1)

	result, noDataError := client.Get(context.Background(), ip).Result()
	if noDataError == redis.Nil {
		if err2 := client.Set(context.Background(), ip, quota, 30 * 60 * time.Second).Err(); err2 != nil{
			err = err2
			return
		}
		result = strconv.Itoa(quota)
	}

	appQuota, _ := strconv.Atoi(result)
	if appQuota <= 0 {
		ttl, _ := client.TTL(context.Background(), ip).Result()
		limit = ttl / (time.Nanosecond * time.Minute)
		err = errors.New("rate limit exceeded")
	}

	return
}