# Go言語で作るScraper
WantedlyでGo言語の仕事を拾ってCSVファイルを作ってくれるScraper
- Go言語意外にのプログラミング言語も検索及びcsvファイルのダウンロード可能

## 作った理由
**Go言語の勉強**のため
- 最初は基本的なTop to Bottom方式の書き方
- 以後go routineを利用したマルチスレッド的な書き方の勉強
  - go routineの強力さを感じたが、何だかんだgo routineをつけていいのか疑問
  - 性能テスト(bench)でバランスを探し出すのが重要だと感じた。

## 内容
- go言語の基本パッケージにcsv関連の強力なパッケージがある、`"encoding/csv"`
- jqueryと似ているgoqueryがある、`"github.com/PuerkitoBio/goquery"`
- 静的ファイルのServingするとき`http.ServeFile`を使う
  - ファイルの名前がuriに基づいて勝手に生成されちゃう `/scrape` -> `scrape.csv`