package myredis

import "errors"

func Get(key interface{}) (value interface{}, e error) {
	r := MainRds.Get()
	defer r.Close()
	return r.Do("GET", key)
}

func Set(key interface{}, value interface{}) (e error) {
	r := MainRds.Get()
	defer r.Close()
	_, e = r.Do("SET", key, value)
	return
}

func SetEx(key interface{}, seconds int, value interface{}) (e error) {
	r := MainRds.Get()
	defer r.Close()
	_, e = r.Do("SETEX", key, seconds, value)
	return

}
func SetNx(key interface{}, value interface{}) (val int, e error) {
	r := MainRds.Get()
	defer r.Close()
	val, e = Int(r.Do("SETNX", key, value))
	return
}

func Incr(key interface{}) (int64, error) {
	r := MainRds.Get()
	defer r.Close()
	return Int64(r.Do("INCR", key))
}

func IncrBy(key interface{}, value interface{}) (int64, error) {
	r := MainRds.Get()
	defer r.Close()
	return Int64(r.Do("INCRBY", key, value))
}

// 批量获取
func MGet(keys ...interface{}) (value interface{}, e error) {
	r := MainRds.Get()
	defer r.Close()
	return r.Do("MGET", keys...)
}

/*
批量设置
< key value > 序列
*/
func MSet(kvs ...interface{}) (value interface{}, e error) {
	if len(kvs)%2 != 0 {
		return nil, errors.New("invalid arguments number")
	}
	r := MainRds.Get()
	defer r.Close()
	return r.Do("MSET", kvs...)
}
