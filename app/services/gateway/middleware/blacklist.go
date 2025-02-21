package middleware

import (
	"context"
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/bitdance-panic/gobuy/app/models"
	"github.com/bitdance-panic/gobuy/app/services/gateway/biz/dal/tidb"
	"github.com/bitdance-panic/gobuy/app/services/gateway/biz/dao"
	"github.com/bitdance-panic/gobuy/app/utils"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/hertz-contrib/jwt"
)

var (
	BlacklistCache sync.Map     // 内存缓存：map[identifier]expiryTime
	CacheMutex     sync.RWMutex // 缓存读写锁
	LastSyncTime   time.Time    // 最后同步时间
)

// 初始化黑名单缓存
func InitBlacklistCache() {
	// 首次启动时全量加载
	loadBlacklistFromDB()

	// 定时同步（每5分钟同步一次）
	go func() {
		ticker := time.NewTicker(5 * time.Minute)
		defer ticker.Stop()
		for range ticker.C {
			loadBlacklistFromDB()
		}
	}()
}

// 从数据库加载黑名单
func loadBlacklistFromDB() {
	// 通过加锁来确保后续操作的线程安全，防止多个 goroutine 同时修改黑名单缓存（blacklistCache）时出现竞态条件
	CacheMutex.Lock()
	// 在函数结束时，释放锁。defer 关键字确保即使发生错误，锁也会被释放
	defer CacheMutex.Unlock()

	var entries []models.Blacklist
	if err := tidb.DB.Where("(expires_at > ? OR expires_at IS NULL) AND is_deleted = 0", time.Now()).Find(&entries).Error; err != nil {
		hlog.Errorf("黑名单同步失败: %v", err)
		return
	}

	// 更新缓存
	tempMap := make(map[string]time.Time)
	for _, entry := range entries {
		tempMap[entry.Identifier] = entry.ExpiresAt
	}
	BlacklistCache = sync.Map{} // 清空缓存
	for k, v := range tempMap {
		BlacklistCache.Store(k, v)
	}

	LastSyncTime = time.Now()
	hlog.Infof("黑名单缓存同步完成，当前条目数: %d", len(tempMap))
}

// 检查是否在黑名单中
func IsBlocked(identifier string) bool {
	// 先检查缓存
	if expiry, ok := BlacklistCache.Load(identifier); ok {
		// 永久封禁或未过期
		if expiry.(time.Time).IsZero() || expiry.(time.Time).After(time.Now()) {
			return true
		}
		// 已过期，从缓存移除
		BlacklistCache.Delete(identifier)
	}
	return false
}

// 黑名单中间件
func BlacklistMiddleware() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		// 检查是否跳过认证(白名单接口)
		if skip, exists := c.Get("skip_auth"); exists && skip.(bool) {
			c.Next(ctx)
			return
		}

		// 获取用户标识（示例：优先取用户ID，未登录则取IP）
		var identifier string
		if claims := jwt.ExtractClaims(ctx, c); claims != nil {
			userID := claims[IdentityKey].(float64)
			userIDStr := strconv.Itoa(int(userID))
			identifier = fmt.Sprintf("user:%v", userIDStr)
		} else {
			identifier = fmt.Sprintf("ip:%s", c.ClientIP())
		}

		// 检查黑名单
		if IsBlocked(identifier) {
			hlog.Warnf("拒绝黑名单用户访问: %s", identifier)
			utils.FailFull(c, consts.StatusForbidden, "您的账户已被封禁", nil)
			c.Abort()
			return
		}

		c.Next(ctx)
	}
}

// 自动清理过期条目
func StartBlacklistCleanupTask() {
	go func() {
		ticker := time.NewTicker(1 * time.Hour)
		defer ticker.Stop()
		for range ticker.C {
			// 删除过期条目
			if err := dao.ClearExpBlacklist(tidb.DB); err != nil {
				hlog.Errorf("黑名单清理失败: %v", err)
			}
		}
	}()
}
