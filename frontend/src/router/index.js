import { createRouter, createWebHistory } from 'vue-router'
import DashboardView from '../views/DashboardView.vue'
import TopArtistsView from '../views/TopArtistsView.vue'
import TopSongsView from '../views/TopSongsView.vue'
import TopAlbumsView from '../views/TopAlbumsView.vue'
import EvolutionSearchView from '../views/EvolutionSearchView.vue'
import WrappedView from '../views/WrappedView.vue'

const router = createRouter({
    history: createWebHistory(import.meta.env.BASE_URL),
    routes: [
        {
            path: '/',
            name: 'home',
            component: DashboardView
        },
        {
            path: '/artists',
            name: 'artists',
            component: TopArtistsView
        },
        {
            path: '/songs',
            name: 'songs',
            component: TopSongsView
        },
        {
            path: '/albums',
            name: 'albums',
            component: TopAlbumsView
        },
        {
            path: '/evolution',
            name: 'evolution',
            component: EvolutionSearchView
        },
        {
            path: '/wrapped',
            name: 'wrapped',
            component: WrappedView
        }
    ]
})

export default router
