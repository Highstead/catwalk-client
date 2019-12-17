package main

import (
	"fmt"

	client "github.com/highstead/catwalk-client"
	log "github.com/sirupsen/logrus"
)

func main() {
	log.SetLevel(log.DebugLevel)
	svc := client.NewCatwalkClient()

	result, err := svc.GetModel("sales_rollup_v7")
	if err == nil {
		fmt.Println("**Result**\n", result)
		return
	}
	log.WithError(err).Errorln("failed to make request")

}
