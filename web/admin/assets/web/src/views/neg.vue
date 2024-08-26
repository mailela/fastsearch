<template>
  <el-container>
    <el-main style="padding-top:0">
      <div id="dashboard">
        <el-card>
          <template #header>
            <div class="header">
              <div style="display:flex;justify-content:space-between">
                <el-button icon="plus" type="success" style="margin-right:10px" @click="addNeg()">添加</el-button>
                <div style="display:block">
                  <el-input @keyup.enter="search()" v-model="params.query" placeholder="输入关键字,Enter">
                    <template #append>
                      <el-button type="primary" @click="search()">查询</el-button>
                    </template>
                  </el-input>
                  <el-input style="float:left;margin-right:10px;" @keyup.enter="search()" v-model="params.filterExp"
                  placeholder="过滤条件，例如: flag=='1' " />
                  <el-input style="float:left;margin-right:10px;" @keyup.enter="search()" v-model="params.scoreExp"
                  placeholder="排序条件，例如: score+[document.sort]" />
                </div>
              </div>
            </div>
            <div v-if="data && data.pageCount > 1" style="margin-top:10px;">
              <el-pagination @size-change="sizeChange" @current-change="currentChange"
                layout="total, sizes, prev, pager, next, jumper" small="small"
                :page-sizes="[10, 20, 30, 50, 100, 200, 300, 500]" background :page-size="params.limit"
                :current-page="params.page" :total="data.total" />
            </div>
            <div style="display:flex">
                  <el-button  style="margin-right:10px" @click="apply()">应用</el-button>
                </div>
          </template>
          <div v-if="data">
            <el-alert type="success" style="margin-bottom:10px;">
              <div>
                查询耗时：{{ data.time }}ms，找到：{{ data.total }}个结果，每页{{ params.limit }}条，总共{{ data.pageCount }}页
              </div>
              <div>
                <b>分词结果：</b>
              </div>
              <div>
                <el-tag style="margin-right: 10px;" v-for="item in data.words" v-text="item" :key="item"></el-tag>
              </div>
            </el-alert>
          </div>


          <el-table :stripe="true" v-loading="loading" :data="tableData" style="width: 100%" @sort-change="sortChange">

            <el-table-column fixed type="expand" prop="document" label="Document" width="100">
              <template #default="scope">
                <json-viewer :value="scope.row" :expand-depth=5 copyable boxed sort></json-viewer>
              </template>
            </el-table-column>

            <el-table-column prop="id" label="ID" sortable   width="380"/>
            <el-table-column prop="title" label="Title">
              <template #default="scope">
                <span v-html="scope.row.title"></span>
              </template>
            </el-table-column>
            <el-table-column prop="time" label="Time" sortable />
            <el-table-column fixed="right" prop="operation" label="Operation" width="160">
              <template #default="scope">
                <el-link @click="showRow(scope.row)" type="primary" style="margin-right:10px;">显示</el-link>
                <el-link @click="updateRow(scope.row)" type="primary" style="margin-right:10px;">修改</el-link>
                <el-link @click="deleteRow(scope.row)" type="danger">删除</el-link>
              </template>
            </el-table-column>
          </el-table>
          <div v-if="data && data.pageCount > 1" style="margin-top:10px;">
            <el-pagination @size-change="sizeChange" @current-change="currentChange"
              layout="total, sizes, prev, pager, next, jumper" small="small"
              :page-sizes="[10, 20, 30, 50, 100, 200, 300, 500]" background :page-size="params.limit"
              :current-page="params.page" :total="data.total" />
          </div>
        </el-card>
      </div>
    </el-main>
    <el-drawer
    v-model="drawer"
    :title="drawer_data.originalTitle"
  >
    <json-viewer :value="drawer_data" :expand-depth=5 copyable boxed sort></json-viewer>
  </el-drawer>
    <NegDialog :data="dialogData" :db="currentDB" :visible="dialogVisible" @success="indexSuccess()"
      @close="dialogVisible = false"></NegDialog>
  </el-container>
</template>

