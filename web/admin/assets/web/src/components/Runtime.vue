<template>

  <el-row>
    <el-col :span="12">
      <el-row>
        <el-col :span="12">OS</el-col>
        <el-col :span="12"><el-tag v-text="system.os"></el-tag></el-col>

        <el-col :span="12">Arch</el-col>
        <el-col :span="12"><el-tag v-text="system.arch"></el-tag></el-col>

        <el-col :span="12">Go version</el-col>
        <el-col :span="12"><el-tag v-text="system.version"></el-tag></el-col>

        <el-col :span="12">PID</el-col>
        <el-col :span="12"><el-tag v-text="system.pid"></el-tag></el-col>

        <el-col :span="12">Gzip</el-col>
        <el-col :span="12"><el-switch disabled v-model="system.enableGzip"></el-switch></el-col>

      </el-row>
    </el-col>
    <el-col :span="12">
      <el-row>
        <el-col :span="12">Shard</el-col>
        <el-col :span="12"><el-tag v-text="system.shard"></el-tag></el-col>
        <el-col :span="12">Debug</el-col>
        <el-col :span="12"><el-switch disabled v-model="system.debug"></el-switch></el-col>

        <el-col :span="12">GOMAXPROCS</el-col>
        <el-col :span="12"><el-tag v-text="system.gomaxprocs"></el-tag></el-col>
        <el-col :span="12">Goroutines</el-col>
        <el-col :span="12"><el-tag v-text="system.goroutines"></el-tag></el-col>
      </el-row>
    </el-col>
    <el-divider content-position="left">认证</el-divider>
    <el-row>
      <el-col :span="12">状态</el-col>
      <el-col :span="12"><el-switch disabled v-model="system.enableAuth"></el-switch></el-col>
    </el-row>
    <el-alert type="warning" show-icon v-if="!system.enableAuth">
      未开启认证，将会导致数据不安全，建议开启，启动时增加参数：--auth=admin:123456
    </el-alert>

    <el-divider content-position="left">路径</el-divider>
    <el-row>
      <el-col :span="6">数据目录</el-col>
      <el-col :span="18">
        <el-tag style="margin-right: 10px;" type="success"><el-tag v-text="system.dataSize"></el-tag>MB</el-tag>
        <el-tag v-text="system.dataPath"></el-tag>
      </el-col>

      <el-col :span="6">词典目录</el-col>
      <el-col :span="18"><el-tag v-text="system.dictionaryPath"></el-tag></el-col>

      <el-col :span="6">启动路径</el-col>
      <el-col :span="18"><el-tag v-text="system.executable"></el-tag></el-col>
    </el-row>
    <el-divider content-position="left">数据库</el-divider>
    <el-row>
      <el-col :span="20">数据库数量</el-col>
      <el-col :span="4"><el-tag v-text="system.dbs"></el-tag></el-col>
    </el-row>
  </el-row>


</template>

<script>
export default {
  name: 'Runtime',
  props: {
    system: {
      type: Object,
      required: true,
    },
  },
}
</script>

<style scoped>

</style>
