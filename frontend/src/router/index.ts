import { createRouter, createWebHistory } from 'vue-router'

// lazy imports
const Landing = () => import('../pages/landing.vue')
const Register = () => import('../pages/register_parent.vue')
const LoginParent = () => import('../pages/login_parent.vue')
const LoginChild = () => import('../pages/login_child.vue')

const LkParent = () => import('../pages/lk_parent.vue')

const LkSmallChild = () => import('../pages/lk_small_child.vue')
const LkBigChild = () => import('../pages/lk_big_child.vue')

const routes = [
  { path: '/', component: Landing },
  // children
  {
    path: '/child',
      children: [
        { path: 'login', component: LoginChild },
        {
          path: 'small',
          children: [
            { path: '', redirect: '/child/small/home' },  
            { path: 'home', component: LkSmallChild },
            { path: 'tasks', component: () => import('../components/small_child/tasklist.vue') },
            { path: 'wishlist', component: () => import('../components/small_child/wishlist.vue') },
            { path: 'new-wish', component: () => import('../components/small_child/new-wish.vue') }
          ]
        },
        {
          path: 'big',
          children: [
            { path: '', redirect: '/child/big/home' },
            { path: 'home', component: LkBigChild },
            { path: 'tasks', component: () => import('../components/big_child/tasklist.vue') },
            { path: 'wishlist', component: () => import('../components/big_child/wishlist.vue') },
            { path: 'new-wish', component: () => import('../components/big_child/new-wish.vue') }
          ]
        }
      ],
  },

  // parent
  {
    path: '/parent',
    children: [
      { path: '', redirect: '/parent/home' },
      { path: 'register', component: Register },
      { path: 'login', component: LoginParent },
      { path: 'home', component: LkParent },
      { path: 'tasks', component: () => import('../components/parent/tasklist.vue') },
      { path: 'new-task', component: () => import('../components/parent/new-task.vue') },
      { path: 'children', component: () => import('../components/parent/childrenlist.vue') },
      { path: 'wishlist', component: () => import('../components/parent/wishlist.vue') }
    ]
  },

  // 404
  { path: '/:pathMatch(.*)*', redirect: '/' }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

// v1 -- {{}} doesn't work in style

//const routes = [
//  { path: '/', component: Landing },

//  // child
//  {
//    path: '/child',
//    children: [
//      { path: '', redirect: '/child/home' },
//      { path: 'login', component: LoginChild },
//      { path: 'home', component: LkChild },
//      { path: 'tasks', component: () => import('../components/child/tasklist.vue') },
//      { path: 'wishlist', component: () => import('../components/child/wishlist.vue') },
//      { path: 'new-wish', component: () => import('../components/child/new-wish.vue') }
//    ]
//  },

//  // parent
//  {
//    path: '/parent',
//    children: [
//      { path: '', redirect: '/parent/home' },
//      { path: 'register', component: Register },
//      { path: 'login', component: LoginParent },
//      { path: 'home', component: LkParent },
//      { path: 'tasks', component: () => import('../components/parent/tasklist.vue') },
//      { path: 'new-task', component: () => import('../components/parent/new-task.vue') },
//      { path: 'children', component: () => import('../components/parent/childrenlist.vue') },
//      { path: 'wishlist', component: () => import('../components/parent/wishlist.vue') }
//    ]
//  },

//  // 404
//  { path: '/:pathMatch(.*)*', redirect: '/' }
//]

//const router = createRouter({
//  history: createWebHistory(),
//  routes
//})

export default router
