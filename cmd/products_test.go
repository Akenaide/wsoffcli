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
	product, err := extractProductInfo(doc)

	if err != nil {
		t.Error("Got unexpected error: ", err)
	}

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

const productHTMLUnexpectedTitle = `
<div class="entry-content">
<h3>トライアルデッキ 富士見ファンタジア文庫 Vol.2</h3>
<div class="product-detail">
<!-- social-buttons -->
<ul class="social-buttons">
<li class="twitter"><iframe id="twitter-widget-1" scrolling="no" frameborder="0" allowtransparency="true" allowfullscreen="true" class="twitter-share-button twitter-share-button-rendered twitter-tweet-button" title="X Post Button" src="https://platform.twitter.com/widgets/tweet_button.2f70fb173b9000da126c79afe2098f02.ja.html#dnt=false&amp;hashtags=ws2tcg&amp;id=twitter-widget-1&amp;lang=ja&amp;original_referer=https%3A%2F%2Fws-tcg.com%2Fproducts%2Ff-td2%2F&amp;size=m&amp;text=%E3%83%88%E3%83%A9%E3%82%A4%E3%82%A2%E3%83%AB%E3%83%87%E3%83%83%E3%82%AD%20%E5%AF%8C%E5%A3%AB%E8%A6%8B%E3%83%95%E3%82%A1%E3%83%B3%E3%82%BF%E3%82%B8%E3%82%A2%E6%96%87%E5%BA%AB%20Vol.2%EF%BD%9C%E3%83%B4%E3%82%A1%E3%82%A4%E3%82%B9%E3%82%B7%E3%83%A5%E3%83%B4%E3%82%A1%E3%83%AB%E3%83%84%EF%BD%9CWei%CE%B2%20Schwarz&amp;time=1721622285925&amp;type=share&amp;url=https%3A%2F%2Fws-tcg.com%2Fproducts%2Ff-td2%2F" style="position: static; visibility: visible; width: 77px; height: 20px;" data-url="https://ws-tcg.com/products/f-td2/"></iframe>
<script>!function(d,s,id){var js,fjs=d.getElementsByTagName(s)[0];if(!d.getElementById(id)){js=d.createElement(s);js.id=id;js.src="//platform.twitter.com/widgets.js";fjs.parentNode.insertBefore(js,fjs);}}(document,"script","twitter-wjs");</script></li>
<li class="facebook"><div class="fb-like fb_iframe_widget" data-href="https://ws-tcg.com/products/f-td2/" data-send="false" data-layout="button_count" data-width="100" data-show-faces="true" fb-xfbml-state="rendered" fb-iframe-plugin-query="app_id=&amp;container_width=0&amp;href=https%3A%2F%2Fws-tcg.com%2Fproducts%2Ff-td2%2F&amp;layout=button_count&amp;locale=ja_JP&amp;sdk=joey&amp;send=false&amp;show_faces=true&amp;width=100"><span style="vertical-align: bottom; width: 130px; height: 28px;"><iframe name="f2f1fb507ca4c5188" width="100px" height="1000px" data-testid="fb:like Facebook Social Plugin" title="fb:like Facebook Social Plugin" frameborder="0" allowtransparency="true" allowfullscreen="true" scrolling="no" allow="encrypted-media" src="https://www.facebook.com/plugins/like.php?app_id=&amp;channel=https%3A%2F%2Fstaticxx.facebook.com%2Fx%2Fconnect%2Fxd_arbiter%2F%3Fversion%3D46%23cb%3Dfd9ad7880f4ecd05a%26domain%3Dws-tcg.com%26is_canvas%3Dfalse%26origin%3Dhttps%253A%252F%252Fws-tcg.com%252Ffa6738fd72072fd38%26relation%3Dparent.parent&amp;container_width=0&amp;href=https%3A%2F%2Fws-tcg.com%2Fproducts%2Ff-td2%2F&amp;layout=button_count&amp;locale=ja_JP&amp;sdk=joey&amp;send=false&amp;show_faces=true&amp;width=100" style="border: none; visibility: visible; width: 130px; height: 28px;" class=""></iframe></span></div></li>
</ul>
<!-- /social-buttons -->
<div class="alignright">
<img width="400" height="400" src="https://ws-tcg.com/wordpress/wp-content/uploads/TD_NOW-PRINTING.png" class="attachment-post-thumbnail size-post-thumbnail wp-post-image" alt="" srcset="https://ws-tcg.com/wordpress/wp-content/uploads/TD_NOW-PRINTING.png 400w, https://ws-tcg.com/wordpress/wp-content/uploads/TD_NOW-PRINTING-150x150.png 150w, https://ws-tcg.com/wordpress/wp-content/uploads/TD_NOW-PRINTING-300x300.png 300w" sizes="(max-width: 400px) 100vw, 400px"></div>
<p class="release"><strong>2024/10/25(Fri) 発売</strong><br>
【
タイトル区分：富士見ファンタジア文庫/ 】
</p>
<div class="section">カード50枚入り構築済みデッキ<br>
クイックマニュアル・デッキ解説書・プレイマット同梱<br>
<br>
1個<br>
希望小売価格1,650円(税込)<br>
<br>
1ボックス 6個入り<br>
希望小売価格9,900円(税込)</div>
<div class="section">富士見ファンタジア文庫の人気ライトノベルが<br>
ヴァイスシュヴァルツに再び参戦！！<br>
買ってすぐに遊べるトライアルデッキ！<br>
これからヴァイスシュヴァルツをはじめる方におすすめの商品です！<br>
<br>
<b>■箔押しサインカード情報（順不同）</b><br>
【冴えない彼女の育てかた】<br>
澤村・スペンサー・英梨々 役：大西沙織さん<br>
<br>
【デート・ア・ライブ】<br>
五河琴里 役：竹達彩奈さん<br>
<br>
【ロクでなし魔術講師と禁忌教典】<br>
ルミア＝ティンジェル 役：宮本侑芽さん<br>
<br>
<b>■収録作品一覧(順不同)</b><br>
・VTuberなんだが配信切り忘れたら伝説になってた<br>
・キミと僕の最後の戦場、あるいは世界が始まる聖戦<br>
・スパイ教室<br>
・スレイヤーズ	<br>
・デート・ア・ライブ  	<br>
・ハイスクールD×D<br>
・フルメタル・パニック!<br>
・ロクでなし魔術講師と禁忌教典 <br>
・冴えない彼女の育てかた	<br>
・生徒会の一存<br>
・転生王女と天才令嬢の魔法革命<br>
</div>
<div class="section"></div>
</div>

</div>
`

func TestExtractProductInfoUnexpectedTitle(t *testing.T) {
	reader := strings.NewReader(productHTMLUnexpectedTitle)
	doc, err := goquery.NewDocumentFromReader(reader)
	if err != nil {
		t.Error("Error parsing HTML:", err)
		return
	}
	if _, err := extractProductInfo(doc); err == nil {
		t.Error("Didn't get expected error")
	}
}
