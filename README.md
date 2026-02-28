# torchi

> Push notification via HTTP API

`torchi`는 HTTP API 호출만으로 모바일/데스크톱에 즉시 푸시 알림을 보내고,
알림에서 바로 **액션(수락 / 거절)** 을 처리할 수 있는 가볍고 개발자 친화적인 푸시 알림 서비스입니다.

서비스 바로가기 -> [torchi.app](https://torchi.app)

## Key Features

- **PWA Based** — 별도 앱 설치 없이 웹에서 바로 푸시 구독
- **Developer Friendly** — 엔드포인트 토큰 하나로 `curl` 한 줄이면 알림 발송
- **Interactive Actions** — 알림에서 수락/거절 등 액션을 받고 API 응답으로 결과 반환

## Quick Start

### 1. 알림 보내기

```bash
curl -X POST https://torchi.app/api/v1/push/{YOUR_TOKEN} \
  -d '배포가 완료되었습니다'
```

### 2. 액션 알림 (승인/거절)

```bash
curl -X POST https://torchi.app/api/v1/push/{YOUR_TOKEN}/ask \
  -d 'msg=프로덕션 배포를 승인하시겠습니까?&actions=승인,거절&timeout=300'
```

응답으로 사용자의 리액션(`승인` 또는 `거절`)이 반환됩니다.
타임아웃 시 `408 Request Timeout`을 반환합니다.

## Architecture

```bash
internal/
├── api/             # HTTP 핸들러, 라우터, 미들웨어
├── core/            # 앱 부트스트랩 (Uber FX)
├── domain/          # 도메인 로직
│   ├── auth/        # 인증 (GitHub OAuth, 게스트)
│   ├── user/        # 사용자 관리
│   ├── token/       # PWA 푸시 구독 토큰
│   ├── endpoint/    # 외부 서비스 엔드포인트 (API 키)
│   ├── push/        # 푸시 알림 발송
│   ├── notifications/ # 알림 이력 및 리액션
│   └── sse/         # SSE 실시간 브로커
├── infrastructure/  # PostgreSQL 연결
└── pkg/             # 공통 패키지 (config, log, jwt)
```

## Tech Stack

- **Language:** Go
- **Router:** chi
- **DI:** Uber FX
- **Database:** PostgreSQL (pgx)
- **Auth:** JWT , GitHub OAuth
- **Push:** Web Push API (VAPID)

## Development

### Prerequisites

- Go 1.25+
- Docker (PostgreSQL)

### Setup

```bash
# 1. PostgreSQL 실행
docker compose up -d

# 2. 환경변수 설정
cp .env.example .env  # 필요한 값 채우기

# 3. 서버 실행
go run . serve
```

### Environment Variables

| Variable               | Description                         |
| ---------------------- | ----------------------------------- |
| `STAGE`                | `dev` / `prod`                      |
| `SERVICE_PORT`         | 서버 포트 (기본 25565)              |
| `DB_URL`               | PostgreSQL 연결 URL                 |
| `JWT_SECRET`           | JWT 서명 키                         |
| `FRONT_URL`            | 프론트엔드 URL (CORS)               |
| `VAPID_PUBLIC_KEY`     | Web Push 공개키                     |
| `VAPID_PRIVATE_KEY`    | Web Push 비밀키                     |
| `GITHUB_CLIENT_ID`     | GitHub OAuth Client ID              |
| `GITHUB_CLIENT_SECRET` | GitHub OAuth Client Secret          |
| `LOG_LEVEL`            | `debug` / `info` / `warn` / `error` |

### VAPID Key 생성

```bash
go run . genkey
```
