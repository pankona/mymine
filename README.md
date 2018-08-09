# mymine

## これは何をするかと言うと

* Redmineで自分が担当者になっているチケットの一覧を表示します。

## Install

* `go get github.com/pankona/mymine`

## Usage

以下の環境変数を設定します。

* MYMINE_REDMINE_URL
  * Redmine へのURL

* MYMINE_REDMINE_API_KEY
  * Redmine の REST API を利用するための API KEY

* 以下のコマンドを実行すると、自分担当のチケット一覧が表示されます。

```bash
$ mymine
```

* `-o {ticket num}` のオプションをつけると、指定されたチケットをブラウザで開きます。

```bash
$ mymine -o 12345
```

## License

* MIT

## Contribution

* Any contribution is welcome!

## Author

Yosuke Akatsuka (a.k.a pankona)
