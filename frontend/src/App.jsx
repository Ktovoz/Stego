import * as React from 'react';
import { createPortal } from 'react-dom';
import { ThemeProvider, useTheme } from './components/theme-provider';
import { I18nProvider, useI18n } from './i18n/i18n-context';
import { LiquidBackground } from './components/liquid-background';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from './components/ui/card';
import { Button } from './components/ui/button';
import { Input } from './components/ui/input';
import { Label } from './components/ui/label';
import { Progress } from './components/ui/progress';
import { Badge } from './components/ui/badge';
import { Alert, AlertDescription } from './components/ui/alert';
import { DatePicker } from './components/ui/date-picker';
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from './components/ui/table';
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from './components/ui/dropdown-menu';
import {
  GetAppInfo,
  GetConfig,
  SaveConfig,
  StartEncrypt,
  CancelEncrypt,
  StartDecrypt,
  CancelDecrypt,
  StartGenerateCarrier,
  CancelGenerate,
  OpenDirectoryDialog,
  OpenFileDialog,
  GetLogs,
  GetLogsCount,
  ExportLogsToFile,
  ClearLogs,
  LogUserAction,
} from '../wailsjs/go/main/App';
import { EventsOn, WindowMinimise, Quit } from '../wailsjs/runtime/runtime';

const logAction = (module, action, details = '') => {
  try {
    LogUserAction(module, action, details);
  } catch (_) {
  }
};

const Icons = {
  Moon: () => (
    <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
      <path d="M12 3a6 6 0 0 0 9 9 9 9 0 1 1-9-9Z" />
    </svg>
  ),
  Sun: () => (
    <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
      <circle cx="12" cy="12" r="4" />
      <path d="M12 2v2" />
      <path d="M12 20v2" />
      <path d="m4.93 4.93 1.41 1.41" />
      <path d="m17.66 17.66 1.41 1.41" />
      <path d="M2 12h2" />
      <path d="M20 12h2" />
      <path d="m6.34 17.66-1.41 1.41" />
      <path d="m19.07 4.93-1.41 1.41" />
    </svg>
  ),
  Lock: () => (
    <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
      <rect width="18" height="11" x="3" y="11" rx="2" ry="2" />
      <path d="M7 11V7a5 5 0 0 1 10 0v4" />
    </svg>
  ),
  Unlock: () => (
    <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
      <rect width="18" height="11" x="3" y="11" rx="2" ry="2" />
      <path d="M7 11V7a5 5 0 0 1 9.9-1" />
    </svg>
  ),
  Image: () => (
    <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
      <rect width="18" height="18" x="3" y="3" rx="2" ry="2" />
      <circle cx="9" cy="9" r="2" />
      <path d="m21 15-3.086-3.086a2 2 0 0 0-2.828 0L6 21" />
    </svg>
  ),
  Settings: () => (
    <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
      <path d="M12.22 2h-.44a2 2 0 0 0-2 2v.18a2 2 0 0 1-1 1.73l-.43.25a2 2 0 0 1-2 0l-.15-.08a2 2 0 0 0-2.73.73l-.22.38a2 2 0 0 0 .73 2.73l.15.1a2 2 0 0 1 1 1.72v.51a2 2 0 0 1-1 1.74l-.15.09a2 2 0 0 0-.73 2.73l.22.38a2 2 0 0 0 2.73.73l.15-.08a2 2 0 0 1 2 0l.43.25a2 2 0 0 1 1 1.73V20a2 2 0 0 0 2 2h.44a2 2 0 0 0 2-2v-.18a2 2 0 0 1 1-1.73l.43-.25a2 2 0 0 1 2 0l.15.08a2 2 0 0 0 2.73-.73l.22-.39a2 2 0 0 0-.73-2.73l-.15-.09a2 2 0 0 1-1-1.74v-.47a2 2 0 0 1 1-1.74l.15-.09a2 2 0 0 0 .73-2.73l-.22-.39a2 2 0 0 0-2.73-.73l-.15.08a2 2 0 0 1-2 0l-.43-.25a2 2 0 0 1-1-1.73V4a2 2 0 0 0-2-2z" />
      <circle cx="12" cy="12" r="3" />
    </svg>
  ),
  Info: () => (
    <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
      <circle cx="12" cy="12" r="10" />
      <path d="M12 16v-4" />
      <path d="M12 8h.01" />
    </svg>
  ),
  FileText: () => (
    <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
      <path d="M14.5 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V7.5L14.5 2z" />
      <polyline points="14 2 14 8 20 8" />
      <line x1="16" x2="8" y1="13" y2="13" />
      <line x1="16" x2="8" y1="17" y2="17" />
      <line x1="10" x2="8" y1="9" y2="9" />
    </svg>
  ),
  Download: () => (
    <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
      <path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4" />
      <polyline points="7 10 12 15 17 10" />
      <line x1="12" x2="12" y1="15" y2="3" />
    </svg>
  ),
  Trash2: () => (
    <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
      <polyline points="3 6 5 6 21 6" />
      <path d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2" />
      <line x1="10" x2="10" y1="11" y2="17" />
      <line x1="14" x2="14" y1="11" y2="17" />
    </svg>
  ),
  Home: () => (
    <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
      <path d="m3 9 9-7 9 7v11a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2z" />
      <polyline points="9 22 9 12 15 12 15 22" />
    </svg>
  ),
  FolderOpen: () => (
    <svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
      <path d="m6 14 1.45-2.9A2 2 0 0 1 9.24 10H20a2 2 0 0 1 1.94 2.5l-1.55 6a2 2 0 0 1-1.94 1.5H4a2 2 0 0 1-2-2V5c0-1.1.9-2 2-2h3.93a2 2 0 0 1 1.66.9l.82 1.2a2 2 0 0 0 1.66.9H18a2 2 0 0 1 2 2v2" />
    </svg>
  ),
  Eye: () => (
    <svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
      <path d="M2 12s3-7 10-7 10 7 10 7-3-7-10-7Z" />
      <circle cx="12" cy="12" r="3" />
    </svg>
  ),
  EyeOff: () => (
    <svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
      <path d="M9.88 9.88a3 3 0 1 0 4.24 4.24" />
      <path d="M10.73 5.08A10.43 10.43 0 0 1 12 5c7 0 10 7 10 7a13.16 13.16 0 0 1-1.67.14" />
      <path d="m6.61 6.61A13.526 13.526 0 0 0 2 12s3 7 10 7c.41 0 .81-.04 1.2-.1" />
      <line x1="2" x2="22" y1="2" y2="22" />
    </svg>
  ),
  RotateCcw: () => (
    <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
      <path d="M3 12a9 9 0 1 0 9-9 9.75 9.75 0 0 0-6.74 2.74L3 8" />
      <path d="M3 3v5h5" />
    </svg>
  ),
  Globe: () => (
    <svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
      <circle cx="12" cy="12" r="10" />
      <path d="M12 2a14.5 14.5 0 0 0 0 20 14.5 14.5 0 0 0 0-20" />
      <path d="M2 12h20" />
    </svg>
  ),
  X: () => (
    <svg xmlns="http://www.w3.org/2000/svg" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
      <path d="M18 6 6 18" />
      <path d="m6 6 12 12" />
    </svg>
  ),
  Minus: () => (
    <svg xmlns="http://www.w3.org/2000/svg" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
      <path d="M5 12h14" />
    </svg>
  ),
  Github: () => (
    <svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
      <path d="M15 22v-4a4.8 4.8 0 0 0-1-3.5c3 0 6-2 6-5.5.08-1.25-.27-2.48-1-3.5.28-1.15.28-2.35 0-3.5 0 0-1-1-3.5 1.5-1.67-1-3.5 1.5-3.5 1.5-2.83-2-3.5-1.5-3.5-1.5-2.5 0-3.5 1.5-3.5 1.5-2.5-1.5-3.5-1.5-3.5 0 0 1-1 3.5 1.5c-.63 1.03-.97 2.23-1 3.5 0 3.5 3 5.5 6 5.5-.39.49-.68 1.05-.85 1.65-.17.6-.22 1.23-.15 1.85v4" />
      <path d="M9 18c-4.51 2-5-2-7-2" />
    </svg>
  ),
};

