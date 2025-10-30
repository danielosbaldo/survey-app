/** @type {import('tailwindcss').Config} */
module.exports = {
  content: [
    "./assets/web/templates/**/*.{gohtml,html}",
    "./internal/handlers/**/*.go",
  ],
  theme: {
    extend: {
      colors: {
        primary: '#D81B60',
        secondary: '#C2185B',
        background: '#F8F8F8',
        text: '#311B0B',
      },
    },
  },
  plugins: [],
}
