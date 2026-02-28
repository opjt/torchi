# Torchi x AI Agent

> AI 에이전트가 작업을 끝내면, 폰으로 알림을 받으세요.

Claude Code, Cursor, Copilot 등 AI 코딩 에이전트에게 긴 작업을 맡기고 자리를 비울 때,
**작업 완료를 푸시 알림으로 받을 수 있습니다.**

별도 설치 없이, `curl` 한 줄이면 됩니다.

## Setup

1. [torchi.app](https://torchi.app) 에서 로그인 (GitHub 또는 게스트)
2. 푸시 알림 구독
3. 엔드포인트 생성 → 토큰 복사

## Usage

### 단방향 알림

AI 에이전트의 CLAUDE.md 또는 시스템 프롬프트에 아래 내용을 추가하세요:

```markdown
작업이 완료되면 아래 명령어로 알림을 보내세요:

curl -s https://torchi.app/api/v1/push/{YOUR_TOKEN} -d "{작업 요약}"
```

### 승인/거절 알림 (Interactive)

배포, 머지 등 중요한 작업 전에 **사용자 승인**을 받을 수 있습니다:

```markdown
중요한 작업 전에 아래 명령어로 승인을 요청하세요:

curl -s "https://torchi.app/api/v1/push/{YOUR_TOKEN}/ask" \
 -d "msg={질문}" \
 -d "actions=승인,거절" \
 -d "timeout=300"

응답값(승인/거절)에 따라 다음 행동을 결정하세요.
```

> **Note:** `timeout`은 사용자의 응답을 기다리는 최대 시간(초)입니다.
> AI 에이전트의 shell timeout은 이보다 약간 길게 설정해야 합니다. (예: timeout=300 → shell timeout 310초)
> 타임아웃 시 HTTP 408이 반환되며, 에이전트는 이를 "사용자 미응답"으로 처리하면 됩니다.

## Example: Claude Code

`CLAUDE.md`에 추가:

```markdown
## Notification

- 작업이 끝나면 `curl -s https://torchi.app/api/v1/push/{YOUR_TOKEN} -d "{요약}"` 으로 알림을 보내세요.
- 사용자 판단이 필요한 경우 아래 형식으로 승인을 요청하세요:
  curl -s "https://torchi.app/api/v1/push/{YOUR_TOKEN}/ask" \
   -d "msg={질문}&actions=승인,거절&timeout=300"
  응답값에 따라 행동하세요.
  shell timeout은 API timeout보다 10초 길게 설정하세요. (예: 310초)
```

## Use Cases

| 상황                  | 알림 타입   |
| --------------------- | ----------- |
| 빌드/테스트 완료      | 단방향 알림 |
| 긴 리팩토링 작업 완료 | 단방향 알림 |
| 프로덕션 배포 전 승인 | 승인/거절   |
| PR 생성 전 확인       | 승인/거절   |
| 에러 발생 보고        | 단방향 알림 |

## Why Torchi?

- **No Install** — PWA 기반, 앱 설치 없이 웹에서 바로 구독
- **One Line** — `curl` 한 줄이면 연동 끝. SDK 불필요
- **Interactive** — 단순 알림이 아닌 승인/거절 액션 지원
- **Agent Friendly** — 어떤 AI 에이전트든 shell 접근만 되면 사용 가능
