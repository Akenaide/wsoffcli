package cmd

import (
	"strings"
	"testing"

	"github.com/PuerkitoBio/goquery"
)

const productHTML = `
<div class="entry-content">
<h3>ブースターパック ラブライブ！スクールアイドルフェスティバル2 MIRACLE LIVE!</h3>
<div class="product-detail">
<!-- social-buttons -->
<ul class="social-buttons">
</ul>
<!-- /social-buttons -->
<div class="alignright">
<img width="400" height="359" src="https://ws-tcg.com/wordpress/wp-content/uploads/20230731171155/ws_ll_sif2_box.png" class="attachment-post-thumbnail size-post-thumbnail wp-post-image" alt="" srcset="https://ws-tcg.com/wordpress/wp-content/uploads/20230731171155/ws_ll_sif2_box.png 400w, https://ws-tcg.com/wordpress/wp-content/uploads/20230731171155/ws_ll_sif2_box-300x269.png 300w" sizes="(max-width: 400px) 100vw, 400px"></div>
<p class="release"><strong>2023/10/27(Fri) 発売</strong><br>
【
タイトル区分：ラブライブ！スクールアイドルフェスティバル2/ 作品番号：SIL,SIS,SIN,SIP,LSF】
</p>
<div class="section">1パック9枚入り<br>
希望小売価格440円(税込)<br>
<br>
1ボックス　16パック入り<br>
希望小売価格7,040円(税込)<br>
<br>
カード種類数(予定)：ノーマル138種＋パラレル142種</div>
<div class="section">■レアリティ(予定)<br>
SEC（シークレット）4種<br>
OFR（オーバーフレームレア）4種<br>
SP（スペシャル）39種<br>
RRR（トリプルレア）16種<br>
SR（スーパーレア）79種<br>
--------------------------------<br>
RR(ダブルレア)　13種<br>
R(レア)　35種<br>
U(アンコモン)　37種<br>
C(コモン)　37種<br>
CC(クライマックスコモン)　16種<br>
<br>
■お気に入りのメンバーでデッキ構築ができる！<br>
「μ's」「Aqours」「虹ヶ咲学園スクールアイドル同好会」「Liella!」　<br>
4つのグループのカードを収録！<br>
本商品のカードは、それぞれのグループのメンバー同士を組み合わせてデッキを構築することができます！<br>
これまでの「ラブライブ！シリーズ」デッキも強化するカードも収録！<br>
<br>
<br>
■出演キャストの箔押しサインカードを4種収録！<br>
新田恵海さん（高坂穂乃果役）<br>
伊波杏樹さん（高海千歌役）<br>
大西亜玖璃さん（上原歩夢役）<br>
伊達さゆりさん（澁谷かのん役）<br>
<br>
<br>
■新規描き下ろしイラストを4点収録！<br>
高坂穂乃果／高海千歌／上原歩夢／澁谷かのん<br>
<br>
<br>
■箔押しメンバーサインカードを39種収録！<br>
スクールアイドル39名全員のメンバーサインを収録！<br>
<br>
■ボックス特典<br>
ボックス封入PRカード全4種のうち、いずれか1枚をボックス内に封入！<br>
※再版分への封入は未定です。<br>
<br>
■特製先攻後攻カード<br>
本商品オリジナルデザインの先攻後攻カード（2枚1セット）がまれにBOXに封入！(全2パターン)<br>
※再版分への封入は未定です。<br>
<br>
<br>
■ネオスタンダード区分<br>
・ネオスタンダード区分 「ラブライブ！スクールアイドルフェスティバル2」<br>
このパックに収録のカードは全て<b><u>「SIL/」「SIS/」「SIN/」「SIP/」「LSF/」</u></b>で始まるカード番号となり、「ラブライブ！スクールアイドルフェスティバル2」としてデッキを組むことができます。<br>
※「LSF/」のカードの収録は1種を予定<br>
<br>
各作品毎に既存のネオスタンダード区分と混ぜてデッキを組むことができます。<br>
・「ラブライブ！」<br>
<b><u>「SIL/」</u></b>で始まるカード番号のカードは<b><u>「LL/」「LSF/」</u></b>で始まるカード番号のカードと混ぜてデッキを組むことができます。<br>
<br>
・「ラブライブ！サンシャイン!!」<br>
<b><u>「SIS/」</u></b>で始まるカード番号のカードは<b><u>「LSS/」「LSF/」</u></b>で始まるカード番号のカードと混ぜてデッキを組むことができます。<br>
<br>
・「ラブライブ！虹ヶ咲学園スクールアイドル同好会」<br>
<b><u>「SIN/」</u></b>で始まるカード番号のカードは<b><u>「LNJ/」「LSF/」</u></b>で始まるカード番号のカードと混ぜてデッキを組むことができます。<br>
<br>
・「ラブライブ！スーパースター!!」<br>
<b><u>「SIP/」</u></b>で始まるカード番号のカードは<b><u>「LSP/」「LSF/」</u></b>で始まるカード番号のカードと混ぜてデッキを組むことができます。<br>
<br>
<font color="crimson">※「SIL/」と「LSS/」、「SIS/」と「LSP/」など、異なるシリーズでは混ぜられません</font></div>
<div class="section"><h3>先攻後攻カードのデザインを公開!(10/24)</h3>
<p>先攻後攻カード２枚セット(全2パターン)のうち、１セットがまれにBOXに封入！<br>
※再版分への封入は未定です。</p>
<h4>■μ’s×Aqours</h4>
<p><img width="400" src="https://ws-tcg.com/wordpress/wp-content/uploads/20231024174647/1_sil_sample.png">  <img width="400" src="https://ws-tcg.com/wordpress/wp-content/uploads/20231024174648/2_sis_sample.png"> </p>
<h4>■虹ヶ咲学園スクールアイドル同好会×Liella!</h4>
<p><img width="400" src="https://ws-tcg.com/wordpress/wp-content/uploads/20231024174650/3_sin_sample.png">  <img width="400" src="https://ws-tcg.com/wordpress/wp-content/uploads/20231024174652/4_sip_sample.png"> </p>
<h3>OFR（オーバーフレームレア）のデザインを公開!(8/31)</h3>
<p><img width="250" src="https://ws-tcg.com/wordpress/wp-content/uploads/20230831160010/WS_SIL_W109_068OFR.png">  <img width="250" src="https://ws-tcg.com/wordpress/wp-content/uploads/20230831160008/WS_SIS_W109_062OFR.png"> </p>
<p><img width="250" src="https://ws-tcg.com/wordpress/wp-content/uploads/20230831160013/WS_SIN_W109_063OFR.png">  <img width="250" src="https://ws-tcg.com/wordpress/wp-content/uploads/20230831160005/WS_SIP_W109_065OFR.png"> </p>
<h3>ボックスPRカード4種のデザインを公開(8/31)</h3>
<p><img width="250" src="https://ws-tcg.com/wordpress/wp-content/uploads/20230831160405/WS_SIL_W109_139.png">  <img width="250" src="https://ws-tcg.com/wordpress/wp-content/uploads/20230831160404/WS_SIS_W109_140.png"> </p>
<p><img width="250" src="https://ws-tcg.com/wordpress/wp-content/uploads/20230831160401/WS_SIN_W109_141.png">  <img width="250" src="https://ws-tcg.com/wordpress/wp-content/uploads/20230831160403/WS_SIP_W109_142.png"> </p>
<h3>SEC（シークレット）は出演キャストの箔押しサインカード！さらにイラストは描き下ろしイラスト！(8/31)</h3>
<p><img src="https://ws-tcg.com/wordpress/wp-content/uploads/20230831160527/bp_llsif2.png" alt="" width="1920" height="1080" class="alignnone size-full wp-image-53840"></p>
<h3>収録カードのテキストを一部公開！(7/31更新)</h3>
<p><img src="https://ws-tcg.com/wordpress/wp-content/uploads/20230731171315/ws_llsif2_1a.png" alt="" width="1920" height="1080" class="alignnone size-full wp-image-53840"></p>
<p><img src="https://ws-tcg.com/wordpress/wp-content/uploads/20230731171318/ws_llsif2_1b.png" alt="" width="1920" height="1080" class="alignnone size-full wp-image-53840"></p>
<p><img src="https://ws-tcg.com/wordpress/wp-content/uploads/20230731171320/ws_llsif2_2a.png" alt="" width="1920" height="1080" class="alignnone size-full wp-image-53840"></p>
<p><img src="https://ws-tcg.com/wordpress/wp-content/uploads/20230731171323/ws_llsif2_2b.png" alt="" width="1920" height="1080" class="alignnone size-full wp-image-53840"></p>
<p><img src="https://ws-tcg.com/wordpress/wp-content/uploads/20230731171326/ws_llsif2_3a.png" alt="" width="1920" height="1080" class="alignnone size-full wp-image-53840"></p>
<p><img src="https://ws-tcg.com/wordpress/wp-content/uploads/20230731171328/ws_llsif2_3b.png" alt="" width="1920" height="1080" class="alignnone size-full wp-image-53840"></p>
<p><img src="https://ws-tcg.com/wordpress/wp-content/uploads/20230731171330/ws_llsif2_4a.png" alt="" width="1920" height="1080" class="alignnone size-full wp-image-53840"></p>
<p><img src="https://ws-tcg.com/wordpress/wp-content/uploads/20230731171333/ws_llsif2_4b.png" alt="" width="1920" height="1080" class="alignnone size-full wp-image-53840"></p>
</div>
</div>

</div>
`

func TestExtractProductInfo(t *testing.T) {
	reader := strings.NewReader(productHTML)
	doc, err := goquery.NewDocumentFromReader(reader)
	if err != nil {
		t.Error("Error parsing HTML:", err)
		return
	}
	product := extractProductInfo(doc)

	if product.ReleaseDate != "2023/10/27" {
		t.Error("ReleaseDate not good. Found: ", product.ReleaseDate)
	}

	if product.Title != "ブースターパック ラブライブ！スクールアイドルフェスティバル2 MIRACLE LIVE!" {
		t.Error("Title not good. Found: ", product.Title)
	}

	if product.LicenceCode != "SIL,SIS,SIN,SIP,LSF" {
		t.Error("licenceCode not good. Found: ", product.LicenceCode)
	}

	if product.Image != "https://ws-tcg.com/wordpress/wp-content/uploads/20230731171155/ws_ll_sif2_box.png" {
		t.Error("Image not good. Found: ", product.Image)
	}

	if product.SetCode != "W109" {
		t.Error("SetCode not good. Found: ", product.SetCode)
	}
}
