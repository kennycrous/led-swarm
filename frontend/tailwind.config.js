/** @type {import('tailwindcss').Config} */
export default {
  content: ['./src/**/*.{html,js,svelte,ts}'],
  theme: {
    extend: {
      colors: {
        cyber: {
          bg: '#06090e',
          card: 'rgba(15, 23, 42, 0.65)',
          border: 'rgba(56, 189, 248, 0.2)',
          borderHover: 'rgba(168, 85, 247, 0.4)',
          cyan: '#06b6d4',
          magenta: '#a855f7',
          gold: '#f59e0b',
          emerald: '#10b981',
          slate: '#0c1017'
        }
      },
      backdropBlur: {
        glass: '16px'
      },
      boxShadow: {
        neonCyan: '0 0 15px rgba(6, 182, 212, 0.3)',
        neonMagenta: '0 0 15px rgba(168, 85, 247, 0.3)'
      }
    },
  },
  plugins: [],
};
