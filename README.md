# httproute wrapper [![Build Status](https://travis-ci.com/ncs-nozomi-nishinohara/httproute_wrapper.svg?branch=master)](https://travis-ci.com/ncs-nozomi-nishinohara/httproute_wrapper) [![godoc](https://godoc.org/github.com/ncs-nozomi-nishinohara/httprouter_wrapper?status.svg)](https://godoc.org/github.com/ncs-nozomi-nishinohara/httprouter_wrapper)

[httprouter](github.com/julienschmidt/httprouter)をラップして、  
構造体にメソッドを追加し、yaml ファイルに設定していくだけで  
自動的に net/http の web サーバーが起動できる。

## Features

yaml ファイルにポート番号,url,構造体から呼び出すメソッドを記述するだけで  
net/http の web サーバーが起動できる。  
さらに yaml ファイルに attribure を設定し、Readme を出力する設定でアプリケーションを起動すると  
markdown ファイル生成される。  
`/refarence`にアクセスすると markdown ファイルが閲覧可能となる  
Swagger が使用できる場合は Swagger を使用した方が良いと思う。

## Usage

.yaml ファイルと main の実装方法を記載します。

```yaml:test_service.yaml
アプリ名:
  describe: アプリ概要記述
  port: ポート番号 #int or string
  paths:
    /パス:
      methods:
        get:
          func: Index # 任意のメソッド名
          attribute:
            describe: 機能概要記述
            parameter: パラメータ記述
```

```go:main.go
package main

import (
    "fmt"
    "httprouter_wrapper"
    "log"
    "net/http"

    "github.com/julienschmidt/httprouter"
)

// レシーバ定義用構造体
type Routes struct{}

// メソッドの定義(ビジネスロジックを記載)
func (Routes) Index(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
    w.Header().Set("Content-Type", "text/plain; charset=utf-8")
    w.WriteHeader(http.StatusOK)
    fmt.Fprint(w, "Welecom\n")
}

func main() {
    // コンストラクタ
    router := httprouter_wrapper.NewRouterWrapperHandler("test_service.yaml", httprouter_wrapper.ReadMe{
        // Markdownファイルを生成するかどうか
        Write:    true,
        // 生成するMarkdownファイル名
        Filename: "TestReadme.md",
    })
    r := Routes{}
    // 構造体をセット
    router.Router = r
    // Handlerのコンストラクタ
    httprouter_wrapper.New(router)
    // サーバー起動
    log.Fatalln(router.ListenServe())
}
```

.yaml の`patameter`の部分については以下の記載も可能となっています。

```yaml:param_json.yaml
TestService:
  describe: テスト用
  port: 8080 # or "8080"
  paths:
    /test:
      methods:
        post:
          func: Post
          attribute:
            describe: HttpMethod:Postテスト
            parameter: '
            {
              "id": "int",
              "name": "int",
              "objs": []
            }'
```

```yaml:param_table.yaml
TestService:
  describe: テスト用
  port: 8080 # or "8080"
  paths:
    /test:
      methods:
        post:
          func: Post
          attribute:
            describe: HttpMethod:Postテスト
            parameter: '
            {
              "id": "int",
              "name": "int",
              "objs": []
            }'
```

`handler_test.go`で使用された.yaml ファイルの記述例

[test_service.yaml](test_service.yaml)

```yaml:test_service.yaml
TestService:
  describe: テスト用
  port: 8080 # or "8080"
  paths:
    /test:
      methods:
        get:
          func: Get
          attribute:
            describe: HttpMethod:Getテスト
            parameter: |+
              | query name | type | required |
              | :-- | :-- | :-- |
              | id | int | true |
              | name | string | false |
        post:
          func: Post
          attribute:
            describe: HttpMethod:Postテスト
            parameter: '
            {
              "id": "int",
              "name": "int",
              "objs": []
            }'
    /all_method:
      methods:
        get:
          func: AllMethod
          attribute:
            describe: メソッド内で分岐(Get)
        post:
          func: AllMethod
          attribute:
            describe: メソッド内で分岐(Post)
        put:
          func: AllMethod
          attribute:
            describe: メソッド内で分岐(Put)
        delete:
          func: AllMethod
          attribute:
            describe: メソッド内で分岐(Delete)
```

上記の`.yaml`から生成される`Markdown`

[TestReadme.md](TestReadme.md)

````markdown:TestReadme.md
# TestService

## Describe

テスト用

<details><summary>/test</summary>

<details><summary>POST</summary>

- descirbe

HttpMethod:Post テスト

- Parameter

```json:paramete.json
{
  "id": "int",
  "name": "int",
  "objs": []
}
```

</details>

<details><summary>GET</summary>

- descirbe

HttpMethod:Get テスト

- Parameter

| query name | type   | required |
| :--------- | :----- | :------- |
| id         | int    | true     |
| name       | string | false    |

</details>

</details><details><summary>/all_method</summary>

<details><summary>GET</summary>

- descirbe

メソッド内で分岐(Get)

</details>

<details><summary>POST</summary>

- descirbe

メソッド内で分岐(Post)

</details>

<details><summary>PUT</summary>

- descirbe

メソッド内で分岐(Put)

</details>

<details><summary>DELETE</summary>

- descirbe

メソッド内で分岐(Delete)

</details>

</details>
````

## TODO

- `Response`キーの対応

- `parameter`が`json`の時に struct を自動生成する機能

- moc サーバー自動生成(path 内容に応じて.go ファイルを出力)  
  `Response`内容については`Response`キーをセットする
