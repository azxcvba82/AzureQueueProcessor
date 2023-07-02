package utils

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
)

const (
	headerAuthorization      = "Authorization"
	headerCacheControl       = "Cache-Control"
	headerContentEncoding    = "Content-Encoding"
	headerContentDisposition = "Content-Disposition"
	headerContentLanguage    = "Content-Language"
	headerContentLength      = "Content-Length"
	headerContentMD5         = "Content-MD5"
	headerContentType        = "Content-Type"
	headerDate               = "Date"
	headerIfMatch            = "If-Match"
	headerIfModifiedSince    = "If-Modified-Since"
	headerIfNoneMatch        = "If-None-Match"
	headerIfUnmodifiedSince  = "If-Unmodified-Since"
	headerRange              = "Range"
	headerUserAgent          = "User-Agent"
	headerXmsDate            = "x-ms-date"
	headerXmsVersion         = "x-ms-version"
)

// NewSharedKeyCredential creates an immutable SharedKeyCredential containing the
// storage account's name and either its primary or secondary key.
func NewSharedKeyCredential(connString string) (*SharedKeyCredential, error) {

	parts := strings.Split(connString, ";")
	accountName := strings.Split(parts[1], "=")[1]

	match1 := regexp.MustCompile("(.?AccountKey)([^;])")
	accountKey := match1.Split(parts[2], -1)[1]

	bytes, err := base64.StdEncoding.DecodeString(accountKey)
	if err != nil {
		return &SharedKeyCredential{}, err
	}
	return &SharedKeyCredential{accountName: accountName, accountKey: bytes}, nil
}

// SharedKeyCredential contains an account's name and its primary or secondary key.
// It is immutable making it shareable and goroutine-safe.
type SharedKeyCredential struct {
	// Only the NewSharedKeyCredential method should set these; all other methods should treat them as read-only
	accountName string
	accountKey  []byte
}

// AccountName returns the Storage account's name.
func (f SharedKeyCredential) AccountName() string {
	return f.accountName
}

func (f SharedKeyCredential) getAccountKey() []byte {
	return f.accountKey
}

func (f *SharedKeyCredential) HttpGetRequest(httpGet *HttpGet) error {
	httpClient := &http.Client{}

	if httpGet.Timeout > 0 {
		httpClient.Timeout = time.Duration(httpGet.Timeout) * time.Second
	}

	httpGet.Request, httpGet.Error = http.NewRequest("GET", httpGet.URI, nil)
	if nil != httpGet.Error {
		return httpGet.Error
	}

	// Add a x-ms-date header if it doesn't already exist
	if d := httpGet.Request.Header.Get(headerXmsDate); d == "" {
		httpGet.Request.Header[headerXmsDate] = []string{time.Now().UTC().Format(http.TimeFormat)}
	}
	httpGet.Request.Header[headerXmsVersion] = []string{"2020-04-08"}
	httpGet.Request.Header[headerContentLength] = []string{strconv.Itoa(0)}
	stringToSign, err := f.buildStringToSign(httpGet.Request)
	if err != nil {
		return err
	}
	signature := f.ComputeHMACSHA256(stringToSign)
	authHeader := strings.Join([]string{"SharedKey ", f.accountName, ":", signature}, "")
	httpGet.Request.Header[headerAuthorization] = []string{authHeader}

	httpGet.Response, httpGet.Error = httpClient.Do(httpGet.Request)
	if httpGet.Error != nil && httpGet.Response != nil {
		return httpGet.Error
	}
	defer httpGet.Response.Body.Close()
	httpGet.StatusCode = httpGet.Response.StatusCode
	httpGet.ResponseBody, httpGet.Error = ioutil.ReadAll(httpGet.Response.Body)
	if nil != httpGet.Error {
		return httpGet.Error
	}
	return nil

}

