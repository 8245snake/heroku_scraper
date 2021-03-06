# heroku-scraper

## 概要
Herokuで動かしてシェアサイクルの台数情報をスクレイピングをするツール。  
一応WEBサーバーの体をなしており、HTTPリクエストをトリガーにしてスクレイピングを開始する。  
ログイン情報などもHTTPで渡す想定。  

## 準備
スクレイピングしたデータを受けるサーバを用意する。  
データは以下のJSONがPOSTリクエストのペイロードに乗ってやってくる。  
```
{
  "spotinfo": [
    {
      "area": "D1",
      "spot": "10",
      "time": "2020/02/16 15:43:20",
      "count": "18"
    },
    {
      "area": "D1",
      "spot": "10",
      "time": "2020/02/16 15:40:12",
      "count": "10"
    }
  ]
}
```
area,spot,time,countを含む構造体の繰り返しであり、エリアごとに100件を上限としてリクエストされる。  
各パラメータの説明は以下の通り  

|フィールド |意味 |
|----|----|
|area |エリアコード（A,B,C,D,E,H,I,J,K,M） |
|spot |スポットコード（01～3桁の連番） |
|time |時刻 |
|count |台数 |

## 使い方
### 台数スクレイピング
エンドポイント： `/start`  
メソッド： `POST`  
|パラメータ |意味 |備考 |
|---|---|---|
|id |ログインID | |
|password |ログインパスワード | |
|areaID |スクレイピングするエリアのコード（1,2,3,5,6,4,10,12,7,8のカンマ区切り） |省略時は全てスクレイピング |
|env |秘密文字列 |環境変数に設定した場合は不要 |
|address |スクレイピング結果を受けとるコールバックURL |SSL証明書エラーは無視するのでhttpでも可 |


### マスタ更新
エンドポイント： `/master`  
メソッド： `POST`  
パラメータ：台数スクレイピングと同様だが`areaID`は無視して全て対象とする

### リカバリ
何らかの事情でスクレイピング結果の送信に失敗したとき（DBサーバが落ちてるなど）、/tmp フォルダにJSONファイルとして溜めておき、あとから送信するという仕組みがある。  

エンドポイント： `/recover`  
メソッド： `POST`  
|パラメータ |意味 |備考 |
|---|---|---|
|max |一回で処理するファイル数 |0を設定した場合/tmpにファイルがいくつ残っているかを返す |

制限事項として、台数スクレイピングもしくはマスタ更新が一回も行われていない場合はエラーを返す（初期化が行われていないため）




