package filter

import (
	"errors"
)

var (
	// ErrUnknownFilter unknown filter error
	ErrUnknownFilter = errors.New("unknown filter")
)

type FilterType string

const (
	PrePrepareFilter FilterType = "PREPARE"
)

type FilterFactory struct{}

func NewFilterFactory() *FilterFactory {
	return &FilterFactory{}
}

func (f *FilterFactory) newFilter(filterType FilterType) (Filter, error) {

	switch filterType {
	case PrePrepareFilter:
		return newPrepareFilter(), nil
	default:
		return nil, ErrUnknownFilter
	}
}
