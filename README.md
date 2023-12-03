Docker で立てたテスト用の RDB に対する操作が競合しないようにするためのアプローチとして、テストが実行されるパッケージ単位でランダムな名前のデータベースを生成して利用するサンプルコード。

# テストの実行方法

direnv が必要。

```shell
cd example
direnv allow
docker compose up -d
```

ときたま失敗するパターン。
失敗しない場合は `-count 10` の数を増やすと通りづらくなるはず。

```shell
go test -race -shuffle on -count 10 ./...
```

今回紹介するアプローチを使ったパターン。

```shell
TESTDB_ISOLATE=true go test -race -shuffle on -count 10 ./...
```
