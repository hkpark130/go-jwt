# go-jwt

### 먼저 .env 파일을 생성할 필요가 있습니다.
```sh
% cat .env
# .env file

DB_USER="user"
DB_PW="password"
DB_NAME="go-jwt"

SECRET_KEY="park"

CONF_FILE="apache.conf"
% 
```

### 빌드

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

