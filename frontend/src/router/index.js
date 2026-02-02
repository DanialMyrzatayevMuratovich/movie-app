import { createRouter, createWebHistory } from 'vue-router'
import Home from '../views/Home.vue'
import MovieDetails from '../views/MovieDetails.vue'
import SeatSelection from '../views/SeatSelection.vue'
import Profile from '../views/Profile.vue'
import Login from '../views/Login.vue'
import TopUp from '../views/TopUp.vue'
import NotFound from '../views/NotFound.vue'
import ServerError from '../views/ServerError.vue'

const routes = [
  {
    path: '/',
    name: 'Home',
    component: Home
  },
  {
    path: '/movie/:id',
    name: 'MovieDetails',
    component: MovieDetails
  },
  {
    path: '/booking/:showtimeId',
    name: 'SeatSelection',
    component: SeatSelection,
    meta: { requiresAuth: true }
  },
  {
    path: '/profile',
    name: 'Profile',
    component: Profile,
    meta: { requiresAuth: true }
  },
  {
    path: '/topup',
    name: 'TopUp',
    component: TopUp,
    meta: { requiresAuth: true }
  },
  {
    path: '/login',
    name: 'Login',
    component: Login
  },
  {
    path: '/error',
    name: 'ServerError',
    component: ServerError
  },
  {
    path: '/:pathMatch(.*)*',
    name: 'NotFound',
    component: NotFound
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

// Navigation guard для защищенных страниц
router.beforeEach((to, from, next) => {
  const isAuthenticated = !!localStorage.getItem('token')

  if (to.meta.requiresAuth && !isAuthenticated) {
    next('/login')
  } else if (to.name === 'Login' && isAuthenticated) {
    next('/')
  } else {
    next()
  }
})

export default router