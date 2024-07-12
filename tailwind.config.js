/** @type {import('tailwindcss').Config} */
module.exports = {
  mode:'jit',
  content: [
      "./components/**/*.{go,templ}",
      "./node_modules/flowbite/**/*.js"

  ],
  theme: {
    extend: {
      colors:{
        shadow:'#f0f9ff',
        temp:'#e11d48',
      },
      fontFamily: {
        roboto: ["Roboto", "sans-serif"],
      },
    },
  },
  plugins: [
    require('@tailwindcss/forms'),
    require('flowbite/plugin')({
      charts: true,
    }),
  ],
}