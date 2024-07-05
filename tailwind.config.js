/** @type {import('tailwindcss').Config} */
module.exports = {
  mode:'jit',
  content: ["./components/**/*.{go,templ}"],
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
  ],
}