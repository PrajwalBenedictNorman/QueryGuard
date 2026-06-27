package adapter

import (
	"QueryGuard/internal/config"
	"bytes"
	"encoding/json"
	"net/http"
	"strings"
	"fmt"
)

type Adapter struct{
	Config *config.Config
	Client *http.Client
}

func NewAdapter(cfg *config.Config)*Adapter{
	client := &http.Client{}
	return &Adapter{Config: cfg,Client: client}
}

func ReplaceQuery(body *config.RequestBody ,query string)*config.RequestBody{
	if body.Message=="{{query}}"{
		body.Message=query
	}
	return body
}

func extractPath(data map[string]interface{},path string)(string,error){
		path=strings.TrimPrefix(path,"$.")
		keys :=strings.Split(path,".")
		var current interface{} = data
		for _, k := range keys {
    m, ok := current.(map[string]interface{})
    if !ok {
        return "", fmt.Errorf("invalid path")
    }
    val, exists := m[k]
    if !exists {
        return "", fmt.Errorf("key %s not found", k)
    }
    current = val 
}

s, ok := current.(string)
if !ok {
    return "", fmt.Errorf("value is not a string")
}
return s, nil
}
func (a *Adapter)Forward(query string)(string,error){
	body:=ReplaceQuery(&a.Config.Request.Body,query)
	jsonBody,_:=json.Marshal(body)
	req,_:=http.NewRequest(a.Config.Target.Method,a.Config.Target.URL,bytes.NewBuffer(jsonBody))
	for i,k:=range a.Config.Request.Header{
		req.Header.Set(i,k)
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := a.Client.Do(req)
	if err !=nil{
		return "",err
	}
	var response map[string]interface{}
	
	er := json.NewDecoder(resp.Body).Decode(&response)
	if er != nil {
    return "",er
	}
	result, err := extractPath(response, a.Config.Response.Path)
	if err != nil {
    return "", err
}
	return result,nil
}