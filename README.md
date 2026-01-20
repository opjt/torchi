# 🏮 torchi

> On-demand Hook Push

`torchi`는 웹훅(Webhook)이나 API를 통해 모바일로 실시간 푸시 알림을 보내고, 알림에서 바로 **액션(수락/거절/응답 등)** 을 취할 수 있는 가볍고 개발자 친화적인 알림 플랫폼입니다.

**Key Features**

- Instant Push: Webhook 호출 즉시 모바일 푸시 알림 발송
- PWA Based: 별도의 앱 설치 없이 웹에서 바로 시작하는 푸시 환경
- Developer Friendly: 간결한 API 구조와 쉬운 연동
- Interactive Actions: 알림에서 바로 실행하는 수락/거절/응답 액션

## 개발 단계

1. MVP 구현
   - webhook/API -> push 알림
   - PWA로 구현. 추후 앱 개발
   - 간단한 상호작용 구현(accept,reject), polling 방식으로 리턴값 받도록
