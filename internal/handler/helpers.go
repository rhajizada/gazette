package handler

import (
	"errors"
	"fmt"
	"net/url"
	"strconv"
)

const MaxLimit = 100

func getPageParams(v url.Values) (PageParams, error) {
	var params PageParams
	limit := v.Get("limit")
	offset := v.Get("offset")
	limitAsInt64, err := strconv.ParseInt(limit, 10, 32)
	if err != nil {
		return params, errors.New("invalid limit type")
	}
	limitAsInt32 := int32(limitAsInt64)
	if limitAsInt32 > MaxLimit {
		msg := fmt.Sprintf("max limit size is %d", MaxLimit)
		return params, errors.New(msg)
	}
	offsetAsInt64, err := strconv.ParseInt(offset, 10, 32)
	if err != nil {
		return params, errors.New("invalid offset type")
	}
	offsetAsInt32 := int32(offsetAsInt64)

	params.Limit, params.Offset = limitAsInt32, offsetAsInt32
	return params, nil
}
