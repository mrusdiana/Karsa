import { createRouter, createWebHistory } from 'vue-router'
import PublicLayout from '@/layouts/PublicLayout.vue'
import { HOME_BY_ROLE, useAuthStore } from '@/stores/auth'
import { useToastStore } from '@/stores/toast'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  scrollBehavior: (to, from, saved) => saved ?? (to.path === from.path ? false : { top: 0 }),
  routes: [
    {
      path: '/',
      component: PublicLayout,
      children: [
        { path: '', name: 'home', component: () => import('@/views/HomeView.vue') },
        { path: 'lowongan', name: 'jobs', component: () => import('@/views/JobsView.vue') },
        { path: 'lowongan/:slug', name: 'job', component: () => import('@/views/JobDetailView.vue') },
        { path: 'perusahaan/:slug', name: 'company', component: () => import('@/views/CompanyView.vue') },
        { path: 'masuk', name: 'login', component: () => import('@/views/AuthView.vue'), props: { mode: 'login' }, meta: { guest: true } },
        { path: 'daftar', name: 'register', component: () => import('@/views/AuthView.vue'), props: { mode: 'register' }, meta: { guest: true } },
        { path: 'saya/profil', name: 'profile', component: () => import('@/views/seeker/ProfileView.vue'), meta: { role: 'seeker' } },
        { path: 'saya/lamaran', name: 'applications', component: () => import('@/views/seeker/ApplicationsView.vue'), meta: { role: 'seeker' } },
        { path: 'saya/lamaran/:id', name: 'application', component: () => import('@/views/seeker/ApplicationDetailView.vue'), meta: { role: 'seeker' } },
        { path: 'saya/tersimpan', name: 'saved', component: () => import('@/views/seeker/SavedJobsView.vue'), meta: { role: 'seeker' } },
        { path: 'saya/undangan', name: 'invitations', component: () => import('@/views/seeker/InvitationsView.vue'), meta: { role: 'seeker' } },
      ],
    },
    {
      path: '/perusahaan-saya',
      component: () => import('@/layouts/DashboardLayout.vue'),
      meta: { role: 'company' },
      children: [
        { path: '', name: 'company-dashboard', component: () => import('@/views/company/DashboardView.vue') },
        { path: 'profil', name: 'company-profile', component: () => import('@/views/company/ProfileView.vue') },
        { path: 'lowongan', name: 'company-jobs', component: () => import('@/views/company/JobsView.vue') },
        { path: 'lowongan/baru', name: 'company-job-new', component: () => import('@/views/company/JobFormView.vue') },
        { path: 'lowongan/:id/edit', name: 'company-job-edit', component: () => import('@/views/company/JobFormView.vue') },
        { path: 'lowongan/:id/pelamar', name: 'company-applicants', component: () => import('@/views/company/ApplicantsView.vue') },
        { path: 'pelamar/:id', name: 'company-applicant', component: () => import('@/views/company/ApplicantDetailView.vue') },
        { path: 'talenta', name: 'company-talents', component: () => import('@/views/company/TalentsView.vue') },
        { path: 'undangan', name: 'company-invitations', component: () => import('@/views/company/InvitationsView.vue') },
      ],
    },
    {
      path: '/admin',
      component: () => import('@/layouts/DashboardLayout.vue'),
      meta: { role: 'admin' },
      children: [
        { path: '', name: 'admin-dashboard', component: () => import('@/views/admin/DashboardView.vue') },
        { path: 'perusahaan', name: 'admin-companies', component: () => import('@/views/admin/CompaniesView.vue') },
        { path: 'pengguna', name: 'admin-users', component: () => import('@/views/admin/UsersView.vue') },
        { path: 'lowongan', name: 'admin-jobs', component: () => import('@/views/admin/JobsView.vue') },
        { path: 'kategori', name: 'admin-categories', component: () => import('@/views/admin/CategoriesView.vue') },
      ],
    },
    { path: '/:pathMatch(.*)*', component: PublicLayout, children: [{ path: '', name: 'not-found', component: () => import('@/views/NotFoundView.vue') }] },
  ],
})

router.beforeEach(async (to) => {
  const auth = useAuthStore()
  await auth.init()
  if (to.meta.guest && auth.role) return { name: HOME_BY_ROLE[auth.role] }
  if (!to.meta.role) return
  if (!auth.role) return { name: 'login', query: { redirect: to.fullPath } }
  if (auth.role !== to.meta.role) {
    useToastStore().show('Kamu tidak punya akses ke halaman itu.', 'danger')
    return { name: HOME_BY_ROLE[auth.role] }
  }
})

export default router
