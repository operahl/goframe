package myredis

import "errors"

func SAdd(key interface{}, values ...interface{}) (e error) {
	args := make([]interface{}, len(values)+1)
	args[0] = key
	copy(args[1:], values)
	r := MainRds.Get()
	defer r.Close()
	_, e = r.Do("SADD", args...)
	return
}

/*
批量添加到set类型的表中
	args: 必须是<key,id>的列表
*/
func SMultiAdd(args ...interface{}) error {
	if len(args)%2 != 0 {
		return errors.New("invalid arguments number")
	}
	r := MainRds.Get()
	defer r.Close()
	if e := r.Send("MULTI"); e != nil {
		return e
	}
	for i := 0; i < len(args); i += 2 {
		if e := r.Send("SADD", args[i], args[i+1]); e != nil {
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

/*
批量删除set类型表中的元素
	args: 必须是<key,id>的列表
*/
func SMultiRem(args ...interface{}) error {
	if len(args)%2 != 0 {
		return errors.New("invalid arguments number")
	}
	r := MainRds.Get()
	defer r.Close()
	if e := r.Send("MULTI"); e != nil {
		return e
	}
	for i := 0; i < len(args); i += 2 {
		if e := r.Send("SREM", args[i], args[i+1]); e != nil {
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

func SRem(key interface{}, values ...interface{}) (e error) {
	args := make([]interface{}, len(values)+1)
	args[0] = key
	copy(args[1:], values)
	r := MainRds.Get()
	defer r.Close()
	_, e = r.Do("SREM", args...)
	return
}

func SIsMember(key interface{}, value interface{}) (isMember bool, e error) {
	r := MainRds.Get()
	defer r.Close()
	return Bool(r.Do("SISMEMBER", key, value))
}

/*
SMembers获取某个key下的所有元素

参数：
	values: 必须是数组的引用
*/
func SMembers(values interface{}, key interface{}) (e error) {
	r := MainRds.Get()
	defer r.Close()
	vs, e := Values(r.Do("SMEMBERS", key))
	if e != nil {
		return e
	}
	return ScanSlice(vs, values)
}

/*
SCard获取某个key下的元素数量

参数：
	values: 必须是数组的引用
*/
func SCard(key interface{}) (count int, e error) {
	r := MainRds.Get()
	defer r.Close()
	count, e = Int(r.Do("SCARD", key))
	if e != nil {
		return
	}
	return
}

/*
SRandMembers获取某个key下的随机count 个元素

参数：
	values: 必须是数组的引用
*/
func SRandMembers(key interface{}, count int, values interface{}) (e error) {
	r := MainRds.Get()
	defer r.Close()
	vs, e := Values(r.Do("SRANDMEMBER", key, count))
	if e != nil {
		return e
	}
	return ScanSlice(vs, values)
}
