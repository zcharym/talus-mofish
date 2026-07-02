/** Management window route IDs. Domain routes use `{domain}.{page}` prefixes. */
export const ManagementRoute = {
  EnglishImport: 'english.import',
  EnglishReading: 'english.reading',
  EnglishVocabulary: 'english.vocabulary',
  Config: 'config',
  Debug: 'debug',
  About: 'about',
} as const;

export type ManagementRouteId = (typeof ManagementRoute)[keyof typeof ManagementRoute];

export const ENGLISH_LEARNING_ROUTES: readonly ManagementRouteId[] = [
  ManagementRoute.EnglishImport,
  ManagementRoute.EnglishReading,
  ManagementRoute.EnglishVocabulary,
];

export function isEnglishLearningRoute(routeId: string): boolean {
  return (ENGLISH_LEARNING_ROUTES as readonly string[]).includes(routeId);
}

export const PAGE_TITLES: Record<ManagementRouteId, string> = {
  [ManagementRoute.EnglishImport]: 'Import',
  [ManagementRoute.EnglishReading]: 'Reading',
  [ManagementRoute.EnglishVocabulary]: 'Vocabulary',
  [ManagementRoute.Config]: 'Configuration',
  [ManagementRoute.Debug]: 'Debug',
  [ManagementRoute.About]: 'About',
};

export const DEFAULT_MANAGEMENT_ROUTE = ManagementRoute.EnglishVocabulary;
