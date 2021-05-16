package filter

// Filter 优先级
var priority = map[FilterName]int{
	// Pre
	PrePrepareFilter:   0,
	PreIPBlackFilter:   1,
	PreIPWhiteFilter:   2,
	PreAuthFilter:      3,
	PreRateLimitFilter: 4,
	PreHeadersFilter_:  5,

	// Post
	PostHeadersFilter_: 100,
}

type FilterName string

const (
	PrePrepareFilter   FilterName = "PREPARE"
	PreAuthFilter      FilterName = "AUTH"
	PreIPWhiteFilter   FilterName = "IP_WHITE"
	PreIPBlackFilter   FilterName = "IP_BLACK"
	PreHeadersFilter_  FilterName = "PRE_HEADERS"
	PreRateLimitFilter FilterName = "RATE_LIMIT"
	PostHeadersFilter_ FilterName = "POST_HEADERS"
)

func NewFilter(filterName FilterName) Filter {
	switch filterName {
	case PrePrepareFilter:
		return newPrepareFilter()
	case PreHeadersFilter_:
		return newPreHeadersFilter()
	case PreIPBlackFilter:
		return newIPBlackFilter()
	case PreIPWhiteFilter:
		return newIPWhiteFilter()
	case PreAuthFilter:
		return newAuthFilter()
	case PreRateLimitFilter:
		return newRateLimitFilter()
	case PostHeadersFilter_:
		return newPostHeadersFilter()
	default:
		return nil
	}
}
