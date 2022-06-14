package redis

import (
	"encoding/json"
	"log"
	"time"

	"github.com/eric-orenge/message-server/models"
	"gopkg.in/redis.v3"
)

type RedisCtrl struct {
	Address  string
	Password string
	DB       int64
	Client   *redis.Client
}

func (rCtrl *RedisCtrl) Connect() error {
	client := redis.NewClient(&redis.Options{
		Addr:     rCtrl.Address,
		Password: rCtrl.Password,
		DB:       rCtrl.DB,
	})
	_, err := client.Ping().Result()
	if err != nil {
		return err
	}
	rCtrl.Client = client
	return nil
}

func (rCtrl *RedisCtrl) PushMessage(key, value string) error { // key - iphone@+254705207666 -
	err := rCtrl.Client.RPush(key, value).Err()
	if err != nil {
		return err
	}
	return nil
}

func (rCtrl *RedisCtrl) PopMessage(key string) []byte { // key - iphone@+254705207666 -
	val := rCtrl.Client.LPop(key)
	return []byte(val.Val())
}

func (rCtrl *RedisCtrl) SetMessage(key string, message []byte) error { // key - iphone@+254705207666 -
	cacheErr := rCtrl.Client.Set(key, message, 10*time.Second).Err()
	if cacheErr != nil {
		return cacheErr
	}
	return nil
}

func (rCtrl *RedisCtrl) GetMessage(key string) models.Message { // key - iphone@+254705207666 -
	val, err := rCtrl.Client.Get(key).Bytes()
	if err != nil {
		log.Println(err)
		return models.Message{}
	}

	msgData := toJson(val)
	return msgData
}

// Converts from []byte to a json object according to the User struct.
func toJson(val []byte) models.Message {
	message := models.Message{}
	err := json.Unmarshal(val, &message)
	if err != nil {
		panic(err)
	}
	return message
}
