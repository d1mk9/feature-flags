package service

import (
	"encoding/json"
)

type FeatureValue struct {
	Bool   *bool    `json:"bool,omitempty"`
	Number *float64 `json:"number,omitempty"`
	String *string  `json:"string,omitempty"`
}

func (v FeatureValue) Validate() bool {
	n := 0
	if v.Bool != nil {
		n++
	}
	if v.Number != nil {
		n++
	}
	if v.String != nil {
		n++
	}
	return n == 1
}

// В ОТВЕТЕ наружу отдаем примитив
func (v FeatureValue) MarshalJSON() ([]byte, error) {
	switch {
	case v.Bool != nil:
		return json.Marshal(*v.Bool)
	case v.Number != nil:
		return json.Marshal(*v.Number)
	case v.String != nil:
		return json.Marshal(*v.String)
	default:
		// null — если ничего не задано
		return []byte("null"), nil
	}
}

// В ЗАПРОСЕ принимаем и примитив, и объектный вид
func (v *FeatureValue) UnmarshalJSON(b []byte) error {
	// bool
	var vb bool
	if err := json.Unmarshal(b, &vb); err == nil {
		v.Bool, v.Number, v.String = &vb, nil, nil
		return nil
	}
	// number
	var vn float64
	if err := json.Unmarshal(b, &vn); err == nil {
		v.Bool, v.Number, v.String = nil, &vn, nil
		return nil
	}
	// string
	var vs string
	if err := json.Unmarshal(b, &vs); err == nil {
		v.Bool, v.Number, v.String = nil, nil, &vs
		return nil
	}
	// объектный вид {"bool":...} / {"number":...} / {"string":...}
	var obj struct {
		Bool   *bool    `json:"bool"`
		Number *float64 `json:"number"`
		String *string  `json:"string"`
	}
	if err := json.Unmarshal(b, &obj); err != nil {
		return ErrUnsupportedType
	}
	v.Bool, v.Number, v.String = obj.Bool, obj.Number, obj.String
	if !v.Validate() {
		return ErrUnsupportedType
	}
	return nil
}
