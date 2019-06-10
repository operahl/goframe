package myredis

import (
	"errors"
	"fmt"
)

type ItemScore struct {
	Key   string
	Score float64
}

func ZCount(key, min, max interface{}) (count uint32, e error) {
	r := MainRds.Get()
	defer r.Close()
	n, e := Uint64(r.Do("ZCOUNT", key, min, max))
	if e != nil {
		return 0, errors.New(fmt.Sprintf("ZCOUNT error: %v", e.Error()))
	}
	return uint32(n), nil
}

func ZCard(key interface{}) (num uint64, e error) {
	r := MainRds.Get()
	defer r.Close()
	return Uint64(r.Do("ZCARD", key))
}

func ZscanAll(key interface{}) (map[string]string, error) {
	r := MainRds.Get()
	defer r.Close()
	reply, _ := Values(r.Do("ZSCAN", key, 0))
	if len(reply) != 2 {
		return nil, nil
	}
	var destin map[string]string
	destin, _ = StringMap(reply[1], nil)
	return destin, nil
}

//获取SortedSet的成员排名
func ZRank(key interface{}, id interface{}, asc bool) (rank int64, e error) {
	cmd := "ZRANK"
	if !asc {
		cmd = "ZREVRANK"
	}
	r := MainRds.Get()
	defer r.Close()
	return Int64(r.Do(cmd, key, id))
}

func ZRangeByScore(key interface{}, min, max interface{}) (items []ItemScore, e error) {
	r := MainRds.Get()
	defer r.Close()
	values, e := Values(r.Do("ZRANGEBYSCORE", key, min, max, "WITHSCORES"))
	if e != nil {
		return nil, errors.New(fmt.Sprintf("ZRANGEBYSCORE error: %v", e.Error()))
	}
	if e = ScanSlice(values, &items); e != nil {
		return nil, errors.New(fmt.Sprintf("ScanSlice error: %v", e.Error()))
	}
	return
}

func ZRevRangeByScore(key interface{}, min, max interface{}) (items []ItemScore, e error) {
	r := MainRds.Get()
	defer r.Close()
	values, e := Values(r.Do("ZREVRANGEBYSCORE", key, max, min, "WITHSCORES"))
	if e != nil {
		return nil, errors.New(fmt.Sprintf("ZREVRANGEBYSCORE error: %v", e.Error()))
	}
	if e = ScanSlice(values, &items); e != nil {
		return nil, errors.New(fmt.Sprintf("ScanSlice error: %v", e.Error()))
	}
	return
}

func ZRevRangeByScoreLimit(key interface{}, min, max interface{}, pageSize int) (items []ItemScore, e error) {
	r := MainRds.Get()
	defer r.Close()
	values, e := Values(r.Do("ZREVRANGEBYSCORE", key, max, min, "WITHSCORES", "LIMIT", 0, pageSize))
	if e != nil {
		return nil, errors.New(fmt.Sprintf("ZREVRANGEBYSCORE error: %v", e.Error()))
	}
	if e = ScanSlice(values, &items); e != nil {
		return nil, errors.New(fmt.Sprintf("ScanSlice error: %v", e.Error()))
	}
	return
}

//fenye
func ZREVRangeWithScoresPS(key interface{}, cur int, ps int) (items []ItemScore, total int, e error) {
	return zRangeWithScoresPS(key, cur, ps, false)
}
func ZRangeWithScoresPS(key interface{}, cur int, ps int) (items []ItemScore, total int, e error) {
	return zRangeWithScoresPS(key, cur, ps, true)
}
func ZREVRangeWithScores(key interface{}, start, end int) (items []ItemScore, total int, e error) {
	return zRangeWithScores(key, start, end, false)
}
func ZRangeWithScores(key interface{}, start, end int) (items []ItemScore, total int, e error) {
	return zRangeWithScores(key, start, end, true)
}
func ZRangeWithScoresOrder(key interface{}, start, end int, isAsc bool) (items []ItemScore, total int, e error) {
	return zRangeWithScores(key, start, end, isAsc)
}

//分页获取SortedSet的ID集合
func ZRangePS(key interface{}, cur int, ps int, asc bool, results interface{}) (total int, e error) {
	start, end := BuildRange(cur, ps, total)
	return ZRange(key, start, end, asc, results)
}

