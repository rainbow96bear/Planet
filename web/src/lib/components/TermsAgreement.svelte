<script lang="ts">
  let { agreeTerms = $bindable(false), agreePrivacy = $bindable(false) } = $props()

  let termsOpen = $state(false)
  let privacyOpen = $state(false)

  const allChecked = $derived(agreeTerms && agreePrivacy)

  function toggleAll(e: Event) {
    const checked = (e.target as HTMLInputElement).checked
    agreeTerms = checked
    agreePrivacy = checked
  }
</script>

<div class="terms-wrapper">
  <label class="terms-all">
    <input type="checkbox" checked={allChecked} onchange={toggleAll} />
    <div class="terms-all-text">
      <span class="terms-all-label">전체 동의</span>
      <span class="terms-all-desc">이용약관 및 개인정보 처리방침에 모두 동의합니다.</span>
    </div>
  </label>

  <div class="terms-divider"></div>

  <!-- 이용약관 -->
  <div class="terms-item">
    <div class="terms-item-row">
      <label class="terms-check">
        <input type="checkbox" bind:checked={agreeTerms} />
        <span>[필수] 이용약관 동의</span>
      </label>
      <button
        type="button"
        class="terms-toggle"
        class:open={termsOpen}
        onclick={() => (termsOpen = !termsOpen)}
        aria-expanded={termsOpen}
      >
        {termsOpen ? '접기' : '보기'}
      </button>
    </div>

    {#if termsOpen}
      <div class="terms-box" role="region" aria-label="이용약관 내용">
        <h3>Planet 이용약관</h3>
        <p class="terms-date">시행일: 2026년 6월 22일</p>

        <h4>제1조 (서비스 소개)</h4>
        <p>Planet은 일정 기록과 커뮤니케이션을 위한 커뮤니티 서비스입니다.</p>
        <p>본 서비스는 개인 개발자가 운영하는 MVP(시범 서비스)로 제공되며 기능과 운영 방식은 수시로 변경될 수 있습니다.</p>

        <h4>제2조 (서비스 이용)</h4>
        <p>서비스 이용을 위해 Naver, Kakao 등 OAuth 로그인이 필요합니다.</p>
        <p>이용자는 본인의 계정을 사용하여 서비스를 이용해야 하며 타인의 계정을 무단으로 사용해서는 안 됩니다.</p>

        <h4>제3조 (이용자의 의무)</h4>
        <p>이용자는 다음 행위를 해서는 안 됩니다.</p>
        <ul>
          <li>타인을 사칭하거나 타인의 계정을 사용하는 행위</li>
          <li>불법 정보 또는 타인의 권리를 침해하는 내용을 게시하는 행위</li>
          <li>욕설, 혐오 표현, 음란물 등 타인에게 피해를 주는 행위</li>
          <li>서비스 운영을 방해하는 행위</li>
          <li>관련 법령을 위반하는 행위</li>
        </ul>

        <h4>제4조 (게시물)</h4>
        <p>게시물에 대한 책임은 해당 게시물을 작성한 이용자에게 있습니다.</p>
        <p>운영자는 법령 위반, 신고 접수 또는 서비스 운영상 필요하다고 판단되는 게시물을 사전 통지 없이 수정 또는 삭제할 수 있습니다.</p>

        <h4>제5조 (이용 제한)</h4>
        <p>운영자는 약관 위반 또는 서비스 운영에 피해를 주는 행위가 확인될 경우 서비스 이용을 제한하거나 계정을 삭제할 수 있습니다.</p>

        <h4>제6조 (서비스 변경 및 종료)</h4>
        <p>Planet은 MVP 서비스이므로 기능 추가, 변경 또는 제거가 언제든지 이루어질 수 있습니다.</p>
        <p>운영자는 서비스의 일부 또는 전부를 변경하거나 종료할 수 있습니다.</p>
        <p>서비스 종료 시 저장된 데이터는 삭제될 수 있으며 이용자는 중요한 데이터를 별도로 보관해야 합니다.</p>

        <h4>제7조 (면책)</h4>
        <p>Planet은 무료로 제공되는 개인 개발 프로젝트입니다.</p>
        <p>운영자는 서비스의 무중단 제공, 데이터 보존 또는 특정 목적에 대한 적합성을 보장하지 않습니다.</p>
        <p>서버 장애, 네트워크 장애, 해킹, OAuth 제공자의 장애 등 운영자가 통제하기 어려운 사유로 발생한 문제에 대해서는 법령상 허용되는 범위 내에서 책임을 제한합니다.</p>
        <p>이용자가 게시하거나 공유한 정보로 인해 발생하는 문제와 이용자 간 분쟁에 대한 책임은 해당 이용자에게 있습니다.</p>
        <p>다만 운영자의 고의 또는 중대한 과실이 있는 경우에는 관련 법령에 따릅니다.</p>

        <h4>제8조 (약관 변경)</h4>
        <p>운영자는 필요 시 본 약관을 변경할 수 있으며 변경 사항은 서비스 내 공지사항을 통해 안내합니다.</p>

        <h4>제9조 (준거법)</h4>
        <p>본 약관은 대한민국 법률에 따라 해석 및 적용됩니다.</p>
      </div>
    {/if}
  </div>

  <!-- 개인정보 처리방침 -->
  <div class="terms-item">
    <div class="terms-item-row">
      <label class="terms-check">
        <input type="checkbox" bind:checked={agreePrivacy} />
        <span>[필수] 개인정보 처리방침 동의</span>
      </label>
      <button
        type="button"
        class="terms-toggle"
        class:open={privacyOpen}
        onclick={() => (privacyOpen = !privacyOpen)}
        aria-expanded={privacyOpen}
      >
        {privacyOpen ? '접기' : '보기'}
      </button>
    </div>

    {#if privacyOpen}
      <div class="terms-box" role="region" aria-label="개인정보 처리방침 내용">
        <h3>개인정보 처리방침</h3>
        <p class="terms-date">시행일: 2026년 6월 22일</p>

        <p>Planet(이하 "서비스")은 개인정보보호법에 따라 이용자의 개인정보를 보호하기 위해 다음과 같이 개인정보 처리방침을 공개합니다.</p>

        <h4>제1조 (수집하는 개인정보)</h4>
        <ul>
          <li><strong>OAuth 로그인 정보:</strong> Naver, Kakao 등 OAuth 제공자가 제공하는 고유 식별자</li>
          <li><strong>프로필 정보:</strong> 이메일, 닉네임, 프로필 이미지(제공된 경우)</li>
          <li><strong>자동 생성 정보:</strong> 접속 IP, 쿠키, 서비스 이용 기록</li>
        </ul>

        <h4>제2조 (수집 목적)</h4>
        <ul>
          <li>회원 식별 및 로그인 처리</li>
          <li>서비스 제공 및 운영</li>
          <li>부정 이용 방지</li>
          <li>서비스 개선 및 오류 분석</li>
        </ul>

        <h4>제3조 (보유 기간)</h4>
        <p>
          회원 탈퇴 시 개인정보는 삭제하는 것을 원칙으로 합니다.
        </p>

        <p>
          다만 게시물, 댓글 등 다른 이용자와 공유된 콘텐츠는 서비스 운영상 남아있거나 익명화될 수 있습니다.
        </p>
        <ul>
          <li>접속 기록: 3개월</li>
        </ul>

        <h4>제4조 (개인정보 제공)</h4>
        <p>서비스는 이용자의 개인정보를 제3자에게 판매하거나 제공하지 않습니다.</p>
        <p>다만 법령에 따른 요청이 있는 경우 예외로 할 수 있습니다.</p>

        <h4>제5조 (이용자의 권리)</h4>
        <ul>
          <li>회원 탈퇴</li>
          <li>개인정보 삭제 요청</li>
          <li>개인정보 처리 관련 문의</li>
        </ul>

        <h4>제6조 (개인정보 보호)</h4>
        <p>서비스는 개인정보 보호를 위해 합리적인 보안 조치를 적용하기 위해 노력합니다.</p>
        <p>다만 인터넷 환경의 특성상 완전한 보안을 보장할 수는 없습니다.</p>

        <h4>제7조 (처리방침 변경)</h4>
        <p>본 처리방침은 법령 또는 서비스 정책 변경에 따라 수정될 수 있으며 변경 시 서비스 내 공지합니다.</p>
      </div>
    {/if}
  </div>
</div>

<style>
  .terms-wrapper {
      margin: 1.25rem 0 0.5rem;
      background: var(--surface);
      border: 1px solid var(--border);
      border-radius: 12px;
      overflow: hidden;
  }

  /* ---------- 전체 동의 ---------- */

  .terms-all {
      display: flex;
      align-items: center;
      gap: 0.75rem;
      padding: 0.9rem 1rem;
      background: rgba(79, 156, 249, 0.05);
      border-bottom: 1px solid var(--border);
      cursor: pointer;
      user-select: none;
  }

  .terms-all input[type='checkbox'] {
      width: 17px;
      height: 17px;
      flex-shrink: 0;
      accent-color: var(--planet-primary);
      cursor: pointer;
  }

  .terms-all-text {
      display: flex;
      flex-direction: column;
      gap: 2px;
  }

  .terms-all-label {
      font-size: 0.85rem;
      font-weight: 600;
      color: var(--text-primary);
  }

  .terms-all-desc {
      font-size: 0.75rem;
      color: var(--text-secondary);
  }

  .terms-divider {
      height: 1px;
      background: var(--border);
  }

  /* ---------- 개별 항목 ---------- */

  .terms-item {
      border-bottom: 1px solid var(--border);
  }

  .terms-item:last-child {
      border-bottom: none;
  }

  .terms-item-row {
      display: flex;
      align-items: center;
      justify-content: space-between;
      gap: 0.75rem;
      padding: 0.85rem 1rem;
      background: var(--surface);
  }

  .terms-check {
      display: flex;
      align-items: center;
      gap: 0.55rem;
      cursor: pointer;
      user-select: none;
      font-size: 0.82rem;
      color: var(--text-secondary);
  }

  .terms-check input[type='checkbox'] {
      width: 15px;
      height: 15px;
      flex-shrink: 0;
      accent-color: var(--planet-primary);
      cursor: pointer;
  }

  .terms-toggle {
      display: inline-flex;
      align-items: center;
      justify-content: center;
      height: 28px;
      padding: 0 0.8rem;
      border: 1px solid var(--border);
      border-radius: 8px;
      background: var(--surface);
      color: var(--text-secondary);
      font-size: 0.75rem;
      font-weight: 600;
      cursor: pointer;
      white-space: nowrap;
      transition:
          background-color .2s,
          border-color .2s,
          color .2s;
  }

  .terms-toggle:hover {
      border-color: var(--planet-primary);
      color: var(--planet-primary);
  }

  .terms-toggle.open {
      background: var(--planet-primary);
      border-color: var(--planet-primary);
      color: #fff;
  }

  /* ---------- 약관 ---------- */

  .terms-box {
      margin: 0 1rem 1rem;
      padding: 1rem;
      background: var(--bg);
      border: 1px solid var(--border);
      border-radius: 10px;
      max-height: 240px;
      overflow-y: auto;

      font-size: 0.8rem;
      line-height: 1.75;
      color: var(--text-secondary);

      scrollbar-width: thin;
      scrollbar-color: var(--border) transparent;
  }

  .terms-box::-webkit-scrollbar {
      width: 6px;
  }

  .terms-box::-webkit-scrollbar-thumb {
      background: var(--border);
      border-radius: 999px;
  }

  .terms-box h3 {
      margin: 0 0 0.25rem;
      font-size: 0.95rem;
      font-weight: 700;
      color: var(--text-primary);
  }

  .terms-date {
      margin: 0 0 1rem;
      font-size: 0.72rem;
      color: var(--text-muted);
  }

  .terms-box h4 {
      margin: 1rem 0 0.35rem;
      font-size: 0.82rem;
      font-weight: 600;
      color: var(--text-primary);
  }

  .terms-box p {
      margin: 0 0 0.45rem;
  }

  .terms-box ul {
      margin: 0.35rem 0 0.6rem;
      padding-left: 1.2rem;
  }

  .terms-box li {
      margin-bottom: 0.25rem;
  }

  .terms-box strong {
      color: var(--text-primary);
      font-weight: 600;
  }

  /* ---------- Mobile ---------- */

  @media (max-width: 520px) {
      .terms-item-row {
          padding: 0.8rem;
      }

      .terms-box {
          margin: 0 0.8rem 0.8rem;
          padding: 0.9rem;
          max-height: 220px;
      }

      .terms-toggle {
          padding: 0 0.7rem;
          font-size: 0.72rem;
      }
  }
</style>