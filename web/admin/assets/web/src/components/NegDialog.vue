<template>
  <el-dialog draggable v-model="visible" @close="$emit('close')" title="添加负面词">
    <el-alert type="success" style="margin-bottom: 10px;">请保持ID的唯一，如果存在将会更新数据。</el-alert>
    <el-form ref="form" :model="form" :rules="rules" label-width="120px">
      <el-form-item label="ID" prop="id" v-if="form.id">
        <el-input v-model="form.id" :disabled="data!=null" placeholder="索引id"></el-input>
      </el-form-item>
      <el-form-item label="标题" prop="title">
        <el-input type="text" v-model="form.title" placeholder="标题"></el-input>
      </el-form-item>
      <el-form-item label="负面词" prop="text">
        <el-input type="textarea" v-model="form.text" placeholder="请输入负面词" :rows="5"></el-input>
      </el-form-item>
    </el-form>
    <template #footer>
      <el-button type="primary" @click="save()">确定</el-button>
    </template>
  </el-dialog>
</template>

<script>
import api from '../api'
export default {
  name: 'IndexDialog',
  props: {
    visible: {
      type: Boolean,
      default: false,
    },
    data: {
      type: Object,
      default: () => null,
    },
  },
  data() {
    return {
      rules: {
        text: [
          { required: true, message: '请输入负面词', trigger: 'blur' },
        ],
      },
      form: {
        title: '',
        text: '',
      },
    }
  },
  watch: {
    data(val) {
      if (val) {
        this.form = val
        if(val.originalText) {
          this.form.text = val.originalText
        }
        if(val.originalTitle) {
          this.form.title = val.originalTitle
        }
      } else {
        this.form = {
          text: '',
          title: '',
        }
      }
    },
  },
  methods: {
  
    save() {
      this.$refs.form.validate(valid => {
        if (valid) {
          //校验json文档
          let data = {
            id: this.form.id,
            title: this.form.title,
            text: this.form.text,
            flag: this.form.flag,
          }
          api.neg_add(this.db, data).then(({ data }) => {
            console.log(data)
            if (data.state) {
              this.$message.success('添加成功!')
              this.$emit('success')
            } else {
              this.$message.error(data.message)
            }
          })
        }
      })
    },
  },
}
</script>

<style scoped>

</style>