//获取SortedSet的ID集合
func ZRange(key interface{}, start, end int, asc bool, results interface{}) (total int, e error) {
	cmd := "ZRANGE"
	if !asc {
		cmd = "ZREVRANGE"
	}
	r := MainRds.Get()
	defer r.Close()
	total, e = Int(r.Do("ZCARD", key))
	if e != nil {
		return 0, errors.New(fmt.Sprintf("ZCARD error: %v", e.Error()))
	}
	values, e := Values(r.Do(cmd, key, start, end))
	if e != nil {
		return 0, e
	}
	if e = ScanSlice(values, results); e != nil {
		return 0, errors.New(fmt.Sprintf("ScanSlice error: %v", e.Error()))
	}
	return
}

//获取SortedSet的ID集合
func ZRangeOk(key interface{}, start, end int, asc bool, results interface{}) (e error) {
	cmd := "ZRANGE"
	if !asc {
		cmd = "ZREVRANGE"
	}
	r := MainRds.Get()
	defer r.Close()
	values, e := Values(r.Do(cmd, key, start, end))
	if e != nil {
		return e
	}
	if e = ScanSlice(values, results); e != nil {
		return errors.New(fmt.Sprintf("ScanSlice error: %v", e.Error()))
	}
	return
}

//分页获取带积分的SortedSet值
func zRangeWithScoresPS(key interface{}, cur int, ps int, asc bool) (items []ItemScore, total int, e error) {
	if ps > 100 {
		ps = 100
	}
	start, end := BuildRange(cur, ps, total)
	return zRangeWithScores(key, start, end, asc)
}

//获取带积分的SortedSet值
func zRangeWithScores(key interface{}, start, end int, asc bool) (items []ItemScore, total int, e error) {
	cmd := "ZRANGE"
	if !asc {
		cmd = "ZREVRANGE"
	}
	r := MainRds.Get()
	defer r.Close()
	total, e = Int(r.Do("ZCARD", key))
	if e != nil {
		return nil, 0, errors.New(fmt.Sprintf("ZCARD error: %v", e.Error()))
	}
	items = make([]ItemScore, 0, 100)
	values, e := Values(r.Do(cmd, key, start, end, "WITHSCORES"))
	if e != nil {
		return nil, 0, errors.New(fmt.Sprintf("%v error: %v", cmd, e.Error()))
	}
	if e = ScanSlice(values, &items); e != nil {
		return nil, 0, errors.New(fmt.Sprintf("ScanSlice error: %v", e.Error()))
	}
	return
}

func ZRevRange(key interface{}, start, end int)(items []ItemScore, e error){
	r := MainRds.Get()
	items = make([]ItemScore, 0, 100)
	values, e := Values(r.Do("ZREVRANGE", key, start, end, "WITHSCORES"))
	if e != nil {
		return nil,  errors.New(fmt.Sprintf("ZREVRANGE error: %v", e.Error()))
	}
	if e = ScanSlice(values, &items); e != nil {
		return nil,  errors.New(fmt.Sprintf("ScanSlice error: %v", e.Error()))
	}
	return
}

/*
批量添加到sorted set类型的表中

	args: 必须是<key,score,id>的列表
*/
func ZAdd(args ...interface{}) error {
	if len(args)%3 != 0 {
		return errors.New("invalid arguments number")
	}
	r := MainRds.Get()
	defer r.Close()
	_, e := r.Do("ZADD", args[0], args[1], args[2])
	return e
}

/*
interface{}{key1,score1,value1,key2,score2,value2}
*/
func ZMAddMoreKey(args ...interface{}) error {
	if len(args)%3 != 0 {
		return errors.New("invalid arguments number")
	}
	r := MainRds.Get()
	defer r.Close()
	/*
		if e := r.Send("MULTI"); e != nil {
			return e
		}
		for i := 0; i < len(args); i += 3 {
			if e := r.Send("ZADD", args[i], args[i+1], args[i+2]); e != nil {
				r.Send("DISCARD")
				return e
			}
		}
		if _, e := r.Do("EXEC"); e != nil {
			r.Send("DISCARD")
			return e
		}
	*/
	for i := 0; i < len(args); i += 3 {
		r.Send("ZADD", args[i], args[i+1], args[i+2])
	}
	err := r.Flush()
	return err
}

