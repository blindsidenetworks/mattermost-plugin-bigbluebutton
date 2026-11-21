const path = require('node:path');

module.exports = {
	entry: './index.js',
	output: {
		path: path.resolve(__dirname, 'dist'),
		filename: 'bigbluebutton_bundle.js',
		clean: true,
	},
	externals: {
		react: 'React',
		'react-dom': 'ReactDOM',
		redux: 'Redux',
		'react-redux': 'ReactRedux',
	},
	resolve: {extensions: ['.js', '.jsx']},
	module: {
		rules: [{
			test: /\.jsx?$/,
			exclude: /node_modules/,
			use: {
				loader: 'babel-loader',
				options: {
					presets: [
						['@babel/preset-env', {targets: 'defaults', modules: false}],
						['@babel/preset-react', {runtime: 'classic'}],
					],
				},
			},
		}],
	},
};
