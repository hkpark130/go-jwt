# go-jwt

### docker-compose run

```sh
docker-compose build
docker-compose up -d
```

### 테스트 커맨드를 실행할 때는 SECRET_KEY를 설정할 필요가 있습니다.

```sh
SECRET_KEY="{SECRET_KEY}" go test api/handlers/auth/* -v
# 또는
docker-compose run golang go test api/handlers/auth/* -v
```

### 인증 흐름
![jwt](https://img1.daumcdn.net/thumb/R1280x0/?scode=mtistory2&fname=https%3A%2F%2Fblog.kakaocdn.net%2Fdn%2Fbowbru%2FbtrvFBvL3KW%2FbavjwicJAa6KUTysGenEX0%2Fimg.png)

### 앞으로 보완할 것
403(Authorization) 기능도 Payload에 퍼미션 데이터를 넣어서 개발할 예정
