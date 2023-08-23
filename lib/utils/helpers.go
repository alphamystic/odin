package utils
/*
  * Contains helper functions
*/

import (
  "os"
  "io"
  "log"
  "fmt"
  "net"
  "time"
  "os/user"
  "errors"
  "regexp"
  "strings"
  "strconv"
  "runtime"
	"unicode"
  "math/rand"
)

func ReturnErrorPlusMessage( e error,text string) error{
  if e != nil{
    return fmt.Errorf(text + fmt.Sprintf("%s",e))
  }
  return fmt.Errorf("Error can not be nil")
}

func RemoveElementFromArray(arr []string, val string) []string {
  index := -1
  for i, v := range arr {
    if v == val {
      index = i
      break
    }
  }
  if index == -1 {
    return arr // value not found
  }
  return append(arr[:index], arr[index+1:]...)
}

func GetCurrentOS() string{ return runtime.GOOS }

func GetUser()(*user.User,error){
  cur,err := user.Current()
  if err != nil{ return nil,err}
  return cur,nil
}

func StringToInt(str string)(int){
  val,err := strconv.Atoi(str)
  if err != nil{
    Logerror(err)
    return 0
  }
  return val
}

// learn go generics and avoid such issues
//check if int is in array
func ArrayContainsInt(arr []int,element int) bool{
  for _,e := range arr {
    if e == element{
      return true
    }
  }
  return false
}

func ArrayContainsString(arr []string, element string) bool{
  for _,e := range arr {
    if e == element{
      return true
    }
  }
  return false
}
func IntToString(val int) string{  return strconv.Itoa(val) }

//Check if a string is empty returns True if string is a string
func CheckifStringIsEmpty(data string) bool{
  if len(strings.TrimSpace(data)) == 0{
    return false
  }
  if len(data) == 0{
    return false
  }
  return true
}

func TrueRand(len int) string{
  bytes := make([]byte,len)
  for i := 0; i < len; i++{
    bytes[i] = byte(randInt(97,122))
  }
  if !CheckifStringIsEmpty(string(bytes)){
    TrueRand(len)
  }
  return string(bytes)
}

func randInt(min int, max int) int {
  return min + rand.Intn(max-min)
}

func RandString(length int) string{
  var output strings.Builder
  rand.Seed(time.Now().Unix())
  charset := []rune("QWERTYUIOPLKJHGFDSAZXCVBNM123456789qwertyuioplkjhgfdsazxcvbnm")
  for i := 0; i < length; i++{
    random := rand.Intn(len(charset))
    randomChar := charset[random]
    output.WriteRune(randomChar)
  }
  id := output.String()
  id = strings.ToUpper(id)
  if !CheckifStringIsEmpty(id){
    RandString(length)
  }
  return id
}

//Retrns a random string with numbers and letters (caps on)
func RandNoLetter(length int) string{
  var output strings.Builder
  rand.Seed(time.Now().Unix())
  charset := []rune("QWERTYUIOPLKJHGFDSAZXCVBNM123456789")
  for i := 0; i < length; i++{
    random := rand.Intn(len(charset))
    randomChar := charset[random]
    output.WriteRune(randomChar)
  }
  id := output.String()
  id = strings.ToUpper(id)
  if !CheckifStringIsEmpty(id){
    RandNoLetter(length)
  }
  return id
}

//Returns A Random letters
func RandLetters(length int) string{
  var output strings.Builder
  rand.Seed(time.Now().Unix())
  charset := []rune("qwertyuioplkjhgfdsazxcvbnmQWERTYUIOPLKJHGFDSAZXCVBBNM")
  for i := 0; i < length; i++{
    random := rand.Intn(len(charset))
    randomChar := charset[random]
    output.WriteRune(randomChar)
  }
  id := output.String()
  if !CheckifStringIsEmpty(id){
    RandLetters(length)
  }
  return id
}

//Returns a random number in string format
func RandNo(length int) string{
  var output strings.Builder
  rand.Seed(time.Now().Unix())
  charset := []rune("1234567890")
  for i := 0; i < length; i++{
    random := rand.Intn(len(charset))
    randomChar := charset[random]
    output.WriteRune(randomChar)
  }
  id := output.String()
  if !CheckifStringIsEmpty(id){
    RandNo(length)
  }
  return id
}

//Log and error to file Allows format string input
func LogErrorToFile(name string,text ...interface{}) error{
  name = "./.data/logs/"+name+".log"
  f,err := os.OpenFile(name,os.O_RDWR|os.O_CREATE|os.O_APPEND,0666)
  if err != nil{
    return err
  }
  defer f.Close()
  writer := io.MultiWriter(os.Stdout,f)
  log.SetOutput(writer)
  log.Println(text)
  return nil
}

