/** @type {import('tailwindcss').Config} */
export default {
  content: ['./index.html', './src/**/*.{vue,js,ts,jsx,tsx}'],
  theme: {
    extend: {
      colors: {
        // Tea 原调色板（保留兼容）
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
        // === 奢侈品牌色板（和 LiveRoom.vue 全局 CSS 一致）===
        // off-black 主舞台、象牙白面板、香槟金点缀
        ink: {
          900: '#0B0A09',   // off-black — Dior / Hermès 主舞台
          800: '#1A1A18',
          700: '#2A2824',
        },
        ivory: {
          50:  '#FAF7F0',
          100: '#F8F5EF',   // 象牙白 —爱马仕丝巾盒底色
          200: '#F0EBE0',
        },
        gold: {
          DEFAULT: '#C5A572',  // 香槟金 — Rolls-Royce coachline / Château Lafite
          soft:    '#D4B785',
          muted:   '#8A7555',
        },
        sand: {
          DEFAULT: '#8A8578',  // 次要文本色 — 像老杂志的铅字灰
          deep:    '#5A5750',
        },
        cream: '#F7F5F0',
        mattegold: '#B8A47C',
        vermillion: '#C0392B',
      },
      fontFamily: {
        // 全站统一为 Cormorant Garamond — 比 Playfair Display 更纤细优雅，
        // 符合 Hermès / Chanel / 老 Savile Row 的传统审美
        display: ['"Cormorant Garamond"', 'Georgia', 'serif'],
        sans: ['Inter', 'system-ui', 'sans-serif'],
        serif: ['"Cormorant Garamond"', 'Georgia', 'serif'],
        song: ['"Noto Serif SC"', '"思源宋体"', 'serif'],
      },
      letterSpacing: {
        'brand': '0.02em',
        'lux':   '0.22em',   // 奢侈品 SERIF CAPS 间距（LiveRoom 用）
      },
      transitionDuration: {
        'lux': '600ms',      // 奢侈品牌 slow transition（比默认 150ms 慢 4 倍）
      },
      borderRadius: {
        'lux': '2px',        // 奢侈品牌微圆角（对比快消品的 16px）
      },
    },
  },
  plugins: [],
}
