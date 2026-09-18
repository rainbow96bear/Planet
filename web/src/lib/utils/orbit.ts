import type { OrbitSchedule } from '$lib/types/task';

// 겹치는 스케줄 중 이름이 다른 사용자만 모아 아바타 스택을 만든다.
// 같은 사람이 그 시간대에 여러 일정을 가진 경우 아바타는 한 번만 보여준다.
// (Map 사용이 .svelte 파일 안에 있으면 svelte/prefer-svelte-reactivity에 걸리므로
// 이 계산은 리렌더 간 상태를 유지하지 않는 순수 함수라 .ts 유틸로 분리한다.)
export function dedupeParticipants(list: OrbitSchedule[]): OrbitSchedule[] {
	const seen = new Map<string, OrbitSchedule>();
	for (const s of list) {
		if (!seen.has(s.user_id)) seen.set(s.user_id, s);
	}
	return Array.from(seen.values());
}
