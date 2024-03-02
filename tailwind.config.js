/** @type {import('tailwindcss').Config} */
module.exports = {
    content: ["./internal/**/*.{html,js,php,ts,go,py}"],
    theme: {
      extend: {
        backgroundImage: {
          'swoop': "url('/public/svg/swoop.svg')",
        },
        colors: {
            "black": "#030303",
            "red": "#E51636",
            "white": "#FFFFFF",
            "blue": "#0000ff",
            "gray": "#999999",
            "lightgray": "#eeeeee",
            "darkgray": "#555555",
            "green": "#2ad457", // Adding dark green
          },
          
        fontFamily: {
          'playful': ['Chelsea Market', 'cursive'],
        },
        screens: {
          'default': '0px',
          'xs': '480px',    // Extra small screens, like smartphones
          'sm': '640px',    // Small screens, like smartphones
          'md': '768px',    // Medium screens, like tablets
          'lg': '1024px',   // Large screens, like laptops
          'xl': '1280px',
          '2xl': '1536px',  // Extra large screens, like desktops
        },
      },
    },
    plugins: [],
  }
  