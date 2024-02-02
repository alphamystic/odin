package utils

import (
  "errors"
  "github.com/dgrijalva/jwt-go"
)


const (
  NI = "NOT IMPLEMENTED"
)

var (
  TOKENKEY = []byte("This is the token Key")
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

func TokenToArray(token string)([]string,error){
  if IsNI(token){
    return nil,NotImplemented
  }
  // Parse and verify the JWT
  token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
    return []byte(TOKENKEY), nil
  })
  if err != nil {
    return nil,err
  }
  var datum []string
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
  claims := toke.Claims.(jwt.MapClaims)
  claims["data"] = data
  tokenString,err := token.SignedString(TOKENKEY)
  if err != nil {
    return nil,fmt.Errorf("Error creating token: %q",err)
  }
  return tokenString,nil
}

func TokenToMultiple(token string) ([]map[string]string,error){
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
      for key,val := range d {
        datum = append(datum,map[string]string{key,val})
      }
    }
  }
  return datum,nil
}

func TokenToString(token string) (string,error){
  if IsNI(token){
    return nil,NotImplemented
  }
  token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
    return []byte(TOKENKEY), nil
  })
  if err != nil {
    return nil,err
  }
  var datum string
  // Extract the data from the JWT
  if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
    data := claims["data"].([]interface{})
    datum = data.(string)
  }
  return datum,nil
}


func TokenToData(tokenString string) (data interface{},err error){
  if IsNI(token){
    return nil,NotImplemented
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
