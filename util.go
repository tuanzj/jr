package jr

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"reflect"
	// "strconv"
	"strings"
	"time"
)

type M map[string]interface{}

func GetTimestamp() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}

func Random(min, max int) int {
	rand.Seed(time.Now().UTC().UnixNano())
	return rand.Intn(max-min) + min
}

func UniqID() int64 {
	// str := strconv.FormatInt(, 10) +
	// 	strconv.Itoa(Random(1000, 9999))
	// if s, err := strconv.ParseInt(str, 10, 64); err == nil {
	// 	return s
	// }
	return time.Now().UnixNano()/int64(time.Millisecond)*1000 + int64(Random(0, 999))
}

func IsPhoneNumber() bool {
	return true
}

func resolveAddress(addr []string) string {
	switch len(addr) {
	case 0:
		if port := os.Getenv("PORT"); len(port) > 0 {
			return ":" + port
		}
		// if len(os.Args) > 1 {
		// return ":" + os.Args[1]
		// }
		return ":80"
	case 1:
		return addr[0]
	default:
		panic("too much parameters")
	}
}
func ByteToString(b []byte) string {
	// n := bytes.IndexByte(b, 0)
	n := bytes.Index(b, []byte{0})
	return string(b[:n])
}

func Request(url string, body string, method string) (string, error) {
	reader := strings.NewReader(body)
	request, err := http.NewRequest(method, url, reader)
	if err != nil {
		return "", err
	}
	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	bodyResp, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	result := string(bodyResp[:])
	return result, nil
}

func EncodeMD5(args ...string) string {
	hasher := md5.New()
	for i := 0; i < len(args); i++ {
		hasher.Write([]byte(args[i]))
	}
	return hex.EncodeToString(hasher.Sum(nil))
}

func InArray(val interface{}, array interface{}) (exists bool, index int) {
	exists = false
	index = -1

	switch reflect.TypeOf(array).Kind() {
	case reflect.Slice:
		s := reflect.ValueOf(array)

		for i := 0; i < s.Len(); i++ {
			if reflect.DeepEqual(val, s.Index(i).Interface()) == true {
				index = i
				exists = true
				return
			}
		}
	}
	return
}

// func convertToString(b []byte) string {
//     s := make([]string,len(b))
//     for i := range b {
//         s[i] = strconv.Itoa(int(b[i]))
//     }
//     return strings.Join(s,",")
// }
