import {defineConfig} from 'vitest/config';
import {transformAsync} from '@babel/core';

export default defineConfig({
	plugins: [{
		name: 'plugin-jsx',
		enforce: 'pre',
		async transform(code, id) {
			if (id.includes('/node_modules/') || !/\.[jt]sx?$/.test(id)) {
				return null;
			}
			return transformAsync(code, {
				filename: id,
				presets: [['@babel/preset-react', {runtime: 'classic'}]],
				sourceMaps: true,
			});
		},
	}],
	test: {environment: 'jsdom'},
});
