<script setup>
import { computed, onBeforeUnmount, onMounted, ref, watch } from 'vue'
import { useRoute } from 'vue-router'
import { Bookmark, Building2, ChevronDown, FileText, Home, LayoutDashboard, LogIn, LogOut, Mail, Search, UserRound } from 'lucide-vue-next'
import AppLogo from '@/components/AppLogo.vue'
import BottomNav from '@/components/BottomNav.vue'
import ConfirmDialog from '@/components/ConfirmDialog.vue'
import UserAvatar from '@/components/UserAvatar.vue'
import { useSignOut } from '@/composables/useSignOut'
import { HOME_BY_ROLE, useAuthStore } from '@/stores/auth'
import { useNotificationStore } from '@/stores/notifications'

const auth = useAuthStore()
const route = useRoute()
const { confirming, signOut } = useSignOut()
const menuOpen = ref(false)
const menu = ref()
const notifications = useNotificationStore()

const seekerLinks = computed(() => [
  { name: 'profile', label: 'Profil & CV' },
  { name: 'applications', label: 'Lamaran Saya', badge: notifications.applications },
  { name: 'saved', label: 'Tersimpan' },
  { name: 'invitations', label: 'Undangan', badge: notifications.invitations },
])

const isOn = (...names) => names.some((name) => String(route.name).startsWith(name))
const bottomItems = computed(() => {
  const common = [
    { label: 'Beranda', icon: Home, to: { name: 'home' }, active: route.name === 'home' },
    { label: 'Lowongan', icon: Search, to: { name: 'jobs' }, active: isOn('job') },
  ]
  if (auth.role === 'seeker') {
    return [
      ...common,
      { label: 'Lamaran', icon: FileText, to: { name: 'applications' }, active: isOn('application'), badge: notifications.applications },
      { label: 'Undangan', icon: Mail, to: { name: 'invitations' }, active: isOn('invitations'), badge: notifications.invitations },
      { label: 'Profil', icon: UserRound, to: { name: 'profile' }, active: isOn('profile', 'saved') },
    ]
  }
  if (auth.role) return [...common, { label: 'Dashboard', icon: LayoutDashboard, to: { name: HOME_BY_ROLE[auth.role] } }]
  return [
    ...common,
    { label: 'Perusahaan', icon: Building2, to: { name: 'register', query: { role: 'company' } }, active: route.name === 'register' && route.query.role === 'company' },
    { label: 'Masuk', icon: LogIn, to: { name: 'login' }, active: route.name === 'login' || (route.name === 'register' && route.query.role !== 'company') },
  ]
})

function closeMenu(event) {
  if (menu.value && !menu.value.contains(event.target)) menuOpen.value = false
}

watch(() => route.fullPath, () => {
  menuOpen.value = false
  notifications.load()
})
watch(() => auth.role, notifications.load)
onMounted(() => {
  notifications.load()
  document.addEventListener('click', closeMenu)
})
onBeforeUnmount(() => document.removeEventListener('click', closeMenu))
</script>

