package utils

/*
  * Contains helper functions
*/


import (
  "os"
  "io"
  "log"
  "fmt"
  "time"
  "sync"
  "errors"
  "regexp"
  "strconv"
  "strings"
  "math/rand"
  cr"crypto/rand"
	"encoding/base64"
  "github.com/google/uuid"
  "golang.org/x/crypto/bcrypt"
  //"github.com/dgrijalva/jwt-go"
)

func GenerateUUID() string {
  return uuid.New().String()
}

type TimeStamps struct {
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (t *TimeStamps) Touch() {
  currentTime := time.Now()
  formattedTime := currentTime.Format("2006-01-02 15:04:05")
  parsedTime, _ := time.Parse("2006-01-02 15:04:05", formattedTime)
  t.UpdatedAt = parsedTime.UTC()
  if t.CreatedAt.IsZero() {
    t.CreatedAt = t.UpdatedAt
  }
}

func IntToString(val int) string {
	return strconv.Itoa(val)
}

// ArrayContainsInt checks if an integer exists in a slice of integers.
func ArrayContainsInt(array []int, target int) bool {
	for _, value := range array {
		if value == target {
			return true
		}
	}
	return false
}

func GetCurrentTime() string {
  var now = time.Now()
  return now.Format("2006-01-02 15:04:05")
}

// RemoveStringDuplicates removes duplicate strings from a slice and returns a new slice with unique values.
func RemoveStringDuplicates(array []string) []string {
	uniqueMap := make(map[string]bool)
	var result []string
	for _, value := range array {
		if _, exists := uniqueMap[value]; !exists {
			uniqueMap[value] = true
			result = append(result, value)
		}
	}

	return result
}

func CheckIfStringIsDomainName(s string) bool {
	domainRegex := `^([a-zA-Z0-9-]{1,63}\.)+[a-zA-Z]{2,}$`
	re := regexp.MustCompile(domainRegex)
	return re.MatchString(s)
}

func ContainsOnlyNumbers(s string) bool {
	// Use a regular expression to check if the string contains only numbers
	match, _ := regexp.MatchString("^[\\p{N}]+$", s)
	return match
}

func HashPassPin(pin string)(string,error) {
  var hash []byte
  hash,err := bcrypt.GenerateFromPassword([]byte(pin),bcrypt.DefaultCost)
  if err != nil{
    return "",fmt.Errorf("Error generating password hash: %q",err)
  }
  return string(hash),nil
}

var invalidUsernameRe = regexp.MustCompile("[^A-Za-z0-9]")

func ValidateUsername(username string) error {
	username = strings.TrimSpace(username)

	if len(username) < 4 || len(username) > 15 {
		return errors.New("invalid username length")
	}

	loc := invalidUsernameRe.FindStringIndex(username)
	if loc != nil {
		return errors.New("invalid username")
	}
	return nil
}


func GenerateBusinessNumber() string {
  var (
  	mu      sync.Mutex
  	counter int
  )
	mu.Lock()
	defer mu.Unlock()
	// Generate a UUID and convert it to a simpler string format
	id := uuid.New()
	uuidStr := id.String()[0:8]
	// Remove '0' and 'O' characters from the UUID
	uuidStr = strings.ReplaceAll(uuidStr, "0", "")
	uuidStr = strings.ReplaceAll(uuidStr, "O", "")
	// Increment the counter
	counter++
	// Capitalize non-digit characters
	result := make([]byte, 0, len(uuidStr))
	for _, ch := range uuidStr {
		if ch >= '0' && ch <= '9' {
			result = append(result, byte(ch))
		} else {
			result = append(result, byte(ch-32)) // Convert to uppercase
		}
	}
	// Combine UUID and counter to create the business number
	businessNumber := fmt.Sprintf("%s-%d", string(result), counter)
	return businessNumber
}

func IsValidEmail(email string) bool {
	// Check if the email has consecutive dots in the domain part
	if strings.Contains(email, "..") {
		return false
	}
	// Regular expression to validate email addresses
	emailPattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	// Compile the regular expression
	re := regexp.MustCompile(emailPattern)
	// Use MatchString method to check if the email matches the pattern
	return re.MatchString(email)
}

func GenerateCSRFToken(csrfTokenLength int) (string, error) {
	tokenBytes := make([]byte, csrfTokenLength)
	_, err := cr.Read(tokenBytes)
	if err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(tokenBytes), nil
}
//create a strings comparer function
//create a convert to lowercase and uppercase function

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

func GetPhone(input string) (bool) {
  if strings.Contains(input, ".") {
		return false
	}
  // Regular expression to match only numeric characters
	numericPattern := `^[0-9]+$`
	// Compile the regular expression
	re := regexp.MustCompile(numericPattern)
	// Use MatchString method to check if the input matches the pattern
	return re.MatchString(input)
}

func StringToInt(data string) int{
  dt,_ := strconv.Atoi(data)
  return dt
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

/*
type Logger struct {
  Name string
  Text  interface{}
}
*/

var ErrorFileNames = []string{"users_sql","apikey_sql","auth_sql","auth_danger_sql","auth_danger"}

// no need to block if  the file to log to is of a different name hence we use a map
var logMutex sync.Map

func LogToFile(ldr Logger) error{
  date := time.Now().Format("2006-01-02")
	name := fmt.Sprintf("./.data/logs/%s/%s.log",date,ldr.Name)
	// Get or create a mutex for the specific log file name
	mutex, _ := logMutex.LoadOrStore(name, new(sync.Mutex))
	mutex.(*sync.Mutex).Lock()
	defer mutex.(*sync.Mutex).Unlock()

	f, err := os.OpenFile(name, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		LogError(err)
		return err
	}
	defer f.Close()
	writer := io.MultiWriter(os.Stdout, f)
	log.SetOutput(writer)
	log.Println(ldr.Text)
	return nil
}


// This log can cause an error when running multiple go routines (solve it with Logger :) i.e LogToFile(Logger)
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

func LogError(err error){
  log.Println(err)
}
