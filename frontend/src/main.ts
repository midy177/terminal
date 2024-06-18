import {createApp} from 'vue'
import App from './App.vue'
import DevUI from 'vue-devui';
import STable from '@shene/table';
import '@shene/table/dist/index.css';
import './style.css';
import 'vue-devui/style.css';
import '@devui-design/icons/icomoon/devui-icon.css';
import { ThemeServiceInit, galaxyTheme } from 'devui-theme';
const themeService = ThemeServiceInit({ galaxyTheme }, 'galaxyTheme');
// 可以动态切换成 galaxyTheme 追光主题
themeService?.applyTheme(galaxyTheme);

createApp(App).use(STable).use(DevUI).mount('#app')
