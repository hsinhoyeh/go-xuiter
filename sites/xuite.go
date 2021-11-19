package sites

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/gocolly/colly/v2"
	log "github.com/golang/glog"

	goxuiter "github.com/hsinhoyeh/go-xuiter"
)

type XuiteAlbumController struct {
	c                 *goxuiter.CollyController
	destinationPrefix string
	password          string
	client            *http.Client
}

func NewXuiteAlbumController(c *goxuiter.CollyController, destinationPrefix string, password string) *XuiteAlbumController {
	alb := &XuiteAlbumController{
		c:                 c,
		password:          password,
		destinationPrefix: destinationPrefix,
		client: &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			},
		},
	}
	alb.RegisterCallbacks()
	return alb
}

func (x *XuiteAlbumController) RegisterCallbacks() {
	x.c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})
	x.c.OnResponse(func(r *colly.Response) {
		fmt.Println("visited", r.Request.URL)
		//fmt.Printf("body:%s\n", string(r.Body))
	})
	//x.c.OnHTML("p[class=album_info_title]", func(e *colly.HTMLElement) {
	x.c.OnHTML("div[class=album_info]", func(e *colly.HTMLElement) {
		title := e.ChildText("p[class=album_info_title]")
		datestr := e.ChildText("p[class=album_info_date]")
		regex := regexp.MustCompile(`[0-9]{4}-[0-9]{2}-[0-9]{2}`)
		matches := regex.FindAllString(datestr, -1)
		var date string
		if len(matches) != 1 {
			date = time.Now().Format("20060102")
		} else {
			date = matches[0]
		}
		folder := fmt.Sprintf("%s/%s", date, title)
		fmt.Printf("folder:%s\n", folder)
		// handle multiple pages
		for pages := 1; pages < 10; pages++ {
			href := fmt.Sprintf("https:%s*%d?t=%s", e.ChildAttr("a[href]", "href"), pages, folder)
			log.Infof("album title: %+v, href:%s", folder, href)

			u, err := url.Parse(href)
			if err != err {
				log.Errorf("parse url failed:%s\n", href)
				return
			}
			x.c.AddRequest(&colly.Request{
				URL:    u,
				Method: "POST",
				Body:   strings.NewReader(fmt.Sprintf("pwd=%s", x.password)),
			})
		}
	})
	x.c.OnHTML(".photo_item.inline-block", func(e *colly.HTMLElement) {
		q := e.Request.URL.Query()
		myTitle := q["t"][0]
		fileName := stripString(e.Text)
		href := e.ChildAttr("img[src]", "src")
		fullResolutionHref := strings.Replace(href, "_c.jpg", "_x.jpg", -1)

		log.Infof("title:%v, filename:%+v, href:%+v\n", myTitle, fileName, fullResolutionHref)
		if err := goxuiter.SaveFile(x.client, x.destinationPrefix, myTitle, fileName, fullResolutionHref); err != nil {
			log.Error("save file failed, err:%v\n", err)
			return
		}
	})
}

func stripString(str string) string {
	return strings.TrimSpace(str)
}

func (x *XuiteAlbumController) AddSite(siteUrl string) {
	x.c.AddURL(siteUrl)
}

func (x *XuiteAlbumController) Run() {
	x.c.Run()
}