/*
批量添加到sorted set类型的表中

	db: 数据库表ID
	opt: 可选参数，必须是NX|XX|CH|INCR|""中的一个
	args: 必须是<key,score,id>的列表
*/
func ZAddOpt(opt string, args ...interface{}) error {
	if len(args)%3 != 0 {
		return errors.New("invalid arguments number")
	}
	r := MainRds.Get()
	defer r.Close()
	if e := r.Send("MULTI"); e != nil {
		return e
	}
	for i := 0; i < len(args); i += 3 {
		if e := r.Send("ZADD", args[i], opt, args[i+1], args[i+2]); e != nil {
			r.Send("DISCARD")
			return e
		}
	}
	if _, e := r.Do("EXEC"); e != nil {
		return e
	}
	return nil
}

func ZRem(args ...interface{}) (affected int, e error) {
	r := MainRds.Get()
	defer r.Close()
	affected, e = Int(r.Do("zrem", args[0], args[1]))
	return
}

/*
ZRem批量删除sorted set表中的元素

参数:
	args: 必须是<key,id>的列表
返回值：
	affected: 每条命令影响的行数(现在没有加这个,其实可以通过r.Receive()这个来获取的 bywcd)
*/
func ZMRemMoreKey(args ...interface{}) (e error) {
	if len(args)%2 != 0 {
		return errors.New("invalid arguments number")
	}
	r := MainRds.Get()
	defer r.Close()
	/*
		if e := r.Send("MULTI"); e != nil {
			return nil, e
		}
		for i := 0; i < len(args); i += 2 {
			if e := r.Send("ZREM", args[i], args[i+1]); e != nil {
				r.Send("DISCARD")
				return nil, e
			}
		}
		if replies, e := Values(r.Do("EXEC")); e != nil {
			return nil, e
		} else {
			affected = []int{}
			if e := ScanSlice(replies, &affected); e != nil {
				return nil, e
			}
		}
	*/

	for i := 0; i < len(args); i += 2 {
		r.Send("ZREM", args[i], args[i+1])
	}
	e = r.Flush()
	return
}

func ZExists(key interface{}, id interface{}) (bool, error) {
	r := MainRds.Get()
	defer r.Close()
	_, e := Int64(r.Do("ZSCORE", key, id))
	switch e {
	case nil:
		return true, nil
	case ErrNil:
		return false, nil
	default:
		return false, e
	}
}

func ZScore(key interface{}, item interface{}) (score int64, e error) {
	r := MainRds.Get()
	defer r.Close()
	return Int64(r.Do("ZSCORE", key, item))
}

//批量获取有序集合的元素的得分
func ZMultiScore(key interface{}, items ...interface{}) (scores map[interface{}]int64, e error) {
	r := MainRds.Get()
	defer r.Close()
	for _, id := range items {
		if e := r.Send("ZSCORE", key, id); e != nil {
			return nil, e
		}
	}
	r.Flush()
	scores = make(map[interface{}]int64, len(items))
	for _, id := range items {
		score, e := Int64(r.Receive())
		switch e {
		case nil:
			scores[id] = score
		case ErrNil:
		default:
			return nil, e
		}
	}
	return scores, nil
}

//批量判断是否是有序集合中的元素
func ZMultiIsMember(key interface{}, items map[interface{}]bool) error {
	r := MainRds.Get()
	defer r.Close()
	ids := make([]interface{}, 0, len(items))
	for id, _ := range items {
		if e := r.Send("ZSCORE", key, id); e != nil {
			return e
		}
		ids = append(ids, id)
	}
	r.Flush()
	for _, id := range ids {
		_, e := Int64(r.Receive())
		switch e {
		case nil:
			items[id] = true
		case ErrNil:
			items[id] = false
		default:
			return e
		}
	}
	return nil
}

/*
根据score 获取有序集 ZREVRANGEBYSCORE min <=score < max  按照score 从大到小排序, ps 获取条数
*/
func ZREVRangeByScoreWithScores(key interface{}, min, max int64, ps int) (items []ItemScore, e error) {
	return zRangeByScoreWithScores(key, min, max, ps, false)
}

/*
根据score 获取有序集 ZRANGEBYSCORE min <score <= max  按照score 从小到大排序, ps 获取条数
*/
func ZRangeByScoreWithScores(key interface{}, min, max int64, ps int) (items []ItemScore, e error) {
	return zRangeByScoreWithScores(key, min, max, ps, true)
}

