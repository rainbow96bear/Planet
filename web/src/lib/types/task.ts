export interface Task {
	id: string;
	title: string;
	description: string;
	start_at: string;
	end_at: string;
	is_completed: boolean;
	is_public: boolean;
}

export type GetTasksByMonthResponse = Task[];

export interface CreateTaskBody {
	title: string;
	description?: string;
	start_at: string;
	end_at: string;
	is_public: boolean;
}

export interface CreateTaskResponse extends Task {
	created_at: string;
}

export interface ToggleTaskResponse {
	id: string;
	is_completed: boolean;
}

// ── Orbit Schedule ──
// My Schedule(Task)과 구분되는 타입: 작성자 정보가 붙고, is_completed가 없다
// (Orbit Schedule은 "할 일"이 아니라 "남의 공개 일정을 보는 뷰"라 완료 개념이 없음).
export interface OrbitSchedule {
	id: string;
	user_id: string;
	nickname: string;
	profile_image: string;
	title: string;
	start_at: string;
	end_at: string;
}

export type GetOrbitSchedulesByMonthResponse = OrbitSchedule[];

export interface OrbitLaneItem {
	schedule: OrbitSchedule;
	lane: number;
	startMin: number;
	endMin: number;
}

export interface OrbitOverlapGroup {
	rangeStartMin: number;
	rangeEndMin: number;
	schedules: OrbitSchedule[];
}

// UpdateTaskBody — CreateTaskBody와 필드는 동일하지만 별도 타입으로 둔다.
// 의미가 "새 일정 생성"과 "기존 일정 수정"으로 다르고, 나중에 필드가 갈라질
// 여지도 있어서 처음부터 구분해뒀다.
export interface UpdateTaskBody {
	title: string;
	description?: string;
	start_at: string;
	end_at: string;
	is_public: boolean;
}

export type UpdateTaskResponse = Task;