<script lang="ts">
  import { onMount } from 'svelte';
  import { auth } from './stores/auth.ts';
  import { page } from '$app/stores';
  import { goto } from '$app/navigation'; // Corrected import

  let scrolled = false;
  let mobileOpen = false;
  let showUserMenu = false;

  const navLinks = [
    { href: '/', text: '首页' },
    { href: '/about', text: '关于我们' },
    { href: '/shop', text: '精选商店' },
    { href: '/shop#assistant', text: '店员助手' },
    { href: '/shop/sent', text: '代寄/邮路测试' },
    { href: '/about/contact', text: '联系我们' },
  ];

  async function handleLogout() {
    try {
      await fetch('/api/auth/logout', { method: 'POST' });
    } catch {
      /* local logout still proceeds */
    }
    auth.logout();
    showUserMenu = false;
    goto('/');
  }

  onMount(() => {
    auth.initialize();
    const handler = () => { scrolled = window.scrollY > 20; };
    window.addEventListener('scroll', handler);
    return () => window.removeEventListener('scroll', handler);
  });

  $: usernameInitial = ($auth.user?.username ?? 'J').slice(0, 1).toUpperCase();

  let hoverTimeout: ReturnType<typeof setTimeout>;

  function handleAvatarEnter() {
    clearTimeout(hoverTimeout);
    showUserMenu = true;
  }
  function handleAvatarLeave() {
    hoverTimeout = setTimeout(() => {
      showUserMenu = false;
    }, 120);
  }
</script>

