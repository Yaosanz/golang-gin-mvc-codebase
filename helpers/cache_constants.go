package helpers

import "time"

// Cache TTL constants - Best practice: different TTLs for different data types
const (
	// User cache TTL - User data changes less frequently
	UserCacheTTL = 30 * time.Minute

	// User authentication cache - Shorter TTL for security
	UserAuthCacheTTL = 5 * time.Minute

	// Shortlink cache - Longer TTL as it changes rarely
	ShortenLinkCacheTTL = 24 * time.Hour

	// Shortlink code cache - Longer TTL for redirect performance
	ShortenLinkCodeCacheTTL = 48 * time.Hour

	// User shortlinks list cache - Medium TTL
	UserShortenLinksCacheTTL = 10 * time.Minute

	// Permission/Role cache - Medium TTL
	PermissionCacheTTL = 15 * time.Minute
	RoleCacheTTL       = 15 * time.Minute

	// Session cache - Short TTL for security
	SessionCacheTTL = 2 * time.Hour
)

// Cache key prefixes for organization
const (
	KeyPrefixUser            = "user"
	KeyPrefixUserAuth        = "user:auth"
	KeyPrefixUserSession     = "user:session"
	KeyPrefixShortlink       = "shortlink"
	KeyPrefixShortenLinkCode = "shortlink:code"
	KeyPrefixPermission      = "permission"
	KeyPrefixRole            = "role"
	KeyPrefixCache           = "cache"
)

// Cache statistics keys for monitoring
const (
	CacheStatsHits      = "cache:stats:hits"
	CacheStatsMisses    = "cache:stats:misses"
	CacheStatsErrors    = "cache:stats:errors"
	CacheStatsEvictions = "cache:stats:evictions"
)
