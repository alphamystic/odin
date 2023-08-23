package main

import (
"os"
"fmt"
"log"
"encoding/hex"
)


func main(){
src,err := os.ReadFile("demo.exe")
if err != nil{ log.Fatal(err)}
dst := make([]byte,hex.EncodedLen(len(src)))

hex.Encode(dst,src)
fmt.Sprintf("%v\n",dst)
}
