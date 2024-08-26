import axios from 'axios'

const BASE_URL = process.env.NODE_ENV === 'production' ? '/api' : 'http://127.0.0.1:5679/api'

function request(url, method = 'get', data) {
  return axios({
    baseURL: BASE_URL,
    url: url,
    method: method,
    data: data,
  })
}

export default {
  getDatabase() {
    return request('/db/list')
  },
  query(db, params) {

    return request(`/query?database=${db}`, 'post', {
      ...params,
      highlight: params.highlight ? {
        preTag: '<em style=\'color:red\'>',
        postTag: '</em>',
      } : null,
    })
  },

  remove(db, id) {
    return request(`/index/remove?database=${db}`, 'post', { id })
  },
  gc() {
    return request('/gc')
  },
  getStatus() {
    return request('/status')
  },
  addIndex(db, index) {
    return request(`/index?database=${db}`, 'post', index )
  },
  drop(db){
    return request(`/db/drop?database=${db}`)
  },
  create(db){
    return request(`/db/create?database=${db}`)
  }

  ,
  neg_query(db, params) {
    return request(`/negative/query?database=${db}`, 'post', {
      ...params,
      highlight: params.highlight ? {
        preTag: '<em style=\'color:red\'>',
        postTag: '</em>',
      } : null,
    })
  },

  neg_remove(db, id) {
    return request(`/negative/remove?database=${db}`, 'post', { id })
  },
  neg_add(db, data) {
    return request(`/negative/add?database=${db}`, 'post', data )
  },
  neg_batch_add(db, data) {
    return request(`/negative/import?database=${db}`, 'post', data )
  },
  neg_apply(db, data) {
    return request(`/negative/apply?database=${db}`, 'post', data )
  },
}
