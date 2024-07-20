package service

import (
	"assignment3/entity"
	"context"
	"encoding/json"
	"fmt"
	"github.com/redis/go-redis/v9"
	"github.com/teris-io/shortid"
	"log"
	"time"
)

// IUserService mendefinisikan interface untuk layanan
type IShortlinkService interface {
	CreateShortlink(ctx context.Context, shortlink *entity.Shortlink) (entity.Shortlink, error)
	GetLongURL(ctx context.Context, shortLink string) (string, error)
}

// IShortlinkRepository mendefinisikan interface untuk repository
type IShortlinkRepository interface {
	CreateShortlink(ctx context.Context, shortlink *entity.Shortlink) (entity.Shortlink, error)
	GetLongURL(ctx context.Context, shortLink string) (entity.Shortlink, error)
}

type shortlinkService struct {
	shortlinkRepo IShortlinkRepository
	rdb           *redis.Client
}

// key redis
const redisByShortlinkKey = "shortlink:%v"

func NewShortlinkService(shortlinkRepo IShortlinkRepository, rdb *redis.Client) IShortlinkService {
	return &shortlinkService{
		shortlinkRepo: shortlinkRepo,
		rdb:           rdb,
	}
}

func (s *shortlinkService) CreateShortlink(ctx context.Context, param *entity.Shortlink) (entity.Shortlink, error) {
	ids, err := shortid.Generate()
	if err != nil {
		return entity.Shortlink{}, fmt.Errorf("Error generating IDs: %v", err)
	}
	shortlink := &entity.Shortlink{
		Shortlink: ids,
		Url:       param.Url,
	}
	createdShortlink, err := s.shortlinkRepo.CreateShortlink(ctx, shortlink)
	if err != nil {
		return entity.Shortlink{}, fmt.Errorf("gagal membuat shortlink: %v", err)
	}

	//save cache to redis
	redisKey := fmt.Sprintf(redisByShortlinkKey, createdShortlink.Shortlink)
	createdUserJSON, err := json.Marshal(createdShortlink)

	if err != nil {
		log.Println("gagal marshal json")
	}
	if err := s.rdb.Set(ctx, redisKey, string(createdUserJSON), 3600*time.Second).Err(); err != nil {
		log.Println("error when set redis key", redisKey)
	}

	return createdShortlink, nil
}

func (s *shortlinkService) GetLongURL(ctx context.Context, param string) (string, error) {
	var longURL string
	redisKey := fmt.Sprintf(redisByShortlinkKey, param)
	val, err := s.rdb.Get(ctx, redisKey).Result()
	if err == nil {
		var dataURL entity.Shortlink
		log.Println("data tersedia di redis")
		err = json.Unmarshal([]byte(val), &dataURL)
		if err != nil {
			log.Println("gagal marshall data di redis, coba ambil ke database")
		}
		longURL = dataURL.Url
	}

	if err != nil {
		log.Println("data tidak tersedia di redis, ambil dari database")
		dataGetURL, err := s.shortlinkRepo.GetLongURL(ctx, param)
		if err != nil {
			log.Println("error get long url to db")
			return "", fmt.Errorf("gagal get longURL: %v", err)
		}

		//save cache to redis
		createdDataJSON, err := json.Marshal(dataGetURL)
		if err != nil {
			log.Println("gagal marshal json")
		}

		if err := s.rdb.Set(ctx, redisKey, createdDataJSON, 3600*time.Second).Err(); err != nil {
			log.Println("error when set redis key", redisKey)
		}
		longURL = dataGetURL.Url
	}

	return longURL, nil
}
