package loader

import (
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
	"wbtech_l2/16/parser"
)

func Load(queue *[]*parser.Node, visited map[string]struct{}) {
	rawURL := (*queue)[0].Url
	N := (*queue)[0].N

	if N == 0 || len(*queue) == 0 {
		return
	}

	fmt.Printf("--%s--  %s\n", time.Now().Format("2006-01-02 15:04:05"), rawURL)

	var noProtoRawURL string
	if !strings.HasPrefix(rawURL, "https://") && !strings.HasPrefix(rawURL, "http://") && !strings.HasPrefix(rawURL, "//") {
		noProtoRawURL = rawURL
		fmt.Printf("Resolvind %s... ", rawURL)
		rawURL = strings.Join([]string{"//", rawURL}, "")
	} else {
		noProto := strings.Split(rawURL, "//")
		noProtoRawURL = noProto[1]
		fmt.Printf("Resolvind %s... ", noProtoRawURL)
	}

	URL, err := url.Parse(rawURL)
	if err != nil {
		fmt.Printf("failed: Name or service not known.\nwget: unable to resolve host address '%s'\n\n", noProtoRawURL)
		return
	}

	URL.Path = strings.TrimSuffix(URL.Path, "/")
	URL.RawQuery = ""
	URL.Fragment = ""

	if URL.Scheme != "http" && URL.Scheme != "https" {
		URL.Scheme = "https"
	}

	fmt.Println("resolved.")

	client := http.Client{
		Timeout: 10 * time.Second,
	}

	fmt.Printf("Connecting to %s... ", strings.Join([]string{URL.Hostname(), URL.Path}, ""))

	resp, err := client.Get(URL.String())
	if err != nil {
		var netErr net.Error
		if ok := errors.As(err, &netErr); ok {
			if netErr.Timeout() {
				fmt.Print("Read error (Connection timed out) in headers.\nGiving up.\n\n")
				return
			} else {
				if URL.Scheme == "https" {
					URL.Scheme = "http"
				} else {
					URL.Scheme = "https"
				}

				resp, err = client.Get(URL.String())
				if err != nil {
					fmt.Printf("%s\nwget: unable to resolve host adress '%s'\n\n", err.Error(), URL)
					return
				}
			}
		} else {
			fmt.Printf("%s\nwget: unable to resolve host adress '%s'\n\n", err.Error(), URL)
			return
		}
	}

	fmt.Println("connected.")

	fmt.Print("HTTP request sent, awaiting response... ")

	defer func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {
			fmt.Printf("%s\nwget: unable to close response body\n\n", err.Error())
			return
		}
	}(resp.Body)

	body := resp.Body
	contentType := resp.Header.Get("Content-Type")

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Read error (Unsuccessful status code: %d) in headers.\nGiving up.\n\n", resp.StatusCode)
		return
	}

	fmt.Println(resp.StatusCode)

	var length string
	if resp.ContentLength == -1 {
		length = "unspecified"
	} else {
		length = strconv.Itoa(int(resp.ContentLength))
	}

	fmt.Printf("Length: %s [%s]\n", length, contentType)

	splittedPath := strings.Split(URL.Path, "/")
	if len(splittedPath) != 0 && splittedPath[0] == "" {
		splittedPath = splittedPath[1:]
	}

	var dir string
	filePattern := regexp.MustCompile(`.*\.[A-Za-z]{1,4}([#?].+|/)?$`)

	if len(splittedPath) == 0 {
		dir = strings.Join([]string{URL.Hostname(), URL.Path}, "")
		URL.Path = strings.Join([]string{URL.Path, "index.html"}, "/")
	} else if matches := filePattern.MatchString(splittedPath[len(splittedPath)-1]); !matches {
		dir = strings.Join([]string{URL.Hostname(), URL.Path}, "")
		URL.Path = strings.Join([]string{URL.Path, "index.html"}, "/")
	} else {
		noFilePath := strings.Join(splittedPath[:len(splittedPath)-1], "/")
		dir = strings.Join([]string{URL.Hostname(), noFilePath}, "/")
	}

	fileName := strings.Join([]string{URL.Hostname(), URL.Path}, "")

	fmt.Printf("Saving to: '%s'\n", fileName)

	err = os.MkdirAll(dir, 0755)
	if err != nil {
		fmt.Printf("wget: error while creating directory\n\n")
		return
	}

	f, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		fmt.Printf("wget: error while creating file\n\n")
		return
	}
	defer f.Close()

	visited[fileName] = struct{}{}

	data := parser.Parse(URL, queue, visited, body, N)

	start := time.Now()
	size, err := f.Write(data)
	elapsed := time.Since(start)

	if err != nil {
		fmt.Printf("wget: error saving file\n\n")
		return
	}

	fmt.Printf("'%s' 100%%[===================>] %s in %fs\n\n", fileName, convertUnits(size), elapsed.Seconds())
}

func convertUnits[T int | int64](n T) string {
	size := float64(n)
	units := []string{"B", "K", "M"}

	for i := 0; i < len(units); i++ {
		newSize := size / 1024.0
		if newSize >= 1 {
			size = newSize
		} else {
			return strings.Join([]string{fmt.Sprintf("%.1f", size), units[i]}, "")
		}
	}

	return fmt.Sprintf("%.1fG", size)
}
