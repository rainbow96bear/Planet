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
        <p class="terms-date">시행일: 2026년 6월 15일</p>

        <h4>제1조 (목적)</h4>
        <p>이 약관은 Planet(이하 "서비스")이 제공하는 커뮤니티 서비스의 이용과 관련하여 서비스와 이용자의 권리, 의무 및 책임 사항을 규정함을 목적으로 합니다.</p>

        <h4>제2조 (정의)</h4>
        <p>"이용자"란 본 약관에 따라 서비스를 이용하는 회원 및 비회원을 말합니다.<br />"회원"이란 서비스에 개인정보를 제공하여 회원 등록을 한 자로, 서비스를 지속적으로 이용할 수 있는 자를 말합니다.<br />"게시물"이란 회원이 서비스를 이용하면서 게재한 문자, 사진, 파일 등의 정보를 말합니다.</p>

        <h4>제3조 (약관의 효력 및 변경)</h4>
        <p>본 약관은 서비스 화면에 게시함으로써 효력이 발생합니다. 서비스는 합리적인 사유가 발생할 경우 관련 법령에 위배되지 않는 범위에서 약관을 변경할 수 있으며, 변경된 약관은 공지사항을 통해 사전 공지합니다.</p>

        <h4>제4조 (회원가입)</h4>
        <p>이용자는 서비스가 정한 양식에 따라 정보를 기입하고 본 약관에 동의함으로써 회원가입을 신청합니다. 다음 각 호에 해당하는 경우 회원가입 신청이 거절될 수 있습니다.</p>
        <ul>
          <li>실명이 아니거나 타인의 정보를 이용하여 신청한 경우</li>
          <li>이전에 약관 위반으로 서비스 이용이 제한된 경우</li>
          <li>기타 서비스가 정한 기준에 적합하지 않은 경우</li>
        </ul>

        <h4>제5조 (이용자의 의무)</h4>
        <p>이용자는 다음 행위를 하여서는 안 됩니다.</p>
        <ul>
          <li>타인의 개인정보를 무단으로 수집하거나 도용하는 행위</li>
          <li>서비스의 운영을 방해하거나 정상적인 이용을 방해하는 행위</li>
          <li>음란, 폭력, 혐오 등 다른 이용자에게 불쾌감을 줄 수 있는 게시물을 게재하는 행위</li>
          <li>저작권 등 지적재산권을 침해하는 행위</li>
          <li>관련 법령에 위배되는 행위</li>
        </ul>

        <h4>제6조 (서비스 이용 제한)</h4> 
        <p> 서비스는 이용자가 본 약관의 의무를 위반하거나 서비스의 정상적인 운영을 방해한 경우 경고, 일시 정지, 영구 이용 정지 등의 조치를 취할 수 있습니다. </p>
       
        <h4>제7조 (게시물 관리)</h4> 
        <p> 게시물에 대한 책임은 해당 게시물을 작성한 이용자에게 있습니다. </p>
        <p> 서비스는 관련 법령 또는 본 약관에 위반되거나 서비스 운영 목적에 부적절하다고 판단되는 게시물을 사전 통지 없이 수정, 이동 또는 삭제할 수 있습니다. </p>
       
        <h4>제8조 (서비스 변경 및 종료)</h4>
        <p> 서비스는 운영상 또는 기술상 필요에 따라 서비스의 일부 또는 전부를 변경하거나 종료할 수 있습니다. </p>
        <p> 서비스 종료 또는 중대한 변경이 있는 경우 서비스 내 공지사항을 통해 안내합니다. </p>
        
        <h4>제9조 (광고)</h4> 
        <p> 서비스는 운영을 위하여 광고를 게재할 수 있습니다. </p> 
        
        <h4>제10조 (면책조항)</h4>
        <p> 서비스는 천재지변, 해킹, 시스템 장애, 네트워크 장애 등 서비스 운영자가 합리적으로 통제하기 어려운 사유로 인해 서비스를 제공할 수 없는 경우 법령상 허용되는 범위 내에서 책임을 부담하지 않습니다. </p>
        <p> 서비스는 이용자의 귀책 사유로 인한 서비스 이용 장애에 대하여 책임을 지지 않습니다. </p>

        <h4>제11조 (준거법 및 분쟁 해결)</h4>
        <p> 본 약관의 해석 및 적용에 관해서는 대한민국 법을 준거법으로 합니다. </p>
        <p> 서비스와 이용자 간에 발생한 분쟁에 관한 소송은 대한민국 법원의 관할에 따릅니다. </p> 
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
        <p class="terms-date">시행일: 2026년 6월 15일</p>

        <p>Planet(이하 "서비스")은 개인정보보호법에 따라 이용자의 개인정보를 보호하고, 이와 관련한 고충을 신속하게 처리할 수 있도록 다음과 같이 개인정보 처리방침을 수립·공개합니다.</p>

        <h4>제1조 (수집하는 개인정보 항목)</h4>
        <ul>
          <li><strong>필수 항목:</strong> 아이디(username), 닉네임</li>
          <li><strong>소셜 로그인 시 수집:</strong> OAuth 제공자의 고유 식별자, 이메일(제공된 경우)</li>
          <li><strong>자동 생성 정보:</strong> 접속 IP, 쿠키, 서비스 이용 기록, 기기 정보</li>
        </ul>

        <h4>제2조 (개인정보 수집 및 이용 목적)</h4>
        <ul>
          <li>회원 식별 및 서비스 제공</li>
          <li>서비스 이용 관련 공지사항 전달</li>
          <li>부정 이용 방지 및 서비스 안정성 확보</li>
          <li>서비스 개선을 위한 통계 분석</li>
        </ul>
        
        <h4>제3조 (개인정보 보유 및 이용 기간)</h4>
        <p> 서비스는 회원 탈퇴 시 개인정보를 삭제 또는 익명화 처리합니다. </p>
        <p> 다만 서비스 운영, 부정 이용 방지, 분쟁 대응 또는 관련 법령에 따라 보관이 필요한 경우 최소한의 정보를 일정 기간 보관할 수 있습니다. </p>
        <ul>
          <li>접속 기록: 3개월 (통신비밀보호법)</li>
        </ul>

        <h4>제4조 (개인정보의 제3자 제공)</h4>
        <p>서비스는 이용자의 개인정보를 원칙적으로 외부에 제공하지 않습니다. 다만, 이용자가 사전에 동의한 경우 또는 법령의 규정에 의거하거나 수사기관의 요청이 있는 경우에는 예외로 합니다.</p>

        <h4>제5조 (이용자의 권리)</h4>
        <ul>
          <li>개인정보 열람 요청</li>
          <li>개인정보 정정·삭제 요청</li>
          <li>개인정보 처리 정지 요청</li>
          <li>회원 탈퇴를 통한 동의 철회</li>
        </ul>

        <h4>제6조 (개인정보 관련 문의)</h4>
        <p> 개인정보 처리와 관련한 문의, 열람, 정정 또는 삭제 요청은 서비스 내 문의 기능 또는 운영자가 제공하는 연락 수단을 통해 접수할 수 있습니다. </p>

        <h4>제7조 (개인정보 보호)</h4>
        <p> 서비스는 개인정보 보호를 위하여 합리적인 보안 조치를 적용하기 위해 노력합니다. </p>
        
        <h4>제8조 (개인정보 처리방침 변경)</h4>
        <p> 본 처리방침은 법령 또는 정책의 변경에 따라 수정될 수 있으며, 변경 시에는 서비스 내 공지사항을 통해 사전 안내합니다. </p>
      </div>
    {/if}
  </div>
