<script setup>
import { computed, onMounted, ref, watch } from 'vue'
import { useRoute } from 'vue-router'
import { BadgeCheck, Briefcase, Building2, Home, LogOut, Send, ShieldCheck, Tag, Users } from 'lucide-vue-next'
import AppLogo from '@/components/AppLogo.vue'
import BottomNav from '@/components/BottomNav.vue'
import UserAvatar from '@/components/UserAvatar.vue'
import CompanyLogo from '@/components/CompanyLogo.vue'
import ConfirmDialog from '@/components/ConfirmDialog.vue'
import { useSignOut } from '@/composables/useSignOut'
import api from '@/services/api'
import { useAuthStore } from '@/stores/auth'
import { COMPANY_STATUS } from '@/utils/format'

const auth = useAuthStore()
const route = useRoute()
const { confirming, signOut } = useSignOut()
const pendingCompanies = ref(0)

const menus = {
  company: [
    { name: 'company-dashboard', label: 'Dashboard', short: 'Beranda', icon: Home },
    { name: 'company-jobs', label: 'Lowongan', short: 'Lowongan', icon: Briefcase, prefix: ['company-job', 'company-applicant'] },
    { name: 'company-talents', label: 'Cari Talenta', short: 'Talenta', icon: Users },
    { name: 'company-invitations', label: 'Undangan Terkirim', short: 'Undangan', icon: Send },
    { name: 'company-profile', label: 'Profil Perusahaan', short: 'Profil', icon: Building2 },
  ],
  admin: [
    { name: 'admin-dashboard', label: 'Dashboard', short: 'Beranda', icon: Home },
    { name: 'admin-companies', label: 'Verifikasi Perusahaan', short: 'Verifikasi', icon: ShieldCheck, count: true },
    { name: 'admin-users', label: 'Pengguna', short: 'Pengguna', icon: Users },
    { name: 'admin-jobs', label: 'Lowongan', short: 'Lowongan', icon: Briefcase },
    { name: 'admin-categories', label: 'Kategori', short: 'Kategori', icon: Tag },
  ],
}
const items = computed(() => menus[auth.role] ?? [])
const isActive = (item) => route.name === item.name || (item.prefix ?? []).some((prefix) => String(route.name).startsWith(prefix))
const bottomItems = computed(() => items.value.map((item) => ({ label: item.short, icon: item.icon, to: { name: item.name }, active: isActive(item), badge: item.count ? pendingCompanies.value : 0 })))

async function loadPending() {
  if (auth.role !== 'admin') return
  try {
    pendingCompanies.value = (await api.get('/admin/companies', { params: { status: 'pending' } })).data.total
  } catch {
    pendingCompanies.value = 0
  }
}

watch(() => route.fullPath, loadPending)
onMounted(loadPending)
</script>

<template>
  <div class="flex min-h-dvh">
    <aside class="sticky top-0 hidden h-dvh w-60 shrink-0 flex-col gap-1 border-r border-line bg-surface-2 px-3.5 py-5 lg:flex">
      <RouterLink :to="{ name: 'home' }" class="mb-5 px-3" aria-label="Karsa, ke beranda"><AppLogo /></RouterLink>
      <nav class="flex flex-col gap-1" aria-label="Menu dashboard">
        <RouterLink
          v-for="item in items"
          :key="item.name"
          :to="{ name: item.name }"
          class="flex min-h-[42px] items-center gap-3 rounded-btn px-3 text-[16px] font-medium whitespace-nowrap transition"
          :class="isActive(item) ? 'bg-surface text-text shadow-[0_1px_2px_rgb(20_23_31/0.06)] [&>svg]:text-primary' : 'text-muted hover:bg-surface/60 hover:text-text'"
          :aria-current="isActive(item) ? 'page' : undefined"
        >
          <component :is="item.icon" class="size-[18px]" aria-hidden="true" />{{ item.label }}
          <span v-if="item.count && pendingCompanies" class="ml-auto rounded-full bg-warning px-2 text-xs font-bold text-white">{{ pendingCompanies }}</span>
        </RouterLink>
      </nav>
      <button type="button" class="mt-auto flex min-h-[42px] items-center gap-3 rounded-btn px-3 text-[16px] font-medium text-muted hover:text-text" @click="confirming = true"><LogOut class="size-[18px]" aria-hidden="true" />Keluar</button>
      <div class="mt-2 flex items-center gap-2.5 border-t border-line px-3 pt-3">
        <CompanyLogo v-if="auth.company" :name="auth.company.name" :size="34" />
        <UserAvatar v-else :name="auth.user.name" :url="auth.user.avatar_url" :size="34" />
        <div class="min-w-0">
          <div class="truncate text-sm font-semibold">{{ auth.company?.name ?? auth.user.name }}</div>
          <div class="flex items-center gap-1 truncate text-[14px] text-muted">
            <BadgeCheck v-if="auth.isVerified" class="size-3.5 text-primary" aria-hidden="true" />{{ auth.company ? COMPANY_STATUS[auth.company.status].label : auth.user.email }}
          </div>
        </div>
      </div>
    </aside>

    <div class="min-w-0 flex-1 pb-16 lg:pb-0">
      <div class="sticky top-0 z-30 flex h-14 items-center gap-3 border-b border-line bg-surface px-4 lg:hidden">
        <RouterLink :to="{ name: 'home' }" aria-label="Karsa, ke beranda"><AppLogo :size="28" /></RouterLink>
        <span class="ml-auto truncate text-sm font-medium text-muted">{{ auth.company?.name ?? 'Admin' }}</span>
        <button type="button" class="btn-icon -mr-2" aria-label="Keluar" @click="confirming = true"><LogOut class="size-5" aria-hidden="true" /></button>
      </div>
      <main class="mx-auto max-w-[1280px] px-4 py-6 sm:px-8 sm:py-8"><RouterView /></main>
    </div>

    <BottomNav :items="bottomItems" class="lg:hidden" />
    <ConfirmDialog :open="confirming" title="Keluar dari Karsa?" text="Kamu perlu masuk lagi untuk membuka dashboard." confirm-label="Keluar" @confirm="signOut" @cancel="confirming = false" />
  </div>
</template>