const MenuItem = ({ icon: Icon, active, onClick, labelKey }) => {
  const { t } = useI18n();
  return (
    <button
      onClick={onClick}
      aria-current={active ? 'page' : undefined}
      className={`group w-full flex items-center gap-3 px-4 py-4 text-base rounded-xl border border-transparent box-border focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 ring-offset-background overflow-hidden transition-gentle ${
        active
          ? 'bg-foreground/5 text-foreground shadow-glass-sm border-foreground/10 ring-1 ring-foreground/10 dark:bg-primary/20 dark:border-primary/30 dark:ring-primary/35'
          : 'text-muted-foreground hover:bg-foreground/4 hover:shadow-glass-sm hover:border-border/50 hover:text-foreground dark:hover:bg-card/80 dark:hover:backdrop-blur-xl dark:hover:shadow-glass-md'
      }`}
    >
      <span className={active ? 'text-foreground dark:text-primary' : 'text-muted-foreground group-hover:text-foreground'}>
        <Icon />
      </span>
      <span>{t(`menu.${labelKey}`)}</span>
    </button>
  );
};

const EncryptPage = ({ config }) => {
  const { t } = useI18n();
  const [formData, setFormData] = React.useState({
    dataSourcePath: '',
    carrierDir: config.defaultCarrierDir || '',
    carrierImagePath: '',
    outputDir: config.defaultOutputDir || '',
    outputFileName: config.defaultEncryptOutputName || '',
    password: config.defaultEncryptPassword || '',
    scatter: true,
  });
  const [progress, setProgress] = React.useState(0);
  const [status, setStatus] = React.useState('');
  const [statusType, setStatusType] = React.useState('info');
  const [task, setTask] = React.useState(null);
  const [isRunning, setIsRunning] = React.useState(false);
  const [showPassword, setShowPassword] = React.useState(false);

  React.useEffect(() => {
    const handler = (p) => {
      setProgress(p.progress);
      setStatus(p.error || p.message);
      setStatusType(p.error ? 'error' : 'info');
      if (p.done) {
        setIsRunning(false);
        setTask(null);
      }
    };
    EventsOn('encryptProgress', handler);
    return () => {};
  }, []);

  const handleStart = async () => {
    try {
      setProgress(0);
      setStatus(t('encrypt.starting'));
      setStatusType('info');
      setIsRunning(true);
      logAction('encrypt', '开始加密', `数据源: ${formData.dataSourcePath}`);

      const taskId = await StartEncrypt({
        ...formData,
        identifier: 'stego',
      });
      setTask(taskId);
    } catch (e) {
      setStatus(String(e));
      setStatusType('error');
      setIsRunning(false);
    }
  };

  const handleCancel = async () => {
    if (!task) return;
    try {
      await CancelEncrypt(task);
      setStatus(t('encrypt.cancelled'));
      setIsRunning(false);
      setTask(null);
      logAction('encrypt', '取消加密', `任务ID: ${task}`);
    } catch (e) {
      setStatus(String(e));
      setStatusType('error');
    }
  };

  return (
    <div className="h-full">
      <Card className="h-full flex flex-col">
      <CardHeader className="pb-3">
        <CardTitle className="text-lg">{t('encrypt.title')}</CardTitle>
        <CardDescription className="text-xs">{t('encrypt.description')}</CardDescription>
      </CardHeader>
      <CardContent className="flex-1 overflow-auto space-y-3">
        <div className="space-y-1.5">
          <Label htmlFor="enc-data" className="text-xs">{t('encrypt.dataSourceDir')}</Label>
          <div className="flex gap-2">
            <Input
              id="enc-data"
              value={formData.dataSourcePath}
              onChange={(e) => setFormData({ ...formData, dataSourcePath: e.target.value })}
              disabled={isRunning}
              className="h-9 text-sm flex-1"
            />
            <Button
              variant="outline"
              size="sm"
              onClick={async () => {
                const result = await OpenDirectoryDialog("");
                if (result) setFormData({ ...formData, dataSourcePath: result });
              }}
              disabled={isRunning}
              className="h-9 px-3"
              title={t('encrypt.selectDirectory')}
            >
              <Icons.FolderOpen />
            </Button>
          </div>
        </div>

        <div className="space-y-1.5">
          <Label htmlFor="enc-carrierDir" className="text-xs">{t('encrypt.carrierDir')}</Label>
          <div className="flex gap-2">
            <Input
              id="enc-carrierDir"
              placeholder={config.defaultCarrierDir || undefined}
              value={formData.carrierDir}
              onChange={(e) => setFormData({ ...formData, carrierDir: e.target.value })}
              disabled={isRunning}
              className="h-9 text-sm flex-1"
            />
            <Button
              variant="outline"
              size="sm"
              onClick={async () => {
                const result = await OpenDirectoryDialog("");
                if (result) setFormData({ ...formData, carrierDir: result });
              }}
              disabled={isRunning}
              className="h-9 px-3"
              title={t('encrypt.selectDirectory')}
            >
              <Icons.FolderOpen />
            </Button>
          </div>
        </div>

        <div className="space-y-1.5">
          <Label htmlFor="enc-out" className="text-xs">{t('encrypt.outputDir')}</Label>
          <div className="flex gap-2">
            <Input
              id="enc-out"
              placeholder={config.defaultOutputDir || undefined}
              value={formData.outputDir}
              onChange={(e) => setFormData({ ...formData, outputDir: e.target.value })}
              disabled={isRunning}
              className="h-9 text-sm flex-1"
            />
            <Button
              variant="outline"
              size="sm"
              onClick={async () => {
                const result = await OpenDirectoryDialog("");
                if (result) setFormData({ ...formData, outputDir: result });
              }}
              disabled={isRunning}
              className="h-9 px-3"
              title={t('encrypt.selectDirectory')}
            >
              <Icons.FolderOpen />
            </Button>
          </div>
        </div>

        <div className="space-y-1.5">
          <Label htmlFor="enc-pass" className="text-xs">{t('encrypt.password')}</Label>
          <div className="flex gap-2">
            <Input
              id="enc-pass"
              type={showPassword ? 'text' : 'password'}
              placeholder={t('encrypt.passwordPlaceholder')}
              value={formData.password}
              onChange={(e) => setFormData({ ...formData, password: e.target.value })}
              disabled={isRunning}
              className="h-9 text-sm flex-1"
            />
            <Button
              variant="outline"
              size="sm"
              onClick={() => setShowPassword(!showPassword)}
              disabled={isRunning}
              className="h-9 px-3"
              title={showPassword ? t('encrypt.hidePassword') : t('encrypt.showPassword')}
            >
              {showPassword ? <Icons.EyeOff /> : <Icons.Eye />}
            </Button>
          </div>
        </div>

        <div className="space-y-1.5">
          <Label htmlFor="enc-outputName" className="text-xs">{t('encrypt.outputFileName')}</Label>
          <Input
            id="enc-outputName"
            placeholder={config.defaultEncryptOutputName || t('encrypt.outputFileNamePlaceholder')}
            value={formData.outputFileName}
            onChange={(e) => setFormData({ ...formData, outputFileName: e.target.value })}
            disabled={isRunning}
            className="h-9 text-sm"
          />
        </div>

        <div className="flex gap-2">
          <Button onClick={handleStart} disabled={isRunning || !formData.dataSourcePath} size="sm">
            {t('encrypt.start')}
          </Button>
          <Button variant="outline" onClick={handleCancel} disabled={!isRunning} size="sm">
            {t('encrypt.cancel')}
          </Button>
        </div>

        {status && (
          <>
            <Progress value={progress} className="h-1.5" />
            <p className={`text-xs ${statusType === 'error' ? 'text-destructive' : 'text-muted-foreground'}`}>
              {status}
            </p>
          </>
        )}
      </CardContent>
      </Card>
    </div>
  );
};

