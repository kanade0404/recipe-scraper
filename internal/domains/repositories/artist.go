package repositories

import (
	"recipe-scraper/internal/domains/models"
)

// ArtistRepository
/*
アーティスト情報の読み取りと書き込みを行うメソッドの定義です。
実体はinternal/infrastructures/servicesにあります。
*/
type ArtistRepository interface {
	// Save
	/*
		アーティスト情報を登録します。
		すでに存在する場合は更新します。
	*/
	Save(ID int, name string) (*models.Artist, error)
}
