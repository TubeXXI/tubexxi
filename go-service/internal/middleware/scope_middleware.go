package middleware

import (
	"encoding/json"
	"net/url"
	"os"
	"strings"
	"sync"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

const settingsScopeLocalKey = "settings_scope"

var (
	settingsScopeMapOnce sync.Once
	settingsScopeMap     map[string]string
)

type ScopeMiddleware struct {
	ctxinject *ContextMiddleware
	logger    *zap.Logger
}

func NewScopeMiddleware(
	ctxinject *ContextMiddleware,
	logger *zap.Logger,
) *ScopeMiddleware {
	return &ScopeMiddleware{
		ctxinject: ctxinject,
		logger:    logger,
	}
}
func (m *ScopeMiddleware) SettingsScopeMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		if c.Locals(settingsScopeLocalKey) == nil {
			scope := m.ResolveSettingsScope(c)
			if scope != "" {
				c.Locals(settingsScopeLocalKey, scope)
			}
		}
		return c.Next()
	}
}

func (m *ScopeMiddleware) GetSettingsScope(c *fiber.Ctx) string {
	v := c.Locals(settingsScopeLocalKey)
	if v == nil {
		return "default"
	}
	if s, ok := v.(string); ok {
		if s == "" {
			return "default"
		}
		return s
	}
	return "default"
}

func (m *ScopeMiddleware) ResolveSettingsScope(c *fiber.Ctx) string {
	m.loadSettingsScopeMap()

	domain := m.resolveClientDomain(c)
	if domain == "" {
		return "default"
	}

	if len(settingsScopeMap) == 0 {
		return "default"
	}

	if scope, ok := settingsScopeMap[domain]; ok {
		return scope
	}

	for k, scope := range settingsScopeMap {
		if k == "" {
			continue
		}
		if strings.HasPrefix(k, "*.") {
			suffix := strings.TrimPrefix(k, "*.")
			if suffix != "" && (domain == suffix || strings.HasSuffix(domain, "."+suffix)) {
				return scope
			}
			continue
		}
		if strings.HasPrefix(k, ".") {
			suffix := strings.TrimPrefix(k, ".")
			if suffix != "" && (domain == suffix || strings.HasSuffix(domain, "."+suffix)) {
				return scope
			}
			continue
		}
	}

	return "default"
}

func (m *ScopeMiddleware) loadSettingsScopeMap() {
	settingsScopeMapOnce.Do(func() {
		raw := strings.TrimSpace(os.Getenv("SETTINGS_SCOPE_MAP"))
		if raw == "" {
			settingsScopeMap = map[string]string{}
			return
		}

		var payload map[string]string
		if err := json.Unmarshal([]byte(raw), &payload); err != nil {
			settingsScopeMap = map[string]string{}
			return
		}

		n := make(map[string]string, len(payload))
		for k, v := range payload {
			key := m.normalizeDomain(k)
			val := strings.TrimSpace(v)
			if key == "" || val == "" {
				continue
			}
			n[key] = val
		}
		settingsScopeMap = n
	})
}

func (m *ScopeMiddleware) resolveClientDomain(c *fiber.Ctx) string {
	if origin := strings.TrimSpace(c.Get("Origin")); origin != "" {
		if u, err := url.Parse(origin); err == nil {
			if d := m.normalizeDomain(u.Host); d != "" {
				return d
			}
		}
	}

	if referer := strings.TrimSpace(c.Get("Referer")); referer != "" {
		if u, err := url.Parse(referer); err == nil {
			if d := m.normalizeDomain(u.Host); d != "" {
				return d
			}
		}
	}

	if xfHost := strings.TrimSpace(c.Get("X-Forwarded-Host")); xfHost != "" {
		if d := m.normalizeDomain(xfHost); d != "" {
			return d
		}
	}

	if host := strings.TrimSpace(c.Get("Host")); host != "" {
		if d := m.normalizeDomain(host); d != "" {
			return d
		}
	}

	return ""
}

func (m *ScopeMiddleware) normalizeDomain(host string) string {
	host = strings.TrimSpace(host)
	host = strings.Trim(host, "\"'")
	host = strings.ToLower(host)
	if host == "" {
		return ""
	}
	if h, _, ok := strings.Cut(host, ":"); ok {
		host = h
	}

	host = strings.TrimPrefix(host, "www.")

	return host
}
