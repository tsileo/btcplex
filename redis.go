package btcplex

import (
    "github.com/garyburd/redigo/redis"
    "time"
)

func GetRedis(conf *Config) (pool *redis.Pool, err error) {
	pool = &redis.Pool{
             MaxIdle: 3,
             IdleTimeout: 240 * time.Second,
             Dial: func () (redis.Conn, error) {
                 c, err := redis.Dial("tcp", conf.RedisHost)
                 if err != nil {
                     return nil, err
                 }
//                 if _, err := c.Do("AUTH", password); err != nil {
//                     c.Close()
//                     return nil, err
//               }
//                 return c, err
                return c, err
             },
                TestOnBorrow: func(c redis.Conn, t time.Time) error {
                    _, err := c.Do("PING")
                 return err
                },
         }
    return
}
