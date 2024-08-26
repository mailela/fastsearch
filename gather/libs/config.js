const fs = require("fs");
const yaml = require("js-yaml");
var doc;

function replaceData(tpl, data) {
  let result = tpl;

  let blocks = result.match(new RegExp(`(\\{\\{\)(.+)(\\}\\})`, "g"));
  if (blocks) {
    for (const block of blocks) {
      let _block=block;
      for (const key in data) {
        _block = _block.replace(
          new RegExp(`#${key}#`, "g"),
          data[key] || ""
        );
      }
      _block = _block.replace(new RegExp(`#(.+)#`, "g"), "''");
      val = eval(_block);
      console.log(block, "->",_block, "->", val)
      result = result.replace(block, val||"");
    }
  }
  console.log(result);
  return result;
}

function stripTags(html) {
  const tmp = document.createElement("DIV");
  tmp.innerHTML = html;
  return tmp.textContent || tmp.innerText || "";
}

// 注意：上面的方法实际上在Node.js中不会工作，因为它使用了document对象，
// 这是浏览器环境特有的。下面是一个仅使用正则表达式的Node.js兼容版本：

function stripTagsWithRegex(html) {
  return html.replace(/<[^>]*>?/gm, "");
}

function load_file(file, argv) {
  try {
    file = file || "./config.yml";
    let cfg = fs.readFileSync(file, "utf8");
    argv = argv || false;
    if (argv) {
      cfg = replaceData(cfg, argv);
    }
    // 读取YAML文件
    doc = load(cfg);
    return doc;
  } catch (e) {
    console.error("Error loading YAML file:", e.message);
  }
  return null;
}
function load(cfg) {
  try {
    // 读取YAML文件
    doc = yaml.load(cfg);
    return doc;
  } catch (e) {
    console.error("Error loading YAML file:", e.message);
  }
  return null;
}
function data() {
  return doc;
}
module.exports = {
  load,
  load_file,
  data,
  stripTags,
  stripTagsWithRegex,
};
