package initialize

import (
	"context"
	"github.com/chromedp/chromedp"
	"time"
)

// Nadia レシピサイトNadia向けのchromedpの初期化を行う
func Nadia(c context.Context) (ctx context.Context, timeoutCancel context.CancelFunc, allocatorCancel context.CancelFunc, contextCancel context.CancelFunc) {
	return initialize(c, []chromedp.ExecAllocatorOption{
		// ヘッドレスモードを無効にする
		chromedp.Flag("headless", false),
		// GPUを無効にする
		chromedp.Flag("disable-gpu", false),
		// 自動テストソフトウェアによる操作を無効にする
		chromedp.Flag("enable-automation", false),
		// 拡張機能を無効にする
		chromedp.Flag("disable-extensions", false),
		// スクロールバーを非表示にする
		chromedp.Flag("hide-scrollbars", false),
		// サンドボックスを無効にする
		chromedp.Flag("no-sandbox", true),
		// メモリの使用量を減らすために/dev/shmを使用しない
		chromedp.Flag("disable-dev-shm-usage", true),
		chromedp.Flag("blink-settings", "imagesEnabled=false"),
		// 音声をミュートにする
		chromedp.Flag("mute-audio", false)}...)
}

// initialize chromedpの初期化を行う
func initialize(c context.Context, opts ...chromedp.ExecAllocatorOption) (ctx context.Context, timeoutCancel context.CancelFunc, allocatorCancel context.CancelFunc, contextCancel context.CancelFunc) {
	options := append(chromedp.DefaultExecAllocatorOptions[:], opts...)
	// 要素の取得がうまくいかなかった時にずっと待ち続けないようにtimeoutを設定する
	ctx, timeoutCancel = context.WithTimeout(c, 5*time.Minute)
	alc, allocatorCancel := chromedp.NewExecAllocator(ctx, options...)
	ctx, contextCancel = chromedp.NewContext(alc)
	return ctx, timeoutCancel, allocatorCancel, contextCancel
}