//根据积分的SortedSet值
func zRangeByScoreWithScores(key interface{}, min, max int64, ps int, asc bool) (items []ItemScore, e error) {
	var s1, s2, cmd string
	if asc {
		s1 = "(" + ToString(min)
		s2 = ToString(max)
		cmd = "ZRANGEBYSCORE"
	} else {
		s1 = "(" + ToString(max)
		s2 = ToString(min)
		cmd = "ZREVRANGEBYSCORE"
	}
	r := MainRds.Get()
	defer r.Close()
	items = make([]ItemScore, 0, 100)
	values, e := Values(r.Do(cmd, key, s1, s2, "WITHSCORES", "LIMIT", 0, ps))
	if e != nil {
		return nil, errors.New(fmt.Sprintf("%v error: %v", cmd, e.Error()))
	}
	if e = ScanSlice(values, &items); e != nil {
		return nil, errors.New(fmt.Sprintf("ScanSlice error: %v", e.Error()))
	}
	return
}

//根据积分的SortedSet值
func ZRangeByScoreWithScoresStr(key interface{}, min, max string, ps int, asc bool, isScore bool) (items []ItemScore, e error) {
	var s1, s2, cmd string
	if asc {
		s1 = "(" + min
		s2 = max
		cmd = "ZRANGEBYSCORE"
	} else {
		s1 = "(" + max
		s2 = min
		cmd = "ZREVRANGEBYSCORE"
	}
	r := MainRds.Get()
	defer r.Close()
	items = make([]ItemScore, 0, 100)
	var values []interface{}

	if isScore == true {
		values, e = Values(r.Do(cmd, key, s1, s2, "WITHSCORES", "LIMIT", 0, ps))
	} else {
		values, e = Values(r.Do(cmd, key, s1, s2, "LIMIT", 0, ps))
	}
	if e != nil {
		return nil, errors.New(fmt.Sprintf("%v error: %v", cmd, e.Error()))
	}
	if e = ScanSlice(values, &items); e != nil {
		return nil, errors.New(fmt.Sprintf("ScanSlice error: %v", e.Error()))
	}
	return
}

/*
移除有序集 key 中，所有 score 值介于 min 和 max 之间(包括等于 min 或 max )的成员
*/
func ZRemRangeByScore(key interface{}, min, max interface{}) error {
	r := MainRds.Get()
	defer r.Close()
	_, e := r.Do("ZREMRANGEBYSCORE", key, min, max)
	return e
}

/*
合并多个有序集合，其中权重weights 默认为1 ，AGGREGATE 默认使用sum
ZUNIONSTORE destination numkeys key [key ...] [WEIGHTS weight [weight ...]] [AGGREGATE SUM|MIN|MAX]
dest_key：合并目标key
keys: 带合并的keys集合 <key> 的列表
expire: 有效时间 （秒值）
aggregate: 聚合方式： SUM | MIN | MAX
*/
func ZUnionSrore(dest_key interface{}, expire int, keys []interface{}, weights []interface{}, aggregate string) error {
	if len(keys) != len(weights) || len(keys) <= 0 {
		return errors.New("invalid numbers of keys and weights")
	}
	args := make([]interface{}, 0, 2*len(keys)+10)
	args = append(args, dest_key, len(keys))
	args = append(args, keys...)
	args = append(args, "WEIGHTS")
	args = append(args, weights...)
	args = append(args, "AGGREGATE", aggregate)
	r := MainRds.Get()
	defer r.Close()
	if _, e := r.Do("ZUNIONSTORE", args...); e != nil {
		return e
	}
	_, e := r.Do("EXPIRE", dest_key, expire)
	return e
}

//ZIsMember判断是否是有序集合的成员
func ZIsMember(key interface{}, item interface{}) (isMember bool, e error) {
	r := MainRds.Get()
	defer r.Close()
	_, e = Float64(r.Do("ZSCORE", key, item))
	switch e {
	case nil:
		return true, nil
	case ErrNil:
		return false, nil
	default:
		return false, e
	}
}
func ZIncrBy(key interface{}, increment interface{}, id interface{}) (score int64, e error) {
	r := MainRds.Get()
	defer r.Close()
	return Int64(r.Do("ZINCRBY", key, increment, id))
}
