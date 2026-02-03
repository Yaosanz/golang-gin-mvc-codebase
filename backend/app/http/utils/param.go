package utils

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type QueryParams struct {
	Paginate      bool              `json:"paginate"`
	Page          int               `json:"page"`
	Limit         int               `json:"limit"`
	Search        string            `json:"search"`
	SortBy        string            `json:"sort_by"`
	SortDirection string            `json:"sort_direction"`
	Filters       map[string]string `json:"filters"`
	DateFrom      *time.Time        `json:"date_from,omitempty"`
	DateTo        *time.Time        `json:"date_to,omitempty"`
	DateField     string            `json:"date_field,omitempty"` // NEW: column name to filter date
}

// QueryConfig Configurable defaults
type QueryConfig struct {
	Paginate      bool
	Page          int
	Limit         int
	MaxLimit      int
	SortBy        string
	SortDirection string
	DateLayout    string // e.g. "2006-01-02" or RFC3339
	DateField     string
}

// defaultQueryConfig Default configuration
var defaultQueryConfig = QueryConfig{
	Paginate:      true,
	Page:          1,
	Limit:         15,
	MaxLimit:      100,
	SortBy:        "created_at",
	SortDirection: "desc",
	DateLayout:    "2006-01-02",
	DateField:     "created_at",
}

// ParseQueryParams Parse query params from Gin context
func ParseQueryParams(c *gin.Context, cfg *QueryConfig) QueryParams {
	if cfg == nil {
		cfg = &defaultQueryConfig
	}

	q := QueryParams{
		Paginate:      cfg.Paginate,
		Page:          cfg.Page,
		Limit:         cfg.Limit,
		SortBy:        cfg.SortBy,
		SortDirection: cfg.SortDirection,
		Filters:       make(map[string]string),
		DateField:     cfg.DateField, // default field
	}

	// paginate
	if paginate, err := strconv.ParseBool(c.DefaultQuery("paginate", "true")); err == nil {
		q.Paginate = paginate
	}

	// page
	if page, err := strconv.Atoi(c.DefaultQuery("page", strconv.Itoa(cfg.Page))); err == nil {
		if page > 0 {
			q.Page = page
		}
	}

	// limit
	if limit, err := strconv.Atoi(c.DefaultQuery("limit", strconv.Itoa(cfg.Limit))); err == nil && limit > 0 {
		q.Limit = limit
	}

	// search
	q.Search = strings.TrimSpace(c.DefaultQuery("search", ""))

	// sort_by
	if sortBy := strings.TrimSpace(c.Query("sort_by")); sortBy != "" {
		q.SortBy = sortBy
	}

	// sort_direction
	if dir := strings.ToLower(c.Query("sort_direction")); dir == "asc" || dir == "desc" {
		q.SortDirection = dir
	}

	// filters
	for key, values := range c.Request.URL.Query() {
		if strings.HasPrefix(key, "filter[") && strings.HasSuffix(key, "]") {
			filterKey := key[7 : len(key)-1]
			if len(values) > 0 {
				q.Filters[filterKey] = values[0]
			}
		}
	}

	// date_from
	if df := c.Query("date_from"); df != "" {
		if t, err := time.Parse("2006-01-02", df); err == nil {
			q.DateFrom = &t
		}
	}

	// date_to
	if dt := c.Query("date_to"); dt != "" {
		if t, err := time.Parse("2006-01-02", dt); err == nil {
			q.DateTo = &t
		}
	}

	// date_field (configurable)
	if field := c.Query("date_field"); field != "" {
		q.DateField = field
	}

	return q
}

// Sanitize ensures defaults and clamps values
func (q *QueryParams) Sanitize() {
	if q.Page <= 0 {
		q.Page = 1
	}
	if q.Limit <= 0 || q.Limit > 100 {
		q.Limit = 10
	}
	if q.SortDirection != "asc" && q.SortDirection != "desc" {
		q.SortDirection = "desc"
	}
	if q.SortBy == "" {
		q.SortBy = "created_at"
	}
	if q.DateField == "" {
		q.DateField = "created_at"
	}
}

// ApplySearch applies search filtering across multiple fields
func ApplySearch(search string, fields ...string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		// search
		if search != "" && len(fields) > 0 {
			like := "%" + search + "%"
			conds := []string{}
			args := []interface{}{}
			for _, field := range fields {
				conds = append(conds, fmt.Sprintf("%s ILIKE ?", field))
				args = append(args, like)
			}
			db = db.Where(strings.Join(conds, " OR "), args...)
		}

		return db
	}
}

// ApplyFilter applies filtering based on allowed fields
func ApplyFilter(filters map[string]string, allowedFields map[string]bool) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		for field, value := range filters {
			if allowedFields[field] && value != "" {
				values := strings.Split(value, ",")
				if len(values) > 1 {
					// Use IN clause for multiple values
					db = db.Where(fmt.Sprintf("%s IN ?", field), values)
				} else if len(values) == 1 {
					// single value -> normal equality
					db = db.Where(fmt.Sprintf("%s = ?", field), value)
				}
			}
		}

		return db
	}
}

// ApplyDateFilter applies filters/search/pagination/sorting
func ApplyDateFilter(dateFrom, dateTo *time.Time, dateField string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if dateFrom != nil && dateTo != nil {
			// Use BETWEEN if both values are present
			db = db.Where(fmt.Sprintf("%s BETWEEN ? AND ?", dateField), dateFrom, dateTo)
		} else if dateFrom != nil {
			db = db.Where(fmt.Sprintf("%s >= ?", dateField), dateFrom)
		} else if dateTo != nil {
			db = db.Where(fmt.Sprintf("%s <= ?", dateField), dateTo)
		}
		return db
	}
}

// ApplySorting applies sorting based on QueryParams
func ApplySorting(q QueryParams) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		// Default sort direction to "asc" if not provided or invalid
		//if strings.ToLower(q.SortDirection) != "asc" && strings.ToLower(q.SortDirection) != "desc" {
		order := fmt.Sprintf("%s %s", q.SortBy, strings.ToUpper(q.SortDirection))
		db = db.Order(order)
		//}

		return db
	}
}

// ApplyPagination applies pagination if enabled
func ApplyPagination(q QueryParams) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if q.Paginate {
			offset := (q.Page - 1) * q.Limit
			db = db.Offset(offset).Limit(q.Limit)
		}
		return db
	}
}

// GetIDParam extracts and validates ID from URI param
func GetIDParam(ctx *gin.Context, keyName string) (int64, error) {
	idStr := ctx.Param(keyName)
	if idStr == "" {
		return 0, fmt.Errorf("missing id parameter")
	}

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil || id <= 0 {
		return 0, fmt.Errorf("invalid id parameter: %v", err)
	}

	return id, nil
}

// GetSlugParam extracts slug from URI param
func GetSlugParam(ctx *gin.Context, keyName string) (string, error) {
	slug := ctx.Param(keyName)
	if slug == "" {
		return "", fmt.Errorf("missing id parameter")
	}

	return slug, nil
}
