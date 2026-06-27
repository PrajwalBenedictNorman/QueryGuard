package config

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"gopkg.in/yaml.v3"
)

type RequestBody struct {
    Message   string `yaml:"message"`
    SessionId string `yaml:"session_id"`
}


type Config struct{
		Server struct{
			Port int `yaml:"port"`
		} `yaml:"server"`
		Metrics struct {
			CUSUM struct{
				K float32 `yaml:"k"`
				H float32 `yaml:"h"`
				Mean float32 `yaml:"mean"`
				Variance float32 `yaml:"variance"`
			}`yaml:"cusum"`
			EWMA struct{
				Lambda float32 `yaml:"lambda"`
				StdDev float32 `yaml:"std_dev"`
				Mean float32 `yaml:"mean"`
				Threshold float32 `yaml:"threshold"`
			}`yaml:"ewma"`
		}`yaml:"metrics"`
		Target struct{
			URL string `yaml:"url"`
			Method string `yaml:"method"`
		}`yaml:"target"`
		Request struct{
			Header  map[string]string `yaml:"header"`
			Body RequestBody `yaml:"body"`
		}`yaml:"request"`
		Response struct{
				Path string `yaml:"path"`
		}`yaml:"response"`		
		Baseline struct{
			Model string `yaml:"model"`
			Threshold float64 `yaml:"threshold"`
			Store string `yaml:"store"`
		}`yaml:"baseline"`

		Cache struct {
			Size int `yaml:"size"`
		}`yaml:"cache"`

}

func Load() (*Config,error){
	dir, _ := os.Getwd()
	fmt.Println("Looking for config at:", dir+"/config.yaml")
	data,err :=os.ReadFile(dir+"/config.yaml")
	if err !=nil{
		log.Fatal("No config file found")
		return  nil,err
	}
	var cfg Config
	if err := yaml.Unmarshal(data,&cfg);err !=nil{
		return  nil,err
	}
	// Url finding
	if cfg.Target.URL == ""{
		fmt.Print("Enter the model url")
		reader :=bufio.NewReader(os.Stdin)
		url,_:=reader.ReadString('\n')
		cfg.Target.URL=strings.TrimSpace(url)
	}

	// Loading evn values
	for key, value := range cfg.Request.Header {
    cfg.Request.Header[key] = os.Expand(value, os.Getenv)
}
	return &cfg,nil
}