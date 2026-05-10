<script lang="ts">
    import { goto } from '$app/navigation'
    import { page } from '$app/stores'
    import { follow, unfollow } from '$lib/api/user'
    import type { PageData } from './$types'
    import './page.css'

    let { data }: { data: PageData } = $props()

    let q = $state(data.q)
    let users = $state(data.users)
    let followLoading = $state<number | null>(null)

    // 검색어 바뀌면 데이터 동기화
    $effect(() => {
        q = data.q
        users = data.users
    })
    async function handleFollow(userid: number) {
        followLoading = userid
        try {
            await follow(userid)
            users = users.map(u => u.userid === userid ? { ...u, is_following: true } : u)
        } catch (e) {
            console.error(e)
        } finally {
            followLoading = null
        }
    }

    async function handleUnfollow(userid: number) {
        followLoading = userid
        try {
            await unfollow(userid)
            users = users.map(u => u.userid === userid ? { ...u, is_following: false } : u)
        } catch (e) {
            console.error(e)
        } finally {
            followLoading = null
        }
    }
</script>

<div class="search-container">
    {#if data.q}
        <p class="search-result-label">
            <span class="search-keyword">"{data.q}"</span> 검색 결과 · {users.length}명
        </p>
    {/if}

    <div class="user-list">
        {#if users.length === 0 && data.q}
            <div class="empty">
                <p class="empty-text">검색 결과가 없습니다.</p>
            </div>
        {:else}
            {#each users as user}
                <div class="user-card">
                    <a href={`/profile/${user.userid}`} class="user-info">
                        <div class="user-avatar">🪐</div>
                        <div class="user-detail">
                            <span class="user-nickname">{user.nickname}</span>
                            <span class="user-username">@{user.username}</span>
                        </div>
                    </a>

                    {#if $page.data.user && $page.data.user.userid !== user.userid}
                        {#if followLoading === user.userid}
                            <button class="btn-follow" disabled>
                                <span class="spinner"></span>
                            </button>
                        {:else if user.is_following}
                            <button class="btn-follow following" onclick={() => handleUnfollow(user.userid)}>
                                팔로잉
                            </button>
                        {:else}
                            <button class="btn-follow" onclick={() => handleFollow(user.userid)}>
                                팔로우
                            </button>
                        {/if}
                    {/if}
                </div>
            {/each}
        {/if}
    </div>
</div>