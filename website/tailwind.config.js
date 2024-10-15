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
  purge: {
    enabled: process.env.NODE_ENV === "production",
    content: ["./docs/.vitepress/**/*.{js,ts,vue}", "./docs/**/*.md"],
    options: {
      safelist: ["html", "body"],
    },
  },
  content: ["./node_modules/flowbite/**/*.js"],
};
