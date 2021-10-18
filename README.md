# go-jwt

### docker-compose run

```sh
docker-compose build
docker-compose up -d
```

### テストコマンドを実行する時、SECRET_KEYを設定する必要があります。

```sh
SECRET_KEY="{SECRET_KEY}" go test api/handlers/auth/*
```
