<template>
  <div v-loading="loading">
    <div style="text-align: right;">
      <el-label>系统版本:</el-label><el-tag v-text="data.version"></el-tag>
      <el-button @click="getStatus" icon="Refresh" type="primary">刷新</el-button>
    </div>
    <div v-if="data">
      <div class="left">
        <el-card>
            <template #header>内存</template>
            <Memory :memory="data.memory"></Memory>
          </el-card>
          <el-card>
            <template #header>CPU</template>
            <CPU :cpu="data.cpu"></CPU>
          </el-card>
          <el-card>
            <template #header>磁盘</template>
            <Disk :disk="data.disk"></Disk>
          </el-card>
      </div>
      <div class="right">
        <el-card>
          <template #header>配置</template>
          <Runtime :system="data.system"></Runtime>
        </el-card>
      </div>  
    </div>
  </div>
</template>

<script>
import api from '../api'
import ProgressChat from '../components/ProgressChat.vue'
import Memory from '../components/Memory.vue'
import CPU from '../components/CPU.vue'
import Disk from '../components/Disk.vue'
import Runtime from '../components/Runtime.vue'

export default {
  name: 'status',
  components: { Runtime, Disk, CPU, Memory, ProgressChat },
  data() {
    return {
      data: null,
      loading: false,
    }
  },
  created() {
    this.getStatus()
  },
  methods: {
    getStatus() {
      let self = this
      self.loading = true

      api.getStatus().then(({ data }) => {

        if (data.state) {
          this.data = data.data
        } else {
          this.$message.error(data.message)
        }
      }).finally(() => self.loading = false)
    },
  },
}
</script>

<style>
.el-card {
  margin-top: 10px;
}
.left {
  display: flex;
  .el-card {
    margin-left: 10px;
  }
}
.right {
  display: flex;
  margin-left: 10px;
}
</style>