<template>
  <div class="flex min-h-dvh flex-col pb-16 md:pb-0">
    <header class="sticky top-0 z-40 border-b border-line bg-surface/95 backdrop-blur">
      <div class="container-app flex h-16 items-center gap-8 sm:h-[68px]">
        <RouterLink :to="{ name: 'home' }" aria-label="Karsa, ke beranda"><AppLogo /></RouterLink>
        <nav class="hidden gap-6 text-[16px] font-medium md:flex" aria-label="Utama">
          <RouterLink :to="{ name: 'jobs' }" class="text-muted hover:text-text" active-class="!text-text">Cari Lowongan</RouterLink>
          <RouterLink v-if="!auth.role" :to="{ name: 'register', query: { role: 'company' } }" class="text-muted hover:text-text">Untuk Perusahaan</RouterLink>
        </nav>

        <div class="ml-auto flex items-center gap-1">
          <template v-if="!auth.role">
            <RouterLink :to="{ name: 'login', query: route.name === 'job' ? { redirect: route.fullPath } : {} }" class="btn-ghost hidden md:inline-flex">Masuk</RouterLink>
            <RouterLink :to="{ name: 'register' }" class="btn-primary btn-sm md:min-h-11 md:px-4 md:text-sm">Daftar</RouterLink>
          </template>
          <template v-else-if="auth.role === 'seeker'">
            <RouterLink :to="{ name: 'saved' }" class="btn-icon" aria-label="Lowongan tersimpan"><Bookmark class="size-5" aria-hidden="true" /></RouterLink>
            <RouterLink :to="{ name: 'invitations' }" class="btn-icon relative hidden md:inline-flex" :aria-label="`Undangan, ${notifications.invitations} menunggu`">
              <Mail class="size-5" aria-hidden="true" />
              <span v-if="notifications.invitations" class="absolute top-1.5 right-1.5 grid h-[18px] min-w-[18px] place-items-center rounded-full border-2 border-surface bg-danger px-1 text-[11.5px] font-semibold text-white">{{ notifications.invitations }}</span>
            </RouterLink>
            <button type="button" class="btn-icon md:hidden" aria-label="Keluar" @click="confirming = true"><LogOut class="size-5" aria-hidden="true" /></button>
            <div ref="menu" class="relative ml-2 hidden md:block">
              <button type="button" class="flex min-h-11 items-center gap-2 rounded-btn px-2 hover:bg-surface-2" :aria-expanded="menuOpen" aria-haspopup="menu" @click="menuOpen = !menuOpen">
                <span class="relative"><UserAvatar :name="auth.user.name" :url="auth.user.avatar_url" /><span v-if="notifications.applications" class="absolute -top-0.5 -right-0.5 size-3 rounded-full border-2 border-surface bg-danger" aria-label="Ada pembaruan lamaran" /></span>
                <span class="max-w-32 truncate text-[16px] font-medium">{{ auth.user.name.split(' ')[0] }}</span>
                <ChevronDown class="size-4 text-muted" aria-hidden="true" />
              </button>
              <div v-if="menuOpen" class="absolute right-0 mt-2 w-52 rounded-card border border-line bg-surface p-1.5 shadow-card" role="menu">
                <RouterLink v-for="link in seekerLinks" :key="link.name" :to="{ name: link.name }" class="flex min-h-10 items-center justify-between rounded-lg px-3 text-[16px] hover:bg-surface-2" role="menuitem">{{ link.label }}<span v-if="link.badge" class="rounded-full bg-danger px-2 text-xs font-semibold text-white">{{ link.badge }}</span></RouterLink>
                <button type="button" class="flex min-h-10 w-full items-center rounded-lg px-3 text-[16px] text-danger hover:bg-surface-2" role="menuitem" @click="confirming = true">Keluar</button>
              </div>
            </div>
          </template>
          <RouterLink v-else :to="{ name: HOME_BY_ROLE[auth.role] }" class="btn-primary hidden md:inline-flex"><LayoutDashboard class="size-4" aria-hidden="true" />Dashboard</RouterLink>
        </div>
      </div>
    </header>

    <main class="flex-1"><RouterView /></main>

    <footer class="border-t border-line bg-surface py-8 text-sm text-muted">
      <div class="container-app flex flex-col gap-4 md:flex-row md:items-center md:justify-between">
        <div class="flex flex-col gap-2 sm:flex-row sm:items-center sm:gap-4"><AppLogo :size="28" /><span>Langkah pertama kariermu di dunia tech.</span></div>
        <span>Proyek demo untuk pembelajaran. Perusahaan dan lowongan bersifat fiktif.</span>
      </div>
    </footer>

    <BottomNav :items="bottomItems" class="md:hidden" />
    <ConfirmDialog :open="confirming" title="Keluar dari Karsa?" text="Kamu perlu masuk lagi untuk melamar dan melihat lamaranmu." confirm-label="Keluar" @confirm="signOut" @cancel="confirming = false" />
  </div>
</template>
