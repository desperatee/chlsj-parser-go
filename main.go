package main

import (
	"encoding/json"
	"flag"
	"fmt"
	http "github.com/bogdanfinn/fhttp"
	"github.com/charmbracelet/log"
	"io"
	"os"
	"strings"
	"time"
)

type Header struct {
	Name  string `json:"name,omitempty"`
	Value string `json:"value,omitempty"`
}

type CharlesRequest struct {
	Status          string      `json:"status,omitempty"`
	Method          string      `json:"method,omitempty"`
	ProtocolVersion string      `json:"protocolVersion,omitempty"`
	Scheme          string      `json:"scheme,omitempty"`
	Host            string      `json:"host,omitempty"`
	ActualPort      int         `json:"actualPort,omitempty"`
	Path            string      `json:"path,omitempty"`
	Query           interface{} `json:"query,omitempty"`
	Tunnel          bool        `json:"tunnel,omitempty"`
	KeptAlive       bool        `json:"keptAlive,omitempty"`
	WebSocket       bool        `json:"webSocket,omitempty"`
	RemoteAddress   string      `json:"remoteAddress,omitempty"`
	ClientAddress   string      `json:"clientAddress,omitempty"`
	ClientPort      int         `json:"clientPort,omitempty"`
	Times           struct {
		Start           time.Time `json:"start,omitempty"`
		RequestBegin    time.Time `json:"requestBegin,omitempty"`
		RequestComplete time.Time `json:"requestComplete,omitempty"`
		ResponseBegin   time.Time `json:"responseBegin,omitempty"`
		End             time.Time `json:"end,omitempty"`
	} `json:"times,omitempty"`
	Durations struct {
		Total    int         `json:"total,omitempty"`
		DNS      interface{} `json:"dns,omitempty"`
		Connect  interface{} `json:"connect,omitempty"`
		Ssl      interface{} `json:"ssl,omitempty"`
		Request  int         `json:"request,omitempty"`
		Response int         `json:"response,omitempty"`
		Latency  int         `json:"latency,omitempty"`
	} `json:"durations,omitempty"`
	Speeds struct {
		Overall  int `json:"overall,omitempty"`
		Request  int `json:"request,omitempty"`
		Response int `json:"response,omitempty"`
	} `json:"speeds,omitempty"`
	TotalSize int `json:"totalSize,omitempty"`
	Ssl       struct {
		Protocol    string `json:"protocol,omitempty"`
		CipherSuite string `json:"cipherSuite,omitempty"`
	} `json:"ssl,omitempty"`
	Alpn struct {
		Protocol string `json:"protocol,omitempty"`
	} `json:"alpn,omitempty"`
	Request struct {
		Sizes struct {
			Headers int `json:"headers,omitempty"`
			Body    int `json:"body,omitempty"`
		} `json:"sizes,omitempty"`
		MimeType        interface{} `json:"mimeType,omitempty"`
		Charset         interface{} `json:"charset,omitempty"`
		ContentEncoding interface{} `json:"contentEncoding,omitempty"`
		Header          struct {
			Headers []Header `json:"headers,omitempty"`
		} `json:"header,omitempty"`
	} `json:"request,omitempty"`
	Response struct {
		Status int `json:"status,omitempty"`
		Sizes  struct {
			Headers int `json:"headers,omitempty"`
			Body    int `json:"body,omitempty"`
		} `json:"sizes,omitempty"`
		MimeType        string `json:"mimeType,omitempty"`
		Charset         string `json:"charset,omitempty"`
		ContentEncoding string `json:"contentEncoding,omitempty"`
		Header          struct {
			Headers []Header `json:"headers,omitempty"`
		} `json:"header,omitempty"`
		Body struct {
			Text    string `json:"text,omitempty"`
			Charset string `json:"charset,omitempty"`
			Decoded bool   `json:"decoded,omitempty"`
		} `json:"body,omitempty"`
	} `json:"response,omitempty"`
}

func isDupe(arr []string, val string) bool {
	for _, value := range arr {
		if value == val {
			return true
		}
	}
	return false
}

func (req *CharlesRequest) Parse() http.Header {
	var headerOrder []string
	var pseudoHeaderOrder []string
	headers := make(http.Header)
	for i := 0; i < len(req.Request.Header.Headers); i++ {
		if !strings.HasPrefix(req.Request.Header.Headers[i].Name, ":") {
			if strings.ToLower(req.Request.Header.Headers[i].Name) != "content-length" {
				headers.Add(req.Request.Header.Headers[i].Name, strings.ReplaceAll(req.Request.Header.Headers[i].Value, "\"", `\"`))
			}
			if !isDupe(headerOrder, req.Request.Header.Headers[i].Name) {
				headerOrder = append(headerOrder, req.Request.Header.Headers[i].Name)
			}
		} else {
			pseudoHeaderOrder = append(pseudoHeaderOrder, req.Request.Header.Headers[i].Name)
		}
	}
	headers[http.HeaderOrderKey] = headerOrder
	headers[http.PHeaderOrderKey] = pseudoHeaderOrder
	return headers
}
func (req *CharlesRequest) String() string {
	headers := req.Parse()
	str := fmt.Sprintf("var parsed = http.Header{\n	")
	for key, vals := range headers {
		if key == http.HeaderOrderKey {
			str += "http.HeaderOrderKey: {"
		} else if key == http.PHeaderOrderKey {
			str += "http.PHeaderOrderKey: {"
		} else {
			str += "\"" + key + "\": {"
		}
		for _, val := range vals {
			str += "\"" + val + "\","
		}
		str = str[:len(str)-1] + "},\n	"
	}
	return str[:len(str)-1] + "}"
}

func ParseCharlesRequests(filename string) error {
	wd, err := os.Getwd()
	if err != nil {
		return err
	}
	file, err := os.Open(wd + "\\" + filename)
	if err != nil {
		return err
	}
	defer file.Close()
	raw, err := io.ReadAll(file)
	if err != nil {
		return err
	}
	var requests []CharlesRequest
	err = json.Unmarshal(raw, &requests)
	if err != nil {
		return err
	}
	for index, request := range requests {
		log.Info(fmt.Sprintf("Parsing request #%v", index), "url", request.Scheme+"://"+request.Host+request.Path)
		fmt.Println(request.String())
	}
	return nil
}

func main() {
	var filename string
	flag.StringVar(&filename, "session", "", "charles session file (.chlsj)")
	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, "Usage: ", os.Args[0], "-session=\"<filename>\"")
		flag.PrintDefaults()
		fmt.Fprintln(os.Stderr, "Parses the given session and extracts request headers")
	}
	flag.Parse()
	err := ParseCharlesRequests(filename)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error: "+err.Error())
		os.Exit(2)
	}
}
