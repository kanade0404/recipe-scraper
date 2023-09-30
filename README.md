# recipe-scraper

## Getting Started

### 1. Install CLI tools

- [air](https://github.com/cosmtrek/air)
- [direnv](https://github.com/direnv/direnv)
- [psqldef](https://github.com/k0kubun/sqldef/releases)
- [taskfile](https://taskfile.dev/installation/)
- [table-to-go](https://github.com/fraenky8/tables-to-go)
- [gomodifytags](https://github.com/fatih/gomodifytags)

### 2. Setup Commands

```shell
cp .envrc.sample .envrc

direnv allow

docker compose up -d

task migrate

task generate
```


### 3. Run Server

```shell
air .
```


## Architecture

https://github.com/golang-standards/project-layout

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

#### handlers

HTTPリクエストを受け付けてHTTPレスポンスを返す責務を持っています。

#### infrastructures

外部とのデータI/Oに対する責務を持ちます。

#### usecases

handlersから受けたリクエストに対する処理を行うinfrastructuresを束ねる役割を持っています。

### scripts

アプリケーションの処理に直接関わらないスクリプトが入っています。
