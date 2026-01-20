package middle

import (
	"context"
	"net/http"
	"sync"
	"time"
	"torchi/internal/pkg/log"

	"go.uber.org/fx"
	"golang.org/x/time/rate"
)

type LimiterItem struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

type RateLimiterManager struct {
	limiters map[string]*LimiterItem
	mu       sync.Mutex
	count    int
	log      *log.Logger
}

func NewRateLimiterManager(lc fx.Lifecycle, log *log.Logger) *RateLimiterManager {
	mgr := &RateLimiterManager{
		limiters: make(map[string]*LimiterItem),
		count:    4, //1초에 최대 n번 허용
		log:      log,
	}

	// FX Lifecycle을 사용하여 백그라운드 청소 작업 시작 및 정지
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go mgr.cleanupLoop()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			// 필요 시 종료 로그 기록
			return nil
		},
	})

	return mgr
}

// 주기적으로 오래된 리미터를 삭제하는 루프
func (m *RateLimiterManager) cleanupLoop() {
	ticker := time.NewTicker(3 * time.Minute) // 3분마다 검사
	defer ticker.Stop()

	for range ticker.C {
		m.log.Debug("Cleaning up rate limiters...")
		m.mu.Lock()
		for token, item := range m.limiters {
			// 5분 동안 활동이 없으면 메모리에서 제거
			if time.Since(item.lastSeen) > 5*time.Minute {
				delete(m.limiters, token)
			}
		}
		m.mu.Unlock()
	}
}

func (m *RateLimiterManager) GetLimiter(token string) *rate.Limiter {
	m.mu.Lock()
	defer m.mu.Unlock()

	if item, exists := m.limiters[token]; exists {
		m.log.Debug("Reusing rate limiter for token")
		item.lastSeen = time.Now() // 사용 시간 업데이트
		return item.limiter
	}
	m.log.Debug("Creating new rate limiter for token")

	l := rate.NewLimiter(rate.Every(time.Second), m.count)
	m.limiters[token] = &LimiterItem{
		limiter:  l,
		lastSeen: time.Now(),
	}
	return l
}

// 미들웨어 팩토리: Manager를 주입받아 미들웨어를 반환
func RateLimitMiddleware(mgr *RateLimiterManager) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			token := r.URL.Path
			mgr.log.Debug("token", "..", token)
			if token == "" {
				next.ServeHTTP(w, r)
				return
			}

			limiter := mgr.GetLimiter(token)
			if !limiter.Allow() {
				http.Error(w, "Too Many Requests", http.StatusTooManyRequests)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
