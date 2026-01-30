import * as React from 'react';
import { messages, t as translate, tp as translateParams } from './i18n.js';

const I18nContext = React.createContext(undefined);

const LOCALE_STORAGE_KEY = 'stego-locale';

function getInitialLocale() {
  const stored = localStorage.getItem(LOCALE_STORAGE_KEY);
  if (stored && messages[stored]) {
    return stored;
  }
  const systemLang = navigator.language || navigator.userLanguage;
  if (systemLang.startsWith('zh')) {
    return 'zh-CN';
  }
  return 'en-US';
}

function setLocale(locale) {
  localStorage.setItem(LOCALE_STORAGE_KEY, locale);
}

export function I18nProvider({ children }) {
  const [locale, setLocaleState] = React.useState(getInitialLocale);

  const changeLocale = React.useCallback((newLocale) => {
    if (messages[newLocale]) {
      setLocale(newLocale);
      setLocaleState(newLocale);
    }
  }, []);

  const t = React.useCallback((key, params) => {
    if (params) {
      return translateParams(messages, locale, key, params);
    }
    return translate(messages, locale, key);
  }, [locale]);

  const value = React.useMemo(() => ({
    locale,
    changeLocale,
    t,
  }), [locale, changeLocale, t]);

  return <I18nContext.Provider value={value}>{children}</I18nContext.Provider>;
}

export function useI18n() {
  const context = React.useContext(I18nContext);
  if (!context) {
    throw new Error('useI18n must be used within I18nProvider');
  }
  return context;
}

export function withI18n(Component) {
  return function WrappedComponent(props) {
    const i18n = useI18n();
    return <Component {...props} i18n={i18n} />;
  };
}
