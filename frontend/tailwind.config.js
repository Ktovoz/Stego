/** @type {import('tailwindcss').Config} */
export default {
  darkMode: ['class'],
  content: [
    './index.html',
    './src/**/*.{js,ts,jsx,tsx}',
  ],
  theme: {
    extend: {
      borderRadius: {
        lg: 'var(--radius)',
        md: 'calc(var(--radius) - 2px)',
        sm: 'calc(var(--radius) - 4px)',
        xl: 'calc(var(--radius) + 4px)',
        '2xl': 'calc(var(--radius) + 8px)',
        '3xl': 'calc(var(--radius) + 12px)',
      },
      colors: {
        background: 'hsl(var(--background))',
        foreground: 'hsl(var(--foreground))',
        card: {
          DEFAULT: 'hsl(var(--card))',
          foreground: 'hsl(var(--card-foreground))',
        },
        popover: {
          DEFAULT: 'hsl(var(--popover))',
          foreground: 'hsl(var(--popover-foreground))',
        },
        primary: {
          DEFAULT: 'hsl(var(--primary))',
          foreground: 'hsl(var(--primary-foreground))',
        },
        secondary: {
          DEFAULT: 'hsl(var(--secondary))',
          foreground: 'hsl(var(--secondary-foreground))',
        },
        muted: {
          DEFAULT: 'hsl(var(--muted))',
          foreground: 'hsl(var(--muted-foreground))',
        },
        accent: {
          DEFAULT: 'hsl(var(--accent))',
          foreground: 'hsl(var(--accent-foreground))',
        },
        destructive: {
          DEFAULT: 'hsl(var(--destructive))',
          foreground: 'hsl(var(--destructive-foreground))',
        },
        border: 'hsl(var(--border))',
        input: 'hsl(var(--input))',
        ring: 'hsl(var(--ring))',
        chart: {
          '1': 'hsl(var(--chart-1))',
          '2': 'hsl(var(--chart-2))',
          '3': 'hsl(var(--chart-3))',
          '4': 'hsl(var(--chart-4))',
          '5': 'hsl(var(--chart-5))',
        },
        'glass': {
          DEFAULT: 'hsl(var(--glass-bg) / 0.6)',
          50: 'hsl(var(--glass-bg) / 0.5)',
          60: 'hsl(var(--glass-bg) / 0.6)',
          70: 'hsl(var(--glass-bg) / 0.7)',
          80: 'hsl(var(--glass-bg) / 0.8)',
          90: 'hsl(var(--glass-bg) / 0.9)',
        },
        'glass-border': {
          DEFAULT: 'hsl(var(--glass-border) / 0.5)',
          50: 'hsl(var(--glass-border) / 0.5)',
          60: 'hsl(var(--glass-border) / 0.6)',
          70: 'hsl(var(--glass-border) / 0.7)',
          80: 'hsl(var(--glass-border) / 0.8)',
        },
      },
      boxShadow: {
        'glass-sm': '0 1px 2px 0 hsl(var(--glass-shadow) / 0.12), 0 1px 3px -1px hsl(var(--glass-shadow) / 0.12)',
        'glass': '0 4px 6px -1px hsl(var(--glass-shadow) / 0.1), 0 2px 4px -2px hsl(var(--glass-shadow) / 0.1)',
        'glass-lg': '0 8px 16px -4px hsl(var(--glass-shadow) / 0.15), 0 4px 8px -2px hsl(var(--glass-shadow) / 0.1)',
        'glass-xl': '0 16px 32px -8px hsl(var(--glass-shadow) / 0.2), 0 8px 16px -4px hsl(var(--glass-shadow) / 0.15)',
        'glow': '0 0 20px hsl(var(--primary) / 0.3)',
        'glow-lg': '0 0 40px hsl(var(--primary) / 0.4)',
        'inset-glass': 'inset 0 1px 0 hsl(255 255 255 / 0.1)',
      },
      backdropBlur: {
        xs: '2px',
        '3xl': '32px',
        '4xl': '64px',
      },
      backdropFilter: {
        'blur': 'blur(20px)',
        'blur-xl': 'blur(24px)',
        'blur-2xl': 'blur(32px)',
        'blur-3xl': 'blur(40px)',
      },
      animation: {
        'float': 'float 6s ease-in-out infinite',
        'pulse-slow': 'pulse-slow 8s ease-in-out infinite',
      },
      keyframes: {
        float: {
          '0%, 100%': { transform: 'translateY(0)' },
          '50%': { transform: 'translateY(-10px)' },
        },
        'pulse-slow': {
          '0%, 100%': {
            opacity: '0.6',
            transform: 'scale(1)',
          },
          '50%': {
            opacity: '0.8',
            transform: 'scale(1.05)',
          },
        },
      },
    },
  },
  plugins: [
    require("tailwindcss-animate"),
  ],
}
