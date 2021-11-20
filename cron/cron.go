package cron

import (
	"github.com/jasonlvhit/gocron"
	"github.com/xhyonline/xdq/services"
)

func Init() {
	_ = gocron.Every(1).Second().Do(services.ScanBucketForReady)
	_ = gocron.Every(1).Second().Do(services.ConsumeReadyJob)
	<-gocron.Start()
}
