export const getUnreadCount = async () => {
	const res = await fetch('/api/v1/notifications/unread-count');
	if (!res.ok) {
		const err = await res.json();
		throw new Error(err.error ?? '서버 오류');
	}

	return res.json();
};

export const getNotifications = async () => {
	const res = await fetch('/api/v1/notifications');
	if (!res.ok) {
		return [];
	}
	return res.json();
};

export const markAllRead = async () => {
	const res = await fetch('/api/v1/notifications/read-all', {
		method: 'PATCH'
	});
	if (!res.ok) {
		const err = await res.json();
		throw new Error(err.error ?? '서버 오류');
	}
};
