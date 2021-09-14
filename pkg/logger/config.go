package logger

import "github.com/ddelizia/hasura-saas/pkg/env"

func ConfigLogLevel() string {
	return env.GetString("logger.level")
}
