<template>
  <div style="text-align:center;width: 200px;color:white" v-if="!isCollapsed">
    <img :src="logoSrc" alt="logo" style="width:50px;">
    <h1 style="margin-top:0px;">FastSearch</h1>
  </div>
  <el-menu
      :default-active="active+''"
      :collapse="isCollapsed"
      class="el-menu-vertical-demo"
      background-color="#191a22"
      text-color="#fff"
      active-text-color="#ffd04b"
  >
    <router-link v-for="(item,index) in menus" :to="{name:item.name}" :key="item.name">
      <el-menu-item :index="index">
        <Icon :name="item.icon" :color="item.color"/>
        <span v-text="item.label"></span>
      </el-menu-item>
    </router-link>

  </el-menu>
</template>

<script>
import { Document, Service } from '@element-plus/icons-vue'
import menus from '../menus'
import Icon from './Icon.vue'
import logoImg from '../assets/logo.png';
export default {
  name: 'Menu',
  components: { Icon, Service, Document },
  props: {
    isCollapsed: {
      type: Boolean,
      default: false,
    },
  },
  data() {
    return {
      logoSrc: logoImg,
      menus: menus,
    }
  },
  computed: {
    active() {
      let active = 0
      this.menus.forEach((item, index) => {
        if (item.name === this.$route.name) {
          active = index
        }
      })
      return active
    },
  },
  methods: {
    openDocument() {
      window.open('https://gitee.com/rachel_os/fastsearch')
    },
  },
}
</script>

<style scoped>
.el-menu a{
  text-decoration: none;
}
</style>
