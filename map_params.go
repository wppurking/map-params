package map_params

import (
	"fmt"
	"math"
	"reflect"
)

// 参考 gocraft/work 中 job.Args 对 map[string]interface{} 的参数获取提供方便方法

// MapParams 对 map[string]interface{} 的类型别名
type MapParams map[string]interface{}

// 返回操作执行后的 error, 如果没有则为 nil
func (m MapParams) Error() error {
	err, ok := m["$_error"]
	if ok {
		return err.(error)
	}
	return nil
}

// String 返回 key 键的字符串类型的值
func (m MapParams) String(key string) string {
	v, ok := m[key]
	if ok {
		typedV, ok := v.(string)
		if ok {
			return typedV
		}
		m["$_error"] = typecastError("string", key, v)
	} else {
		m["$_error"] = missingKeyError("string", key)
	}
	return ""
}

// Bool 尝试返回布尔值, 如果解析错误返回默认的 false
func (m MapParams) Bool(key string) bool {
	v, ok := m[key]
	if ok {
		typedV, ok := v.(bool)
		if ok {
			return typedV
		}
		m["$_error"] = typecastError("bool", key, v)
	} else {
		m["$_error"] = missingKeyError("bool", key)
	}
	return false
}

// Int64 尝试返回 int64 值(float 类型也会尝试转换成为 int64)
func (m MapParams) Int64(key string) int64 {
	v, ok := m[key]
	if ok {
		rVal := reflect.ValueOf(v)
		if isIntKind(rVal) {
			return rVal.Int()
		} else if isUintKind(rVal) {
			vUint := rVal.Uint()
			if vUint <= math.MaxInt64 {
				return int64(vUint)
			}
		} else if isFloatKind(rVal) {
			vFloat64 := rVal.Float()
			vInt64 := int64(vFloat64)
			if vFloat64 == math.Trunc(vFloat64) && vInt64 <= 9007199254740892 && vInt64 >= -9007199254740892 {
				return vInt64
			}
		}
		m["$_error"] = typecastError("int64", key, v)
	} else {
		m["$_error"] = missingKeyError("int64", key)
	}
	return 0
}

// Float64 尝试返回 float64 值(int 类型也会尝试转换成为 float64)
func (m MapParams) Float64(key string) float64 {
	v, ok := m[key]
	if ok {
		rVal := reflect.ValueOf(v)
		if isIntKind(rVal) {
			return float64(rVal.Int())
		} else if isUintKind(rVal) {
			return float64(rVal.Uint())
		} else if isFloatKind(rVal) {
			return rVal.Float()
		}
		m["$_error"] = typecastError("float64", key, v)
	} else {
		m["$_error"] = missingKeyError("float64", key)
	}
	return 0.0
}

func isIntKind(v reflect.Value) bool {
	k := v.Kind()
	return k == reflect.Int || k == reflect.Int8 || k == reflect.Int16 || k == reflect.Int32 || k == reflect.Int64
}

func isUintKind(v reflect.Value) bool {
	k := v.Kind()
	return k == reflect.Uint || k == reflect.Uint8 || k == reflect.Uint16 || k == reflect.Uint32 || k == reflect.Uint64
}

func isFloatKind(v reflect.Value) bool {
	k := v.Kind()
	return k == reflect.Float32 || k == reflect.Float64
}

func missingKeyError(jsonType, key string) error {
	return fmt.Errorf("looking for a %s in job.Arg[%s] but key wasn't found", jsonType, key)
}

func typecastError(jsonType, key string, v interface{}) error {
	actualType := reflect.TypeOf(v)
	return fmt.Errorf("looking for a %s in job.Arg[%s] but value wasn't right type: %v(%v)", jsonType, key, actualType, v)
}
