package healthx

import (
	"fmt"
	"net"

	"git.edenfarm.id/edenlabs/edenlabs/opt"
)

type HealthCheckStatus struct {
	Redis   CheckStatus `json:"redis"`
	PubSub  CheckStatus `json:"pubsub"`
	MongoDB CheckStatus `json:"mongodb"`
	MySQL   CheckStatus `json:"mysql"`
}

type CheckStatus struct {
	Enabled   bool `json:"enabled"`
	Connected bool `json:"connected"`
}

func Check(opt *opt.Options) (helathCheckStatus HealthCheckStatus, err error) {
	var conn net.Conn

	// 1. Redis
	if opt.Config.Redis.Enabled {
		conn, err = net.Dial("tcp", fmt.Sprintf("%s:%d", opt.Config.Redis.Host, opt.Config.Redis.Port))
		if err != nil {
			helathCheckStatus.Redis.Enabled = true
			helathCheckStatus.Redis.Connected = false
			return
		}
		defer conn.Close()

		helathCheckStatus.Redis.Enabled = true
		helathCheckStatus.Redis.Connected = true
	} else {
		helathCheckStatus.Redis.Enabled = false
		helathCheckStatus.Redis.Connected = false
	}

	// 2. Mongo
	if opt.Config.Mongodb.Enabled {
		conn, err = net.Dial("tcp", fmt.Sprintf("%s:%d", opt.Config.Mongodb.Host, opt.Config.Mongodb.Port))
		if err != nil {
			helathCheckStatus.MongoDB.Enabled = true
			helathCheckStatus.MongoDB.Connected = false
			return
		}
		defer conn.Close()

		helathCheckStatus.MongoDB.Enabled = true
		helathCheckStatus.MongoDB.Connected = true
	} else {
		helathCheckStatus.MongoDB.Enabled = false
		helathCheckStatus.MongoDB.Connected = false
	}

	// 3. Pub/Sub
	// conn, err = net.Dial("tcp", fmt.Sprintf("%s:%d", opt.Config.PubSub, opt.Config.Redis.Port))
	// if err != nil {
	// 	return
	// }
	// defer conn.Close()

	// 4. Mysql
	if opt.Config.Database.Enabled {
		conn, err = net.Dial("tcp", fmt.Sprintf("%s:%d", opt.Config.Database.Write.Host, opt.Config.Database.Write.Port))
		if err != nil {
			helathCheckStatus.MySQL.Enabled = true
			helathCheckStatus.MySQL.Connected = false
			return
		}
		defer conn.Close()

		helathCheckStatus.MySQL.Enabled = true
		helathCheckStatus.MySQL.Connected = true
	} else {
		helathCheckStatus.MySQL.Enabled = false
		helathCheckStatus.MySQL.Connected = false
	}

	return
}
