/** @type {import('tailwindcss').Config} */
module.exports = {
	content: ["./internal/**/*.{html,js,php,ts,go,py}"],
	theme: {
	  extend: {
		backgroundImage: {
		  'swoop': "url('/public/svg/swoop.svg')",
		},
		colors: {
		  "red": "#E51636",
		  "white": "#FFFFFF",
		  "blue": "#0000ff",
		  "gray": "#999999",
		  "lightgray": "#eeeeee",
		},
		fontFamily: {
		  'playful': ['Chelsea Market', 'cursive'],
		},
    screens: {
      'sm': '640px',    // Small screens, like smartphones
      'md': '768px',    // Medium screens, like tablets
      'lg': '1024px',   // Large screens, like laptops
      'xl': '1280px',
      '2xl': '1536px',    // Extra large screens, like desktops
    },
	  },
	},
	plugins: [],
  }