package main

import (
"fmt"
"strings"
"io/ioutil"
)

func main(){
file,err := ioutil.ReadFile("SMALL_SKIRT.exe")
if err != nil{
	panic(err)
}
//fmt.Sprintf("%s",string(file))
nr := strings.NewReader(string(file))
b,err := ioutil.ReadAll(nr)
if err != nil{ panic(err)}
fmt.Printf("%s",b)
}
