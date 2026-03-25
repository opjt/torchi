# Torchi Codemap

## 라우트 구조 (router.go)

```
/api/
├── /health            → handler/health.go:Check
├── /v1 (Rate Limit)   → handler/api.go
│   ├── POST /push/{token}      → Push
│   ├── POST /push/{token}/ask  → Ask (PushAndWait)
│   ├── POST /react/{id}        → React
│   ├── POST /demo              → DemoPush
│   └── POST /test/{token}      → TestPush
├── /auth              → handler/auth.go
│   ├── GET  /github/callback   → OauthGithubCallback
│   ├── POST /guest             → GuestLogin
│   ├── POST /refresh           → Refresh
│   └── POST /logout            → Logout
└── (AuthMiddleware)
    ├── /subscriptions → handler/subscription.go
    │   ├── POST /              → Subscribe
    │   └── DELETE /            → Unsubscribe
    ├── /users         → handler/user.go
    │   ├── GET  /whoami        → Whoami
    │   ├── POST /terms-agree   → TermsAgree
    │   └── DELETE /            → Withdraw
    ├── /endpoints     → handler/endpoint.go
    │   ├── POST /              → Add
    │   ├── GET  /              → GetList
    │   ├── DELETE /{token}     → Delete
    │   ├── POST /{token}/mute  → Mute
    │   └── DELETE /{token}/mute→ Unmute
    ├── /notifications → handler/notification.go
    │   ├── GET  /              → GetList
    │   ├── POST /read-until    → Read
    │   └── DELETE /{id}        → Delete
    └── /sse           → handler/sse.go
        └── GET /notifications  → Stream
```

## 도메인 의존성

```
handler/api.go
  → push/service.go:Push, PushAndWait, React
    → endpoint/service.go:FindByToken
    → token/service.go:FindByUserID
    → notifications/service.go:Register, UpdateStatusSent, SaveReaction
    → sse/broker.go:Publish

handler/auth.go
  → auth/service.go:OauthGithubFlow, GuestLogin, Refresh
    → user/service.go:UpsertByProvider, UpsertGuestUser
    → pkg/token:CreatePairToken, Validate

handler/subscription.go
  → push/service.go:Subscribe, Unsubscribe
    → token/service.go:Register, Unregister

handler/notification.go
  → notifications/service.go:GetListWithCursor, MarkAllAsRead, MarkDelete

handler/endpoint.go
  → endpoint/service.go:Add, GetList, Delete, Mute, Unmute

handler/sse.go
  → sse/broker.go:Subscribe, Unsubscribe
```

## 프론트엔드 ↔ 백엔드 매핑

```
frontend/src/lib/api/
├── endpoints.ts
│   ├── addEndpoint()      → POST /endpoints
│   ├── fetchEndpoints()   → GET /endpoints
│   ├── deleteEndpoint()   → DELETE /endpoints/{token}
│   ├── muteEndpoint()     → POST /endpoints/{token}/mute
│   └── unmuteEndpoint()   → DELETE /endpoints/{token}/mute
├── notifications.ts
│   ├── getNotifications() → GET /notifications?cursor&limit&endpoint_id&query
│   ├── markAsReadUntil()  → POST /notifications/read-until
│   ├── deleteNotification()→ DELETE /notifications/{id}
│   └── postReaction()     → POST /v1/react/{id}
├── push.ts
│   ├── checkSubscription()→ POST /subscriptions/check
│   └── subscribe()        → POST /subscriptions
└── user.ts
    ├── fetchWhoami()      → GET /users/whoami
    ├── agreeToTerms()     → POST /users/terms-agree
    └── withdraw()         → DELETE /users
```

## 페이지 → API 사용

```
routes/app/+page.svelte (알림 목록)
  → notifications.ts: getNotifications, markAsReadUntil, deleteNotification, postReaction
  → endpoints.ts: fetchEndpoints
  → SSE: EventSource('/api/sse/notifications')

routes/app/+layout.svelte (레이아웃)
  → auth/auth.ts: init
  → push.ts: handleSubscribe

routes/app/setting/+page.svelte (설정)
  → endpoints.ts: fetchEndpoints, addEndpoint, deleteEndpoint, muteEndpoint, unmuteEndpoint
  → user.ts: withdraw
  → push.ts: subscribe (테스트 알림)

routes/app/welcome/+page.svelte (약관 동의)
  → user.ts: agreeToTerms

routes/app/guide/+page.svelte (가이드)
  → API 호출 없음
```

## DTO 매핑 (알림)

```
DB: notifications 테이블
  → notifications/entity.go:Noti
  → handler/notification.go:resNoti (JSON 응답)
  → SSE: push/service.go:publishSSE (map[string]interface{})
  ↕
  → frontend/notifications.ts:NotificationApiResponse (API 타입)
  → frontend/notifications.ts:DisplayNotification (UI 타입)
  → transformNotification() 으로 변환
```

## 핵심 흐름: 푸시 발송

```
POST /api/v1/push/{token} -d 'msg=hello'
  1. handler/api.go:Push
  2. push/service.go:Push
     a. endpoint/service.go:FindByToken → endpoint 확인
     b. token/service.go:FindByUserID → 디바이스 토큰 조회
     c. notifications/service.go:Register → DB 기록
     d. push/service.go:publishSSE → sse/broker.go:Publish → handler/sse.go → 클라이언트
     e. push/service.go:pushNotification → webpush-go → 브라우저
     f. notifications/service.go:UpdateStatusSent
```

## 핵심 흐름: 리액션 대기 (ask)

```
POST /api/v1/push/{token}/ask -d 'msg=배포?' -d 'actions=승인,거절' -d 'timeout=300'
  1. handler/api.go:Ask → push/service.go:PushAndWait
  2. 푸시 발송 (위와 동일)
  3. push/waitmap.go:Set(notiID, ch) → 채널 대기
  4. 사용자가 POST /v1/react/{notiID} -d 'reaction=승인'
     → push/service.go:React → waitmap.Get(notiID) → ch <- "승인"
  5. PushAndWait 채널 수신 → 응답 반환 "승인"
  (타임아웃 시 notifications/service.go:UpdateStatusTimeout)
```

## DB 관계

```
users (1) ──→ (*) push_tokens    ON DELETE CASCADE
      (1) ──→ (*) endpoints      ON DELETE CASCADE
      (1) ──→ (*) notifications  ON DELETE CASCADE
endpoints (1) ──→ (*) notifications  ON DELETE SET NULL
```
