import { createApp } from "vue";
import App from "./App.vue";
import "element-plus/dist/index.css";
import ElementPlus from "element-plus";
import zhCn from "element-plus/es/locale/lang/zh-cn";

import * as ElementPlusIconsVue from "@element-plus/icons-vue";
import "element-plus/theme-chalk/src/dark/css-vars.scss";
import router from "./router";

let app = createApp(App);
// 屏蔽错误信息
// app.config.errorHandler = () => null;
// 屏蔽黄色警告信息
app.config.warnHandler = () => null;

app.use(ElementPlus, {
  locale: zhCn,
});
app.use(router);
for (const [key, component] of Object.entries(ElementPlusIconsVue)) {
  app.component(key, component);
  app.use(component);
}


app.config.silent = true;
app.mount("#app");