const DecryptPage = ({ config }) => {
  const { t } = useI18n();
  const [formData, setFormData] = React.useState({
    imagePath: '',
    outputDir: config.defaultOutputDir || '',
    password: config.defaultDecryptPassword || '',
  });
  const [progress, setProgress] = React.useState(0);
  const [status, setStatus] = React.useState('');
  const [statusType, setStatusType] = React.useState('info');
  const [task, setTask] = React.useState(null);
  const [isRunning, setIsRunning] = React.useState(false);
  const [showPassword, setShowPassword] = React.useState(false);

  React.useEffect(() => {
    const handler = (p) => {
      setProgress(p.progress);
      setStatus(p.error || p.message);
      setStatusType(p.error ? 'error' : 'info');
      if (p.done) {
        setIsRunning(false);
        setTask(null);
      }
    };
    EventsOn('decryptProgress', handler);
    return () => {};
  }, []);

  const handleStart = async () => {
    try {
      setProgress(0);
      setStatus(t('decrypt.starting'));
      setStatusType('info');
      setIsRunning(true);
      logAction('decrypt', '开始解密', `图片路径: ${formData.imagePath}`);

      const taskId = await StartDecrypt({
        ...formData,
        identifier: 'stego',
      });
      setTask(taskId);
    } catch (e) {
      setStatus(String(e));
      setStatusType('error');
      setIsRunning(false);
    }
  };

  const handleCancel = async () => {
    if (!task) return;
    try {
      await CancelDecrypt(task);
      setStatus(t('decrypt.cancelled'));
      setIsRunning(false);
      setTask(null);
      logAction('decrypt', '取消解密', `任务ID: ${task}`);
    } catch (e) {
      setStatus(String(e));
      setStatusType('error');
    }
  };

  return (
    <Card className="h-full flex flex-col">
      <CardHeader className="pb-3">
        <CardTitle className="text-lg">{t('decrypt.title')}</CardTitle>
        <CardDescription className="text-xs">{t('decrypt.description')}</CardDescription>
      </CardHeader>
      <CardContent className="flex-1 overflow-auto space-y-3">
        <div className="space-y-1.5">
          <Label htmlFor="dec-img" className="text-xs">{t('decrypt.imagePath')}</Label>
          <div className="flex gap-2">
            <Input
              id="dec-img"
              value={formData.imagePath}
              onChange={(e) => setFormData({ ...formData, imagePath: e.target.value })}
              disabled={isRunning}
              className="h-9 text-sm flex-1"
            />
            <Button
              variant="outline"
              size="sm"
              onClick={async () => {
                const result = await OpenFileDialog("");
                if (result) setFormData({ ...formData, imagePath: result });
              }}
              disabled={isRunning}
              className="h-9 px-3"
              title={t('decrypt.selectFile')}
            >
              <Icons.FolderOpen />
            </Button>
          </div>
        </div>

        <div className="space-y-1.5">
          <Label htmlFor="dec-out" className="text-xs">{t('decrypt.outputDir')}</Label>
          <div className="flex gap-2">
            <Input
              id="dec-out"
              placeholder={config.defaultOutputDir || undefined}
              value={formData.outputDir}
              onChange={(e) => setFormData({ ...formData, outputDir: e.target.value })}
              disabled={isRunning}
              className="h-9 text-sm flex-1"
            />
            <Button
              variant="outline"
              size="sm"
              onClick={async () => {
                const result = await OpenDirectoryDialog("");
                if (result) setFormData({ ...formData, outputDir: result });
              }}
              disabled={isRunning}
              className="h-9 px-3"
              title={t('decrypt.selectDirectory')}
            >
              <Icons.FolderOpen />
            </Button>
          </div>
        </div>

        <div className="space-y-1.5">
          <Label htmlFor="dec-pass" className="text-xs">{t('decrypt.password')}</Label>
          <div className="flex gap-2">
            <Input
              id="dec-pass"
              type={showPassword ? 'text' : 'password'}
              placeholder={t('decrypt.passwordPlaceholder')}
              value={formData.password}
              onChange={(e) => setFormData({ ...formData, password: e.target.value })}
              disabled={isRunning}
              className="h-9 text-sm flex-1"
            />
            <Button
              variant="outline"
              size="sm"
              onClick={() => setShowPassword(!showPassword)}
              disabled={isRunning}
              className="h-9 px-3"
              title={showPassword ? t('decrypt.hidePassword') : t('decrypt.showPassword')}
            >
              {showPassword ? <Icons.EyeOff /> : <Icons.Eye />}
            </Button>
          </div>
        </div>

        <div className="flex gap-2">
          <Button onClick={handleStart} disabled={isRunning || !formData.imagePath} size="sm">
            {t('decrypt.start')}
          </Button>
          <Button variant="outline" onClick={handleCancel} disabled={!isRunning} size="sm">
            {t('decrypt.cancel')}
          </Button>
        </div>

        {status && (
          <>
            <Progress value={progress} className="h-1.5" />
            <p className={`text-xs ${statusType === 'error' ? 'text-destructive' : 'text-muted-foreground'}`}>
              {status}
            </p>
          </>
        )}
      </CardContent>
    </Card>
  );
};

