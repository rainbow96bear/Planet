<script lang="ts">
    import { goto } from '$app/navigation'
    import { page } from '$app/stores'
    import type { PageData } from './$types'
    import './page.css'

    let { data }: { data: PageData } = $props()

    let q = $state(data.q)
    let users = $state(data.users)

    // 검색어 바뀌면 데이터 동기화
    $effect(() => {
        q = data.q
        users = data.users
    })
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
                </div>
            {/each}
        {/if}
    </div>
</div>