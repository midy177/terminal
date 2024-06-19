import {createApp} from 'vue'
import App from './App.vue'
import DevUI from 'vue-devui';
import './style.css';
import 'vue-devui/style.css';
import '@devui-design/icons/icomoon/devui-icon.css';
import 'ant-design-vue/dist/reset.css';
import { ThemeServiceInit, galaxyTheme } from 'devui-theme';
const themeService = ThemeServiceInit({ galaxyTheme }, 'galaxyTheme');
// 可以动态切换成 galaxyTheme 追光主题
themeService?.applyTheme(galaxyTheme);
createApp(App).use(DevUI).mount('#app')