const GeneratePage = ({ config }) => {
  const { t } = useI18n();
  const [formData, setFormData] = React.useState({
    outputDir: config.defaultCarrierDir || '',
    targetMB: 10,
    count: 5,
    noiseEnabled: true,
  });
  const [progress, setProgress] = React.useState(0);
  const [status, setStatus] = React.useState('');
  const [statusType, setStatusType] = React.useState('info');
  const [task, setTask] = React.useState(null);
  const [isRunning, setIsRunning] = React.useState(false);

  React.useEffect(() => {
    const handler = (p) => {
      setProgress(p.progress);
      const current = p.current ?? 0;
      const total = p.total ?? 0;
      const displayCurrent = current > 0 ? current : 0;
      const msg = total > 0 && !p.done ? `${p.message} (${displayCurrent}/${total})` : p.message;
      setStatus(msg);
      setStatusType(p.error ? 'error' : 'info');
      if (p.done) {
        setIsRunning(false);
        setTask(null);
      }
    };
    EventsOn('generateProgress', handler);
    return () => {};
  }, []);

  const handleStart = async () => {
    try {
      setProgress(0);
      setStatus(t('generate.starting'));
      setStatusType('info');
      setIsRunning(true);
      logAction('generate', '开始生成载体', `容量: ${formData.targetMB}MB, 数量: ${formData.count}`);

      const targetBytes = formData.targetMB * 1024 * 1024;

      const taskId = await StartGenerateCarrier({
        outputDir: formData.outputDir,
        targetBytes: targetBytes,
        count: formData.count,
        prefix: 'carrier',
        noiseEnabled: formData.noiseEnabled,
      });
      setTask(taskId);
    } catch (e) {
      setStatus(String(e));
      setStatusType('error');
      setIsRunning(false);
    }
  };

  const handleSelectOutputDir = async () => {
    try {
      const result = await OpenDirectoryDialog("");
      if (result) {
        setFormData({ ...formData, outputDir: result });
      }
    } catch (e) {
      console.error('Failed to open directory dialog:', e);
    }
  };

  const handleCancel = async () => {
    if (!task) return;
    try {
      await CancelGenerate(task);
      setStatus(t('generate.cancelled'));
      setIsRunning(false);
      setTask(null);
      logAction('generate', '取消生成', `任务ID: ${task}`);
    } catch (e) {
      setStatus(String(e));
      setStatusType('error');
    }
  };

  return (
    <Card className="h-full flex flex-col">
      <CardHeader className="pb-3">
        <CardTitle className="text-lg">{t('generate.title')}</CardTitle>
        <CardDescription className="text-xs">{t('generate.description')}</CardDescription>
      </CardHeader>
      <CardContent className="flex-1 overflow-auto space-y-3">
        <div className="space-y-1.5">
          <Label htmlFor="gen-out" className="text-xs">{t('generate.outputDir')}</Label>
          <div className="flex gap-2">
            <Input
              id="gen-out"
              placeholder={config.defaultCarrierDir || undefined}
              value={formData.outputDir}
              onChange={(e) => setFormData({ ...formData, outputDir: e.target.value })}
              disabled={isRunning}
              className="h-9 text-sm flex-1"
            />
            <Button
              variant="outline"
              size="sm"
              onClick={handleSelectOutputDir}
              disabled={isRunning}
              className="h-9 px-3"
              title={t('generate.selectDirectory')}
            >
              <Icons.FolderOpen />
            </Button>
          </div>
        </div>

        <div className="space-y-1.5">
          <Label htmlFor="gen-mb" className="text-xs">{t('generate.targetMB')}</Label>
          <Input
            id="gen-mb"
            type="number"
            placeholder="10"
            value={formData.targetMB}
            onChange={(e) => setFormData({ ...formData, targetMB: parseInt(e.target.value) || 0 })}
            disabled={isRunning}
            className="h-9 text-sm"
          />
        </div>

        <div className="space-y-1.5">
          <Label htmlFor="gen-count" className="text-xs">{t('generate.count')}</Label>
          <Input
            id="gen-count"
            type="number"
            value={formData.count}
            onChange={(e) => setFormData({ ...formData, count: parseInt(e.target.value) || 1 })}
            disabled={isRunning}
            className="h-9 text-sm"
          />
        </div>

        <div className="flex gap-2">
          <Button onClick={handleStart} disabled={isRunning || !formData.outputDir} size="sm">
            {t('generate.start')}
          </Button>
          <Button variant="outline" onClick={handleCancel} disabled={!isRunning} size="sm">
            {t('generate.cancel')}
          </Button>
        </div>

        {status && (
          <>
            <Progress value={progress} className="h-1.5" />
            <p className={`text-xs ${statusType === 'error' ? 'text-destructive' : 'text-muted-foreground'}`}>
              {status}
            </p>
          </>
        )}
      </CardContent>
    </Card>
  );
};

const Modal = ({ show, onClose, title, children, footer }) => {
  if (!show) return null;

  return createPortal(
    <div className="fixed inset-0 z-[9999] flex items-center justify-center" onClick={onClose}>
      <div className="absolute inset-0 bg-black/60 backdrop-blur-sm" />
      <div
        className="relative z-10 w-full max-w-md rounded-2xl border bg-glass/90 backdrop-blur-xl shadow-glass-xl animate-in fade-in zoom-in duration-200"
        onClick={(e) => e.stopPropagation()}
      >
        {title && (
          <div className="px-6 py-4 border-b border-border/50">
            <div className="text-sm font-semibold">{title}</div>
          </div>
        )}
        <div className="p-6 text-sm text-muted-foreground">
          {children}
        </div>
        {footer && (
          <div className="px-6 pb-6 pt-0 flex justify-end gap-2">
            {footer}
          </div>
        )}
      </div>
    </div>,
    document.body
  );
};

