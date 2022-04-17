package main

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/JoungSik/MovieRating_crawler/cmd/models"
	"github.com/araddon/dateparse"
	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/chromedp"
)

func main() {

	var movie = models.Movie{Code: "204138", Page: 1}
	var url = "https://movie.naver.com/movie/bi/mi/pointWriteFormList.naver?code=" + movie.Code + "&type=after&onlyActualPointYn=N&onlySpoilerPointYn=N&order=sympathyScore&page=" + strconv.Itoa(movie.Page)

	ctx, cancel := chromedp.NewContext(
		context.Background(),
		// chromedp.WithDebugf(log.Printf),
	)
	defer cancel()

	var reples []models.Reple

	var maxPages string
	var nodes []*cdp.Node
	var ratings []*cdp.Node
	var messages []*cdp.Node
	var dates []*cdp.Node

	err := chromedp.Run(ctx,
		chromedp.Navigate(url),
		chromedp.Text("div.score_total > strong.total > em", &maxPages),
		chromedp.Nodes("div.score_result > ul > li", &nodes),
		chromedp.Nodes("div.score_result > ul > li > div.star_score > em", &ratings),
		chromedp.Nodes("div.score_result > ul > li > div.score_reple > p > span", &messages),
		chromedp.Nodes("div.score_result > ul > li > div.score_reple > dl > dt > em", &dates),
	)
	if err != nil {
		log.Fatal(err)
	}

	// r := funk.Filter([]int{1, 2, 3, 4}, func(x int) bool {
	// 	return x%2 == 0
	// })

	for _, n := range messages {
		if n.AttributeValue("id") != "" {
			reples = append(reples, models.Reple{
				Message: strings.TrimSpace(n.Children[0].NodeValue),
			})
		}
	}

	fmt.Println(len(reples))

	for i, n := range ratings {
		value, err := strconv.Atoi(strings.TrimSpace(n.Children[0].NodeValue))
		if err == nil {
			reples[i].Score = value
		}
	}

	for i, n := range dates {
		if n.Children != nil {
			date := strings.Replace(n.Children[0].NodeValue, ".", "/", -1)
			t, err := dateparse.ParseLocal(date)
			if err != nil {
				log.Fatal(err.Error())
			}
			fmt.Println(i, date, t)
		}
	}

	for i, n := range reples {
		fmt.Println(i, n.Score, n.Message)
	}
}
