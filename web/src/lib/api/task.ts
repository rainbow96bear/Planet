import type {
	CreateTaskBody,
	CreateTaskResponse,
	UpdateTaskBody,
	UpdateTaskResponse,
	GetTasksByMonthResponse,
	GetOrbitSchedulesByMonthResponse
} from '$lib/types/task';

export const getTasksByMonth = async (
	userid: string,
	year: number,
	month: number
): Promise<GetTasksByMonthResponse> => {
	const res = await fetch(`/api/v1/users/${userid}/tasks?year=${year}&month=${month}`);

	if (!res.ok) {
		const err = await res.json();
		throw new Error(err.error ?? '서버 오류');
	}
	const tasks = await res.json();

	return tasks;
};

// Orbit Schedule은 "본인 것만" 조회 가능하다 (Backend가 이미 강제하지만,
// Frontend에서도 다른 유저 프로필에서는 이 함수 자체를 호출하지 않도록 한다 — Phase 5/6에서 처리).
export const getOrbitSchedulesByMonth = async (
	userid: string,
	year: number,
	month: number
): Promise<GetOrbitSchedulesByMonthResponse> => {
	const res = await fetch(`/api/v1/users/${userid}/orbit-schedules?year=${year}&month=${month}`);

	if (!res.ok) {
		const err = await res.json();
		throw new Error(err.error ?? '서버 오류');
	}
	const schedules = await res.json();

	return schedules;
};

export const createTask = async (body: CreateTaskBody): Promise<CreateTaskResponse> => {
	const res = await fetch('/api/v1/tasks', {
		method: 'POST',
		headers: { 'Content-Type': 'application/json' },
		body: JSON.stringify(body)
	});

	if (!res.ok) {
		const err = await res.json();
		throw new Error(err.error ?? '서버 오류');
	}

	return res.json();
};

export const updateTask = async (
	taskId: string,
	body: UpdateTaskBody
): Promise<UpdateTaskResponse> => {
	const res = await fetch(`/api/v1/tasks/${taskId}`, {
		method: 'PATCH',
		headers: { 'Content-Type': 'application/json' },
		body: JSON.stringify(body)
	});

	if (!res.ok) {
		const err = await res.json();
		throw new Error(err.error ?? '서버 오류');
	}

	return res.json();
};

export const deleteTask = async (taskId: string): Promise<void> => {
	const res = await fetch(`/api/v1/tasks/${taskId}`, {
		method: 'DELETE'
	});

	if (!res.ok) {
		const err = await res.json();
		throw new Error(err.error ?? '서버 오류');
	}
};

export const toggleTask = async (taskId: string): Promise<void> => {
	const res = await fetch(`/api/v1/tasks/${taskId}/toggle`, {
		method: 'POST'
	});

	if (!res.ok) {
		throw new Error('Failed to toggle task');
	}
};
