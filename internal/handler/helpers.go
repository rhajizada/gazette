package handler

import (
	"errors"
	"net/url"
	"strconv"

	"github.com/rhajizada/gazette/internal/repository"
)

func getListFeedsParams(v url.Values) (repository.ListFeedsParams, error) {
	var params repository.ListFeedsParams
	limit := v.Get("limit")
	offset := v.Get("offset")
	limitAsInt64, err := strconv.ParseInt(limit, 10, 32)
	if err != nil {
		return params, errors.New("invalid limit type")
	}
	limitAsInt32 := int32(limitAsInt64)
	offsetAsInt64, err := strconv.ParseInt(offset, 10, 32)
	if err != nil {
		return params, errors.New("invalid offset type")
	}
	offsetAsInt32 := int32(offsetAsInt64)

	params.Limit, params.Offset = limitAsInt32, offsetAsInt32
	return params, nil
}
