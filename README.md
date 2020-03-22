# go-notif-breakout

## Usage
- 環境変数: DISCORD_ID, DISCORD_TOKEN   
Discord webhookのid=数字の羅列, token=文字の羅列を環境変数にセット  
``` bash
// ex.
$ export DISCORD_ID=23940293402
$ export DISCORD_TOKEN=adfad4980alwejfaafd43kafregjeihga
$ go build
$ nohup ./go-notif-breakout -product <binance_product_code> -term <minutes> &
```

## Author
[@_numbP](https://twitter.com/_numbP)  

