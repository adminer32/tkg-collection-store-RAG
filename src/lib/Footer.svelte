<script lang="ts">
  // ── 可覆盖 slot 内容，或直接使用默认样式 ──────────────────────────────────

  const colLinks = [
    { title: '服务', links: [
      { text: '明信片代寄', href: '/shop/sent?service=postcard-forward' },
      { text: '邮票零售', href: '/shop' },
      { text: '信件定制', href: '/shop/sent?addon=letter-custom' },
      { text: '礼盒套装', href: '/shop/sent?addon=gift-box' },
    ] },
    { title: '收藏', links: [
      { text: '日本系列', href: '/shop' },
      { text: '欧洲系列', href: '/shop' },
      { text: '复古系列', href: '/shop' },
      { text: '限量珍藏', href: '/shop' },
    ] },
    { title: '帮助', links: [
      { text: '店员助手', href: '/shop#assistant' },
      { text: '邮路测试', href: '/shop/sent?service=mail-route-test' },
      { text: '联系我们', href: '/about/contact' },
      { text: '个人中心', href: '/userspace' },
    ] },
  ];

  const socials = ['微博', '微信', '小红书', 'B站'];
</script>

<!-- ===== FOOTER ===== -->
<footer id="about" style="background-color:#0b1829;" class="text-slate-400 relative">

  <!-- Airmail diagonal stripe border top -->
  <div class="airmail-stripe w-full h-3 absolute top-0 left-0"></div>

  <div class="max-w-7xl mx-auto px-6 pt-20 pb-12">
    <div class="grid md:grid-cols-2 lg:grid-cols-5 gap-12 mb-12">

      <!-- Brand col -->
      <div class="lg:col-span-2 space-y-4">
        <div class="flex items-center gap-2">
            <img src="favicon.svg" class="w-10 h-10 text-white" alt="logo" />
          <span class="font-extrabold text-xl text-white">
            <span class="text-rose-400">TKG·</span><span class="text-purple-400">Collection Store</span>
          </span>
        </div>

        <p class="text-sm leading-relaxed max-w-xs" style="color:#8ca0bc;">
          高木桑收藏品店 —— 专注于明信片、信件与邮票的收藏与代寄，让每一份珍藏都承载温度与故事。
        </p>

        <!-- Postmark decoration + socials -->
        <div class="flex items-center gap-3">
          <svg width="48" height="48" viewBox="0 0 48 48" fill="none" class="opacity-40 flex-shrink-0">
            <circle cx="24" cy="24" r="22" stroke="#7c3aed" stroke-width="1.5"/>
            <circle cx="24" cy="24" r="17" stroke="#7c3aed" stroke-width="0.8" fill="none"/>
            <line x1="4" y1="24" x2="44" y2="24" stroke="#7c3aed" stroke-width="1.2"/>
            <text x="24" y="21" text-anchor="middle" font-size="5" fill="#7c3aed" font-family="monospace">TKG STORE</text>
            <text x="24" y="29" text-anchor="middle" font-size="4.5" fill="#7c3aed" font-family="monospace">2026</text>
          </svg>
          <div class="flex gap-2">
            {#each socials as platform}
              <a
                href="#social"
                class="social-btn"
                aria-label={platform}
              >
                {platform[0]}
              </a>
            {/each}
          </div>
        </div>
      </div>

      <!-- Link columns -->
      {#each colLinks as col}
        <div class="space-y-4">
          <h4 class="font-bold text-white text-sm">{col.title}</h4>
          <ul class="space-y-2">
            {#each col.links as lnk}
              <li>
                <a href={lnk.href} class="footer-link">{lnk.text}</a>
              </li>
            {/each}
          </ul>
        </div>
      {/each}
    </div>

    <!-- Bottom bar -->
    <div class="bottom-bar">
      <div class="flex flex-col gap-2 text-sm" style="color:#8ca0bc;">
        <p>© 2026 广州市白云区三元里高木桑收藏品店 保留所有权利。</p>
        <div class="beian-info">
          <a
            href="https://beian.miit.gov.cn/#/Integrated/index"
            target="_blank"
            rel="noopener noreferrer"
            class="footer-link"
          >粤ICP备2026021900号</a>
          <span class="beian-separator"> </span>
          <a href="https://beian.mps.gov.cn/#/query/webSearch?code=44011102484680" class="footer-link beian-link">
            <img src="/cops/isb.png" alt="公安备案" class="cops-icon">
            <span>粤公网安备44011102484680号</span>
          </a>
        </div>
        <p>本店服务通过中国邮政实现，不提供实际寄送服务。</p>
      </div>

      <div class="flex items-center gap-4 flex-shrink-0">
        <!-- Service status badge (slot-able) -->
        <slot name="status">
          <div class="flex items-center gap-2 text-xs">
            <span class="w-2 h-2 rounded-full bg-red-400 animate-pulse"></span>
            <span class="text-red-400">代寄服务测试中</span>
          </div>
        </slot>

        <!-- 网站更新暗号 -->
        <div class="airmail-badge px-3 py-1 rounded text-[10px] font-mono font-bold tracking-widest" style="color:#0b1829;">
          0316 1725
        </div>
      </div>
    </div>
  </div>
</footer>

<style>
  /* ── Airmail diagonal stripe border ──────────────────────── */
  .airmail-stripe {
    background: repeating-linear-gradient(
      -45deg,
      #c8102e 0px, #c8102e 6px,
      #ffffff 6px, #ffffff 12px,
      #002868 12px, #002868 18px,
      #ffffff 18px, #ffffff 24px
    );
    opacity: 0.7;
  }

  /* ── Airmail label badge ──────────────────────────────────── */
  .airmail-badge {
    background: repeating-linear-gradient(
      -45deg,
      #c8102e 0px, #c8102e 4px,
      #ffffff 4px, #ffffff 8px,
      #002868 8px, #002868 12px,
      #ffffff 12px, #ffffff 16px
    );
    letter-spacing: 0.12em;
  }

  /* ── Social buttons ──────────────────────────────────────── */
  .social-btn {
    width: 2.25rem; height: 2.25rem;
    border-radius: 0.75rem;
    background: #142236;
    display: flex; align-items: center; justify-content: center;
    font-size: 0.75rem; font-weight: 700;
    color: #8ca0bc;
    text-decoration: none;
    transition: all 0.2s;
  }
  .social-btn:hover {
    background: linear-gradient(135deg, #f43f5e, #a855f7);
    color: white;
  }

  /* ── Footer link ─────────────────────────────────────────── */
  .footer-link {
    font-size: 0.875rem;
    color: #8ca0bc;
    text-decoration: none;
    transition: color 0.15s;
  }
  .footer-link:hover { color: #f43f5e; }

  /* ── Bottom bar ──────────────────────────────────────────── */
  .bottom-bar {
    border-top: 1px solid #1a2e44;
    padding-top: 2rem;
    display: flex;
    flex-direction: column;
    gap: 1rem;
  }
  @media (min-width: 768px) {
    .bottom-bar {
      flex-direction: row;
      align-items: flex-end;
      justify-content: space-between;
    }
  }

  /* ── 备案信息 ──────────────────────────────────────────── */
  .beian-info {
    display: flex;
    align-items: center;
    gap: 0.5rem;
  }

  .beian-separator {
    color: #5a7a95;
  }

  .beian-link {
    display: flex;
    align-items: center;
    gap: 0.375rem;
  }

  .cops-icon {
    width: 16px;
    height: 16px;
    flex-shrink: 0;
  }
</style>
