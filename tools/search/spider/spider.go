package spider

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"strings"

	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/components/tool/utils"
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/go-rod/stealth"
	"github.com/go-shiori/go-readability"
)

type SpiderParams struct {
	URL string `json:"url" jsonschema:"description=用于指定要爬取的目标 URL"`
}

func SpiderFunc(ctx context.Context, config *SpiderConfig) func(ctx context.Context, params *SpiderParams) (string, error) {
	return func(ctx context.Context, params *SpiderParams) (string, error) {
		u := launcher.New().Bin(config.Bin).MustLaunch()
		browser := rod.New().ControlURL(u).MustConnect()
		defer browser.Close()
		// 启动一个新的页面
		page, err := stealth.Page(browser)
		if err != nil {
			return "", err
		}
		defer page.Close()
		// 导航到目标 URL
		err = page.Navigate(params.URL)
		if err != nil {
			return "", err
		}
		page.MustElement("body").MustWaitVisible()
		html, err := page.HTML()
		if err != nil {
			return "", err
		}
		parsedURL, err := url.Parse(params.URL)
		if err != nil {
			fmt.Fprintf(os.Stderr, "解析 URL 失败: %v\n", err)
			return "", err
		}
		article, err := readability.FromReader(strings.NewReader(html), parsedURL)
		if err != nil {
			return "", err
		}
		return article.TextContent, nil // 返回纯文本
	}
}

func NewSpiderTool(ctx context.Context, config *SpiderConfig) (tool.InvokableTool, error) {
	spiderTool, err := utils.InferTool(
		"spider", // tool name
		`Web spider are used to scrape websites to obtain detailed content; 
		they can be used to crawl dynamic JavaScript web pages, eg: https://www.baidu.com`, // tool description
		SpiderFunc(ctx, config))
	if err != nil {
		return nil, err
	}
	return spiderTool, nil
}
