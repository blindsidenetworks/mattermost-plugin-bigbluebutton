module.exports = {
    entry: './index.js',
    output: {
        filename: 'dist/bigbluebutton_bundle.js'
    },
    externals: {
      react: 'React',
      redux: 'Redux',
      'react-redux': 'ReactRedux',
      'react-bootstrap': 'ReactBootstrap',

    },
    module: {
        loaders: [
            {
                test: /\.(js|jsx)?$/,
                loader: 'babel-loader',
                exclude: /(node_modules|non_npm_dependencies)/,
                query: {
                    presets: [
                        'react',
                        ['es2015', {modules: false}],
                        'stage-0'
                    ],
                    plugins: ['transform-runtime']
                }
            }
        ]
    }
};
