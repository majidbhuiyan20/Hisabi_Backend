package middleware

import (
	"hisabi.com/m/utils"
	"net/http"
	"strings"
	"sync"
	"time"
)

const (
	maxAttempts = 10
	windowTime  = 5 * time.Minute
	blockTime   = 5 * time.Minute
)

type attemptRecord struct {
	count        int
	windowStart  time.Time
	blockedUntil time.Time
	mu           sync.Mutex
}

var (
	records   = make(map[string]*attemptRecord)
	recordsMu sync.RWMutex
)

func LoginRateLimit(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		ip := getClientIP(r)

		// Get or create record for this IP
		recordsMu.Lock()
		rec, exists := records[ip]
		if !exists {
			rec = &attemptRecord{windowStart: time.Now()}
			records[ip] = rec
		}
		recordsMu.Unlock()

		rec.mu.Lock()
		defer rec.mu.Unlock()

		now := time.Now()

		// ── Check if currently blocked ────────────────────
		if now.Before(rec.blockedUntil) {
			remaining := rec.blockedUntil.Sub(now).Round(time.Second)
			w.Header().Set("Retry-After", remaining.String())
			w.Header().Set("X-RateLimit-Limit", "10")
			w.Header().Set("X-RateLimit-Remaining", "0")
			utils.JSONStatus(w, http.StatusTooManyRequests,
				false,
				"Too many login attempts. Please try again in "+remaining.String(),
				map[string]interface{}{
					"retry_after_seconds": int(remaining.Seconds()),
				},
			)
			return
		}

		// ── Reset window if 5 minutes have passed ─────────
		if now.Sub(rec.windowStart) > windowTime {
			rec.count = 0
			rec.windowStart = now
		}

		// ── Check attempt count ───────────────────────────
		if rec.count >= maxAttempts {
			// Block the IP for 5 minutes
			rec.blockedUntil = now.Add(blockTime)
			rec.count = 0 // reset for next window after block
			rec.windowStart = now

			w.Header().Set("Retry-After", blockTime.String())
			w.Header().Set("X-RateLimit-Limit", "10")
			w.Header().Set("X-RateLimit-Remaining", "0")
			utils.JSONStatus(w, http.StatusTooManyRequests,
				false,
				"Too many login attempts. Your IP has been blocked for 5 minutes.",
				map[string]interface{}{
					"retry_after_seconds": int(blockTime.Seconds()),
				},
			)
			return
		}

		// ── Increment attempt count ───────────────────────
		rec.count++
		remaining := maxAttempts - rec.count

		// Set rate limit headers so client can see status
		w.Header().Set("X-RateLimit-Limit", "10")
		w.Header().Set("X-RateLimit-Remaining", string(rune('0'+remaining)))

		next.ServeHTTP(w, r)
	})
}

func getClientIP(r *http.Request) string {
	if forwarded := r.Header.Get("X-Forwarded-For"); forwarded != "" {
		parts := strings.Split(forwarded, ",")
		return strings.TrimSpace(parts[0])
	}

	if realIP := r.Header.Get("X-Real-IP"); realIP != "" {
		return strings.TrimSpace(realIP)
	}

	ip := r.RemoteAddr
	if idx := strings.LastIndex(ip, ":"); idx != -1 {
		ip = ip[:idx]
	}
	return ip
}
