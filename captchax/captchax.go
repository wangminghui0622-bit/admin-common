// Package captchax 封装 base64Captcha（5 位数字）+ Redis store，
// key 格式与 simple-admin 一致：CAPTCHA:<id>。
package captchax

import (
	"time"

	"github.com/mojocn/base64Captcha"
	"github.com/zeromicro/go-zero/core/stores/redis"
)

// Store 实现 base64Captcha.Store 接口，验证码存 Redis。
type Store struct {
	redis *redis.Redis
}

func NewStore(r *redis.Redis) *Store {
	return &Store{redis: r}
}

// Set 把验证码答案写入 Redis，5 分钟后自动失效。
func (s *Store) Set(id, value string) error {
	return s.redis.Setex("CAPTCHA:"+id, value, 300)
}

// Get 取出答案；clear 为 true 时取完即删（一次性验证码）。
func (s *Store) Get(id string, clear bool) string {
	key := "CAPTCHA:" + id
	v, err := s.redis.Get(key)
	if err != nil {
		return ""
	}
	if clear {
		_, _ = s.redis.Del(key)
	}
	return v
}

// Verify 校验答案。
func (s *Store) Verify(id, answer string, clear bool) bool {
	return s.Get(id, clear) == answer
}

// Captcha 数字验证码生成器。
type Captcha struct {
	c *base64Captcha.Captcha
}

func New(r *redis.Redis) *Captcha {
	// 5 位数字，宽 240 高 80
	driver := base64Captcha.NewDriverDigit(80, 240, 5, 0.7, 80)
	return &Captcha{c: base64Captcha.NewCaptcha(driver, NewStore(r))}
}

// Generate 返回 (captchaId, base64 data uri, 答案, err)。
func (c *Captcha) Generate() (string, string, string, error) {
	return c.c.Generate()
}

// Verify 校验并清除验证码。
func (c *Captcha) Verify(id, answer string, clear bool) bool {
	return c.c.Verify(id, answer, clear)
}

// Set 兼容 go-zero captcha 的 Store 签名（过期时间入参，预留）。
func (s *Store) SetWithExpire(id, value string, expire time.Duration) error {
	return s.redis.Setex("CAPTCHA:"+id, value, int(expire.Seconds()))
}
