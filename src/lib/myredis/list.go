package myredis

//批量插入队尾
func RPush(key interface{}, values ...interface{}) (value interface{}, e error) {
	if len(values) == 0 {
		return
	}
	r := MainRds.Get()
	defer r.Close()
	vs := []interface{}{key}
	vs = append(vs, values...)
	return r.Do("RPUSH", vs...)
}

//批量插入队头
func LPush(key interface{}, values ...interface{}) (value interface{}, e error) {
	if len(values) == 0 {
		return
	}
	r := MainRds.Get()
	defer r.Close()
	vs := []interface{}{key}
	vs = append(vs, values...)
	return r.Do("LPUSH", vs...)
}

/*
获取队列数据
*/
func LRange(key interface{}, start, stop interface{}) (value interface{}, e error) {
	r := MainRds.Get()
	defer r.Close()
	return r.Do("LRANGE", key, start, stop)
}

/*
获取队列长度，如果key不存在，length=0，不会报错。
*/
func LLen(key interface{}) (length int, e error) {
	r := MainRds.Get()
	defer r.Close()
	return Int(r.Do("LLEN", key))
}

/*
弹出队列数据
*/
func LPop(key interface{}) (value interface{}, e error) {
	r := MainRds.Get()
	defer r.Close()
	return r.Do("LPOP", key)
}

/*
keys: 键列表
*/
func BRPop(timeout interface{}, keys ...interface{}) (value interface{}, e error) {
	if len(keys) > 0 {
		keys = append(keys, timeout)
	}
	r := MainRds.Get()
	defer r.Close()
	return r.Do("BRPOP", keys...)
}

func RPop(key interface{}) (value interface{}, e error) {
	r := MainRds.Get()
	defer r.Close()
	return r.Do("RPOP", key)
}

/*
keys: 键列表
*/
func BLpop(timeout interface{}, keys ...interface{}) (value interface{}, e error) {
	if len(keys) > 0 {
		keys = append(keys, timeout)
	}
	r := MainRds.Get()
	defer r.Close()
	return r.Do("BLPOP", keys...)
}

/*
剪裁
*/
func LTrim(key interface{}, start, end int) (e error) {
	r := MainRds.Get()
	defer r.Close()
	_, e = r.Do("LTRIM", key, start, end)
	return
}

/*
删除
*/
func LRem(key interface{}, count int, value interface{}) (e error) {
	r := MainRds.Get()
	defer r.Close()
	_, e = r.Do("LREM", key, count, value)
	return
}

/*
索引元素
*/
func LIndex(key interface{}, index int) (value interface{}, e error) {
	r := MainRds.Get()
	defer r.Close()
	value, e = r.Do("LINDEX", key, index)
	return
}

/*
获取头部元素
*/
func Front(key interface{}) (value interface{}, e error) {
	return LIndex(key, 0)
}

/*
获取尾部元素
*/
func Back(key interface{}) (value interface{}, e error) {
	return LIndex(key, -1)
}