</div>

<style>
  .terms-wrapper {
    margin: 1.25rem 0 0.25rem;
    border: 1px solid var(--border);
    border-radius: 10px;
    overflow: hidden;
  }

  /* 전체 동의 */
  .terms-all {
    display: flex;
    align-items: center;
    gap: 10px;
    padding: 13px 14px;
    background: color-mix(in srgb, var(--planet-primary) 6%, var(--surface));
    border-bottom: 1px solid color-mix(in srgb, var(--planet-primary) 20%, transparent);
    cursor: pointer;
    user-select: none;
  }

  .terms-all input[type='checkbox'] {
    width: 17px;
    height: 17px;
    accent-color: var(--planet-primary);
    cursor: pointer;
    flex-shrink: 0;
  }

  .terms-all-text {
    display: flex;
    flex-direction: column;
    gap: 1px;
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

  /* 개별 항목 */
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
    padding: 11px 14px;
    background: var(--surface);
  }

  .terms-check {
    display: flex;
    align-items: center;
    gap: 8px;
    cursor: pointer;
    user-select: none;
    font-size: 0.82rem;
    color: var(--text-secondary);
  }

  .terms-check input[type='checkbox'] {
    width: 15px;
    height: 15px;
    accent-color: var(--planet-primary);
    cursor: pointer;
    flex-shrink: 0;
  }

  .terms-toggle {
    font-size: 0.72rem;
    color: var(--planet-primary);
    background: none;
    border: 1px solid color-mix(in srgb, var(--planet-primary) 35%, transparent);
    border-radius: 4px;
    padding: 3px 10px;
    cursor: pointer;
    white-space: nowrap;
    transition:
      background-color 0.15s,
      color 0.15s;
  }

  .terms-toggle:hover,
  .terms-toggle.open {
    background: var(--planet-primary);
    color: #fff;
  }

  /* 약관 본문 */
  .terms-box {
    margin: 0 14px 12px;
    padding: 12px 14px;
    background: var(--bg);
    border: 1px solid var(--border);
    border-radius: 8px;
    max-height: 200px;
    overflow-y: auto;
    font-size: 0.78rem;
    line-height: 1.75;
    color: var(--text-secondary);
    scrollbar-width: thin;
    scrollbar-color: var(--border) transparent;
  }

  .terms-box::-webkit-scrollbar {
    width: 4px;
  }

  .terms-box::-webkit-scrollbar-thumb {
    background: var(--border);
    border-radius: 4px;
  }

  .terms-box h3 {
    font-size: 0.85rem;
    font-weight: 700;
    color: var(--text-primary);
    margin: 0 0 2px;
  }

  .terms-box .terms-date {
    font-size: 0.72rem;
    color: var(--text-muted);
    margin: 0 0 10px;
  }

  .terms-box h4 {
    font-size: 0.78rem;
    font-weight: 600;
    color: var(--text-primary);
    margin: 12px 0 3px;
  }

  .terms-box p {
    margin: 0 0 4px;
  }

  .terms-box ul {
    margin: 3px 0 4px;
    padding-left: 16px;
  }

  .terms-box ul li {
    margin-bottom: 2px;
  }

  .terms-box strong {
    color: var(--text-primary);
    font-weight: 600;
  }
</style>