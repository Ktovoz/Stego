import zhCN from './locales/zh-CN.json';
import enUS from './locales/en-US.json';

export const messages = {
  'zh-CN': zhCN,
  'en-US': enUS,
};

export const defaultLocale = 'zh-CN';

export const languages = [
  { code: 'zh-CN', name: '中文' },
  { code: 'en-US', name: 'English' },
];

export function t(messages, locale, key) {
  const keys = key.split('.');
  let result = messages[locale];

  for (const k of keys) {
    if (result && typeof result === 'object') {
      result = result[k];
    } else {
      result = undefined;
      break;
    }
  }

  if (result === undefined) {
    result = messages[defaultLocale];
    for (const k of keys) {
      if (result && typeof result === 'object') {
        result = result[k];
      } else {
        result = key;
        break;
      }
    }
  }

  return result || key;
}

export function tp(messages, locale, key, params) {
  let result = t(messages, locale, key);

  if (params && typeof result === 'string') {
    Object.keys(params).forEach(param => {
      result = result.replace(`{${param}}`, params[param]);
    });
  }

  return result;
}
