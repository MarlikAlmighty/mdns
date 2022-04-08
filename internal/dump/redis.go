package dump

import (
	"context"
	"encoding/json"
	"github.com/MarlikAlmighty/mdns/internal/data"
	"github.com/MarlikAlmighty/mdns/internal/gen/models"
	"github.com/go-redis/redis/v8"
)

// Store client redis
type Store struct {
	Client   *redis.Client
	Resolver data.Resolver
}

// New simple constructor
func New(redisUrl string, data data.Resolver) (*Store, error) {
	opt, err := redis.ParseURL(redisUrl)
	if err != nil {
		return nil, err
	}
	r := redis.NewClient(&redis.Options{
		Addr:     opt.Addr,
		Password: opt.Password,
		DB:       opt.DB,
	})
	if _, err = r.Ping(context.Background()).Result(); err != nil {
		return nil, err
	}
	return &Store{
		Client:   r,
		Resolver: data,
	}, nil
}

func (s *Store) InitMaps(key string) error {
	res, err := s.Pop(key)
	if err != nil {
		return err
	}
	if len(res) > 0 {
		mp := make(map[string]models.DNSEntry)
		if err = json.Unmarshal([]byte(res), &mp); err != nil {
			return err
		}
		s.Resolver.SetMap(mp)
	}
	return nil
}

// Pop return value of key
func (s *Store) Pop(value string) (string, error) {
	res, err := s.Client.Get(context.Background(), value).Result()
	if err == redis.Nil {
		return "", nil
	} else if err != nil {
		return "", err
	}
	return res, nil
}

// Push insert dump into database
func (s *Store) Push(key string, value []byte) error {
	if err := s.Client.Set(context.Background(), key, value, 0).Err(); err != nil {
		return err
	}
	return nil
}

// Shutdown Saving the dump and closing the connection
func (s *Store) Shutdown(redisKey string) error {
	mp := s.Resolver.GetMap()
	b, err := json.Marshal(mp)
	if err != nil {
		return err
	}
	if err = s.Push(redisKey, b); err != nil {
		return err
	}
	if err = s.Client.Close(); err != nil {
		return err
	}
	return nil
}
