package utils
/*
  This package implements simple functions to help in form filling and implementing htpp clients
*/

import (
  "net"
  "net/http"
)

func ServerAddHeaderVal(res http.ResponseWriter, name,val string) {
  res.Header().Set(name,val)
}

func ClientAddHeaderVal(req *http.Request,name,val string){
  req.Header.Set(name,val)
}

func CheckIfStringIsIp(s string) bool {
	return net.ParseIP(s) != nil
}


/*
import (
  "net/url"
  "fmt"
  //"time"
  "bytes"
  //"strings"
  "regexp"
  "net/url"
  "net/http"
  "io/ioutil"
)
/*
type Transporter struct {
  http.RoundTripper
  headers [][]string
}

func (tr *Transporter) RoundTrip(req *http.Request) (*http.Response,error){
  for _, header := range tr.headers {
    req.Header.Add(header[0],header[1])
  }
  return tr.RoundTripper.RoundTrip(req)
}
// returns a http client
func GetHTTPClient(req *http.Request,headers [][]string) *http.Client {
    if len(headers) > 0 {
        tr := &Transporter{
            DisableKeepAlives: true,
            DialContext: (&net.Dialer{
                Timeout:   30 * time.Second,
                KeepAlive: 30 * time.Second,
                DualStack: true,
            }).DialContext,
            TLSHandshakeTimeout:   10 * time.Second,
            ResponseHeaderTimeout: 10 * time.Second,
            ExpectContinueTimeout: 1 * time.Second,
            headers: headers,
        }
        client := &http.Client{Transport:tr}
    }
    return client
}

// returns a https client
func HTTPSClient(headers [][]string) *http.Client {
    client := &http.Client{}
    if len(headers) > 0 {
        client.Transport = &http.Transport{
            TLSClientConfig: &tls.Config{
                InsecureSkipVerify: true,
            },
            DisableKeepAlives: true,
            DialContext: (&net.Dialer{
                Timeout:   30 * time.Second,
                KeepAlive: 30 * time.Second,
                DualStack: true,
            }).DialContext,
            TLSHandshakeTimeout:   10 * time.Second,
            ResponseHeaderTimeout: 10 * time.Second,
            ExpectContinueTimeout: 1 * time.Second,
        }
        for _, header := range headers {
            client.Transport.Header().Set(header[0], header[1])
        }
    }
    return client
}
*/
/*
headers := [][]string{
    {"X-Forwarded-For", "localhost"},
    {"User-Agent", "MyCustomClient/1.0"},
}
client := utils.HTTPClient(headers)


// header checker receives a https response and prints out/returns all headers
func GetRespHeaders(resp *http.Response) []string {
    var headers []string
    for name, values := range resp.Header {
        for _, value := range values {
            headers = append(headers, fmt.Sprintf("%s: %s", name, value))
        }
    }
    return headers
}

// header checker receives a https request and prints out/returns
func GetReqHeaders(req *http.Request) []string {
    var headers []string
    for name, values := range req.Header {
        for _, value := range values {
            headers = append(headers, fmt.Sprintf("%s: %s", name, value))
        }
    }
    return headers
}
/*
req, err := http.NewRequest("GET", "http://example.com", nil)
if err != nil {
    // handle error
}
client := &http.Client{}
resp, err := client.Do(req)
if err != nil {
    // handle error
}
defer resp.Body.Close()
headers := GetHeaders(resp)
fmt.Println(headers)

func FillFormURLParam(url string, form url.Values,client *http.Client) (string, string, error) {
    // Encode the form data as a query string
    query := form.Encode()
    // Create a new POST request with the form data as the URL query
    req, err := http.NewRequest("POST", url+"?"+query, nil)
    if err != nil {
        return "", "", err
    }
    // Set the content type to "application/x-www-form-urlencoded" (We can set more headers using HTTPSClient )
    req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
    // Use a client to send the request
    resp, err := client.Do(req)
    if err != nil {
        return "", "", err
    }
    defer resp.Body.Close()
    // Read the response body
    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        return "", "", err
    }
    return string(body), resp.Request.URL.String(), nil
}

// this client can be from an ninitial request or in a new an initialized request
func FillFormBody(url string, form url.Values,client *http.Client) (string, string, error) {
    // Encode the form data as a reader
    reader := strings.NewReader(form.Encode())
    // Use a client to send the request
    resp, err := client.Do(req)
    if err != nil {
        return "", "", err
    }
    defer resp.Body.Close()
    // Read the response body
    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        return "", "", err
    }
    return string(body), resp.Request.URL.String(), nil
}

/*
func FillFormBody(url string, form url.Values) (*http.Response, error) {
    // Encode the form data as a reader
    reader := strings.NewReader(form.Encode())

    // Create a new POST request with the form data as the body
    req, err := http.NewRequest("POST", url, reader)
    if err != nil {
        return nil, err
    }

    // Set the content type to "application/x-www-form-urlencoded"
    req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

    // Use a client to send the request
    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    // Read the response body
    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        return nil, err
    }

    // Parse the response body as HTML
    doc, err := goquery.NewDocumentFromReader(bytes.NewReader(body))
    if err != nil {
        return nil, err
    }

    // Find the form element in the HTML
    formEl := doc.Find("form")
    if formEl.Length() == 0 {
        return nil, errors.New("form not found")
    }

    // Get the action and method attributes of the form
    action, exists := formEl.Attr("action")
    if !exists {
        return nil, errors.New("form action not found")
    }
    method, exists := formEl.Attr("method")
    if !exists {
        return nil, errors.New("form method not found")
    }

    // Initialize a new form data object
    newForm := url.Values{}

    // Find all the input elements in the form
    formEl.Find("input").Each(func(i int, s *goquery.Selection) {
        // Get the name and value attributes of the input element
        name, exists := s.Attr("name")
        if !exists {
            return
        }
        value, exists := s.Attr("value")
        if !exists {
            value = ""
        }

        // Add the name-value pair to the form data object
        newForm.Set(name, value)
    })

    // Set the method and form data in the new POST request
    req, err = http.NewRequest(method, action, strings.NewReader(newForm.Encode()))
    if err != nil {
        return nil, err
    }
    req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

    // Send the new POST request and return the response
    return client.Do(req)
}

func FillFormBody(url net.Url,client *http.Client) (*http.Response, error){
  resp,err := client.Get(url)
  if err != nil {
    return nil,err
  }
  defer resp.Body.Close()
  body,err := ioutil.ReadAll(resp.Body)
  if err != nil {
    return nil,err
  }
  formValues := map[string]string{}
	re := regexp.MustCompile(`name="(.*?)" value="(.*?)"`)
	matches := re.FindAllStringSubmatch(string(body), -1)
	for _, match := range matches {
		formValues[match[1]] = match[2]
	}
	// Set the form data
	formData := url.Value
	for key, value := range formValues {
		formData.Set(key, value)
	}
  req,err := http.NewRequest("POST","https://target/form",bytes.NewBufferString(formData.Encode()))
  if err != nil {
    return nil,err
  }
  req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
  return  client.Do(req)
}
*/
