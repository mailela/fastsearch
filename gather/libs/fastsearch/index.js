const axios = require("axios");
const CryptoJS = require("crypto-js");

function Client(addr, database, auth) {
  this.addr = addr;

  this.database = database;
  this.auth = auth;

  this.request = function (url, data) {
    // basic auth
    let self = this;

    //数据库
    if (self.database) {
      if (url.indexOf("?") === -1) {
        url = url + "?database=" + self.database;
      } else {
        url = url + "&database=" + self.database;
      }
    }

    let apiUrl = self.addr + url;
    // console.log(apiUrl,data,auth)
    return axios({
      method: "post",
      url: apiUrl,
      data: data,
      auth: self.auth ? self.auth : null,
    });
  };

  /**
   *  添加文档
   * @param id 文档id
   * @param text 文本内容
   * @param document 文档内容
   * @returns {Promise<AxiosResponse<any>>}
   */
  this.addDocument = function (id, title, text, flag, tags, document) {
    return this.request("/index", {
      id: id,
      title: title,
      text: `${text}`,
      flag: `${flag}`,
      tags: `${tags}`,
      document: document,
    });
  };

  /**
   * 查询索引
   * @param query 关键词
   * @param page 页码
   * @param limit 每页数量
   * @param order 排序，ASC、DESC
   * @param highlight 高亮，{"preTag":"<span style='color:red'>","postTag":"</span>"}
   * @returns {Promise<AxiosResponse<any>>}
   */
  this.query = function (
    query,
    page = 1,
    limit = 10,
    order = "desc",
    highlight = null
  ) {
    return this.request("/query", {
      query: query,
      page: page,
      limit: limit,
      order: order,
      highlight: highlight,
    });
  };

  /**
   * 删除索引
   * @param id 文档id
   * @returns {Promise<AxiosResponse<any>>}
   */
  this.removeDocument = function (id) {
    return this.request("/index/remove", {
      id: id,
    });
  };
}

module.exports = Client;
