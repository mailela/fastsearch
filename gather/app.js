const cfg = require("./libs/config");
const Client = require("./libs/fastsearch");
const mysql = require("mysql2/promise");
const CryptoJS = require("crypto-js");

// 使用yargs获取命令行参数
const yargs = require("yargs/yargs");
const { hideBin } = require("yargs/helpers");

// 定义命令行接口
const argv = yargs(hideBin(process.argv)).argv;
const styles=require("./libs/colors")

// 使用示例：node yourScript.js --name="John Doe" --age=30
function echo(msg, flag, style) {
  c = style || styles["green"];
  c = c[0];
  flag = flag || "";
  console.log(c, `${flag}${msg}`);
}
function error(...msg) {
  echo(msg, "[ERROR]", styles.red);
}
function info(...msg) {
  echo(msg, "[INFO]", styles.green);
}
function sleep(ms) {
  return new Promise((resolve) => setTimeout(resolve, ms));
}
function replaceData(tpl, data) {
  let result = tpl;
  for (const key in data) {
    result = result.replace(new RegExp(`\\{${key}\\}`, "g"), data[key] || "");
  }
  return result;
}
async function fetchData(file) {
  let config = cfg.load_file(file, argv);
  openClient(config);
  let connection;
  try {
    // 创建连接
    console.log(config);
    let dbc = config["database"];
    connection = await mysql.createConnection({
      host: dbc["host"],
      port: dbc["port"],
      user: dbc["user"],
      password: dbc["password"],
      database: dbc["database"],
    });
    let i = 0;
    let page = config.source["page"] || 10;
    let table = config.source["table"];
    let where = config.source["where"] || " 1 =1";
    console.log("表名:", table);
    let map = config.source["map"];
    let document_map = map["document"];
    while (true) {
      // 执行SQL查询
      let sql = `select * from ${table} where ${where} limit ${
        page * i
      },${page}`;
      const [rows, fields] = await connection.execute(sql);
      echo(sql);
      if (rows.length == 0) break;
      for (var k = 0; k < rows.length; k++) {
        let item = rows[k];
        // console.log(item);
        let document = {};
        for (let key in document_map) {
          document[key] = item[document_map[key]];
        }
        let key = replaceData(map.key, item);
        let title = replaceData(map.title, item);
        let text = cfg.stripTagsWithRegex(replaceData(map.text, item));
        let flag = replaceData(map.flag, item);
        let tags = replaceData(map.tags, item);
        let state = replaceData(map.state, item);
        let state_flag = map.state_flag || 30;
        info(key, title, flag, tags, state);
        result = false;
        if (state === state_flag) {
           await addDocument(client, key, title, text, flag, tags, document)
            .then(function (value) {
              result = value;
            })
        } else {
          remoteDocument(client, key)
            .then(function (value) {
              // result = value;
            })
            result = true;
        }
        if (result == false) {
          k--;
          error("中断SLEEP",result);
          await sleep(1000);
        }
      }
      if (config["debug"]) {
        error("中断");
        break;
      }
      i++;
      
      echo("第" + i + "页");
    }
    await connection.end();
  } catch (err) {
    console.error("数据库连接或查询出错:", err);
  }
  // 关闭连接
  // await connection.end();
}

function openClient(config) {
  //指定数据库和开启认证
  let db = config["target"]["db"];
  let api_url = config["target"]["api"];
  let username = config["target"]["username"];
  let password = config["target"]["password"];
  info(db, api_url, username, password);
  client = new Client(api_url, db, {
    username: username,
    password: password,
  });
}
// console.log(api_url, db, username, password)
//添加索引
function addDocument(client, id, title, text, flag, tags, document) {
  return new Promise((resolve, reject) => {
    client
      .addDocument(id, title, text, flag, tags, document || {})
      .then((r) => {
        console.log("索引返回：", r.data);
        resolve(r.data.state||false);
      })
      .catch((e) => {
        // console.log(document)
        error(e.message);
        resolve(false);
      });
  });
}
function remoteDocument(client, id) {
  return new Promise((resolve, reject) => {
    client
      .removeDocument(id)
      .then((r) => {
        console.log("删除索引返回：", r.data);
        resolve(true);
      })
      .catch((e) => {
        error(e.message);
        resolve(false);
      });
  });
}

fetchData(argv.file);
