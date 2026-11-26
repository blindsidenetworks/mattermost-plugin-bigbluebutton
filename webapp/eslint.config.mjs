import js from '@eslint/js';
import globals from 'globals';

export default [
	{ignores: ['dist/**', 'node_modules/**']},
	js.configs.recommended,
	{
		files: ['**/*.{js,jsx,mjs}'],
		languageOptions: {
			parserOptions: {ecmaFeatures: {jsx: true}},
			globals: {...globals.browser, ...globals.node},
		},
		rules: {
			'no-unused-vars': ['error', {args: 'none', caughtErrors: 'none', varsIgnorePattern: '^React$'}],
			quotes: ['error', 'single', {avoidEscape: true}],
			semi: ['error', 'always'],
		},
	},
];
