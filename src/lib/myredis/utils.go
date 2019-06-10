package myredis

import (
	redigo "github.com/garyburd/redigo/redis"
)

func Values(reply interface{}, err error) ([]interface{}, error) {
	return redigo.Values(reply, err)
}
func Scan(src []interface{}, dest ...interface{}) ([]interface{}, error) {
	return redigo.Scan(src, dest...)
}

func ScanSlice(src []interface{}, dest interface{}, fieldNames ...string) error {
	return redigo.ScanSlice(src, dest, fieldNames...)
}
func ScanStruct(src []interface{}, dest interface{}) error {
	return redigo.ScanStruct(src, dest)
}
func Bytes(reply interface{}, err error) ([]byte, error) {
	return redigo.Bytes(reply, err)
}
func String(reply interface{}, err error) (string, error) {
	return redigo.String(reply, err)
}
func Strings(reply interface{}, err error) ([]string, error) {
	return redigo.Strings(reply, err)
}
func Bool(reply interface{}, err error) (bool, error) {
	return redigo.Bool(reply, err)
}
func Int(reply interface{}, err error) (int, error) {
	return redigo.Int(reply, err)
}
func Int64(reply interface{}, err error) (int64, error) {
	return redigo.Int64(reply, err)
}
func Uint64(reply interface{}, err error) (uint64, error) {
	return redigo.Uint64(reply, err)
}
func Uint32(reply interface{}, err error) (uint32, error) {
	v, e := redigo.Uint64(reply, err)
	if e != nil {
		return 0, e
	}
	return uint32(v), nil
}

func Float64(reply interface{}, err error) (float64, error) {
	return redigo.Float64(reply, err)
}

func Int64Map(reply interface{}, err error) (map[string]int64, error) {
	return redigo.Int64Map(reply, err)
}

//func Int32Int32Map(reply interface{}, err error) (map[int32]int32, error) {
//	return redigo.Int32Int32Map(reply, err)
//}

func IntMap(reply interface{}, err error) (map[string]int, error) {
	return redigo.IntMap(reply, err)
}

func StringMap(reply interface{}, err error) (map[string]string, error) {
	return redigo.StringMap(reply, err)
}

func Exists(key interface{}) (bool, error) {
	r := MainRds.Get()
	defer r.Close()
	return redigo.Bool(r.Do("EXISTS", key))
}

func Del(key interface{}) error {
	r := MainRds.Get()
	defer r.Close()
	_, e := r.Do("DEL", key)
	return e
}

/*
设置key的有效时间
*/
func Expire(expire int, key interface{}) (e error) {
	r := MainRds.Get()
	defer r.Close()
	_, e = r.Do("EXPIRE", key, expire)
	return

}

/*
MultiExpire 批量设置key的有效时间

	db: 数据库表ID
	expire:缓存失效时间(秒值)
	args:key的列表
*/
func MultiExpire(db, expire int, args ...interface{}) (e error) {
	if len(args) <= 0 {
		return
	}
	r := MainRds.Get()
	defer r.Close()
	if e := r.Send("MULTI"); e != nil {
		return e
	}
	for _, key := range args {
		if e := r.Send("EXPIRE", key, expire); e != nil {
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
设置key的到期时间

	db: 数据库表ID
	expireat:缓存失效的到期时间(unix 时间戳)
	key: 键值
*/
func Expireat(db int, expireat int64, key interface{}) (ret int, e error) {
	r := MainRds.Get()
	defer r.Close()
	ret, e = redigo.Int(r.Do("EXPIREAT", key, expireat))
	return

}

func Multi(db int, cmd func(con redigo.Conn) error) error {
	r := MainRds.Get()
	defer r.Close()
	if e := r.Send("MULTI"); e != nil {
		return e
	}
	if e := cmd(r); e != nil {
		r.Send("DISCARD")
		return e
	}
	if _, e := r.Do("EXEC"); e != nil {
		return e
	}
	return nil
}

type ScanResult struct {
	Cursor string
	Keys   []string
}

func KeyScan(db int, cursor string) (r ScanResult, e error) {
	rs := MainRds.Get()
	defer rs.Close()
	reply, e := redigo.Values(rs.Do("SCAN", cursor, "COUNT", 1000))
	if e != nil {
		return
	}
	if len(reply) == 2 {
		r.Cursor, e = String(reply[0], nil)
		r.Keys, e = Strings(reply[1], nil)
	}
	// fmt.Println(fmt.Sprintf("reply %v", reply))
	// e = ScanStruct(reply, &r)
	return
}

func KeyScanWithPattern(db int, cursor string, pattern string) (r ScanResult, e error) {
	rs := MainRds.Get()
	defer rs.Close()
	reply, e := redigo.Values(rs.Do("SCAN", cursor, "MATCH", pattern, "COUNT", 100))
	if e != nil {
		return
	}
	if len(reply) == 2 {
		r.Cursor, e = String(reply[0], nil)
		r.Keys, e = Strings(reply[1], nil)
	}
	return
}

/*
获取key的有效时间
*/
func TTL(db int, key interface{}) (expire int, e error) {
	r := MainRds.Get()
	defer r.Close()
	expire, e = redigo.Int(r.Do("TTL", key))
	return

}

func Keys(pattern string) (keys []string, e error) {
	r := MainRds.Get()
	defer r.Close()
	keys, e = redigo.Strings(r.Do("KEYS", pattern))
	return
}

func Publish(db int, channel, value interface{}) error {
	r := MainRds.Get()
	defer r.Close()
	_, err := r.Do("PUBLISH", channel, value)
	if err != nil {
		return err
	}
	return nil
}

// 生成redis分页
func BuildRange(cur, ps, total int) (int, int) {
	begin := 0
	if cur > 1 {
		begin = (cur - 1) * ps
	}
	end := begin + ps - 1
	if total > 0 && end >= total {
		end = total - 1
	}
	return begin, end
}
