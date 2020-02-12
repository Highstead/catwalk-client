package main

import (
	client "github.com/highstead/catwalk-client"
	log "github.com/sirupsen/logrus"
)

func main() {
	log.SetLevel(log.DebugLevel)
	svc := client.NewCatwalkClient()

	//models := []string{"sales_rollup_v7", "orders_rollup_v7"}
	//models := []string{"sales_rollup_v7"}
	models := []string{"marketing_activity_daily_v3"}
	result, err := svc.GetModels(models)
	if err == nil {
		log.Println("**Result**\n", result)
		return
	}
	log.WithError(err).Errorln("failed to make request")

}
