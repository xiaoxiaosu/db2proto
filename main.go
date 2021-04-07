package main

import (
	"log"
	"xiaoxiaosu.com/db2pb/pkg/cmd"
)

func main() {
	err := cmd.Execute()
	if err != nil {
		log.Fatalf("cmd.Execute err:%v", err)
	}
}