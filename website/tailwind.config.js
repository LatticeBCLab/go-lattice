module.exports = {
  plugins: [
    require("flowbite/plugin"),
    function ({ addUtilities }) {
      const newUtilities = {
        ".no-underline-important": {
          "text-decoration": "none !important",
        },
      };
      addUtilities(newUtilities, ["responsive", "hover"]);
    },
  ],
  content: {
    files: [
      "./docs/.vitepress/**/*.{js,ts,vue}",
      "./docs/**/*.md",
      "./docs/components/**/*.vue",
      "./node_modules/flowbite/**/*.js",
    ],
    safelist: ["html", "body"],
  },
};
