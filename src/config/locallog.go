package config

import (
	"github.com/cihub/seelog"
)

func InitLocalLog() error {
	filePath := "conf/seelog.xml"
	logger, err := seelog.LoggerFromConfigAsFile(filePath)

	if err != nil {
		return err
	} else {
		seelog.UseLogger(logger)
		seelog.Flush()
		return nil
	}
}
