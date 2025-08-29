package main

import (
    "context"
    "fmt"
    "log"
    "time"
    
    "github.com/chromedp/chromedp"
)

func main() {
    ctx, cancel := chromedp.NewContext(context.Background())
    defer cancel()

    var allLinks []string

    err := chromedp.Run(ctx,
        chromedp.Navigate("https://startups.gallery"),
        chromedp.Sleep(5*time.Second),
        chromedp.Evaluate(`Array.from(document.querySelectorAll('a')).map(a => a.href)`, &allLinks),
    )

    if err != nil {
        log.Fatal(err)
    }

    fmt.Println("All links found:")
    for i, link := range allLinks {
        fmt.Printf("%d: %s\n", i+1, link)
    }
    
    // ioutil.WriteFile("output.html", []byte(html), 0644)
}
