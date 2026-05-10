package utils

import (
	"encoding/json"
	"fmt"
	"gin-advanced/core"
	"time"
)

const (
	// Redis keys
	UserInfoKeyPrefix = "user:info:"   // 用户信息缓存
	TokenKeyPrefix    = "user:token:"  // Token映射用户ID
	UserExpire        = 24 * time.Hour // 缓存过期时间
)

// UserCache 用户缓存结构
type UserCache struct {
	ID           uint     `json:"id"`
	Username     string   `json:"username"`
	Nickname     string   `json:"nickname"`
	Roles        []string `json:"roles"`
	Permissions  []string `json:"permissions"`
}

// SetUserCache 设置用户缓存
func SetUserCache(userID uint, username, nickname string, roles, permissions []string) error {
	cache := UserCache{
		ID:           userID,
		Username:     username,
		Nickname:     nickname,
		Roles:        roles,
		Permissions:  permissions,
	}
	data, err := json.Marshal(cache)
	if err != nil {
		return err
	}
	key := fmt.Sprintf("%s%d", UserInfoKeyPrefix, userID)
	return core.RDB.Set(core.Ctx, key, data, UserExpire).Err()
}

// SetTokenCache 设置Token映射
func SetTokenCache(token string, userID uint) error {
	key := fmt.Sprintf("%s%s", TokenKeyPrefix, token)
	return core.RDB.Set(core.Ctx, key, userID, UserExpire).Err()
}

// GetUserIDByToken 通过Token获取用户ID
func GetUserIDByToken(token string) (uint, error) {
	key := fmt.Sprintf("%s%s", TokenKeyPrefix, token)
	userID, err := core.RDB.Get(core.Ctx, key).Uint64()
	if err != nil {
		return 0, err
	}
	return uint(userID), nil
}

// GetUserCache 获取用户缓存
func GetUserCache(userID uint) (*UserCache, error) {
	key := fmt.Sprintf("%s%d", UserInfoKeyPrefix, userID)
	data, err := core.RDB.Get(core.Ctx, key).Bytes()
	if err != nil {
		return nil, err
	}
	var cache UserCache
	if err := json.Unmarshal(data, &cache); err != nil {
		return nil, err
	}
	return &cache, nil
}

// GetUsernameByCache 通过缓存获取用户名
func GetUsernameByCache(userID uint) (string, error) {
	cache, err := GetUserCache(userID)
	if err != nil {
		return "", err
	}
	return cache.Username, nil
}

// DelUserCache 删除用户缓存
func DelUserCache(userID uint) error {
	key := fmt.Sprintf("%s%d", UserInfoKeyPrefix, userID)
	return core.RDB.Del(core.Ctx, key).Err()
}

// DelTokenCache 删除Token缓存
func DelTokenCache(token string) error {
	key := fmt.Sprintf("%s%s", TokenKeyPrefix, token)
	return core.RDB.Del(core.Ctx, key).Err()
}
