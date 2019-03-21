package cmd

import (
	"log"
	"strings"
	"testing"

	"github.com/PuerkitoBio/goquery"
)

func TestExtractData(t *testing.T) {
	chara := `
	<th><a href="/cardlist/?cardno=BD/W63-036SPMa&amp;l"><img src="https://s3-ap-northeast-1.amazonaws.com/static.ws-tcg.com/wordpress/wp-content/cardimages/b/bd_w63/bd_w63_036spma.gif" alt="“私達、参上っ！”上原ひまり"/></a></th>
	<td>
	<h4><a href="/cardlist/?cardno=BD/W63-036SPMa&amp;l"><span>
	“私達、参上っ！”上原ひまり</span>(<span>BD/W63-036SPMa</span>)</a> -「バンドリ！ ガールズバンドパーティ！」Vol.2<br/></h4>
	<span class="unit">
	サイド：<img src="https://s3-ap-northeast-1.amazonaws.com/static.ws-tcg.com/wordpress/wp-content/cardimages/_partimages/w.gif"/></span>
	<span class="unit">種類：キャラ</span>
	<span class="unit">レベル：2</span><br/>
	<span class="unit">色：<img src="https://s3-ap-northeast-1.amazonaws.com/static.ws-tcg.com/wordpress/wp-content/cardimages/_partimages/green.gif"/></span>
	<span class="unit">パワー：6000</span>
	<span class="unit">ソウル：<img src="https://s3-ap-northeast-1.amazonaws.com/static.ws-tcg.com/wordpress/wp-content/cardimages/_partimages/soul.gif"/><img src="https://s3-ap-northeast-1.amazonaws.com/static.ws-tcg.com/wordpress/wp-content/cardimages/_partimages/soul.gif"/></span>
	<span class="unit">コスト：1</span><br/>
	<span class="unit">レアリティ：SPMa</span>
	<span class="unit">トリガー：<img src="https://s3-ap-northeast-1.amazonaws.com/static.ws-tcg.com/wordpress/wp-content/cardimages/_partimages/soul.gif"/></span>
	<span class="unit">特徴：<span>音楽・Afterglow</span></span><br/>
	<span class="unit">フレーバー：-</span><br/>
	<br/>
	<span>【永】 あなたのターン中、他のあなたの「“止まらずに、前へ”美竹蘭」がいるなら、このカードのパワーを＋6000。<br/>【自】［(1)］ このカードがアタックした時
	、あなたはコストを払ってよい。そうしたら、そのアタック中、あなたはトリガーステップにトリガーチェックを2回行う。</span>
	</td>
	`

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(chara))

	if err != nil {
		log.Fatal(err)
	}

	card := ExtractData(doc.Clone())
	if card.JpName != "“私達、参上っ！”上原ひまり" {
		t.Errorf("got %v: expected “私達、参上っ！”上原ひまり", card.JpName)
	}
	if card.Set != "BD" {
		t.Errorf("got %v: expected BD", card.Set)
	}
	if card.SetName != "BD" {
		t.Errorf("got %v: expected BD", card.SetName)
	}
	if card.Side != "W" {
		t.Errorf("got %v: expected W", card.Side)
	}
	if card.Release != "63" {
		t.Errorf("got %v: expected 63", card.Release)
	}
	if card.ID != "036SPMa" {
		t.Errorf("got %v: expected 036SPMa", card.ID)
	}
	if card.Level != "2" {
		t.Errorf("got %v: expected 2", card.Level)
	}
	if card.Color != "green" {
		t.Errorf("got %v: expected green", card.Color)
	}
	if card.Power != "6000" {
		t.Errorf("got %v: expected 6000", card.Power)
	}
	if card.Soul != "2" {
		t.Errorf("got %v: expected 2", card.Soul)
	}
	if card.Cost != "1" {
		t.Errorf("got %v: expected 1", card.Cost)
	}
	if card.Rarity != "SPMa" {
		t.Errorf("got %v: expected SPMa", card.Rarity)
	}
	t.Errorf("lyay")
}
