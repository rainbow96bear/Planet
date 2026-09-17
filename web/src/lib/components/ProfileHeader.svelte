<script lang="ts">
	import { resolve } from '$app/paths';

	let {
		profileUser,
		orbitCount,
		gravityCount,
		isOwner,
		isOrbiting,
		orbitLoading,
		onEnterOrbit,
		onLeaveOrbit
	}: {
		profileUser: {
			nickname: string;
			username: string;
			profile_image?: string;
		};
		orbitCount: number;
		gravityCount: number;
		isOwner: boolean;
		isOrbiting: boolean;
		orbitLoading: boolean;
		onEnterOrbit: () => void;
		onLeaveOrbit: () => void;
	} = $props();
</script>

<div class="profile-header">
	<div class="profile-image">
		{#if profileUser.profile_image}
			<img
				src={profileUser.profile_image}
				alt="{profileUser.nickname}님의 프로필 이미지"
				class="profile-image-photo"
			/>
		{:else}
			<span class="profile-image-placeholder">🪐</span>
		{/if}
	</div>
	<div class="profile-info">
		<h1 class="profile-nickname">{profileUser.nickname}</h1>
		<span class="profile-username">@{profileUser.username}</span>

		<div class="orbit-stats">
			<span class="orbit-stat">
				<strong>{orbitCount}</strong>
				<span class="orbit-stat-label">Orbit</span>
			</span>
			<span class="orbit-stat">
				<strong>{gravityCount}</strong>
				<span class="orbit-stat-label">Gravity</span>
			</span>
		</div>
	</div>

	<div class="profile-actions">
		{#if isOwner}
			<a href={resolve('/settings/profile')} class="action-btn secondary">
				<svg
					width="14"
					height="14"
					viewBox="0 0 24 24"
					fill="none"
					stroke="currentColor"
					stroke-width="2"
					stroke-linecap="round"
					stroke-linejoin="round"
				>
					<path d="M11 4H4a2 2 0 0 0-2 2v14a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2v-7" />
					<path d="M18.5 2.5a2.121 2.121 0 0 1 3 3L12 15l-4 1 1-4 9.5-9.5z" />
				</svg>
				프로필 수정
			</a>
		{:else if orbitLoading}
			<button class="action-btn primary" disabled aria-label="처리 중">
				<span class="spinner"></span>
			</button>
		{:else if isOrbiting}
			<button class="action-btn orbiting" onclick={onLeaveOrbit}>
				<span class="btn-text-default">In Orbit</span>
				<span class="btn-text-hover">Leave Orbit</span>
			</button>
		{:else}
			<button class="action-btn primary" onclick={onEnterOrbit}> Enter Orbit </button>
		{/if}
	</div>
</div>

<style>
	.profile-header {
		display: flex;
		align-items: center;
		gap: var(--space-xl);
		margin-bottom: var(--space-3xl);
	}
	.profile-image {
		display: flex;
		align-items: center;
		justify-content: center;
		flex-shrink: 0;
		width: 72px;
		height: 72px;
		overflow: hidden;
		background: var(--surface-alt);
		border: 1px solid var(--border);
		border-radius: 50%;
	}
	.profile-image-photo {
		display: block;
		width: 100%;
		height: 100%;
		object-fit: cover;
	}
	.profile-image-placeholder {
		font-size: 2.5rem;
	}
	.profile-info {
		flex: 1;
		min-width: 0;
	}
	.profile-nickname {
		margin: 0;
		color: var(--text-primary);
		font-size: 1.35rem;
		font-weight: 700;
		letter-spacing: -0.02em;
	}
	.profile-username {
		margin-top: var(--space-2xs);
		color: var(--text-secondary);
		font-size: 0.9rem;
	}
	/* ---------- Orbit / Gravity stats ---------- */
	.orbit-stats {
		display: flex;
		align-items: center;
		gap: var(--space-md);
		margin-top: var(--space-xs);
	}
	.orbit-stat {
		display: inline-flex;
		align-items: center;
		gap: 5px;
		color: var(--text-secondary);
		font-size: 0.8rem;
	}
	.orbit-stat strong {
		color: var(--text-primary);
		font-weight: 700;
	}
	.orbit-stat-label {
		color: var(--text-muted);
		font-size: 0.72rem;
	}
	/* ---------- Actions ---------- */
	.profile-actions {
		display: flex;
		align-items: center;
		gap: var(--space-sm);
		flex-shrink: 0;
	}
	.action-btn {
		display: inline-flex;
		align-items: center;
		justify-content: center;
		gap: var(--space-xs);
		height: 40px;
		padding: 0 18px;
		border-radius: var(--radius-md);
		border: 1px solid transparent;
		cursor: pointer;
		font-size: 0.875rem;
		font-weight: 600;
		font-family: inherit;
		text-decoration: none;
		white-space: nowrap;
		transition:
			background var(--transition-fast),
			border-color var(--transition-fast),
			color var(--transition-fast);
	}
	.action-btn:disabled {
		opacity: 0.5;
		cursor: not-allowed;
	}
	.action-btn.primary {
		background: var(--planet-primary);
		border-color: var(--planet-primary);
		color: var(--text-on-dark);
	}
	.action-btn.primary:hover:not(:disabled) {
		background: var(--planet-primary-hover);
	}
	.action-btn.secondary {
		background: var(--surface);
		border-color: var(--border);
		color: var(--text-secondary);
	}
	.action-btn.secondary:hover:not(:disabled) {
		background: var(--surface-alt);
		color: var(--text-primary);
	}
	.action-btn.orbiting {
		position: relative;
		background: var(--surface);
		border-color: var(--border);
		color: var(--text-secondary);
	}
	.action-btn.orbiting .btn-text-hover {
		display: none;
	}
	.action-btn.orbiting:hover:not(:disabled) {
		background: rgba(var(--danger-rgb), 0.08);
		border-color: var(--danger);
		color: var(--danger);
	}
	.action-btn.orbiting:hover:not(:disabled) .btn-text-default {
		display: none;
	}
	.action-btn.orbiting:hover:not(:disabled) .btn-text-hover {
		display: inline;
	}
	/* ---------- Spinner ---------- */
	.spinner {
		width: 14px;
		height: 14px;
		border: 2px solid currentColor;
		border-top-color: transparent;
		border-radius: 50%;
		animation: spin 0.6s linear infinite;
	}
	@keyframes spin {
		to {
			transform: rotate(360deg);
		}
	}
	/* ---------- Responsive ---------- */
	@media (max-width: 520px) {
		.profile-header {
			flex-wrap: wrap;
			gap: var(--space-md);
		}
		.profile-actions {
			width: 100%;
		}
		.action-btn {
			flex: 1;
		}
	}
</style>