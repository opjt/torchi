#!/bin/sh

# env.js 파일 경로 
OUTPUT_FILE="/usr/share/nginx/html/env.js"

# 컨테이너 실행 시 주입된 환경변수를 JS 파일로 추출
echo "window.__APP_CONFIG__ = {" > $OUTPUT_FILE
echo "  VAPID_KEY: \"${VAPID_KEY}\"," >> $OUTPUT_FILE
echo "  GITHUB_CLIENT_ID: \"${GITHUB_CLIENT_ID}\"" >> $OUTPUT_FILE
echo "};" >> $OUTPUT_FILE

echo "Runtime config generated in env.js"

# Nginx 실행
exec "$@"