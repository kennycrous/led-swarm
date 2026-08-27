import eslintPluginSvelte from 'eslint-plugin-svelte';

export default [
  ...eslintPluginSvelte.configs['flat/recommended'],
  {
    rules: {
      'svelte/no-at-html-tags': 'warn',
      'svelte/valid-compile': 'error'
    }
  },
  {
    ignores: ['build/', '.svelte-kit/', 'dist/', 'node_modules/']
  }
];
