module.exports = {
  options: {
    entry: 'index.jsx',
  },
  use: [
    ['neutrino-preset-react', {
      devServer: { port: process.env.PORT || 3000 }
    }],
    (neutrino) => neutrino.config
      .entry('vendor')
        .add('react')
        .add('react-dom')
  ]
};