const LogsPage = () => {
  const { t } = useI18n();
  const [logs, setLogs] = React.useState([]);
  const [totalCount, setTotalCount] = React.useState(0);
  const [levelFilter, setLevelFilter] = React.useState('ALL');

  const [startDate, setStartDate] = React.useState('');
  const [endDate, setEndDate] = React.useState('');
  const [exportFormat, setExportFormat] = React.useState('txt');
  const [showClearConfirm, setShowClearConfirm] = React.useState(false);
  const [isClearing, setIsClearing] = React.useState(false);
  const [isResetting, setIsResetting] = React.useState(false);

  const loadLogs = React.useCallback(async () => {
    try {
      let start = 0;
      let end = 0;

      if (startDate) {
        const [year, month, day] = startDate.split('-').map(Number);
        const startDateObj = new Date(year, month - 1, day, 0, 0, 0);
        start = Math.floor(startDateObj.getTime() / 1000);
      }

      if (endDate) {
        const [year, month, day] = endDate.split('-').map(Number);
        const endDateObj = new Date(year, month - 1, day, 23, 59, 59);
        end = Math.floor(endDateObj.getTime() / 1000);
      }

      const result = await GetLogs(levelFilter, start, end, 100, 0);
      setLogs(result || []);
      const count = await GetLogsCount();
      setTotalCount(count || 0);
    } catch (e) {
      console.error('Failed to load logs:', e);
      setLogs([]);
      setTotalCount(0);
    }
  }, [levelFilter, startDate, endDate]);

  const handleExport = async () => {
    try {
      let start = 0;
      let end = 0;

      if (startDate) {
        const [year, month, day] = startDate.split('-').map(Number);
        const startDateObj = new Date(year, month - 1, day, 0, 0, 0);
        start = Math.floor(startDateObj.getTime() / 1000);
      }

      if (endDate) {
        const [year, month, day] = endDate.split('-').map(Number);
        const endDateObj = new Date(year, month - 1, day, 23, 59, 59);
        end = Math.floor(endDateObj.getTime() / 1000);
      }

      const savePath = await ExportLogsToFile(exportFormat, start, end);
      if (!savePath) return;
      logAction('logs', '导出日志', `格式: ${exportFormat}, 日期范围: ${startDate || '全部'} - ${endDate || '全部'}`);

    } catch (e) {
      console.error('Failed to export logs:', e);
    }
  };



  const handleClear = () => {
    setShowClearConfirm(true);
  };

  const confirmClear = async () => {
    if (isClearing) return;
    setIsClearing(true);
    try {
      await ClearLogs();
      logAction('logs', '清空日志', '');
      setShowClearConfirm(false);
      loadLogs();
    } catch (e) {
      console.error('Failed to clear logs:', e);
    } finally {
      setIsClearing(false);
    }
  };

  const cancelClear = () => {
    if (isClearing) return;
    setShowClearConfirm(false);
  };

  const handleResetDate = async () => {
    setIsResetting(true);
    await new Promise(resolve => setTimeout(resolve, 300));
    setStartDate('');
    setEndDate('');
    logAction('logs', '重置日期筛选', '');
    setIsResetting(false);
  };

  React.useEffect(() => {
    loadLogs();
  }, [loadLogs]);

  const getLevelColor = (level) => {
    switch (level) {
      case 'ERROR': return 'text-destructive';
      case 'WARN': return 'text-orange-500';
      case 'INFO': return 'text-blue-500';
      default: return 'text-muted-foreground';
    }
  };

  const formatDate = (timestamp) => {
    const date = new Date(timestamp);
    return date.toLocaleString('zh-CN', {
      year: 'numeric',
      month: '2-digit',
      day: '2-digit',
      hour: '2-digit',
      minute: '2-digit',
      second: '2-digit',
    });
  };

  return (
    <Card className="h-full flex flex-col">
      <CardHeader className="pb-3">
        <CardTitle className="text-lg">{t('logs.title')}</CardTitle>
        <CardDescription className="text-xs">{t('logs.description')}</CardDescription>
      </CardHeader>
      <CardContent className="flex-1 overflow-auto space-y-3">
        <div className="flex gap-2 flex-wrap">
          <div className="flex-1 min-w-[120px]">
            <Label htmlFor="log-level" className="text-xs">{t('logs.level')}</Label>
            <select
              id="log-level"
              value={levelFilter}
              onChange={(e) => {
                setLevelFilter(e.target.value);
                logAction('logs', '更改级别筛选', `级别: ${e.target.value}`);
              }}
              className="w-full h-9 px-2 text-sm border rounded-md bg-background"
            >
              <option value="ALL">{t('logs.levels.all')}</option>
              <option value="ERROR">{t('logs.levels.error')}</option>
              <option value="WARN">{t('logs.levels.warn')}</option>
              <option value="INFO">{t('logs.levels.info')}</option>
            </select>
          </div>
          <div className="flex-1 min-w-[140px]">
            <Label htmlFor="log-start" className="text-xs">{t('logs.startDate')}</Label>
            <DatePicker
              id="log-start"
              value={startDate}
              onChange={(val) => {
                setStartDate(val);
                logAction('logs', '更改开始日期', val);
              }}
              placeholder={t('logs.selectStartDate')}
              className="w-full"
            />
          </div>
          <div className="flex-1 min-w-[140px]">
            <Label htmlFor="log-end" className="text-xs">{t('logs.endDate')}</Label>
            <DatePicker
              id="log-end"
              value={endDate}
              onChange={(val) => {
                setEndDate(val);
                logAction('logs', '更改结束日期', val);
              }}
              placeholder={t('logs.selectEndDate')}
              className="w-full"
            />
          </div>
          <div className="flex-1 min-w-[100px]">
            <Label htmlFor="log-format" className="text-xs">{t('logs.exportFormat')}</Label>
            <select
              id="log-format"
              value={exportFormat}
              onChange={(e) => {
                setExportFormat(e.target.value);
                logAction('logs', '更改导出格式', e.target.value);
              }}
              className="w-full h-9 px-2 text-sm border rounded-md bg-background"
            >
              <option value="txt">{t('logs.formats.text')}</option>
              <option value="json">{t('logs.formats.json')}</option>
            </select>
          </div>
        </div>

        <div className="flex gap-2">
          <Button onClick={handleResetDate} variant="outline" size="sm" disabled={isResetting}>
            <Icons.RotateCcw className={isResetting ? 'animate-spin' : ''} />
            <span className="ml-1">{t('logs.resetDate')}</span>
          </Button>
          <Button onClick={handleExport} variant="outline" size="sm">
            <Icons.Download />
            <span className="ml-1">{t('logs.export')}</span>
          </Button>
          <Button onClick={handleClear} variant="outline" size="sm" className="text-destructive hover:text-destructive">
            <Icons.Trash2 />
            <span className="ml-1">{t('logs.clear')}</span>
          </Button>
          <div className="flex-1" />
          <span className="text-xs text-muted-foreground self-center">{t('logs.total', { count: totalCount })}</span>
        </div>

        <div className="border rounded-md overflow-auto flex-1" style={{ maxHeight: 'calc(100vh - 400px)' }}>
          {logs.length === 0 ? (
            <div className="flex items-center justify-center h-32 text-sm text-muted-foreground">
              {t('logs.noLogs')}
            </div>
          ) : (
            <Table>
              <TableHeader>
                <TableRow>
                  <TableHead className="w-[140px]">{t('logs.columns.time')}</TableHead>
                  <TableHead className="w-[60px]">{t('logs.columns.level')}</TableHead>
                  <TableHead className="w-[80px]">{t('logs.columns.module')}</TableHead>
                  <TableHead>{t('logs.columns.message')}</TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                {logs.map((log) => (
                  <TableRow key={log.id}>
                    <TableCell className="text-xs">{formatDate(log.timestamp)}</TableCell>
                    <TableCell className={`text-xs font-medium ${getLevelColor(log.level)}`}>
                      {log.level}
                    </TableCell>
                    <TableCell className="text-xs">{log.module}</TableCell>
                    <TableCell className="text-xs">
                      {log.message}
                      {log.details && <span className="text-muted-foreground ml-2">- {log.details}</span>}
                    </TableCell>
                  </TableRow>
                ))}
              </TableBody>
            </Table>
          )}
        </div>
      </CardContent>

      <Modal
        show={showClearConfirm}
        onClose={cancelClear}
        title={t('logs.clear')}
        footer={
          <>
            <Button variant="secondary" onClick={cancelClear} disabled={isClearing}>
              取消
            </Button>
            <Button
              variant="destructive"
              onClick={confirmClear}
              disabled={isClearing}
              className="hover:bg-destructive/90 hover:text-destructive-foreground"
            >
              清空
            </Button>
          </>
        }
      >
        {t('logs.confirmClear')}
      </Modal>
    </Card>
  );
};

