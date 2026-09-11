export interface Feed {
	feed_id: string;
	actor_id: string;
	actor_nickname: string;
	type: 'task.created' | 'task.completed';
	task_id: string;
	task_title?: string;
	like_count: number;
	cheer_count: number;
	is_liked: boolean;
	is_cheered: boolean;
	created_at: string;
}