<script>
import api from '../api'
import jsonViewer from 'vue-json-viewer'
import NegDialog from '../components/NegDialog.vue'
export default {
  name: 'dashboard',
  components: { NegDialog, jsonViewer },
  data() {
    return {
      drawer: false,
      drawer_data:"",
      dbs: {},
      currentDB: '',
      loading: false,
      dialogVisible: false,
      dialogData: null,
      params: {
        query: '',
        filterExp: "",
        scoreExp: "",
        page: 1,
        limit: 10,
        highlight: true,
        order: 'desc',
      },
      data: null,
    }
  },
  watch: {
    currentDB(val) {
      this.search()
    },
    'params.highlight'(val) {
      this.search()
    },
  },
  computed: {
    databases() {
      let rs = []
      for (let db in this.dbs) {
        rs.push({
          label: db,
          value: db,
        })
      }
      return rs
    },
    db() {
      if (!this.currentDB || !this.dbs) return null
      return this.dbs[this.currentDB]
    },
    tableData() {
      if (!this.data) return null
      return this.data.documents
    },
  },
  mounted() {
    this.params.query = this.queryURLParams("query") || '';
    this.params.filterExp = this.queryURLParams("filterExp") || '';
    this.params.filterExp = this.queryURLParams("filterExp") || '';
  },
  created() {
    this.search()
  },
  methods: {
    queryURLParams: function (paramName,url) {
      url=url||window.location;
      // 正则表达式模式，用于匹配URL中的参数部分。正则表达式的含义是匹配不包含 ?、&、= 的字符作为参数名，之后是等号和不包含 & 的字符作为参数值
      let pattern = /([^?&=]+)=([^&]+)/g;
      let params = {};

      // match用于存储正则匹配的结果
      let match;
      // while 循环和正则表达式 exec 方法来迭代匹配URL中的参数
      while ((match = pattern.exec(url)) !== null) {
        // 在字符串url中循环匹配pattern，并对每个匹配项进行解码操作，将解码后的键和值分别存储在key和value变量中
        let key = decodeURIComponent(match[1]);
        let value = decodeURIComponent(match[2]);

        if (params[key]) {
          if (Array.isArray(params[key])) {
            params[key].push(value);
          } else {
            params[key] = [params[key], value];
          }
        } else {
          // 参数名在 params 对象中不存在，将该参数名和对应的值添加到 params 对象中
          params[key] = value;
        }
      }

      if (!paramName) {
        // 没有传入参数名称, 返回含有所有参数的对象params
        return params;
      } else {
        if (params[paramName]) {
          return params[paramName];
        } else {
          return '';
        }
      }
    },
    sortChange({ column, prop, order }) {
      if (order === 'ascending') {
        this.params.order = 'ASC'
      } else {
        this.params.order = 'DESC'
      }
      this.search()
    },
    addNeg() {
      this.dialogData = null
      this.dialogVisible = true
    },
    indexSuccess() {
      this.search()
      this.dialogVisible = false
    },
    selectDB(db) {
      for (let key in this.dbs) {
        this.dbs[key].active = false
      }
      db.active = true
      this.$forceUpdate()
      this.currentDB = db.DatabaseName
    },
    sizeChange(size) {
      this.params.limit = size
      this.$nextTick(() => this.search())
    },
    currentChange(page) {
      this.params.page = page
      this.$nextTick(() => this.search())
    },
    search() {
      this.loading = true
      api.neg_query(this.currentDB, this.params).then(res => {
        if (!res.data.state) {
          this.$message.error(res.message)
        }
        this.data = res.data.data
      }).finally(() => {
        this.loading = false
      })
    },
    apply() {
      this.loading = true
      api.neg_apply(this.currentDB, this.params).then(res => {
        console.log(res)
        this.$message.info(res.data.message)
      }).finally(() => {
        this.loading = false
      })
    },
    updateRow(row) {
      this.dialogData = row
      this.dialogVisible = true
    },
    showRow(row) {
      this.drawer_data = row
      this.drawer = true
    },
    deleteRow(row) {
      let self = this
      //删除
      //弹出确认框
      this.$confirm('确定删除吗？', '提示', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning',
      }).then(() => {
        api.neg_remove(this.currentDB, row.id).then(res => {
          console.log(res.data)
          if (!res.data.state) {
            this.$message.error(res.data.message)
            return
          }
          this.$message.success('删除成功')
          self.search()
        })
      }).catch(() => {
        this.$message({
          type: 'info',
          message: '已取消删除',
        })
      })
    },
  }
  ,
}
</script>

<style scoped>
.header {
  display: flex;
  justify-content: space-between;
}

.database .item {
  line-height: 35px;
  border-bottom: 1px solid var(--el-card-border-color);
  cursor: pointer;
  transition: background-color .3s;
  padding: 0 5px;
}

.database .item:hover {
  background-color: #f3f6f9;
}

.database .item .name {
  margin-left: 10px;
}


.database .active {
  background-color: #fef5ea;
  color: var(--el-color-primary);
}
</style>
