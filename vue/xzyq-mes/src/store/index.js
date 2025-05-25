import Vue from 'vue'
import Vuex from 'vuex'
import axios from 'axios'

Vue.use(Vuex)

export default new Vuex.Store({
  state: {
    token: localStorage.getItem('token') || '',
    user: JSON.parse(localStorage.getItem('user')) || null,
    organizations: [],
    projects: [],
    products: [],
    users: [],
    roles: [],
    logs: []
  },
  mutations: {
    SET_TOKEN(state, token) {
      state.token = token
      localStorage.setItem('token', token)
    },
    SET_USER(state, user) {
      state.user = user
      localStorage.setItem('user', JSON.stringify(user))
    },
    SET_ORGANIZATIONS(state, organizations) {
      state.organizations = organizations
    },
    SET_PROJECTS(state, projects) {
      state.projects = projects
    },
    SET_PRODUCTS(state, products) {
      state.products = products
    },
    SET_USERS(state, users) {
      state.users = users
    },
    SET_ROLES(state, roles) {
      state.roles = roles
    },
    SET_LOGS(state, logs) {
      state.logs = logs
    },
    CLEAR_AUTH(state) {
      state.token = ''
      state.user = null
      localStorage.removeItem('token')
      localStorage.removeItem('user')
    }
  },
  actions: {
    async login({ commit }, credentials) {
      try {
        const response = await axios.post('/auth/login', credentials)
        const { token, user } = response.data
        commit('SET_TOKEN', token)
        commit('SET_USER', user)
        axios.defaults.headers.common['Authorization'] = `Bearer ${token}`
        return response
      } catch (error) {
        throw error
      }
    },
    async fetchOrganizations({ commit }) {
      try {
        const response = await axios.get('http://localhost:8080/api/organizations')
        commit('SET_ORGANIZATIONS', response.data)
        return response
      } catch (error) {
        throw error
      }
    },
    async fetchProjects({ commit }) {
      try {
        const response = await axios.get('http://localhost:8080/api/projects')
        commit('SET_PROJECTS', response.data)
        return response
      } catch (error) {
        throw error
      }
    },
    async fetchProducts({ commit }) {
      try {
        const response = await axios.get('http://localhost:8080/api/products')
        commit('SET_PRODUCTS', response.data)
        return response
      } catch (error) {
        throw error
      }
    },
    async fetchUsers({ commit }) {
      try {
        const response = await axios.get('http://localhost:8080/api/users')
        commit('SET_USERS', response.data)
        return response
      } catch (error) {
        throw error
      }
    },
    async fetchRoles({ commit }) {
      try {
        const response = await axios.get('http://localhost:8080/api/roles')
        commit('SET_ROLES', response.data)
        return response
      } catch (error) {
        throw error
      }
    },
    async fetchLogs({ commit }) {
      try {
        const response = await axios.get('http://localhost:8080/api/logs')
        commit('SET_LOGS', response.data)
        return response
      } catch (error) {
        throw error
      }
    },
    logout({ commit }) {
      commit('CLEAR_AUTH')
      delete axios.defaults.headers.common['Authorization']
    }
  }
})