func (f *SharedKeyCredential) HttpPostRequest(httpPost *HttpPost) error {
	httpClient := &http.Client{}

	if httpPost.Timeout > 0 {
		httpClient.Timeout = time.Duration(httpPost.Timeout) * time.Second
	}

	httpPost.Request, httpPost.Error = http.NewRequest("POST", httpPost.URI, bytes.NewBuffer(httpPost.RequestBody))
	if nil != httpPost.Error {
		return httpPost.Error
	}

	// Add a x-ms-date header if it doesn't already exist
	if d := httpPost.Request.Header.Get(headerXmsDate); d == "" {
		httpPost.Request.Header[headerXmsDate] = []string{time.Now().UTC().Format(http.TimeFormat)}
	}
	httpPost.Request.Header[headerXmsVersion] = []string{"2020-04-08"}
	httpPost.Request.Header[headerContentLength] = []string{strconv.Itoa(len(string(httpPost.RequestBody)))}
	stringToSign, err := f.buildStringToSign(httpPost.Request)
	if err != nil {
		return err
	}
	signature := f.ComputeHMACSHA256(stringToSign)
	authHeader := strings.Join([]string{"SharedKey ", f.accountName, ":", signature}, "")
	httpPost.Request.Header[headerAuthorization] = []string{authHeader}

	httpPost.Response, httpPost.Error = httpClient.Do(httpPost.Request)
	if httpPost.Error != nil && httpPost.Response != nil {
		return httpPost.Error
	}
	defer httpPost.Response.Body.Close()
	httpPost.StatusCode = httpPost.Response.StatusCode
	httpPost.ResponseBody, httpPost.Error = ioutil.ReadAll(httpPost.Response.Body)
	if nil != httpPost.Error {
		return httpPost.Error
	}
	return nil

}

func (f *SharedKeyCredential) HttpPutRequest(httpPost *HttpPost, params ...map[string][]string) error {
	httpClient := &http.Client{}

	var header map[string][]string
	header = make(map[string][]string)
	if len(params) > 0 {
		header = params[0]
	}

	if httpPost.Timeout > 0 {
		httpClient.Timeout = time.Duration(httpPost.Timeout) * time.Second
	}

	httpPost.Request, httpPost.Error = http.NewRequest("PUT", httpPost.URI, bytes.NewBuffer(httpPost.RequestBody))
	if nil != httpPost.Error {
		return httpPost.Error
	}

	// Add a x-ms-date header if it doesn't already exist
	if d := httpPost.Request.Header.Get(headerXmsDate); d == "" {
		httpPost.Request.Header[headerXmsDate] = []string{time.Now().UTC().Format(http.TimeFormat)}
	}
	httpPost.Request.Header[headerXmsVersion] = []string{"2020-04-08"}
	httpPost.Request.Header[headerContentLength] = []string{strconv.Itoa(len(string(httpPost.RequestBody)))}

	for k, v := range header {
		httpPost.Request.Header[k] = v
	}

	stringToSign, err := f.buildStringToSign(httpPost.Request)
	if err != nil {
		return err
	}
	signature := f.ComputeHMACSHA256(stringToSign)
	authHeader := strings.Join([]string{"SharedKey ", f.accountName, ":", signature}, "")
	httpPost.Request.Header[headerAuthorization] = []string{authHeader}

	httpPost.Response, httpPost.Error = httpClient.Do(httpPost.Request)
	if httpPost.Error != nil && httpPost.Response != nil {
		return httpPost.Error
	}
	defer httpPost.Response.Body.Close()
	httpPost.StatusCode = httpPost.Response.StatusCode
	httpPost.ResponseBody, httpPost.Error = ioutil.ReadAll(httpPost.Response.Body)
	if nil != httpPost.Error {
		return httpPost.Error
	}
	return nil

}

func (f *SharedKeyCredential) HttpDeleteRequest(httpGet *HttpGet) error {
	httpClient := &http.Client{}

	if httpGet.Timeout > 0 {
		httpClient.Timeout = time.Duration(httpGet.Timeout) * time.Second
	}

	httpGet.Request, httpGet.Error = http.NewRequest("DELETE", httpGet.URI, nil)
	if nil != httpGet.Error {
		return httpGet.Error
	}

	// Add a x-ms-date header if it doesn't already exist
	if d := httpGet.Request.Header.Get(headerXmsDate); d == "" {
		httpGet.Request.Header[headerXmsDate] = []string{time.Now().UTC().Format(http.TimeFormat)}
	}
	httpGet.Request.Header[headerXmsVersion] = []string{"2020-04-08"}
	httpGet.Request.Header[headerContentLength] = []string{strconv.Itoa(0)}
	stringToSign, err := f.buildStringToSign(httpGet.Request)
	if err != nil {
		return err
	}
	signature := f.ComputeHMACSHA256(stringToSign)
	authHeader := strings.Join([]string{"SharedKey ", f.accountName, ":", signature}, "")
	httpGet.Request.Header[headerAuthorization] = []string{authHeader}

	httpGet.Response, httpGet.Error = httpClient.Do(httpGet.Request)
	if httpGet.Error != nil && httpGet.Response != nil {
		return httpGet.Error
	}
	defer httpGet.Response.Body.Close()
	httpGet.StatusCode = httpGet.Response.StatusCode
	httpGet.ResponseBody, httpGet.Error = ioutil.ReadAll(httpGet.Response.Body)
	if nil != httpGet.Error {
		return httpGet.Error
	}
	return nil

}

