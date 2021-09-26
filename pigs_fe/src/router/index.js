import {createRouter, createWebHistory} from 'vue-router'
import env from "@/store/env";
// import Watermark from '@/assets/watermark';
import NProgress from 'nprogress'

const routes = [
    {
        path: '/',
        name: 'Home',
        component: () => import('../Layout.vue'),
        children: [
            {
                path: '',
                name: 'Index',
                component: () => import('../views/Index.vue'),
                meta: {
                    title: '仪表盘',
                }
            },
            {
                path: 'k8s/cluster',
                name: 'K8S',
                component: () => import('../views/container/K8SClusterManage.vue'),
                meta: {
                    title: '集群管理',
                    module: "容器管理"
                },
                children: []
            },
        ]
    },
    {
        path: '/user/login',
        name: 'Login',
        component: () => import('../views/user/Login.vue'),
        meta: {
            title: '用户登录',
            module: "用户登录"
        }
    },
    {
        path: '/403',
        name: '403',
        component: () => import('../views/error/403.vue'),
        meta: {
            title: '当前无权限操作, 请联系管理员!'
        }
    },
    {
        path: '/refresh',
        name: 'refresh',
        component: () => import('../views/error/refresh.vue'),
        meta: {
            title: '网络已断开'
        }
    },
    {
        path: '/:catchAll(.*)',
        component: () => import('../views/error/404.vue'),
        meta: {
            title: 'Not Found 当前页面未找到!'
        }
    },
]

const router = createRouter({
    history: createWebHistory(),
    routes
})

// 导航守卫
router.beforeEach((to, from, next) => {
    if (to.path === '/user/forgotPwd'){
        return next()
    }

    if (!localStorage.getItem("onLine")) {
        if (to.path !== '/user/login') {
            return next('/user/login')
        }
    }

    // 进度条start
    NProgress.start();
    document.title = to.meta.title + " - " + env.Title;
    next()
});

router.afterEach(() => {
    NProgress.done();
    // 注释水印
    // const water = localStorage.getItem("email")
    // if (typeof (water) === "undefined" || !water){
    //     // Watermark.set("小飞猪 - 运维平台")
    // }else{
    //     Watermark.set(localStorage.getItem("email") + "(" + localStorage.getItem("name") + ")");
    // }
});
export default router