<header class="sticky top-0 z-50 bg-white/80 backdrop-blur-lg border-b border-slate-100/50">
  <nav class="container mx-auto flex items-center justify-between px-4 py-3">
    <!-- Logo -->
    <a href="/" class="flex items-center gap-2 group">
      <img src="/favicon.svg" alt="logo" class="h-8 w-8 group-hover:rotate-12 transition-transform duration-300" />
      <span class="text-lg font-bold text-slate-800 tracking-tight">
        高木桑<span class="text-rose-500">.</span>Store
      </span>
    </a>

    <!-- Desktop Nav -->
    <div class="hidden items-center gap-8 text-sm font-medium text-slate-600 md:flex">
      {#each navLinks as link}
        <a href={link.href} data-sveltekit-reload class="relative hover:text-rose-500 transition-colors py-2 group">
          {link.text}
          <span class="absolute bottom-0 left-0 w-full h-0.5 bg-rose-500 scale-x-0 group-hover:scale-x-100 transition-transform origin-left rounded-full"></span>
        </a>
      {/each}

      {#if $auth.isAuthenticated}
        <div
            class="relative py-2"
            role="group"
            aria-label="User info"
            on:mouseenter={handleAvatarEnter}
            on:mouseleave={handleAvatarLeave}
        >
          <button
            class="user-menu-trigger"
            class:active={showUserMenu}
            on:click={() => showUserMenu = !showUserMenu}
          >
            <span class="user-avatar">{usernameInitial}</span>
            <span class="user-name">{$auth.user?.username}</span>
            <svg class="w-3.5 h-3.5 text-slate-400 transition-transform duration-200" class:rotate-180={showUserMenu} viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round">
              <path d="M6 9l6 6 6-6"/>
            </svg>
          </button>

          {#if showUserMenu}
            <div class="user-menu-card">
              <div class="p-4 border-b border-slate-50 bg-slate-50/50">
                <p class="text-[10px] uppercase tracking-wider text-slate-400 font-bold mb-0.5">Signed in as</p>
                <p class="text-sm font-bold text-slate-800 truncate">{$auth.user?.username}</p>
              </div>
              <div class="p-1.5 flex flex-col gap-0.5">
                <a href="/userspace" class="menu-item" on:click={() => showUserMenu = false}>
                  <span class="text-lg">👤</span>
                  个人中心
                </a>
                <a href="/userspace" class="menu-item" on:click={() => showUserMenu = false}>
                  <span class="text-lg">📦</span>
                  我的申请
                </a>
                <div class="h-px bg-slate-100 my-1 mx-2"></div>
                <button class="menu-item text-rose-500 hover:bg-rose-50!" on:click={handleLogout}>
                  <span class="text-lg">🚪</span>
                  退出登录
                </button>
              </div>
            </div>
          {/if}
        </div>
      {:else}
        <a href="/auth?from={$page.url.pathname}" data-sveltekit-reload class="inline-block rounded-full bg-gradient-to-r from-rose-500 to-purple-600 px-5 py-2 font-semibold text-white">
          登录/注册
        </a>
      {/if}
    </div>

    <!-- Mobile hamburger -->
    <button
      class="md:hidden p-2 rounded-lg text-slate-600 hover:text-rose-500"
      on:click={() => mobileOpen = !mobileOpen}
      aria-label="菜单"
    >
      <svg class="w-6 h-6" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24">
        {#if mobileOpen}
          <path stroke-linecap="round" stroke-linejoin="round" d="M6 18L18 6M6 6l12 12"/>
        {:else}
          <path stroke-linecap="round" stroke-linejoin="round" d="M4 6h16M4 12h16M4 18h16"/>
        {/if}
      </svg>
    </button>
  </nav>

  <!-- mobile nav -->
  <div class:hidden={!mobileOpen} class="border-t md:hidden">
    <div class="container mx-auto px-4">
      <ul class="flex flex-col gap-4 py-4 text-sm font-medium text-slate-600">
        {#each navLinks as link}
          <li><a href={link.href} on:click={() => mobileOpen=false} data-sveltekit-reload class="hover:text-rose-500">{link.text}</a></li>
        {/each}
        {#if $auth.isAuthenticated}
          <li>
            <div class="px-2 py-2 border-t border-slate-100 mt-2">
              <span class="block text-xs text-slate-400 mb-1">已登录</span>
              <div class="font-bold text-rose-500 mb-2">{$auth.user?.username}</div>
              <a href="/userspace" class="block py-1 hover:text-rose-600" on:click={() => mobileOpen=false}>个人中心</a>
              <button class="block w-full text-left py-1 text-slate-500 hover:text-rose-600" on:click={() => {handleLogout(); mobileOpen=false;}}>退出登录</button>
            </div>
          </li>
        {:else}
          <li>
            <a href="/auth?from={$page.url.pathname}" data-sveltekit-reload class="inline-block rounded-full bg-gradient-to-r from-rose-500 to-purple-600 px-5 py-2 font-semibold text-white">
              登录/注册
            </a>
          </li>
        {/if}
      </ul>
    </div>
  </div>
</header>

<style>

  .user-menu-trigger {
    display: flex;
    align-items: center;
    gap: 0.6rem;
    padding: 0.25rem 0.75rem 0.25rem 0.25rem;
    border-radius: 9999px;
    background: transparent;
    border: 1px solid transparent;
    color: #334155;
    font-weight: 600;
    font-size: 0.9rem;
    transition: all 0.2s ease;
    cursor: pointer;
  }

  .user-menu-trigger:hover, .user-menu-trigger.active {
    background: #f8fafc;
    border-color: #e2e8f0;
    box-shadow: 0 2px 4px rgba(0,0,0,0.02);
  }

  .user-avatar {
    width: 2rem;
    height: 2rem;
    border-radius: 50%;
    background: linear-gradient(135deg, #f43f5e, #a855f7);
    color: white;
    display: flex;
    align-items: center;
    justify-content: center;
    font-size: 0.85rem;
    font-weight: 700;
    box-shadow: 0 2px 5px rgba(244, 63, 94, 0.2);
  }

  .user-name {
    max-width: 100px;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
  }

  .user-menu-card {
    position: absolute;
    top: 100%;
    right: 0;
    margin-top: 0.25rem;
    width: 220px;
    background: white;
    border-radius: 1rem;
    border: 1px solid #f1f5f9;
    box-shadow: 0 10px 30px -5px rgba(0, 0, 0, 0.08), 0 4px 6px -2px rgba(0, 0, 0, 0.02);
    transform-origin: top right;
    animation: slideDown 0.2s cubic-bezier(0.16, 1, 0.3, 1);
    overflow: hidden;
  }

  .menu-item {
    display: flex;
    align-items: center;
    gap: 0.75rem;
    width: 100%;
    padding: 0.6rem 0.75rem;
    border-radius: 0.5rem;
    font-size: 0.85rem;
    font-weight: 500;
    color: #475569;
    text-align: left;
    transition: all 0.15s;
  }
  .menu-item:hover {
    background: #f8fafc;
    color: #0f172a;
  }

  @keyframes slideDown {
    from { opacity: 0; transform: translateY(-8px) scale(0.96); }
    to { opacity: 1; transform: translateY(0) scale(1); }
  }
</style>
