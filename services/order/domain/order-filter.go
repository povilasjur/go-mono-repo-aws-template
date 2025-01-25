package domain

import (
	"common"
	"net/url"
	"time"
)

type OrderFilter struct {
	// in: query
	// required: false
	Id string `json:"id"`

	// in: query
	// required: false
	Name string `json:"name"`

	// in: query
	// required: false
	CreatedFrom *time.Time `json:"createdFrom"`

	// in: query
	// required: false
	CreatedTo *time.Time `json:"createdTo"`
}

func ParseOrderFilter(queryParams url.Values) (*OrderFilter, error) {
	createdFrom, err := common.GetTimestampFromQueryParams(queryParams, "createdFrom")
	if err != nil {
		return nil, err
	}
	createdTo, err := common.GetTimestampFromQueryParams(queryParams, "createdTo")
	if err != nil {
		return nil, err
	}
	return &OrderFilter{
		Id:          common.GetFilterByName("id", queryParams),
		Name:        common.GetFilterByName("name", queryParams),
		CreatedFrom: createdFrom,
		CreatedTo:   createdTo,
	}, nil
}
