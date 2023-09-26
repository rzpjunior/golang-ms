package global

import (
	"context"
	"fmt"
	"time"
)

// untuk melakukan pengecekan apakah key yang dikirimkan untuk melimit
// akses ada di redis atau tidak.
func IsLimit(key string, period time.Duration, maxRetry int) bool {
	var c int64
	dbredis := Setup.Common.Redisx
	if _, e := dbredis.Ping(context.TODO()); e == nil {
		//tidak ada di redis
		if !dbredis.CheckCacheByKey(context.TODO(), key) {
			dbredis.SetCache(context.TODO(), key, c+1, period)
			return true
		}

		dbredis.GetCache(context.TODO(), key, &c)
		if c != int64(maxRetry) {
			dbredis.SetCache(context.TODO(), key, c+1, period)
			return true
		}
	} else {
		fmt.Println("=======================================================================")
		fmt.Println("=======================================================================")
		fmt.Println(e)
		fmt.Println("=======================================================================")
		fmt.Println("=======================================================================")
	}
	return false
}

// untuk melakukan pengecekan apakah key yang dikirimkan untuk melimit
// akses ada di redis atau tidak.
func IsLimitReqOTP(key string, key2 string, period time.Duration, period2 time.Duration, maxRetry int, dataType string) int {
	var c int64
	dbredis := Setup.Common.Redisx

	//tidak ada di redis set ke redis
	if dbredis.CheckCacheByKey(context.TODO(), key) {
		dbredis.GetCache(context.TODO(), key, &c)
	}

	if c >= int64(maxRetry) {
		return 0
	}

	// saat regis, tidak terkena validasi menunggu 1 menit
	if dbredis.CheckCacheByKey(context.TODO(), key2) && dataType == "2" {
		dbredis.SetCache(context.TODO(), key, c+1, period)
		return 1
	}

	if !dbredis.CheckCacheByKey(context.TODO(), key2) {
		dbredis.SetCache(context.TODO(), key2, 1, period2)
		dbredis.SetCache(context.TODO(), key, c+1, period)
		return 1
	}

	return 2
}

// untuk mengecek apakah counter di redis sudah memenuhi max retry
// tanpa mengubah value di redis
func IsMaxLimit(key string, maxRetry int) bool {
	var c int64
	dbredis := Setup.Common.Redisx

	dbredis.GetCache(context.TODO(), key, &c)

	//sudah max jadi ditolak untuk coba lagi
	if c == int64(maxRetry) {
		return false
	}

	return true
}

// resetCache : handler to reset cache
func resetCache(key string) (e error) {
	dbredis := Setup.Common.Redisx

	e = dbredis.DeleteAll(context.TODO())
	return e
}
