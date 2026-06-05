import { createElement, type ReactNode } from 'react';
import {
  notifications,
  type NotificationData,
} from '@mantine/notifications';
import {
  IconAlertTriangle,
  IconCircleCheck,
  IconInfoCircle,
  IconX,
} from '@tabler/icons-react';

export type NotifyVariant = 'success' | 'warning' | 'failed' | 'info';

export type NotifyOptions = Partial<
  Omit<NotificationData, 'title' | 'message' | 'color' | 'icon'>
> & {
  title?: string;
  message?: ReactNode;
};

const variantConfig: Record<
  NotifyVariant,
  { color: string; icon: typeof IconCircleCheck }
> = {
  success: { color: 'green', icon: IconCircleCheck },
  warning: { color: 'yellow', icon: IconAlertTriangle },
  failed: { color: 'red', icon: IconX },
  info: { color: 'blue', icon: IconInfoCircle },
};

const defaultOptions: Partial<NotificationData> = {
  withBorder: true,
  radius: 'sm',
  autoClose: 5000,
};

function show(
  variant: NotifyVariant,
  title: string,
  message?: ReactNode,
  options?: NotifyOptions,
) {
  const { color, icon: Icon } = variantConfig[variant];

  return notifications.show({
    ...defaultOptions,
    title,
    message: message ?? '',
    color,
    icon: createElement(Icon, { size: 18 }),
    ...options,
  });
}

export const notify = {
  success(title: string, message?: ReactNode, options?: NotifyOptions) {
    return show('success', title, message, options);
  },

  warning(title: string, message?: ReactNode, options?: NotifyOptions) {
    return show('warning', title, message, options);
  },

  failed(title: string, message?: ReactNode, options?: NotifyOptions) {
    return show('failed', title, message, options);
  },

  info(title: string, message?: ReactNode, options?: NotifyOptions) {
    return show('info', title, message, options);
  },

  show(options: NotifyOptions & { title: string; message: ReactNode }) {
    return notifications.show({
      ...defaultOptions,
      ...options,
    });
  },

  hide(id: string) {
    return notifications.hide(id);
  },

  update(id: string, options: NotifyOptions & { message: ReactNode }) {
    return notifications.update({ id, ...defaultOptions, ...options });
  },

  clean() {
    notifications.clean();
  },
};
