package myredis

import (
	"errors"
	"github.com/garyburd/redigo/redis"
)

func HMSet(key string, args map[string]string) (e error) {
	r := MainRds.Get()
	defer r.Close()

	_, e = r.Do("HMSET", redis.Args{}.Add(key).AddFlat(args)...)
	return
}

/*
HSet批量设置HashSet中的值

	args: 必须是<key,id,value>的列表
*/
func HMultiSet(args ...interface{}) (e error) {
	if len(args)%3 != 0 {
		return errors.New("invalid arguments number")
	}
	r := MainRds.Get()
	defer r.Close()
	if e := r.Send("MULTI"); e != nil {
		return e
	}
	for i := 0; i < len(args); i += 3 {
		if e := r.Send("HSET", args[i], args[i+1], args[i+2]); e != nil {
			r.Send("DISCARD")
			return e
		}
	}
	if _, e := r.Do("EXEC"); e != nil {
		r.Send("DISCARD")
		return e
	}
	return nil
}

func HSet(key interface{}, id interface{}, value interface{}) (e error) {
	r := MainRds.Get()
	defer r.Close()
	_, e = r.Do("HSET", key, id, value)
	return
}

func HGet(key interface{}, name interface{}) (value interface{}, e error) {
	r := MainRds.Get()
	defer r.Close()
	return r.Do("HGET", key, name)
}

/*
HMGet针对同一个key获取hashset中的部分元素的值

参数：
	args: 第一个值必须是key，后续的值都是id
	values: 必须是数组的引用，如果某个id不存在，会把对应数据类型的零值放在数组对应位置上
*/
func HMGet(values interface{}, key interface{}, ids ...interface{}) (e error) {
	if len(ids) == 0 {
		return
	}
	args := make([]interface{}, len(ids)+1)
	args[0] = key
	copy(args[1:], ids)
	r := MainRds.Get()
	defer r.Close()
	vs, e := Values(r.Do("HMGET", args...))
	if e != nil {
		return e
	}
	return ScanSlice(vs, values)
}

/*
HMultiGet批量获取HashSet中多个key中ID的值

参数：
	args: 必须是<key,id>的列表
返回值：
	values: 一个两层的map，第一层的key是参数中的key，第二层的key是参数中的id
*/
func HMultiGet(args ...interface{}) (values map[interface{}]map[interface{}]interface{}, e error) {
	if len(args)%2 != 0 {
		return nil, errors.New("invalid arguments number")
	}
	r := MainRds.Get()
	defer r.Close()
	for i := 0; i < len(args); i += 2 {
		if e := r.Send("HGET", args[i], args[i+1]); e != nil {
			return nil, e
		}
	}
	r.Flush()
	values = make(map[interface{}]map[interface{}]interface{}, len(args))
	for i := 0; i < len(args); i += 2 {
		v, e := r.Receive()
		switch e {
		case nil:
			idm, ok := values[args[i]]
			if !ok {
				idm = make(map[interface{}]interface{})
				values[args[i]] = idm
			}
			idm[args[i+1]] = v
		case ErrNil:
		default:
			return nil, e
		}
	}
	return values, nil
}

/*
HDel批量删除某个Key中的元素
	args: 第一个必须是key，后面的都是id
*/
func HDel(args ...interface{}) (e error) {
	if len(args) <= 1 {
		return nil
	}
	r := MainRds.Get()
	defer r.Close()
	_, e = r.Do("HDEL", args...)
	return
}

/*
HGetAll针对同一个key获取hashset中的所有元素的值
*/
func HGetAll(key interface{}) (reply interface{}, e error) {
	r := MainRds.Get()
	defer r.Close()
	return r.Do("HGETALL", key)
}

/*
MHGetAll批量获取多个key所有的字段
*/
func MHGetAll(args ...interface{}) (reply map[interface{}][]interface{}, e error) {
	r := MainRds.Get()
	defer r.Close()
	for _, key := range args {
		if e = r.Send("HGETALL", key); e != nil {
			return
		}
	}
	r.Flush()
	reply = make(map[interface{}][]interface{})
	for _, key := range args {
		r, e := Values(r.Receive())
		if e != nil {
			return reply, e
		}
		reply[key] = r
	}
	return
}

func HIncrBy(key interface{}, field interface{}, increment int64) (reply int64, e error) {
	r := MainRds.Get()
	defer r.Close()
	return Int64(r.Do("HINCRBY", key, field, increment))
}

func HExists(key string, field interface{}) (bool, error) {
	r := MainRds.Get()
	defer r.Close()
	return Bool(r.Do("hexists", key, field))
}
func HKeys(key interface{}) ([]string, error) {
	r := MainRds.Get()
	defer r.Close()
	return Strings(r.Do("hkeys", key))
}

func Hlen(key string) (reply int64, e error) {
	r := MainRds.Get()
	defer r.Close()
	return Int64(r.Do("hlen", key))
}