const SettingsPage = ({ config, onConfigChange }) => {
  const { t } = useI18n();
  const [formData, setFormData] = React.useState({
    defaultCarrierDir: config.defaultCarrierDir || '',
    defaultOutputDir: config.defaultOutputDir || '',
    defaultEncryptOutputName: config.defaultEncryptOutputName || '',
    defaultEncryptPassword: config.defaultEncryptPassword || '',
    defaultDecryptPassword: config.defaultDecryptPassword || '',
  });
  const [status, setStatus] = React.useState('');
  const [showEncPassword, setShowEncPassword] = React.useState(false);
  const [showDecPassword, setShowDecPassword] = React.useState(false);

  const handleSave = async () => {
    try {
      setStatus(t('settings.saving'));
      await SaveConfig(formData);
      setStatus(t('settings.saved'));
      onConfigChange(formData);
      logAction('config', '保存设置', '');
    } catch (e) {
      setStatus(String(e));
    }
  };

  return (
    <Card className="h-full flex flex-col">
      <CardHeader className="pb-3">
        <CardTitle className="text-lg">{t('settings.title')}</CardTitle>
        <CardDescription className="text-xs">{t('settings.description')}</CardDescription>
      </CardHeader>
      <CardContent className="flex-1 overflow-auto space-y-3">
        <div className="space-y-1.5">
          <Label htmlFor="cfg-carrierDir" className="text-xs">{t('settings.defaultCarrierDir')}</Label>
          <div className="flex gap-2">
            <Input
              id="cfg-carrierDir"
              value={formData.defaultCarrierDir}
              onChange={(e) => setFormData({ ...formData, defaultCarrierDir: e.target.value })}
              className="h-9 text-sm flex-1"
            />
            <Button
              variant="outline"
              size="sm"
              onClick={async () => {
                const result = await OpenDirectoryDialog("");
                if (result) setFormData({ ...formData, defaultCarrierDir: result });
              }}
              className="h-9 px-3"
              title={t('settings.selectDirectory')}
            >
              <Icons.FolderOpen />
            </Button>
          </div>
        </div>

        <div className="space-y-1.5">
          <Label htmlFor="cfg-outputDir" className="text-xs">{t('settings.defaultOutputDir')}</Label>
          <div className="flex gap-2">
            <Input
              id="cfg-outputDir"
              value={formData.defaultOutputDir}
              onChange={(e) => setFormData({ ...formData, defaultOutputDir: e.target.value })}
              className="h-9 text-sm flex-1"
            />
            <Button
              variant="outline"
              size="sm"
              onClick={async () => {
                const result = await OpenDirectoryDialog("");
                if (result) setFormData({ ...formData, defaultOutputDir: result });
              }}
              className="h-9 px-3"
              title={t('settings.selectDirectory')}
            >
              <Icons.FolderOpen />
            </Button>
          </div>
        </div>

        <div className="space-y-1.5">
          <Label htmlFor="cfg-outputName" className="text-xs">{t('settings.defaultEncryptOutputName')}</Label>
          <Input
            id="cfg-outputName"
            value={formData.defaultEncryptOutputName}
            onChange={(e) => setFormData({ ...formData, defaultEncryptOutputName: e.target.value })}
            className="h-9 text-sm"
          />
        </div>

        <div className="space-y-1.5">
          <Label htmlFor="cfg-encPass" className="text-xs">{t('settings.encryptPassword')}</Label>
          <div className="flex gap-2">
            <Input
              id="cfg-encPass"
              type={showEncPassword ? 'text' : 'password'}
              value={formData.defaultEncryptPassword}
              onChange={(e) => setFormData({ ...formData, defaultEncryptPassword: e.target.value })}
              className="h-9 text-sm flex-1"
            />
            <Button
              variant="outline"
              size="sm"
              onClick={() => setShowEncPassword(!showEncPassword)}
              className="h-9 px-3"
              title={showEncPassword ? t('settings.hidePassword') : t('settings.showPassword')}
            >
              {showEncPassword ? <Icons.EyeOff /> : <Icons.Eye />}
            </Button>
          </div>
        </div>

        <div className="space-y-1.5">
          <Label htmlFor="cfg-decPass" className="text-xs">{t('settings.decryptPassword')}</Label>
          <div className="flex gap-2">
            <Input
              id="cfg-decPass"
              type={showDecPassword ? 'text' : 'password'}
              value={formData.defaultDecryptPassword}
              onChange={(e) => setFormData({ ...formData, defaultDecryptPassword: e.target.value })}
              className="h-9 text-sm flex-1"
            />
            <Button
              variant="outline"
              size="sm"
              onClick={() => setShowDecPassword(!showDecPassword)}
              className="h-9 px-3"
              title={showDecPassword ? t('settings.hidePassword') : t('settings.showPassword')}
            >
              {showDecPassword ? <Icons.EyeOff /> : <Icons.Eye />}
            </Button>
          </div>
        </div>

        <div className="flex gap-2 items-center">
          <Button onClick={handleSave} size="sm">{t('settings.save')}</Button>
          {status && <p className="text-xs text-muted-foreground">{status}</p>}
        </div>
      </CardContent>
    </Card>
  );
};

