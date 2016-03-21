package helpers

import (
    "fmt"
    "encoding/json"
    "os"
)
//Configuration struct mapped to configuration json:
//{
//  "GithubPersonalToken": "whatever",
//  "Port": 3000    
//}
type Configuration struct{
    GithubPersonalToken string
    Port                int
}

//LoadConfigFromFile creating config struct from file with specified path.
func LoadConfigFromFile(path string) (*Configuration, error){
    conf:= &Configuration{}
    file, err := os.Open(path)
    if err != nil{
        return nil,err
    }
    decoder := json.NewDecoder(file)
    err = decoder.Decode(conf)
    return conf, nil
}
//SaveConfigToFile saving modified struct to file with specified path.
func (conf *Configuration) SaveConfigToFile(path string) error{
    file, err := os.Open(path)
    if err !=nil {
        fmt.Println("error:",err)
               return err
    }
    encoder := json.NewEncoder(file)
    err = encoder.Encode(conf)
    if err != nil {
        fmt.Println("error:", err)
        return err
    }
    return nil
    
}
