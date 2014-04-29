package main

import (
	"fmt"
	"github.com/ChimeraCoder/anaconda"
	rss "github.com/haarts/go-pkg-rss"
	"log"
	"net/url"
	"os"
	"regexp"
	"time"
	"strings"
)

const timeout = 50

var first = map[string]bool{}

func main() {
	fmt.Println("-------------------------")
	fmt.Printf("Init from ")
	fmt.Println(time.Now().Format("2006-01-02  15:04:05"))
	log.SetOutput(os.Stdout)
	fmt.Printf("Second\n")
	//go PollFeed("http://blog.golang.org/feed.atom", itemHandlerGoBlog)
	//go PollFeed("https://news.ycombinator.com/rss", itemHandlerHackerNews)
	//PollFeed("http://www.reddit.com/r/golang.rss", itemHandlerReddit)
	PollFeed("http://allafrica.com/tools/headlines/rdf/latest/headlines.rdf", itemHandlerAllafrica)
	//go PollFeed("http://www.biztechafrica.com/feed/rss", itemHandlerBiztechafrica)
	//go PollFeed("http://feeds.bbci.co.uk/news/world/africa/rss.xml", itemHandlerBBCAfrica)
	fmt.Printf("Start to get the RSS\n")
}

func PollFeed(uri string, itemHandler rss.ItemHandler) {
	feed := rss.New(timeout, true, chanHandler, itemHandler)

	for {
		if err := feed.Fetch(uri, nil); err != nil {
			fmt.Fprintf(os.Stderr, "[e] %s: %s", uri, err)
			return
		}
		fmt.Printf("We are waitting for the result\n")
		<-time.After(time.Duration(feed.SecondsTillUpdate() * 1e9))
	}
}

func chanHandler(feed *rss.Feed, newchannels []*rss.Channel) {
	//noop
}

func genericItemHandler(feed *rss.Feed, ch *rss.Channel, newItems []*rss.Item, individualItemHandler func(*rss.Item)) {
	log.Printf("%d new item(s) in %s\n", len(newItems), feed.Url)
	for _, item := range newItems {
		individualItemHandler(item)
	}
}

func itemHandlerHackerNews(feed *rss.Feed, ch *rss.Channel, newItems []*rss.Item) {
	f := func(item *rss.Item) {
		if match, _ := regexp.MatchString(`\w Go( |$|\.)`, item.Title); match {
			short_title := item.Title
			if len(short_title) > 100 {
				short_title = short_title[:99] + "…"
			}
			PostTweet(short_title + " " + item.Links[0].Href + " #hackernews")
		}
	}

	if _, ok := first["hn"]; !ok {
		first["hn"] = false
	} else {
		genericItemHandler(feed, ch, newItems, f)
	}
}

func itemHandlerGoBlog(feed *rss.Feed, ch *rss.Channel, newItems []*rss.Item) {
	f := func(item *rss.Item) {
		short_title := item.Title
		if len(short_title) > 100 {
			short_title = short_title[:99] + "…"
		}
		PostTweet(short_title + " " + item.Links[0].Href + " #go_blog")
	}

	if _, ok := first["go"]; !ok {
		first["go"] = false
	} else {
		genericItemHandler(feed, ch, newItems, f)
	}
}

func itemHandlerReddit(feed *rss.Feed, ch *rss.Channel, newItems []*rss.Item) {
	f := func(item *rss.Item) {
		re := regexp.MustCompile(`([^"]+)">\[link\]`)
		matches := re.FindStringSubmatch(item.Description)
		if len(matches) == 2 {
			short_title := item.Title
			if len(short_title) > 100 {
				short_title = short_title[:99] + "…"
			}
			PostTweet(short_title + " " + matches[1] + " #reddit")
		}
	}

	if _, ok := first["reddit"]; !ok {
		first["reddit"] = false
	} else {
		genericItemHandler(feed, ch, newItems, f)
	}
}

func itemHandlerAllafrica(feed *rss.Feed, ch *rss.Channel, newItems []*rss.Item) {
        fmt.Printf("Allafrica start\n")
        defer fmt.Printf("Allafrica end\n")
	f := func(item *rss.Item) {
		short_title := item.Title
		if len(short_title) > 100 {
			short_title = short_title[:99] + "…"
		}
		tag_country := strings.Replace(strings.Split(item.Title, ":")[0], " ", "", -1)
/*		PostTweet(short_title + " " + item.Links[0].Href + " #allafrica")*/
PostTweet(short_title + " " + item.Links[0].Href + " #afmobi" + " #" + tag_country + " #allafrica")
	}

	if _, ok := first["allafrica"]; !ok {
		first["allafrica"] = false
	} else {
		genericItemHandler(feed, ch, newItems, f)
	}
}

func itemHandlerBiztechafrica(feed *rss.Feed, ch *rss.Channel, newItems []*rss.Item) {
        fmt.Printf("Biztechafrica start\n")
        defer fmt.Printf("Biztechafrica end\n")
	f := func(item *rss.Item) {
		short_title := item.Title
		if len(short_title) > 100 {
			short_title = short_title[:99] + "…"
		}
		PostTweet(short_title + " " + item.Links[0].Href + " #afmobi" + " #biztechafrica")
	}

	if _, ok := first["biztechafrica"]; !ok {
		first["biztechafrica"] = false
	} else {
		genericItemHandler(feed, ch, newItems, f)
	}
}

func itemHandlerBBCAfrica(feed *rss.Feed, ch *rss.Channel, newItems []*rss.Item) {
				fmt.Printf("BBCAfrica start\n")
				defer fmt.Printf("BBCAfrica end\n")
	f := func(item *rss.Item) {
		short_title := item.Title
		if len(short_title) > 100 {
			short_title = short_title[:99] + "…"
		}
		PostTweet(short_title + " " + item.Links[0].Href + " #afmobi" + " #bbcafrica")
	}

	if _, ok := first["bbcafrica"]; !ok {
		first["bbcafrica"] = false
	} else {
		genericItemHandler(feed, ch, newItems, f)
	}
}

func PostTweet(tweet string) {
	fmt.Printf("Run to post the tweet\n")
	defer fmt.Printf("Finish post tweet\n")
	anaconda.SetConsumerKey(ReadConsumerKey())
	anaconda.SetConsumerSecret(ReadConsumerSecret())
	api := anaconda.NewTwitterApi(ReadAccessToken(), ReadAccessTokenSecret())

	v := url.Values{}
	_, err := api.PostTweet(tweet, v)
	if err != nil {
		log.Printf("Error posting tweet: %s", err)
	}
}
