package common

import "time"

func GetTimeUnix() int64{
	return time.Now().Unix()
}


