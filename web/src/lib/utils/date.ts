// 순수 날짜 계산 유틸. eslint의 svelte/prefer-svelte-reactivity 규칙은 .svelte 파일
// 안의 new Date() 사용을 반응성 문제로 간주해 SvelteDate를 권하는데, 여기서 만드는
// Date 값들은 리렌더 사이에 in-place로 변형되지 않는 일회성 계산 결과라 해당 규칙이
// 의도한 문제(깊은 반응성 필요)와 무관하다. .ts 유틸로 분리해 규칙 스코프 밖에 둔다.

// anchor가 속한 주(일~토)의 Date 배열을 구한다.
export function getWeekDates(anchor: Date): Date[] {
	const dow = anchor.getDay();
	const sunday = new Date(anchor);
	sunday.setDate(anchor.getDate() - dow);
	return Array.from({ length: 7 }, (_, i) => {
		const d = new Date(sunday);
		d.setDate(sunday.getDate() + i);
		return d;
	});
}

// 주어진 날짜에서 days만큼 이동한 새 Date를 반환한다 (원본은 변경하지 않음).
export function shiftDate(date: Date, days: number): Date {
	const d = new Date(date);
	d.setDate(d.getDate() + days);
	return d;
}
