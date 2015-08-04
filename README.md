# これは何をするかと言うと

* Redmineで自分が担当者になっているチケットの一覧を表示します。

# Install

* `go get github.com/pankona/mymine`

# Usage

* 以下の環境変数を設定します。

  * REDMINE_URL
    * redmine へのURL

  * REDMINE_API_KEY
    * redmine の REST API を利用するための API KEY

  * 以下のコマンドを実行すると、自分担当のチケット一覧が表示されます。

    * `$ mymine`

  * `-o {ticket num}` のオプションをつけると、指定されたチケットをブラウザで開きます。

    * `$ mymine -o 12345`

# License

* MIT

# Contribution

* Any contribution is welcome!
