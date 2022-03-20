package dump

import (
	"context"
	"github.com/MarlikAlmighty/mdns/internal/config"
	"github.com/go-redis/redis/v8"
	"log"
)

// Store client redis
type Store struct {
	Client *redis.Client
}

// New simple constructor
func New(c *config.Configuration) (*Store, error) {
	opt, err := redis.ParseURL(c.RedisURl)
	if err != nil {
		return nil, err
	}
	r := redis.NewClient(&redis.Options{
		Addr:     opt.Addr,
		Password: opt.Password,
		DB:       opt.DB,
	})
	return &Store{
		Client: r,
	}, nil
}

// Pop return value of key
func (s *Store) Pop(value string) ([]byte, error) {
	res, err := s.Client.Get(context.Background(), value).Result()
	if err == redis.Nil {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	return []byte(res), nil
}

// Push insert dump into database
func (s *Store) Push(key string, value []byte) {
	if err := s.Client.Set(context.Background(), key, value, 0).Err(); err != nil {
		log.Fatalln(err)
	}
}

// Close shutdown connect
func (s *Store) Close() error {
	if err := s.Client.Close(); err != nil {
		return err
	}
	return nil
}
