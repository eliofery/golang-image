/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ["./internal/resources/views/**/*.{gohtml,html}"],
  theme: {
    extend: {},
  },
  plugins: ["prettier-plugin-tailwindcss", "@awmottaz/prettier-plugin-void-html"],
}
