// Package rediscli implements a Redis client.
package rediscli

import (
	"time"

	"github.com/go-redis/cache"
	"github.com/go-redis/redis"
	"github.com/vmihailenco/msgpack"
	// "gitlab2.trumptech.com/wct-global/backend/pkg/qeutil"
)

// Client defines a redis database connection and its serialize codec.
type Client struct {
	client *redis.Client // Client defines how to connect the cache sever.
	codec  *cache.Codec  // Codec defines how the model serialize in cache.
}

// NewRedisCli initialises a new Client with msgpack codec.
func NewRedisCli(opt *redis.Options) *Client {
	var cli Client
	cli.client = redis.NewClient(opt)
	cli.codec = &cache.Codec{
		Redis:     cli.client,
		Marshal:   func(v interface{}) ([]byte, error) { return msgpack.Marshal(v) },
		Unmarshal: func(b []byte, v interface{}) error { return msgpack.Unmarshal(b, v) },
	}
	return &cli
}

// Set object into redis client.
func (cli *Client) Set(key string, obj interface{}, exp time.Duration) error {
	return cli.codec.Set(&cache.Item{
		Key:        key,
		Object:     obj,
		Expiration: exp,
	})
}

// Get object into redis client.
func (cli *Client) Get(key string, obj interface{}) error {
	return cli.codec.Get(key, obj)
}

// Close redis connection.
func (cli *Client) Close(key string, obj interface{}) error {
	return cli.client.Close()
}

// Do the command in redis connection.
func (cli *Client) Do(args ...interface{}) error {
	return cli.client.Process(redis.NewCmd(args...))
}

// UnlinkKeys remove all keys in redis that matching the given conditions.
func (cli *Client) UnlinkKeys(keys []string) error {
	for i := range keys {
		// Scan and unlink key
		iter := cli.client.Scan(0, keys[i], 0).Iterator()
		for iter.Next() {
			err := cli.client.Unlink(iter.Val()).Err()
			if err != nil {
				return err
			}
		}
		if err := iter.Err(); err != nil {
			return err
		}
	}

	return nil
}

// Client return the redis connection client.
func (cli *Client) Client() *redis.Client {
	return cli.client
}
