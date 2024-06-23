package repositories

import (
	"context"
	"os"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/redis/go-redis/v9"
	"ingenhouzs.com/chesshouzs/go-game/interfaces"
	"ingenhouzs.com/chesshouzs/go-game/models"
)

type Repository struct {
	postgres *gorm.DB
	redis    *redis.Client
}

func ConnectPostgreSQL(psql models.SqlConnection) (*gorm.DB, error) {
	var err error

	connectionString := "host=" + psql.Host + " user=" + psql.User + " password=" + psql.Password + " dbname=" + psql.Database + " port=" + psql.Port + " sslmode=disable TimeZone=UTC+7"
	db, err := gorm.Open(psql.Driver, connectionString)
	if err != nil {
		return db, err
	}

	db.SingularTable(true)
	if os.Getenv("LOG_GORM") == "ON" {
		db.LogMode(true)
	}

	return db, nil
}

func ConnectRedis(r models.RedisConnection) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     r.Host + ":" + r.Port,
		Password: r.Password,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := client.Ping(ctx).Result()
	if err != nil {
		return nil, err
	}

	return client, nil
}

func NewRepository(postgres *gorm.DB, redis *redis.Client) interfaces.Repository {
	return &Repository{postgres, redis}
}
