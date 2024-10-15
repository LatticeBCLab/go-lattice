import "./tailwind.postcss";
import { onMounted } from "vue";
import DefaultTheme from "vitepress/theme";

export default {
  ...DefaultTheme,
  enhanceApp({ app, router, siteData }) {
    if (typeof window !== "undefined") {
      onMounted(() => {
        import("../../../node_modules/flowbite/dist/flowbite.min.js");
      });
    }
  },
};
