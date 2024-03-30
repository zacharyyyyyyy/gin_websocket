package tools

import "gin_websocket/lib/logger"

func RecoverFunc() {
	if v := recover(); v != nil {
		logger.Runtime.Error(v.(string))
	}
}
