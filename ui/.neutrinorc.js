module.exports = {
  options: {
    entry: 'index.jsx',
  },
  use: [
    [
      'neutrino-preset-react',
      {
        html: {
          title: 'QTUM Portal'
        },
      },
    ],
    (neutrino) => neutrino.config
      .entry('vendor')
      .add('react')
      .add('react-dom')
  ]
};