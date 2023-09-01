package jq

import "encoding/json"

type JQResponseType int

const (
	Empty JQResponseType = iota
	String
	Map
	ArrayMap
	Array
	Unknown
)

type JQResponse struct {
	Type JQResponseType
	Raw  []byte
	Data any
}

func NewResponse(data []byte) (resp *JQResponse) {
	var _obj string
	var _map map[string]any
	var _mapArray []map[string]any
	var _array []any

	resp = &JQResponse{
		Type: Unknown,
		Raw:  data,
		Data: data,
	}

	if string(data) == "" {
		resp.Type = Empty
		return
	}

	if err := json.Unmarshal(data, &_obj); err == nil {
		resp.Data = _obj
		if _obj == "" {
			resp.Type = Empty
			return
		}
		resp.Type = String
		return
	}

	if err := json.Unmarshal(data, &_map); err == nil {
		resp.Data = _map
		resp.Type = Map
		return
	}

	if err := json.Unmarshal(data, &_mapArray); err == nil {
		resp.Data = _mapArray
		resp.Type = ArrayMap
		return
	}

	if err := json.Unmarshal(data, &_array); err == nil {
		resp.Data = _array
		resp.Type = Array
		return
	}

	return
}
