package memory

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

type Redis struct {
	redisClient *redis.Client
}

func New(redisClient *redis.Client) *Redis {
	return &Redis{
		redisClient: redisClient,
	}
}

var ctx = context.Background()

func (r *Redis) Test(key string, value interface{}, exp time.Duration) (*string, error) {
	data, err := r.redisClient.Get(ctx, key).Result()
	if err == redis.Nil {
		fmt.Printf("GET %s does not exist. \n", key)
		err := r.redisClient.Set(ctx, key, value, exp).Err()
		if err != nil {
			return nil, err
		}
	} else if err != nil {
		return nil, err
	} else {
		fmt.Printf("GET Cache %s:%v \n", key, data)
		return &data, nil
	}
	return nil, nil
}

func (r *Redis) GetKey(key string) (*string, error) {
	data, err := r.redisClient.Get(ctx, key).Result()
	if err == redis.Nil {
		//fmt.Println("key2 does not exist")
		fmt.Printf("%s does not exist. \n", key)

	} else if err != nil {
		return nil, err
	} else {
		//fmt.Println("key2", data)
		fmt.Printf("%s:%v \n", key, data)
		return &data, nil
	}
	return nil, nil
}

func (r *Redis) SetKey(key string, value interface{}, exp time.Duration) error {

	err := r.redisClient.Set(ctx, key, value, exp).Err()
	if err != nil {
		return err
	}

	return nil
}

func (r *Redis) Delete(key string) error {
	err := r.redisClient.Del(ctx, key).Err()
	if err != nil {
		return err
	}
	return nil
}
