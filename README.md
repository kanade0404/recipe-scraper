# recipe-scraper

## Getting Started

### 1. Install CLI tools

各ライブラリのドキュメントを見て適宜installしてください。

#### [air](https://github.com/cosmtrek/air)

golangアプリケーションのホットリロードを行ってくれるライブラリです。

#### [direnv](https://github.com/direnv/direnv)

envファイルの環境変数をそのディレクトリでのみ環境変数に書き込んでくれるライブラリです。

#### [psqldef](https://github.com/k0kubun/sqldef/releases)

sqlでのスキーマ定義に応じて宣言的にテーブルの作成・削除・変更を実行してくれるライブラリです。

- [taskfile](https://taskfile.dev/installation/)

Go Moduleにはnpm scriptsのようなタスクランナーの仕組みがないので、色々なタスクを定義して実行してくれるライブラリを入れています。

- [table-to-go](https://github.com/fraenky8/tables-to-go)

データベースのテーブル定義からGoの構造体をタグ付きで出力してくれるライブラリです。
このライブラリでdomains/modelsのファイルを生成しています。

- [gomodifytags](https://github.com/fatih/gomodifytags)

table-to-goで生成した構造体のタグを書き換えられるライブラリです。
JSONにparseする際に必要になるjsonタグをつけるために入れています。

### 2. Setup Commands

アプリケーションの実行前に以下コマンドを実行してください。

```shell
cp .envrc.sample .envrc

direnv allow

docker compose up -d

task migrate

task generate
```


### 3. Run Server

Goアプリケーションを以下コマンドで起動します。

```shell
air .
```


## Architecture

ディレクトリ構造は大枠ではgolang-standards/project-layoutに準拠しています。

https://github.com/golang-standards/project-layout

internal内の構造はlayered architectureぽくなっています。

大きくはdomain layerとinfrastructure layer、usecase layer、interfaces layerがあります。
依存のの向きは

infrastructure→domain←usecase←interfaces

になっています。

domain layerに当たるのはinternal/domains、infrastructure layerはinternal/infrastructures、usecase layerはinternal/usecase、interfaces layerはinternal/handlersです。


### cmd

コマンドラインツールのエントリーポイントが入っています。
ここのmain.goを実行することでアプリケーションが起動します。

### database

PostgreSQLの初期化やマイグレーションを行うためのSQLファイルが入っています。

### internal

ここに他の実装が入っています。

#### config

環境変数や設定を定義しています。

#### domains

アプリケーションのコアとなるデータモデルとその操作を定義しています。

domainsはどこに対しても依存しません。

#### handlers

HTTPリクエストを受け付けてusecaseにメインロジックを依頼して、返ってきた結果に応じてHTTPレスポンスを返す責務を持っています。

handlersはusecaseに対して依存します。

#### infrastructures

外部とのデータI/Oに対する責務を持ちます。

infrastructuresはdomainに対して依存します。

#### usecases

handlersから受けたリクエストに対する処理を行うinfrastructuresを束ねる役割を持っています。

usecasesはdomainに対して依存します。

### scripts

アプリケーションの処理に直接関わらないスクリプトが入っています。
