package redis

import "github.com/Edu4rdoNeves/ingestor-magalu/utils"

func (c *RedisClient) IncrementCounter(key string, value float64) error {
	convertValue := utils.FloatToInt64(value)
	return c.Client.IncrBy(key, convertValue).Err()
}

func (c *RedisClient) GetValue(key string) (string, error) {
	val, err := c.Client.Get(key).Result()
	if err != nil {
		return "", err
	}
	return val, nil
}

func (c *RedisClient) DeleteKey(key string) error {
	return c.Client.Del(key).Err()
}

func (c *RedisClient) GetKeysByPattern(pattern string) ([]string, error) {
	return c.Client.Keys(pattern).Result()
}
