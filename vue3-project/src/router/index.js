import { createRouter, createWebHistory } from 'vue-router'
import LoginPage from '../views/LoginPage.vue'
import HomePage from '../views/HomePage.vue'
import UserPage from '../views/UserPage.vue'

const routes = [
  {
    path: '/login',
    name: 'LoginPage',
    component: LoginPage
  },
  {
    path: '/',
    name: 'HomePage',
    component: HomePage,
    meta: { requiresAuth: true },
    children: [
      {
        path: 'users',
        name: 'UserPage',
        component: UserPage,
        meta: { requiresAuth: true }
      },
      {
        path: 'organizations',
        name: 'Organizations',
        component: () => import('../views/OrganizationPage.vue'),
        meta: { requiresAuth: true }
      },
      {
        path: 'organizations/:id',
        name: 'OrganizationDetail',
        component: () => import('../views/OrganizationDetailPage.vue'),
        meta: { requiresAuth: true }
      },
      {
        path: 'object-classes',
        name: 'ObjectClasses',
        component: () => import('../views/ObjectClassPage.vue'),
        meta: { requiresAuth: true }
      },
      {
        path: 'object-classes/:id',
        name: 'ObjectClassDetail',
        component: () => import('../views/ObjectClassDetailPage.vue'),
        meta: { requiresAuth: true }
      }
    ]
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

router.beforeEach((to, from, next) => {
  const isAuthenticated = localStorage.getItem('token')
  
  if (to.matched.some(record => record.meta.requiresAuth) && !isAuthenticated) {
    next('/login')
  } else if (to.path === '/login' && isAuthenticated) {
    next('/')
  } else {
    next()
  }
})

export default router