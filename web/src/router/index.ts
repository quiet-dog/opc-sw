import { createRouter, createWebHashHistory, createWebHistory, RouteRecordRaw } from 'vue-router'


const route: RouteRecordRaw[] = [
    {
        path: "/",
        name: "home",
        component: () => import("../views/home/index.vue"),
    },
    {
        path: "/service",
        name: "service",
        component: () => import("../views/service/index.vue"),
    }
]
const router = createRouter({
    routes: route,
    history: createWebHashHistory()
})
export default router