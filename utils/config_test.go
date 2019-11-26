package utils

import (
	"midmsg/log"
	"testing"
)

func TestGetSubscribeByKey(t *testing.T) {
	list := GetSubscribeByKey("test.service")
	log.Info(list)
}