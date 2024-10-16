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
    // config the files to be scanned by tailwind
    files: [
      "./docs/.vitepress/**/*.{js,ts,vue}",
      "./docs/**/*.md",
      "./docs/components/**/*.vue",
      "./node_modules/flowbite/**/*.js",
    ],
    safelist: ["html", "body"],
  },
  theme: {
    extend: {
      colors: {
        primary: {
          600: "#1e5ae6",
          700: "#1a46d0",
          800: "#1b39a4",
        },
      },
    },
  },
};