const AboutPage = ({ appInfo }) => {
  const { t } = useI18n();
  return (
    <Card className="h-full flex flex-col">
      <CardHeader className="pb-4">
        <CardTitle className="text-lg">{t('about.title')}</CardTitle>
        <CardDescription className="text-xs">{t('about.description')}</CardDescription>
      </CardHeader>
      <CardContent className="flex-1 overflow-auto">
        <div className="space-y-6">
          <div className="text-center space-y-4">
            <div>
              <h2 className="text-3xl font-bold">{appInfo.name || t('app.name')}</h2>
              <p className="text-sm text-muted-foreground mt-2">{t('app.tagline')}</p>
            </div>
            <Badge variant="secondary" className="text-sm px-3 py-1">
              v{appInfo.version || '-'}
            </Badge>
          </div>

          <div className="border-t" />

          <div className="space-y-4">
            <p className="text-xs font-medium text-muted-foreground uppercase tracking-wide">{t('about.features.title')}</p>
            <div className="space-y-3 text-sm">
              <div className="flex gap-2">
                <div className="w-1.5 h-1.5 rounded-full bg-primary mt-1.5 flex-shrink-0" />
                <div>
                  <p className="font-medium">{t('about.features.lsb.name')}</p>
                  <p className="text-xs text-muted-foreground">{t('about.features.lsb.desc')}</p>
                </div>
              </div>
              <div className="flex gap-2">
                <div className="w-1.5 h-1.5 rounded-full bg-primary mt-1.5 flex-shrink-0" />
                <div>
                  <p className="font-medium">{t('about.features.scatter.name')}</p>
                  <p className="text-xs text-muted-foreground">{t('about.features.scatter.desc')}</p>
                </div>
              </div>
              <div className="flex gap-2">
                <div className="w-1.5 h-1.5 rounded-full bg-primary mt-1.5 flex-shrink-0" />
                <div>
                  <p className="font-medium">{t('about.features.aes.name')}</p>
                  <p className="text-xs text-muted-foreground">{t('about.features.aes.desc')}</p>
                </div>
              </div>
              <div className="flex gap-2">
                <div className="w-1.5 h-1.5 rounded-full bg-primary mt-1.5 flex-shrink-0" />
                <div>
                  <p className="font-medium">{t('about.features.rs.name')}</p>
                  <p className="text-xs text-muted-foreground">{t('about.features.rs.desc')}</p>
                </div>
              </div>
              <div className="flex gap-2">
                <div className="w-1.5 h-1.5 rounded-full bg-primary mt-1.5 flex-shrink-0" />
                <div>
                  <p className="font-medium">{t('about.features.interleave.name')}</p>
                  <p className="text-xs text-muted-foreground">{t('about.features.interleave.desc')}</p>
                </div>
              </div>
              <div className="flex gap-2">
                <div className="w-1.5 h-1.5 rounded-full bg-primary mt-1.5 flex-shrink-0" />
                <div>
                  <p className="font-medium">{t('about.features.integrity.name')}</p>
                  <p className="text-xs text-muted-foreground">{t('about.features.integrity.desc')}</p>
                </div>
              </div>
              <div className="flex gap-2">
                <div className="w-1.5 h-1.5 rounded-full bg-primary mt-1.5 flex-shrink-0" />
                <div>
                  <p className="font-medium">{t('about.features.carrier.name')}</p>
                  <p className="text-xs text-muted-foreground">{t('about.features.carrier.desc')}</p>
                </div>
              </div>
              <div className="flex gap-2">
                <div className="w-1.5 h-1.5 rounded-full bg-primary mt-1.5 flex-shrink-0" />
                <div>
                  <p className="font-medium">{t('about.features.capacity.name')}</p>
                  <p className="text-xs text-muted-foreground">{t('about.features.capacity.desc')}</p>
                </div>
              </div>
              <div className="flex gap-2">
                <div className="w-1.5 h-1.5 rounded-full bg-primary mt-1.5 flex-shrink-0" />
                <div>
                  <p className="font-medium">{t('about.features.largeFile.name')}</p>
                  <p className="text-xs text-muted-foreground">{t('about.features.largeFile.desc')}</p>
                </div>
              </div>
            </div>
          </div>

          <div className="border-t" />

          <div className="space-y-4">
            <div className="space-y-1">
              <p className="text-xs font-medium text-muted-foreground uppercase tracking-wide">{t('about.buildInfo')}</p>
              <p className="text-sm">{appInfo.buildDate || '-'}</p>
            </div>

            <div className="space-y-1">
              <p className="text-xs font-medium text-muted-foreground uppercase tracking-wide">{t('about.buildHash')}</p>
              <p className="text-sm font-mono bg-muted/50 px-2 py-1 rounded inline-block">{appInfo.buildHash || '-'}</p>
            </div>

            <div className="space-y-1">
              <p className="text-xs font-medium text-muted-foreground uppercase tracking-wide">{t('about.author')}</p>
              <p className="text-sm font-semibold">{appInfo.author || '-'}</p>
            </div>

            <div className="space-y-1">
              <p className="text-xs font-medium text-muted-foreground uppercase tracking-wide">{t('about.blog')}</p>
              <a
                href="https://www.ktovoz.com"
                target="_blank"
                rel="noopener noreferrer"
                className="text-sm text-primary hover:underline flex items-center gap-1.5"
              >
                <Icons.Globe />
                <span>www.ktovoz.com</span>
              </a>
            </div>

            <div className="space-y-1">
              <p className="text-xs font-medium text-muted-foreground uppercase tracking-wide">{t('about.github')}</p>
              <a
                href={appInfo.github || '#'}
                target="_blank"
                rel="noopener noreferrer"
                className="text-sm text-primary hover:underline flex items-center gap-1.5"
              >
                <Icons.Github />
                <span>github.com/Ktovoz</span>
              </a>
            </div>
          </div>

          <div className="pt-2">
            <Alert>
              <AlertDescription className="text-center text-xs">
                {t('about.disclaimer')}
              </AlertDescription>
            </Alert>
          </div>
        </div>
      </CardContent>
    </Card>
  );
};

const HomePage = () => {
  const { t } = useI18n();
  return (
    <div className="h-full relative bg-transparent">
      <div className="relative h-full flex flex-col items-center justify-center p-10 -mt-10">
        <div className="mb-10 relative flex-shrink-0 animate-float">
          <div className="absolute inset-0 bg-primary/20 blur-3xl rounded-full opacity-70 animate-pulse-slow" />
          <div className="relative w-24 h-24 bg-gradient-to-br from-primary to-[hsl(var(--chart-2))] rounded-2xl flex items-center justify-center shadow-glow-lg transition-transform-gentle hover:scale-105">
            <svg xmlns="http://www.w3.org/2000/svg" width="48" height="48" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="1.5" strokeLinecap="round" strokeLinejoin="round" className="text-primary-foreground">
              <rect width="18" height="18" x="3" y="3" rx="2" ry="2" />
              <circle cx="9" cy="9" r="2" />
              <path d="m21 15-3.086-3.086a2 2 0 0 0-2.828 0L6 21" />
            </svg>
          </div>
        </div>

        <h1 className="text-6xl font-bold mb-6 text-foreground tracking-tight flex-shrink-0 leading-tight">
          <span className="bg-gradient-to-r from-primary via-primary/80 to-[hsl(var(--chart-2))] bg-clip-text text-transparent">
            {t('app.name')}
          </span>
        </h1>

        <p className="text-xl text-muted-foreground text-center max-w-md leading-relaxed flex-shrink-0">
          {t('home.title')}
          <br />
          <span className="text-base">{t('home.subtitle')}</span>
        </p>
      </div>
    </div>
  );
};

