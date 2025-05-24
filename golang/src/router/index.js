import Vue from 'vue'
import VueRouter from 'vue-router'
import Login from '../views/Login.vue'
import Register from '../views/Register.vue'
import Home from '../views/Home.vue'

Vue.use(VueRouter)

const routes = [
  {
    path: '/',
    redirect: '/login'
  },
  {
    path: '/login',
    name: 'Login',
    component: Login
  },
  {
    path: '/register',
    name: 'Register',
    component: Register
  },
  {
    path: '/home',
    name: 'Home',
    component: Home,
    children: [
      {
        path: '/logs',
        name: 'Logs',
        component: () => import('../views/Logs.vue')
      },
      {
        path: '/users',
        name: 'Users',
        component: () => import('../views/Users.vue')
      },
      {
        path: '/tenants',
        name: 'Tenants',
        component: () => import('../views/Tenants.vue')
      }
    ]
  }
]

const router = new VueRouter({
  mode: 'history',
  base: process.env.BASE_URL,
  routes
})

// 添加路由守卫
router.beforeEach((to, from, next) => {
  const publicPages = ['/login', '/register']
  const authRequired = !publicPages.includes(to.path)
  const loggedIn = sessionStorage.getItem('username')

  if (authRequired && !loggedIn) {
    return next('/login')
  }

  next()
})

export default router