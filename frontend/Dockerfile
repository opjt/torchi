# Build stage
FROM --platform=$BUILDPLATFORM node:20-alpine AS build-stage

WORKDIR /app

# pnpm 설치
RUN corepack enable && corepack prepare pnpm@latest --activate

# 의존성 파일만 먼저 복사 (캐싱 최적화)
COPY package.json pnpm-lock.yaml ./

# pnpm 캐시 마운트 사용하여 설치 (빌드 시간 단축)
RUN --mount=type=cache,id=pnpm,target=/pnpm/store \
    pnpm install --frozen-lockfile

# 소스 코드 복사, 빌드 실행
COPY . .
RUN pnpm run build

# Production stage
FROM --platform=linux/amd64 nginx:alpine

COPY --from=build-stage /app/build /usr/share/nginx/html

# Nginx 설정 파일 복사 (SPA 라우팅 지원)
COPY nginx.conf /etc/nginx/conf.d/default.conf

# bootstrap 스크립트 복사 및 실행 권한 부여
COPY bootstrap.sh /bootstrap.sh

EXPOSE 80

# 컨테이너 시작 시 실행될 스크립트 설정
ENTRYPOINT ["/bootstrap.sh"]

CMD ["nginx", "-g", "daemon off;"]