const ThemeToggle = () => {
  const { theme, setTheme } = useTheme();
  const [isDark, setIsDark] = React.useState(false);

  React.useEffect(() => {
    const isDarkTheme = theme === 'dark' || (theme === 'system' && window.matchMedia('(prefers-color-scheme: dark)').matches);
    setIsDark(isDarkTheme);
  }, [theme]);

  const handleToggle = () => {
    const newTheme = isDark ? 'light' : 'dark';
    setTheme(newTheme);
    logAction('ui', '切换主题', `切换到: ${newTheme}`);
  };

  return (
    <Button
      variant="ghost"
      size="icon"
      onClick={handleToggle}
      className="h-10 w-10 hover:bg-accent"
    >
      {isDark ? <Icons.Sun /> : <Icons.Moon />}
    </Button>
  );
};

const LanguageToggle = () => {
  const { locale, changeLocale } = useI18n();

  const languages = [
    { code: 'zh-CN', name: '中文', flag: '🇨🇳' },
    { code: 'en-US', name: 'English', flag: '🇺🇸' },
  ];

  const handleSelect = (langCode) => {
    changeLocale(langCode);
    logAction('ui', '切换语言', `切换到: ${langCode}`);
  };

  return (
    <DropdownMenu>
      <DropdownMenuTrigger asChild>
        <button
          className="h-10 w-10 flex items-center justify-center hover:bg-accent rounded transition-colors outline-none focus-visible:outline-none focus-visible:ring-0"
          title="Language / 语言"
        >
          <Icons.Globe />
        </button>
      </DropdownMenuTrigger>
      <DropdownMenuContent align="end" className="w-48">
        {languages.map((lang) => (
          <DropdownMenuItem
            key={lang.code}
            onClick={() => handleSelect(lang.code)}
            className={locale === lang.code ? 'bg-accent' : ''}
          >
            <span className="text-xl mr-2">{lang.flag}</span>
            <span className="text-base">{lang.name}</span>
          </DropdownMenuItem>
        ))}
      </DropdownMenuContent>
    </DropdownMenu>
  );
};

const WindowControls = () => {
  const handleMinimize = () => {
    try {
      WindowMinimise();
      logAction('ui', '最小化窗口', '');
    } catch (e) {
      console.error('Minimize failed:', e);
    }
  };

  const handleClose = () => {
    try {
      logAction('ui', '关闭应用程序', '');
      Quit();
    } catch (e) {
      console.error('Quit failed:', e);
    }
  };

  return (
    <div className="flex items-center gap-1" style={{ '--wails-draggable': 'no-drag' }}>
      <LanguageToggle />
      <ThemeToggle />
      <button
        onClick={handleMinimize}
        className="h-10 w-12 flex items-center justify-center hover:bg-accent hover:text-accent-foreground rounded transition-colors"
        title="最小化"
      >
        <Icons.Minus />
      </button>
      <button
        onClick={handleClose}
        className="h-10 w-12 flex items-center justify-center hover:bg-destructive hover:text-destructive-foreground rounded transition-colors"
        title="关闭"
      >
        <Icons.X />
      </button>
    </div>
  );
};

function AppContent() {
  const { t } = useI18n();
  const [appInfo, setAppInfo] = React.useState({ name: '', version: '', buildDate: '', author: '', repository: '', contact: '' });
  const [config, setConfig] = React.useState({});
  const [currentPage, setCurrentPage] = React.useState('home');

  React.useEffect(() => {
    const loadData = async () => {
      try {
        const info = await GetAppInfo();
        setAppInfo(info);
        const cfg = await GetConfig();
        setConfig(cfg);
      } catch (e) {
        console.error('Failed to load app info:', e);
      }
    };
    loadData();
  }, []);

  const handleConfigChange = (newConfig) => {
    setConfig(newConfig);
  };

  const renderPage = () => {
    switch (currentPage) {
      case 'home':
        return <HomePage />;
      case 'encrypt':
        return <EncryptPage config={config} />;
      case 'decrypt':
        return <DecryptPage config={config} />;
      case 'generate':
        return <GeneratePage config={config} />;
      case 'logs':
        return <LogsPage />;
      case 'settings':
        return <SettingsPage config={config} onConfigChange={handleConfigChange} />;
      case 'about':
        return <AboutPage appInfo={appInfo} />;
      default:
        return <HomePage />;
    }
  };

  return (
    <div className="h-screen w-screen flex flex-col bg-transparent overflow-hidden relative">
      <LiquidBackground />
      <div className="relative z-10 h-full w-full flex flex-col shell-surface">
        <header
          className="h-14 flex items-center justify-between px-4 select-none bg-transparent shell-divider-b"
          style={{ '--wails-draggable': 'drag' }}
        >
        <button
          onClick={() => {
            logAction('ui', '返回主页', '');
            setCurrentPage('home');
          }}
          className="flex items-center gap-2 font-semibold text-xl hover:text-primary transition-colors"
          style={{ '--wails-draggable': 'no-drag' }}
        >
          <Icons.Home />
          <span>{t('app.name')}</span>
        </button>
        <WindowControls />
      </header>

        <div className="flex-1 min-h-0 flex overflow-hidden">
          <aside className="w-52 min-h-0 flex flex-col bg-transparent shell-divider-r">
            <nav className="flex-1 min-h-0 p-3 space-y-1.5 overflow-auto">
            <MenuItem
              icon={Icons.Home}
              labelKey="home"
              active={currentPage === 'home'}
              onClick={() => setCurrentPage('home')}
            />
            <MenuItem
              icon={Icons.Lock}
              labelKey="encrypt"
              active={currentPage === 'encrypt'}
              onClick={() => setCurrentPage('encrypt')}
            />
            <MenuItem
              icon={Icons.Unlock}
              labelKey="decrypt"
              active={currentPage === 'decrypt'}
              onClick={() => setCurrentPage('decrypt')}
            />
            <MenuItem
              icon={Icons.Image}
              labelKey="generate"
              active={currentPage === 'generate'}
              onClick={() => setCurrentPage('generate')}
            />
            <MenuItem
              icon={Icons.FileText}
              labelKey="logs"
              active={currentPage === 'logs'}
              onClick={() => setCurrentPage('logs')}
            />
            <MenuItem
              icon={Icons.Settings}
              labelKey="settings"
              active={currentPage === 'settings'}
              onClick={() => setCurrentPage('settings')}
            />
            <MenuItem
              icon={Icons.Info}
              labelKey="about"
              active={currentPage === 'about'}
              onClick={() => setCurrentPage('about')}
            />
            </nav>
          </aside>

          <main className={`flex-1 min-h-0 overflow-hidden bg-transparent ${currentPage === 'home' ? 'p-0' : 'p-4'}`}>
            {currentPage === 'home' ? (
              renderPage()
            ) : (
              <div className="h-full rounded-3xl shell-frame p-0 overflow-hidden">
                {renderPage()}
              </div>
            )}
          </main>
        </div>
      </div>
    </div>
  );
}

function App() {
  return (
    <ThemeProvider defaultTheme="dark" storageKey="stego-ui-theme">
      <I18nProvider>
        <AppContent />
      </I18nProvider>
    </ThemeProvider>
  );
}

export default App;