func CheckIfStringIsIp(name string) bool{
  ip := net.ParseIP(name)
  if ip != nil{
    return true
  }
  return false
}

func CheckIfStringIsDomainName(name string)bool{
	// Regular expression to match domain names
	domainRegex := regexp.MustCompile(`^(([a-zA-Z0-9]|[a-zA-Z0-9][a-zA-Z0-9\-]*[a-zA-Z0-9])\.)*([A-Za-z0-9]|[A-Za-z0-9][A-Za-z0-9\-]*[A-Za-z0-9])$`)
	if domainRegex.MatchString(name) {
		return true
	}
	return false
}

func RemoveStringDuplicates(strings []string) []string {
	// Create a map to store unique strings
	unique := make(map[string]bool)
	// Iterate through the input array
	for _, s := range strings {
		unique[s] = true
	}
	// Create a new slice to store the unique strings
	result := []string{}
	// Iterate through the map and append the keys (unique strings) to the result slice
	for key := range unique {
		result = append(result, key)
	}
	return result
}


func checkEmailValid(email string) error {
	// check email syntax is valid
	//func MustCompile(str string) *Regexp
	emailRegex, err := regexp.Compile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	if err != nil {
		fmt.Println(err)
		return errors.New("sorry, something went wrong")
	}
	rg := emailRegex.MatchString(email)
	if rg != true {
		return errors.New("Email address is not a valid syntax, please check again")
	}
	// check email length
	if len(email) < 4 {
		return errors.New("Email length is too short")
	}
	if len(email) > 253 {
		return errors.New("Email length is too long")
	}
	return nil
}

func checkEmailDomain(email string) error {
	i := strings.Index(email, "@")
	host := email[i+1:]
	// func LookupMX(name string) ([]*MX, error)
	_, err := net.LookupMX(host)
	if err != nil {
		err = errors.New("Could not find email's domain server, please chack and try again")
		return err
	}
	return nil
}


func checkUsernameCriteria(username string) error {
	// check username for only alphaNumeric characters
	var nameAlphaNumeric = true
	for _, char := range username {
		// func IsLetter(r rune) bool, func IsNumber(r rune) bool
		if !unicode.IsLetter(char) && !unicode.IsNumber(char) {
			nameAlphaNumeric = false
		}
	}
	if nameAlphaNumeric != true {
		// func New(text string) error
		return errors.New("Username must only contain letters and numbers")
	}
	// check username length
	var nameLength bool
	if 5 <= len(username) && len(username) <= 50 {
		nameLength = true
	}
	if nameLength != true {
		return errors.New("Username must be longer than 4 characters and less than 51")
	}
	return nil
}

func checkPasswordCriteria(password string) error {
	var err error
	// variables that must pass for password creation criteria
	var pswdLowercase, pswdUppercase, pswdNumber, pswdSpecial, pswdLength, pswdNoSpaces bool
	pswdNoSpaces = true
	for _, char := range password {
		switch {
		// func IsLower(r rune) bool
		case unicode.IsLower(char):
			pswdLowercase = true
		// func IsUpper(r rune) bool
		case unicode.IsUpper(char):
			pswdUppercase = true
			err = errors.New("Pa")
		// func IsNumber(r rune) bool
		case unicode.IsNumber(char):
			pswdNumber = true
		// func IsPunct(r rune) bool, func IsSymbol(r rune) bool
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			pswdSpecial = true
		// func IsSpace(r rune) bool, type rune = int32
		case unicode.IsSpace(int32(char)):
			pswdNoSpaces = false
		}
	}
	// check password length
	if 11 < len(password) && len(password) < 60 {
		pswdLength = true
	}
	// create error for any criteria not passed
	if !pswdLowercase || !pswdUppercase || !pswdNumber || !pswdSpecial || !pswdLength || !pswdNoSpaces {
		switch false {
		case pswdLowercase:
			err = errors.New("Password must contain atleast one lower case letter")
		case pswdUppercase:
			err = errors.New("Password must contain atleast one uppercase letter")
		case pswdNumber:
			err = errors.New("Password must contain atleast one number")
		case pswdSpecial:
			err = errors.New("Password must contain atleast one special character")
		case pswdLength:
			err = errors.New("Passward length must atleast 12 characters and less than 60")
		case pswdNoSpaces:
			err = errors.New("Password cannot have any spaces")
		}
		return err
	}
	return nil
}
