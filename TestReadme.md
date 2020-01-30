# TestService

## Describe

テスト用

<details><summary>/test</summary>

<details><summary>POST</summary>

- descirbe

HttpMethod:Postテスト

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

HttpMethod:Getテスト

- Parameter

| query name | type | required |
| :-- | :-- | :-- |
| id | int | true |
| name | string | false |


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