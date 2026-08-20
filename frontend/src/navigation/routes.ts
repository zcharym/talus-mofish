/** Management window route IDs. Domain routes use `{domain}.{page}` prefixes. */
export const ManagementRoute = {
  EnglishImport: 'english.import',
  EnglishReading: 'english.reading',
  EnglishVocabulary: 'english.vocabulary',
  ObsidianNotes: 'obsidian.notes',
  ObsidianSearch: 'obsidian.search',
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

export const OBSIDIAN_ROUTES: readonly ManagementRouteId[] = [
  ManagementRoute.ObsidianNotes,
  ManagementRoute.ObsidianSearch,
];

export function isEnglishLearningRoute(routeId: string): boolean {
  return (ENGLISH_LEARNING_ROUTES as readonly string[]).includes(routeId);
}

export function isObsidianRoute(routeId: string): boolean {
  return (OBSIDIAN_ROUTES as readonly string[]).includes(routeId);
}

export const PAGE_TITLES: Record<ManagementRouteId, string> = {
  [ManagementRoute.EnglishImport]: 'Import',
  [ManagementRoute.EnglishReading]: 'Reading',
  [ManagementRoute.EnglishVocabulary]: 'Vocabulary',
  [ManagementRoute.ObsidianNotes]: 'Notes',
  [ManagementRoute.ObsidianSearch]: 'Search',
  [ManagementRoute.Config]: 'Configuration',
  [ManagementRoute.Debug]: 'Debug',
  [ManagementRoute.About]: 'About',
};

export const DEFAULT_MANAGEMENT_ROUTE = ManagementRoute.EnglishVocabulary;
