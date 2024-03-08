package redis

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/yanlihongaichila/framework/nacos"
	"gopkg.in/yaml.v2"
	"time"
)

type RedisConfig struct {
	Redis struct {
		Host string `yaml:"host"`
		Port string `yaml:"port"`
		Db   int    `yaml:"db"`
	} `yaml:"redis"`
}

func getRedisConfig(group, server string) (RedisConfig, error) {
	config, err := nacos.GetConfig(group, server)
	if err != nil {
		return RedisConfig{}, err
	}
	redisCon := RedisConfig{}

	err = yaml.Unmarshal([]byte(config), &redisCon)
	if err != nil {
		return RedisConfig{}, err
	}

	return redisCon, nil
}

// redis连接
func withRedis(group, server string, hand func(cli *redis.Client) error) error {
	//获取Redis配置,连接Redis
	redisConfig, err := getRedisConfig(group, server)
	rCfg := redisConfig.Redis
	if err != nil {
		return err
	}

	redisClient := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%v:%v", rCfg.Host, rCfg.Port),
		DB:   rCfg.Db,
	})

	//使用完之后需要关闭
	defer redisClient.Close()

	err = hand(redisClient)
	if err != nil {
		return err
	}

	return nil
}

// redis的简单操作
func GetRedisInfo(ctx context.Context, group, server, key string) (string, error) {
	var result string
	var err error
	err = withRedis(group, server, func(cli *redis.Client) error {
		result, err = cli.Get(ctx, key).Result()
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return "", err
	}

	return result, nil
}

func SetRedisInfo(ctx context.Context, group, server, key string, value any, time time.Duration) bool {
	var err error
	err = withRedis(group, server, func(cli *redis.Client) error {
		err = cli.Set(ctx, key, value, time).Err()
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return false
	}

	return true
}

func ExistInfo(ctx context.Context, group, server, key string) bool {
	var err error
	err = withRedis(group, server, func(cli *redis.Client) error {
		err = cli.Exists(ctx, key).Err()
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return false
	}
	return true
}