func (f SharedKeyCredential) ComputeHMACSHA256(message string) (base64String string) {
	h := hmac.New(sha256.New, f.accountKey)
	h.Write([]byte(message))
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}

func (f *SharedKeyCredential) buildCanonicalizedResource(u *url.URL) (string, error) {
	// https://docs.microsoft.com/en-us/rest/api/storageservices/authentication-for-the-azure-storage-services
	cr := bytes.NewBufferString("/")
	cr.WriteString(f.accountName)

	if len(u.Path) > 0 {
		// Any portion of the CanonicalizedResource string that is derived from
		// the resource's URI should be encoded exactly as it is in the URI.
		// -- https://msdn.microsoft.com/en-gb/library/azure/dd179428.aspx
		cr.WriteString(u.EscapedPath())
	} else {
		// a slash is required to indicate the root path
		cr.WriteString("/")
	}

	// params is a map[string][]string; param name is key; params values is []string
	params, err := url.ParseQuery(u.RawQuery) // Returns URL decoded values
	if err != nil {
		return "", errors.New("parsing query parameters must succeed, otherwise there might be serious problems in the SDK/generated code")
	}

	if len(params) > 0 { // There is at least 1 query parameter
		paramNames := []string{} // We use this to sort the parameter key names
		for paramName := range params {
			paramNames = append(paramNames, paramName) // paramNames must be lowercase
		}
		sort.Strings(paramNames)

		for _, paramName := range paramNames {
			paramValues := params[paramName]
			sort.Strings(paramValues)

			// Join the sorted key values separated by ','
			// Then prepend "keyName:"; then add this string to the buffer
			cr.WriteString("\n" + paramName + ":" + strings.Join(paramValues, ","))
		}
	}
	return cr.String(), nil
}

func buildCanonicalizedHeader(headers http.Header) string {
	cm := map[string][]string{}
	for k, v := range headers {
		headerName := strings.TrimSpace(strings.ToLower(k))
		if strings.HasPrefix(headerName, "x-ms-") {
			cm[headerName] = v // NOTE: the value must not have any whitespace around it.
		}
	}
	if len(cm) == 0 {
		return ""
	}

	keys := make([]string, 0, len(cm))
	for key := range cm {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	ch := bytes.NewBufferString("")
	for i, key := range keys {
		if i > 0 {
			ch.WriteRune('\n')
		}
		ch.WriteString(key)
		ch.WriteRune(':')
		ch.WriteString(strings.Join(cm[key], ","))
	}
	return ch.String()
}

func (f *SharedKeyCredential) buildStringToSign(request *http.Request) (string, error) {
	// https://docs.microsoft.com/en-us/rest/api/storageservices/authentication-for-the-azure-storage-services
	headers := request.Header
	contentLength := headers.Get(headerContentLength)
	if contentLength == "0" {
		contentLength = ""
	}

	canonicalizedResource, err := f.buildCanonicalizedResource(request.URL)
	if err != nil {
		return "", err
	}

	stringToSign := strings.Join([]string{
		request.Method,
		headers.Get(headerContentEncoding),
		headers.Get(headerContentLanguage),
		contentLength,
		headers.Get(headerContentMD5),
		headers.Get(headerContentType),
		"", // Empty date because x-ms-date is expected (as per web page above)
		headers.Get(headerIfModifiedSince),
		headers.Get(headerIfMatch),
		headers.Get(headerIfNoneMatch),
		headers.Get(headerIfUnmodifiedSince),
		headers.Get(headerRange),
		buildCanonicalizedHeader(headers),
		canonicalizedResource,
	}, "\n")
	return stringToSign, nil
}
