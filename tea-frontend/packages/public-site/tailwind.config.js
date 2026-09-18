/** @type {import('tailwindcss').Config} */
export default {
  content: ['./index.html', './src/**/*.{vue,js,ts,jsx,tsx}'],
  theme: {
    extend: {
      colors: {
        // Tea 调色板 — 和 style.css 里的 CSS 变量保持一致
        tea: {
          50:  '#fdf8f0',
          100: '#f7ecd8',
          200: '#efd8b0',
          300: '#e3bd7d',
          400: '#d4a04f',
          500: '#c48a2f',
          600: '#a66e23',
          700: '#855620',
          800: '#6d451f',
          900: '#5a3a1c',
          950: '#311d0c',
        },
        cream: '#F7F5F0',
        ink: '#1A1A1A',
        mattegold: '#B8A47C',
        vermillion: '#C0392B',
      },
      fontFamily: {
        display: ['"Playfair Display"', 'Georgia', 'serif'],
        sans: ['Inter', 'system-ui', 'sans-serif'],
        serif: ['"Playfair Display"', 'Georgia', 'serif'],
        song: ['"Noto Serif SC"', '"思源宋体"', 'serif'],
      },
      letterSpacing: {
        'brand': '0.02em',
      },
    },
  },
  plugins: [],
}
