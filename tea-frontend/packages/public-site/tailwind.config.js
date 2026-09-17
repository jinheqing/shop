/** @type {import('tailwindcss').Config} */
export default {
  content: ['./index.html', './src/**/*.{vue,js,ts,jsx,tsx}'],
  theme: {
    extend: {
      colors: {
        cream: '#F7F5F0',
        ink: '#1A1A1A',
        mattegold: '#B8A47C',
        vermillion: '#C0392B',
      },
      fontFamily: {
        display: ['"Playfair Display"', 'Georgia', 'serif'],
        sans: ['Inter', 'system-ui', 'sans-serif'],
        song: ['"思源宋体"', '"Noto Serif SC"', 'serif'],
      },
      letterSpacing: {
        'brand': '0.02em',
      },
    },
  },
  plugins: [],
}
