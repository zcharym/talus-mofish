import { useEffect, useState } from 'react';
import { SystemService } from '../utils/api';

export function usePlatform(): string | null {
  const [platform, setPlatform] = useState<string | null>(null);

  useEffect(() => {
    SystemService.Platform()
      .then((value) => setPlatform(value || null))
      .catch(() => setPlatform(null));
  }, []);

  return platform;
}

export function isMacPlatform(platform: string | null): boolean {
  return platform === 'mac';
}
