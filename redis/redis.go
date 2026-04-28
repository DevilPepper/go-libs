package redis

import (
	"context"

	"github.com/gomodule/redigo/redis"
)

type Redis struct {
	Pool   *redis.Pool
	ctx    context.Context
	prefix string
}

func (r *Redis) Get(key string) (interface{}, error) {
	conn, err := r.Pool.GetContext(r.ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	return conn.Do("GET", r.prefix+key)
}

func (r *Redis) GetBytes(key string) ([]byte, error) {
	b, e := r.Get(key)
	return redis.Bytes(b, e)
}

func (r *Redis) Set(key string, value interface{}) error {
	conn, err := r.Pool.GetContext(r.ctx)
	if err != nil {
		return err
	}
	defer conn.Close()

	_, err = conn.Do("SET", r.prefix+key, value)
	return err
}

func (r *Redis) SetWithExpireTimestamp(key string, value interface{}, expire int64) error {
	conn, err := r.Pool.GetContext(r.ctx)
	if err != nil {
		return err
	}
	defer conn.Close()

	_, err = conn.Do("SET", r.prefix+key, value, "EXAT", expire)
	return err
}

func (r *Redis) SetWithDuration(key string, value interface{}, duration int64) error {
	conn, err := r.Pool.GetContext(r.ctx)
	if err != nil {
		return err
	}
	defer conn.Close()

	_, err = conn.Do("SET", r.prefix+key, value, "EX", duration)
	return err
}

func (r *Redis) Delete(key string) error {
	conn, err := r.Pool.GetContext(r.ctx)
	if err != nil {
		return err
	}
	defer conn.Close()

	_, err = conn.Do("DEL", r.prefix+key)
	return err
}

func GetRedis(ctx context.Context, prefix string) *Redis {
	return &Redis{
		Pool:   GetRedisPool(),
		ctx:    ctx,
		prefix: prefix,
	}
}
