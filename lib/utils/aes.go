package utils

import (
  "io"
  "fmt"
  "bytes"
  "errors"
  "crypto/aes"
  "crypto/rand"
  "crypto/cipher"
)

func pad(buf []byte) []byte{
  length := len(buf)
  padding := aes.BlockSize - (length % aes.BlockSize)
  if padding == 0{
    padding = aes.BlockSize
  }
  padded := make([]byte,length + padding)
  copy(padded,buf)
  copy(padded[length:],bytes.Repeat([]byte{byte(padding)},padding))
  return padded
}

func unpad(buf []byte) []byte{
  padding := int(buf[len(buf)-1])
  return buf[:len(buf)-padding]
}

type RansomSim struct{
  Data []byte
  Key []byte
  CipherText []byte
}

func (rs *RansomSim) Encrypt() (error){
  var (
		iv         []byte
		block      cipher.Block
		mode       cipher.BlockMode
		err        error
  )
  if block,err = aes.NewCipher(rs.Key);err != nil{
    return fmt.Errorf("Error creatng block (Encrypt). %v",err)
  }
  iv = make([]byte,aes.BlockSize)
  if _,err := io.ReadFull(rand.Reader,iv); err != nil{
    return fmt.Errorf("Error reading unique blocksize. %v",err)
  }
  mode = cipher.NewCBCEncrypter(block,iv)
  rs.Data = pad(rs.Data)
  rs.CipherText = make([]byte,aes.BlockSize+len(rs.Data))
  copy(rs.CipherText,iv)
  mode.CryptBlocks(rs.CipherText[aes.BlockSize:],rs.Data)
  return nil
}

func (rs *RansomSim) Decrypt() (error){
  var (
		iv         []byte
		block      cipher.Block
		mode       cipher.BlockMode
		err        error
  )
  if len(rs.CipherText) < aes.BlockSize {
    return errors.New("Invalid cyphertext length: too short.")
  }
  if len(rs.CipherText) % aes.BlockSize != 0 {
    return errors.New("Invalid ciphertext length: Not a multiple of blocksize")
  }
  iv = rs.CipherText[:aes.BlockSize]
  rs.CipherText = rs.CipherText[aes.BlockSize:]
  if block,err = aes.NewCipher(rs.Key); err != nil{
    return fmt.Errorf("Error creatng block.(Decrypt): %v",err)
  }
  mode = cipher.NewCBCDecrypter(block,iv)
  rs.Data = make([]byte,len(rs.CipherText))
  mode.CryptBlocks(rs.Data,rs.CipherText)
  rs.Data = unpad(rs.Data)
  return nil
}

func (rs *RansomSim) RunSimulator(action bool) error{
  if _,err := io.ReadFull(rand.Reader,rs.Key); err != nil {
    return fmt.Errorf("Error reading key. %v",err)
  }
  if action{
    //encrypt
    if err := rs.Encrypt();err != nil{
      return fmt.Errorf("[ENCRYPTION]: %v",err)
    }
    PrintTextInASpecificColorInBold("yellow",fmt.Sprintf("KEY        = %x\n",rs.Key))
  } else {
    //assume its a decrypt
    if err := rs.Decrypt(); err != nil{
      return fmt.Errorf("[DECRYPTION]: %v",err)
    }
    PrintTextInASpecificColorInBold("yellow","Successfully decrypted......")
  }
  return nil
}
