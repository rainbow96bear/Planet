// Task/OrbitSchedule의 시작·종료 시각을 "하루 중 분(minute)" 단위로 계산하는 공용 유틸.
// MonthCalendar(미니 타임트랙), DayView(Orbit 타임라인), +page.svelte(Orbit 상세 모달
// 겹침 계산) 여러 곳에서 같은 기준으로 시간을 다뤄야 해서 한 곳에 모아둔다.

export const DAY_MINUTES = 24 * 60;

export function startMinutes(iso: string): number {
	const d = new Date(iso);
	return d.getHours() * 60 + d.getMinutes();
}

// 자정을 넘는 일정은 그날 안에서는 24:00로 잘라 표시한다
// (실제 end_at 데이터는 그대로 유지 — 문서 5번 원칙).
export function endMinutesClamped(startIso: string, endIso: string): number {
	const start = new Date(startIso);
	const end = new Date(endIso);
	if (
		end.getFullYear() !== start.getFullYear() ||
		end.getMonth() !== start.getMonth() ||
		end.getDate() !== start.getDate()
	) {
		return DAY_MINUTES;
	}
	return end.getHours() * 60 + end.getMinutes();
}
