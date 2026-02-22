package main

import (
	"golang.org/x/net/html"
    "net/http"
    "strings"
    "io"
    "fmt"
    "os"
    "log"
)

func HandleRequest(base string)  {

	var (
		process func(*html.Node)
        urlTovisit []string
        ScrapUrlTovisit []string
    )

    urlTovisit = []string{base} 
    

	process = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "a" {
			for _, ele := range n.Attr {
				if ele.Key == "href" && ele.Val != "../" && ! strings.Contains(ele.Val, "README") {
                    ScrapUrlTovisit = append(ScrapUrlTovisit,base + ele.Val)
				}

                if ele.Key == "href" &&  strings.Contains(ele.Val, "README"){ 
                   url := base + ele.Val
                   resp, err := http.Get(url)
                    if err != nil {
                        log.Fatal(url, err)
                    }

                    bodyBytes, err := io.ReadAll(resp.Body)
                    if err != nil {
                        log.Fatal(url, err)
                    }
                    if strings.Contains(string(bodyBytes), "flag") {
                            fmt.Printf("url: %s\nFound flag : %s\n", url ,string(bodyBytes))
                            os.Exit(0)
                    } else {
                            fmt.Printf("flag not found: %s -> %s\n", url, string(bodyBytes))
                    }
                    resp.Body.Close()
                }
			}
		}


		// traverse the child nodes
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			process(c)
		}
	}


    for len(urlTovisit) > 0 { 


        for _, url := range(urlTovisit){
            base = url
            resp, err := http.Get(url)
            if err != nil {
                log.Fatal(err)
            }


            defer resp.Body.Close()
            doc, err := html.Parse(resp.Body)
            if err != nil {
                return 
            }
            process(doc)
        }
        urlTovisit = ScrapUrlTovisit
        ScrapUrlTovisit = []string{}
    }
}

func main() {

    if len(os.Args) !=  2 {
        log.Fatal("Invalid agrgs")
    }
    HandleRequest(os.Args[1])
}
