import {createApp} from 'vue'
import App from './App.vue'
import DevUI from 'vue-devui';
import '@shene/table/dist/index.css';
import './style.css';
import 'vue-devui/style.css';
import '@devui-design/icons/icomoon/devui-icon.css';
import { ThemeServiceInit, galaxyTheme } from 'devui-theme';
const themeService = ThemeServiceInit({ galaxyTheme }, 'galaxyTheme');
// 可以动态切换成 galaxyTheme 追光主题
themeService?.applyTheme(galaxyTheme);
// import STable from '@shene/table';

createApp(App).use(DevUI).mount('#app')
