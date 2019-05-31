package limit

import (
	"fmt"
	"time"

	"encoding/hex"

	"github.com/jiachuhuang/concurrentcache"
)

var cache, _ = concurrentcache.NewConcurrentCache(256, 10240)

func AccessLimit(key string, limit_num int, limit_time int) bool {
	old_num, _ := cache.Get(key)
	var old_ok_num = 0
	if old_num != nil {
		old_ok_num = old_num.(int)
	}
	fmt.Println("old:")
	fmt.Println(old_ok_num)
	if old_ok_num >= limit_num {
		//c.AbortWithStatus(401)
		//c.JSON(401, gin.H{})
		return true
	} else {
		new_num := old_ok_num + 1
		cache.Set(key, new_num, time.Duration(limit_time)*time.Second)
	}
	return false
}

func AeskeyCache(aeskey string) (key []byte) {
	key_tmp, _ := cache.Get("aeskey")
	if key_tmp == nil {
		key, _ = hex.DecodeString(aeskey)
		cache.Set("aeskey", key, time.Duration(100000))
	} else {
		key = key_tmp.([]byte)
	}
	return
}
