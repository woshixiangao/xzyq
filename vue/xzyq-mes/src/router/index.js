import Vue from 'vue'
import VueRouter from 'vue-router'

Vue.use(VueRouter)

const routes = [
  {
    path: '/',
    redirect: '/login'
  },
  {
    path: '/login',
    name: 'Login',
    component: () => import('../views/Login.vue')
  },
  {
    path: '/register',
    name: 'Register',
    component: () => import('../views/Register.vue')
  },
  {
    path: '/dashboard',
    name: 'Dashboard',
    component: () => import('../views/Dashboard.vue'),
    children: [
      {
        path: 'organizations',
        name: 'Organizations',
        component: () => import('../views/organization/List.vue')
      },
      {
        path: 'projects',
        name: 'Projects',
        component: () => import('../views/project/List.vue')
      },
      {
        path: 'products',
        name: 'Products',
        component: () => import('../views/product/List.vue')
      },
      {
        path: 'users',
        name: 'Users',
        component: () => import('../views/user/List.vue')
      },
      {
        path: 'roles',
        name: 'Roles',
        component: () => import('../views/role/List.vue')
      },
      {
        path: 'logs',
        name: 'Logs',
        component: () => import('../views/log/List.vue')
      }
    ]
  }
]

const router = new VueRouter({
  routes
})

// 路由守卫
router.beforeEach((to, from, next) => {
  const token = localStorage.getItem('token')
  if (to.path === '/login' || to.path === '/register') {
    next()
  } else if (!token) {
    next('/login')
  } else {
    next()
  }
})

export default router