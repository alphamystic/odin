package utils

import (
  "fmt"
  "errors"
  "github.com/dgrijalva/jwt-go"
)


const (
  TOKENKEY = "This is the token Key"
  NI = "NOT IMPLEMENTED"
)

var (
  NotImplemented = errors.New("Requested data is not inmplemeted and set to null")
)


var IsNI = func(s string) bool{
  if s == NI {
    return true
  }
  return false
}

func ArrayToToken(arr []string) (string,error){
  token := jwt.New(jwt.SigningMethodHS256)
  claims := token.Claims.(jwt.MapClaims)
  claims["data"] = arr
  // Sign the token with your secret key
  tokenString, err := token.SignedString(TOKENKEY)
  if err != nil {
    return "",err
  }
  return tokenString,nil
}

func TokenToArray(tokenString string)([]string,error){
  if IsNI(tokenString){
    return nil,NotImplemented
  }
  var datum []string
  // Parse and verify the JWT
  token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
    return []byte(TOKENKEY), nil
  })
  if err != nil {
    return datum,err
  }
  // Extract the data from the JWT
  if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
    data := claims["data"].([]interface{})
    for _, d := range data {
        datum = append(datum, d.(string))
    }
  }
  return datum,nil
}

func MultipleToToken(data interface{}) (string,error){
  token := jwt.New(jwt.SigningMethodHS256)
  claims := token.Claims.(jwt.MapClaims)
  claims["data"] = data
  tokenString,err := token.SignedString(TOKENKEY)
  if err != nil {
    return "",fmt.Errorf("Error creating token: %q",err)
  }
  return tokenString,nil
}

func TokenToMultiple(tokenString string) ([]map[string]string,error){
  token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
    return []byte(TOKENKEY), nil
  })
  if err != nil {
    return nil,err
  }
  var datum []map[string]string
  // Extract the data from the JWT
  if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
    data := claims["data"].([]interface{})
    for _, d := range data {
      for key,val := range d.(string) {
        datum = append(datum,map[string]string{key,val})
      }
    }
  }
  return datum,nil
}

func TokenToString(tokenString string) (string,error){
  if IsNI(tokenString){
    return "",NotImplemented
  }
  token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
    return []byte(TOKENKEY), nil
  })
  if err != nil {
    return "",err
  }
  var datum string
  // Extract the data from the JWT
  if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
    data := claims["data"].(interface{})
    datum = data.(string)
  }
  return datum,nil
}


func TokenToData(tokenString string) (data interface{},err error){
  if IsNI(tokenString){
    return data,NotImplemented
  }
  token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
    return []byte(TOKENKEY), nil
  })
  if err != nil {
    return nil,err
  }
  var datum interface{}
  // Extract the data from the JWT
  if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
    data := claims["data"].([]interface{})
    datum = data
  }
  return datum,nil
}
