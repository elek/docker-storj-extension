import {createApp} from 'vue/dist/vue.esm-bundler.js'
import {createRouter, createWebHashHistory} from 'vue-router'
import Push from "./components/Push.vue";
import Pull from "./components/Pull.vue";
import Main from "./components/Main.vue";
import Configure from "./components/Configure.vue";

const routes = [
    {path: '/', component: Main},
    {path: '/push', component: Push},
    {path: '/pull', component: Pull},
    {path: '/config', component: Configure}
]


const router = createRouter({
    history: createWebHashHistory(),
    routes,
})

// 5. Create and mount the root instance.
const app = createApp({}).use(router).mount('#app')
