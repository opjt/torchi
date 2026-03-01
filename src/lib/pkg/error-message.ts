export const ERROR_MESSAGES = {
	AUTH_REQUIRED: '로그인이 필요한 서비스입니다.',
	INVALID_PARAMETER: '입력값이 올바르지 않습니다.',
	USER_NOT_FOUND: '존재하지 않는 사용자입니다.',
	NETWORK_ERROR: '네트워크 연결이 원활하지 않습니다.',
	DEFAULT: '알 수 없는 오류가 발생했습니다. 잠시 후 다시 시도해주세요.',
} as const;

export type ErrorCode = keyof typeof ERROR_MESSAGES;

// 에러 코드에 맞는 메시지를 찾아 반환하거나 기본 메시지를 반환
export const getErrorMessage = (code?: string, defaultMsg?: string): string => {
	if (!code) return defaultMsg || ERROR_MESSAGES.DEFAULT;
	return ERROR_MESSAGES[code as ErrorCode] || defaultMsg || ERROR_MESSAGES.DEFAULT;
};
