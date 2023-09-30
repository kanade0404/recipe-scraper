package run

import (
	"context"
	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
)

// Nadia
/*
レシピサイトNadiaのスクレイピングを行う
*/

func Nadia(ctx context.Context, actions []chromedp.Action) error {
	return run(ctx, append(chromedp.Tasks{chromedp.Tasks{network.SetBlockedURLS([]string{
		// ロードに時間がかかるので広告系のJSのリクエストをブロックする
		"https://ib.adnxs.com/*",
		"https://securepubads.g.doubleclick.net/tag/js/gpt.js",
		"https://c.amazon-adsystem.com/aax2/apstag.js",
		"https://www.googletagmanager.com/gtm.js?id=*",
		"https://secure.cdn.fastclick.net/js/pubcid/latest/pubcid.min.js",
		"https://flux-cdn.com/client/oceans/nadia.min.js",
		"https://static.adsafeprotected.com/main.*.js",
		"https://bs.serving-sys.com/Serving/adServer.bs?*",
		"https://www.clarity.ms/*",
		"https://static.adsafeprotected.com/sca.*.js",
		"https://www.googletagservices.com/activeview/js/current/rx_lidar.js*",
		"https://bam.nr-data.net*",
	})}}, actions...))
}

// run
/*
chromedpでスクレイピングを行う
*/
func run(ctx context.Context, actions []chromedp.Action) error {
	return chromedp.Run(ctx, actions...)
}
