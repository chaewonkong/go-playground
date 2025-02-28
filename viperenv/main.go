package main

import (
	"log"

	"github.com/spf13/viper"
)

func main() {
	v := viper.New()
	v.AutomaticEnv()

	p := v.GetString("SERVER_PORT")
	s := v.GetString("SECRET")
	log.Printf("port: %s, secret: %s\n", p, s)
}

/*
Pros and cons

os.GetEnv와 뭐가 다른가?
*/
