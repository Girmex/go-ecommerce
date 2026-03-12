package main

import (
	"log"

	"github.com/Girmex/go-ecommerce/config"
	"github.com/Girmex/go-ecommerce/internal/api"
)

func main(){

	cfg, err := config.SetupEnv()

	if err != nil{
     log.Fatalf("config file is not loaded correctly %v\n",err)
	}
	api.StartServer(cfg)
}