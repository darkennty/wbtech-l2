package parser

import (
	"fmt"
	"io"
	"net/url"
	"os"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type Node struct {
	Url string
	N   int
}

func Parse(URL *url.URL, queue *[]*Node, visited map[string]struct{}, body io.ReadCloser, N int) []byte {
	noProtoURL := *URL
	noProtoURL.Scheme = ""

	domain := noProtoURL.Hostname()
	path := strings.Split(noProtoURL.Path, "/")
	hypertextPattern := regexp.MustCompile(".+(.html|.php)$")

	if matches := hypertextPattern.MatchString(path[len(path)-1]); matches {
		doc, err := goquery.NewDocumentFromReader(body)
		if err != nil {
			fmt.Printf("wget: error while parsing document\n\n")
			return nil
		}

		currDir, _ := os.Getwd()

		uniqueLinks := make(map[string]struct{})

		getGoqueryEachFunc := func(a string) func(i int, s *goquery.Selection) {
			return func(i int, s *goquery.Selection) {
				attr, exists := s.Attr(a)
				if exists && attr != "" && attr != URL.String() && attr != noProtoURL.String() {
					if attr[0] == '#' {
						return
					}

					if attr[0] == '/' {
						attr = strings.Join([]string{domain, attr}, "")
					}

					if !strings.HasPrefix(attr, "https://") && !strings.HasPrefix(attr, "http://") && !strings.HasPrefix(attr, "//") {
						attr = strings.Join([]string{"//", attr}, "")
					}

					link, err := url.Parse(attr)
					if err != nil {
						fmt.Printf("wget: unable to resolve hosta address '%s'\n\n", a)
						return
					}

					if link.Host != domain {
						return
					}

					attr = strings.Join([]string{link.Hostname(), link.Path}, "")
					local := attr

					newPath := strings.Split(local, "/")
					line := newPath[len(newPath)-1]

					filePattern := regexp.MustCompile(`.*\.[A-Za-z]{1,4}([#?].+|/)?$`)
					matches = filePattern.MatchString(line)
					if newPath[len(newPath)-1] == domain || !matches {
						local = strings.Join([]string{local, "index.html"}, "/")
					}

					local = strings.Join([]string{currDir, local}, "/")
					local = strings.Replace(local, "//", "/", -1)
					local = strings.Replace(local, "\\", "/", -1)
					local = strings.Join([]string{"file:///", local}, "")

					if link.RawQuery != "" {
						local = strings.Join([]string{local, link.RawQuery}, "?")
					}
					if link.Fragment != "" {
						local = strings.Join([]string{local, link.Fragment}, "#")
					}

					s.SetAttr(a, local)

					if strings.TrimRight(attr, "/") == domain {
						_, ok1 := uniqueLinks[attr]
						_, ok2 := visited[attr]

						if !ok1 && !ok2 {
							*queue = append(*queue, &Node{Url: attr, N: N - 1})
							uniqueLinks[attr] = struct{}{}
						}
					} else {
						_, ok1 := uniqueLinks[attr]
						_, ok2 := visited[attr]

						if !ok1 && !ok2 {
							*queue = append(*queue, &Node{Url: attr, N: N - 1})
							uniqueLinks[attr] = struct{}{}
						}
					}
				}
			}
		}

		doc.Find("link").Each(getGoqueryEachFunc("href"))
		doc.Find("script").Each(getGoqueryEachFunc("src"))
		doc.Find("a").Each(getGoqueryEachFunc("href"))
		doc.Find("img").Each(getGoqueryEachFunc("src"))

		html, _ := doc.Html()

		return []byte(html)
	} else {
		data, err := io.ReadAll(body)
		if err != nil {
			fmt.Printf("wget: error saving file\n\n")
			return nil
		}
		return data
	}
}
