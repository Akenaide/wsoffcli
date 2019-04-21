package cmd

import (
	"log"
	"strings"
	"testing"

	"github.com/PuerkitoBio/goquery"
)

func equalSlice(sliceA []string, sliceB []string) bool {
	if len(sliceA) != len(sliceB) {
		log.Printf("wrong len sliceA %v, len sliceB %v", len(sliceA), len(sliceB))
		return false
	}

	for i := range sliceA {
		if sliceA[i] != sliceB[i] {
			log.Printf("wrong value sliceA %v, len sliceB %v", sliceA[i], sliceB[i])
			return false
		}
	}
	return true
}

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
	<span class="unit">トリガー：<img src="https://s3-ap-northeast-1.amazonaws.com/static.ws-tcg.com/wordpress/wp-content/cardimages/_partimages/soul.gif"/>
	<img src="https://s3-ap-northeast-1.amazonaws.com/static.ws-tcg.com/wordpress/wp-content/cardimages/_partimages/bounce.gif"/>
	<img src="https://s3-ap-northeast-1.amazonaws.com/static.ws-tcg.com/wordpress/wp-content/cardimages/_partimages/shot.gif"/>
	<img src="https://s3-ap-northeast-1.amazonaws.com/static.ws-tcg.com/wordpress/wp-content/cardimages/_partimages/treasure.gif"/>
	<img src="https://s3-ap-northeast-1.amazonaws.com/static.ws-tcg.com/wordpress/wp-content/cardimages/_partimages/standby.gif"/>
	<img src="https://s3-ap-northeast-1.amazonaws.com/static.ws-tcg.com/wordpress/wp-content/cardimages/_partimages/salvage.gif"/>
	<img src="https://s3-ap-northeast-1.amazonaws.com/static.ws-tcg.com/wordpress/wp-content/cardimages/_partimages/gate.gif"/>
	<img src="https://s3-ap-northeast-1.amazonaws.com/static.ws-tcg.com/wordpress/wp-content/cardimages/_partimages/draw.gif"/>
	</span>
	<span class="unit">特徴：<span>音楽・Afterglow</span></span><br/>
	<span class="unit">フレーバー：-</span><br/>
	<br/>
	<span>【永】 あなたのターン中、他のあなたの「“止まらずに、前へ”美竹蘭」がいるなら、このカードのパワーを＋6000。<br/>【自】［(1)］ このカードがアタックした時 、あなたはコストを払ってよい。そうしたら、そのアタック中、あなたはトリガーステップにトリガーチェックを2回行う。</span>
	</td>
	`

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(chara))
	var expectedTrigger = []string{"SOUL", "RETURN", "SHOT", "TREASURE", "STANDBY", "COMEBACK", "GATE", "DRAW"}
	var expectedTrait = []string{"音楽", "Afterglow"}
	var expectedAbility = []string{
		"【永】 あなたのターン中、他のあなたの「“止まらずに、前へ”美竹蘭」がいるなら、このカードのパワーを＋6000。",
		"【自】［(1)］ このカードがアタックした時 、あなたはコストを払ってよい。そうしたら、そのアタック中、あなたはトリガーステップにトリガーチェックを2回行う。",
	}

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
	if card.Colour != "GREEN" {
		t.Errorf("got %v: expected GREEN", card.Colour)
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
	if card.CardType != "CH" {
		t.Errorf("got %v: expected CH", card.CardType)
	}
	if card.Rarity != "SPMa" {
		t.Errorf("got %v: expected SPMa", card.Rarity)
	}
	if !equalSlice(card.Trigger, expectedTrigger) {
		t.Errorf("got %v: expected %v", card.Trigger, expectedTrigger)
	}
	if !equalSlice(card.SpecialAttrib, expectedTrait) {
		t.Errorf("got %v: expected %v", card.SpecialAttrib, expectedTrait)
	}
	if !equalSlice(card.Ability, expectedAbility) {
		t.Errorf("got \n %v: expected \n %v", card.Ability, expectedAbility)
	}
}

func TestExtractDataEvent(t *testing.T) {
	chara := `
	<th><a href="/cardlist/?cardno=BD/W63-022&amp;l"><img src="https://s3-ap-northeast-1.amazonaws.com/static.ws-tcg.com/wordpress/wp-content/cardimages/b/bd_w63/bd_w63_022.gif" alt="ミッシェルからの伝言"></a></th>
	<td>
	<h4><a href="/cardlist/?cardno=BD/W63-022&amp;l"><span class="highlight_target">
	ミッシェルからの伝言</span>(<span class="highlight_target">BD/W63-022</span>)</a> -「バンドリ！ ガールズバンドパーティ！」Vol.2<br></h4>
	<span class="unit">
	サイド：<img src="https://s3-ap-northeast-1.amazonaws.com/static.ws-tcg.com/wordpress/wp-content/cardimages/_partimages/w.gif"></span>
	<span class="unit">種類：イベント</span>
	<span class="unit">レベル：1</span><br>
	<span class="unit">色：<img src="https://s3-ap-northeast-1.amazonaws.com/static.ws-tcg.com/wordpress/wp-content/cardimages/_partimages/yellow.gif"></span>
	<span class="unit">パワー：-</span>
	<span class="unit">ソウル：-</span>
	<span class="unit">コスト：0</span><br>
	<span class="unit">レアリティ：U</span>
	<span class="unit">トリガー：－</span>
	<span class="unit">特徴：<span class="highlight_target">-・-</span></span><br>
	<span class="unit">フレーバー：美咲「あはは……ありがとう、はぐみ」</span><br>
	<br>
	<span class="highlight_target">このカードは、あなたの《ハロー、ハッピーワールド！》のキャラが2枚以下なら、手札からプレイできない。<br>あなたは自分の山札の上から2枚を、控え室に置き、自分の控え室のレベルＸ以下のキャラを1枚選び、手札に戻す。Ｘはそれらのカードのレベルの合計に等しい。（クライマックスのレベルは0として扱う）</span>
	</td>
	`

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(chara))
	var expectedTrigger = []string{}

	if err != nil {
		log.Fatal(err)
	}

	card := ExtractData(doc.Clone())
	if !equalSlice(card.Trigger, expectedTrigger) {
		t.Errorf("got %v: expected %v", card.Trigger, expectedTrigger)
	}

	if card.CardType != "EV" {
		t.Errorf("got %v: expected EV", card.CardType)
	}

	if !equalSlice(card.SpecialAttrib, []string{}) {
		t.Errorf("got %v: expected empty", card.SpecialAttrib)
	}

	if card.Soul != "0" {
		t.Errorf("got %v: expected ''", card.Soul)
	}

	if card.Power != "0" {
		t.Errorf("got %v: expected 0", card.Power)
	}
}

func TestExtractDataCX(t *testing.T) {
	chara := `
	<th><a href="/cardlist/?cardno=BD/W63-025&amp;l"><img src="https://s3-ap-northeast-1.amazonaws.com/static.ws-tcg.com/wordpress/wp-content/cardimages/b/bd_w63/bd_w63_025.gif" alt="キラキラのお日様"></a></th>
	<td>
	<h4><a href="/cardlist/?cardno=BD/W63-025&amp;l"><span class="highlight_target">
	キラキラのお日様</span>(<span class="highlight_target">BD/W63-025</span>)</a> -「バンドリ！ ガールズバンドパーティ！」Vol.2<br></h4>
	<span class="unit">
	サイド：<img src="https://s3-ap-northeast-1.amazonaws.com/static.ws-tcg.com/wordpress/wp-content/cardimages/_partimages/w.gif"></span>
	<span class="unit">種類：クライマックス</span>
	<span class="unit">レベル：-</span><br>
	<span class="unit">色：<img src="https://s3-ap-northeast-1.amazonaws.com/static.ws-tcg.com/wordpress/wp-content/cardimages/_partimages/yellow.gif"></span>
	<span class="unit">パワー：-</span>
	<span class="unit">ソウル：-</span>
	<span class="unit">コスト：-</span><br>
	<span class="unit">レアリティ：CR</span>
	<span class="unit">トリガー：<img src="https://s3-ap-northeast-1.amazonaws.com/static.ws-tcg.com/wordpress/wp-content/cardimages/_partimages/soul.gif"><img src="https://s3-ap-northeast-1.amazonaws.com/static.ws-tcg.com/wordpress/wp-content/cardimages/_partimages/bounce.gif"></span>
	<span class="unit">特徴：<span class="highlight_target">-・-</span></span><br>
	<span class="unit">フレーバー：楽しい気持ちは誰かといると生まれるものってこと！</span><br>
	<br>
	<span class="highlight_target">【永】 あなたのキャラすべてに、パワーを＋1000し、ソウルを＋1。<br>（<img src="https://s3-ap-northeast-1.amazonaws.com/static.ws-tcg.com/wordpress/wp-content/cardimages/_partimages/bounce.gif">：このカードがトリガーした時、あなたは相手のキャラを1枚選び、手札に戻してよい）</span>
	</td>
	`

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(chara))

	if err != nil {
		log.Fatal(err)
	}

	card := ExtractData(doc.Clone())

	if card.CardType != "CX" {
		t.Errorf("got %v: expected CX", card.CardType)
	}

	if card.Soul != "0" {
		t.Errorf("got %v: expected ''", card.Soul)
	}

	if card.Level != "0" {
		t.Errorf("got %v: expected 0", card.Level)
	}

	if card.Cost != "0" {
		t.Errorf("got %v: expected 0", card.Cost)
	}

	if strings.Contains(card.Ability[1], "img") {
		t.Errorf("got img tag in %v", card.Ability)
	}
}
