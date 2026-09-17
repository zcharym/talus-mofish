export {
  AuthService,
  ChatService,
  CloudflareService,
  ConfigService,
  EnglishService,
  FeedsService,
  ObsidianService,
  SudokuService,
  SystemService,
} from '../../../bindings/github.com/songwei.ma/talus-mofish/backend/services';

export {
  App as AppConfig,
  Cloudflare,
  Feeds,
  OAuth,
  Obsidian,
  Sudoku,
  Vocabulary,
  VocabularyUpdate,
  Card,
  CardContentUpdate,
  Deck,
  Article,
  AnkiImportRecord,
  ArticlePageResult,
  ArticleSummary,
  VocabularyPageResult,
  ChatSession,
  ChatMessage,
  OverlayChatTurnRequest,
  StartChatTurnResult,
  SudokuGame,
  SudokuCheckResult,
  SudokuSession,
  ConfigChanged,
} from '../../../bindings/github.com/songwei.ma/talus-mofish/backend/types/models';

export { Config as AIConfig, Provider } from '../../../bindings/github.com/songwei.ma/talus-mofish/backend/utils/aiclient/models';

export { Note, SearchHit, FileEntry } from '../../../bindings/github.com/songwei.ma/talus-mofish/backend/obsidian/models';

export {
  AnkiDeckPreview,
  AnkiPreview,
  ImportDeckConfig,
  ImportResult,
} from '../../../bindings/github.com/songwei.ma/talus-mofish/backend/english/content/models';

export type { Dashboard as CloudflareDashboardSnapshot, Worker as CloudflareWorker } from '../../../bindings/github.com/songwei.ma/talus-mofish/backend/cloudflare/models';

export type FeedSource = {
  id: string;
  kind: string;
  title: string;
  url?: string;
  remoteId?: string;
  enabled?: boolean;
  lastError?: string;
  lastFetchedAt?: string;
};

export type FeedItem = {
  id: string;
  sourceId?: string;
  sourceKind?: string;
  sourceTitle?: string;
  title: string;
  url: string;
  author?: string;
  summary?: string;
  thumbnailUrl?: string;
  publishedAt?: string;
  saved?: boolean;
  read?: boolean;
};

export type FeedInbox = {
  sources?: FeedSource[];
  items?: FeedItem[];
  filter?: string;
  fetchedAt?: string;
};

export { toApiError, isRuntimeError } from './errors';
export {
  onConfigChanged,
  onAgentStreamChunk,
  onAgentTurnDone,
  onAgentTurnError,
  onAgentTurnCancelled,
  readEventData,
  ConfigChangedEvent,
} from './events';
export type {
  ConfigChangedPayload,
  AgentStreamChunkEvent,
  AgentTurnDoneEvent,
  AgentTurnErrorEvent,
  AgentTurnCancelledEvent,
} from './events';
