import type { ReactionType, ToggleReactionBody } from '$lib/types/reaction';

export const addTaskReaction = async (taskId: string, type: ReactionType) => {
	const res = await fetch(`/api/v1/tasks/${taskId}/reactions`, {
		method: 'POST',
		headers: {
			'Content-Type': 'application/json'
		},
		body: JSON.stringify({
			type
		} satisfies ToggleReactionBody)
	});

	if (!res.ok) {
		const text = await res.text();
		const err = text ? JSON.parse(text) : {};
		throw new Error(err.error ?? '서버 오류');
	}
};

export const removeTaskReaction = async (taskId: string, type: ReactionType) => {
	const res = await fetch(`/api/v1/tasks/${taskId}/reactions`, {
		method: 'DELETE',
		headers: {
			'Content-Type': 'application/json'
		},
		body: JSON.stringify({
			type
		} satisfies ToggleReactionBody)
	});

	if (!res.ok) {
		const text = await res.text();
		const err = text ? JSON.parse(text) : {};
		throw new Error(err.error ?? '서버 오류');
	}
};
