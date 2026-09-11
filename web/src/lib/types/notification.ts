export interface Notification {
	id: string;
	type: 'orbit_entered' | 'comment' | 'reaction';
	message: string;
	is_read: boolean;
	created_at: string;
	actor_id: string;
	actor_nickname: string;
}
