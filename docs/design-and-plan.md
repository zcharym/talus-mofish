# Talus Echo — English Learning App: Design & Plan

> **Project:** Talus Echo (aka Talus MoFish)  
> **Stack:** Wails v3 (Go 1.24+ / React 18 + TypeScript + Mantine) / SQLite (modernc.org/sqlite)  
> **Date:** 2026-06-05  

---

## Table of Contents

1. [Executive Summary](#1-executive-summary)
2. [Reference Analysis](#2-reference-analysis)
3. [Product Vision & Philosophy](#3-product-vision--philosophy)
4. [Architecture Overview](#4-architecture-overview)
5. [Core Modules](#5-core-modules)
6. [Spaced Repetition System (SRS) Design](#6-spaced-repetition-system-srs-design)
7. [Data Model](#7-data-model)
8. [UI/UX Design](#8-uiux-design)
9. [Implementation Phases](#9-implementation-phases)
10. [Technical Decisions](#10-technical-decisions)
11. [Open Questions](#11-open-questions)

---

## 1. Executive Summary

Talus Echo is a cross-platform desktop English learning application targeting intermediate Chinese-speaking programmers. It combines four proven approaches into one cohesive product:

| Source | Idea Absorbed |
|--------|---------------|
| **Echo-Loop** | Core SRS engine architecture, flashcard pipeline design |
| **LibreLingo** | Structured course hierarchy, multi-type challenges, gamification |
| **A Programmer's Guide to English** | First-principles methodology, "build a parser" mental model, corpus-first approach |
| **Read Frog** | Immersive reading, AI-powered contextual explanations, real-content learning |

The application is already scaffolded with Wails v3 + React (Mantine) + SQLite. The navbar defines the major feature areas: **Reading**, **Recite Words**, **Vocabulary**, **Listening**, **Grammar** — each currently a placeholder. This document turns those placeholders into a buildable roadmap.

---

## 2. Reference Analysis

### 2.1 Echo-Loop

> Repository: `github.com/echo-loop/Echo-Loop`

The original Echo-Loop project is the namesake and spiritual ancestor of this app. Its core contribution is a **pure-client SRS (Spaced Repetition System)** engine modeled on SM-2 (the algorithm behind Anki/SuperMemo), designed specifically for English-Chinese vocabulary acquisition. Key architectural ideas:

- **Queue-based review pipeline**: Each `due` flashcard enters a review queue; user ratings (`again`, `hard`, `good`, `easy`) determine next interval.
- **Inter-day & intra-day spacing**: New cards appear the same day, then graduated to longer intervals.
- **Leech detection**: Cards repeatedly failed are flagged and set aside to avoid clogging the queue.
- **JSON-driven card format** — front (Chinese prompt) / back (English word + example sentence + phonetic).

**Ideas adopted:**
- SM-2-based interval calculator as the SRS core (Go implementation)
- Card data structure with ease factor, interval, repetitions count
- Review queue scheduling with graduated intervals

### 2.2 LibreLingo

> Repository: `github.com/kantord/LibreLingo`

LibreLingo is a community-owned, open-source Duolingo alternative (AGPLv3, Svelte + PouchDB). Its teaching methodology provides:

- **Course → Module → Skill → Challenge** hierarchy — clean, composable curriculum design.
- **5 challenge types** (in difficulty order): Cards (image-word), Options (MCQ), Short Input (typing), Listening (audio), Chips (drag-to-order syntax).
- **Stale skill tracking**: Skills not practiced recently are visually flagged and pushed into the review queue.
- **Automatic challenge generation** from YAML-defined skill content — avoids hard-coding exercises.
- **Gamification**: Progress bars, streaks, XP-like tracking (though minimalist compared to Duolingo).
- **Cross-device sync** via PouchDB/CouchDB.

**Ideas adopted:**
- Course/module/skill/challenge hierarchy for structured learning paths
- Multi-type challenge system (not just flashcards)
- Stale skill detection → automatic review scheduling
- YAML/JSON-based course content (contributor-friendly)

### 2.3 A Programmer's Guide to English

> Repository: `github.com/yujiangshui/A-Programmers-Guide-to-English`

This is a methodology guide (16k+ stars) written by a Chinese programmer who went from CET-4 442 to IELTS 6.5. Its core philosophy:

- **Language is a skill, not a subject** — input-driven acquisition, not grammar memorization.
- **"Build an English-recognizing program"** — think of your brain as a parser: input processing → pattern matching → output generation.
- **Three pillar training:**
  1. **Corpus building** — systematic comprehensible input (reading + listening at i+1 level)
  2. **Pronunciation & listening** — IPA-based, not Chinese approximation; minimal pair drills
  3. **English thinking** — bypass native-language translation; think directly in English
- **No quick fixes** — the guide warns upfront of time commitment and rejects "magic method" marketing.
- **Emphasis on level-appropriate materials** — most advice is misleading because it comes from advanced learners.

**Ideas adopted:**
- Three-pillar curriculum design (corpus/listening/thinking)
- i+1 level recommendation engine for reading/listening content
- IPA-based pronunciation module with minimal pair drills
- "Build-a-parser" mental model as app's guiding metaphor
- Time-budgeted study plans (configurable daily goals)

### 2.4 Read Frog

> Repository: `github.com/mengxi-ream/read-frog`

Read Frog (陪读蛙) is an open-source AI browser extension (5,500+ stars, 20K weekly active users) for immersive language learning through reading:

- **Immersive bilingual overlay**: Translations appear inline preserving page layout — no context switching.
- **Context-aware AI translation**: Full-page context used for domain-appropriate translation (not word-by-word).
- **Selection translation + explanation**: Highlight any text → AI explanation + TTS + contextual word meaning.
- **Level-adaptive**: Explanations adjust depth to learner proficiency.
- **Multi-provider AI**: 20+ providers (OpenAI, Claude, DeepSeek, local Ollama, etc.) — user brings their own key.
- **Cost optimization**: Batch requests group many translations into one API call, saving up to 70%.

**Ideas adopted:**
- Immersive reading mode with AI-powered contextual explanations
- Bring-your-own-API-key model for AI features
- Selection-based lookup with TTS and level-adaptive detail
- Batch translation/monolingual explanation to manage API costs
- YouTube subtitle translation for listening practice

---

## 3. Product Vision & Philosophy

### 3.1 Vision Statement

> *"Talus Echo is a programmer's English learning workbench — combining the rigor of spaced repetition with the richness of real-world content, powered by AI that adapts to your level."*

### 3.2 Design Principles

1. **Input-first, output-second** — Prioritize comprehensible input (reading + listening). Speaking/writing come after sufficient corpus.
2. **SRS as the memory backbone** — Everything that needs memorization (words, phrases, grammar patterns) goes through the spaced repetition engine.
3. **Real content over contrived exercises** — Where possible, use real articles, YouTube videos, and code documentation rather than artificial sentences.
4. **Level-adaptive AI** — AI features (explanations, translations, sentence generation) adapt depth to the learner's current level.
5. **Contributor-friendly content** — Course content lives in plain YAML/JSON files; community contributions should be straightforward.
6. **Offline-first, sync-optional** — Core SRS works fully offline. AI features and content downloads are additive.
7. **For programmers, by programmers** — Mental models (parsers, queues, intervals) resonate with the target audience. The app should feel like a well-engineered tool, not a gamified toy.

### 3.3 Target Persona

- **Chinese software developer** (native Mandarin speaker)
- Current English level: CET-4 to CET-6 (IELTS 4.5–6.0)
- Familiar with CLI tools, understands algorithmic concepts
- Has 30–60 min/day to study, prefers desktop (Mac/Windows)
- Values efficiency and measurable progress over gamification fluff

---

## 4. Architecture Overview

```
┌─────────────────────────────────────────────────────────────┐
│                   React + Mantine UI                        │
│  ┌──────┐ ┌──────┐ ┌──────┐ ┌────────┐ ┌────────┐ ┌──────┐ │
│  │Reading│ │Recite│ │Vocab │ │Listening│ │Grammar │ │Config│ │
│  └──┬───┘ └──┬───┘ └──┬───┘ └───┬────┘ └───┬────┘ └──┬───┘ │
│     └────────┴────────┴─────────┴──────────┴─────────┘      │
│                         │  Wails Bindings                    │
├─────────────────────────┼───────────────────────────────────┤
│              Go Backend (Wails v3 Services)                  │
│  ┌───────────────────────────────────────────────────────┐   │
│  │                    AppService                         │   │
│  │  ┌──────────┐ ┌────────────┐ ┌──────────────────┐    │   │
│  │  │ SRS Engine│ │ AI Service │ │ Content Manager  │    │   │
│  │  │ (SM-2)   │ │ (OpenAI/   │ │ (Course/Lesson   │    │   │
│  │  │          │ │  Claude/   │ │  Loader)         │    │   │
│  │  │          │ │  Ollama)   │ │                  │    │   │
│  │  └──────────┘ └────────────┘ └──────────────────┘    │   │
│  │  ┌──────────┐ ┌────────────┐ ┌──────────────────┐    │   │
│  │  │ Config   │ │ Stats/Track│ │ Import/Export    │    │   │
│  │  │ Service  │ │ Service    │ │ Service          │    │   │
│  │  └──────────┘ └────────────┘ └──────────────────┘    │   │
│  └───────────────────────────────────────────────────────┘   │
│                         │                                     │
│              sqlc-generated store layer                       │
├───────────────────────────────────────────────────────────────┤
│                      SQLite (modernc.org/sqlite)              │
│  ┌────────┐ ┌──────────┐ ┌────────┐ ┌────────┐ ┌─────────┐  │
│  │settings│ │cards     │ │reviews │ │courses │ │vocabulary│  │
│  └────────┘ └──────────┘ └────────┘ └────────┘ └─────────┘  │
│  ┌────────┐ ┌──────────┐ ┌────────┐ ┌─────────────────────┐  │
│  │stats   │ │articles  │ │streaks │ │...                  │  │
│  └────────┘ └──────────┘ └────────┘ └─────────────────────┘  │
└───────────────────────────────────────────────────────────────┘
```

### 4.1 Key Technology Decisions

| Layer | Technology | Rationale |
|-------|-----------|-----------|
| Desktop shell | Wails v3 | Native desktop performance, Go backend, Web UI |
| UI framework | React 18 + Mantine 7 | Mature ecosystem, excellent component library |
| Database | SQLite (modernc.org/sqlite) | Zero CGO, embedded, portable; already in project |
| Query layer | sqlc | Type-safe Go codegen from SQL; already configured |
| AI API | OpenAI-compatible HTTP client | Multi-provider (OpenAI, Claude, DeepSeek, Ollama) |
| Storage | Local filesystem | Config JSON + SQLite DB at OS user data dir |

---

## 5. Core Modules

### 5.1 SRS Engine (`internal/srs/`)

The heart of the application. See [Section 6](#6-spaced-repetition-system-srs-design) for detailed algorithm design.

**Key components:**
- `card.go` — Card data structure (front, back, ease, interval, repetitions, due date, lapsed status)
- `scheduler.go` — Review queue: which cards are due now, new card introduction rate
- `calculator.go` — SM-2 interval calculation with configurable parameters
- `leech.go` — Leech detection and handling

### 5.2 AI Service (`internal/ai/`)

A multi-provider AI client for contextual explanations, translation, and content generation.

**Key components:**
- `client.go` — HTTP client with provider abstraction (OpenAI, Claude, DeepSeek, Ollama, etc.)
- `explanations.go` — Word/phrase explanation generator (level-adaptive)
- `translator.go` — Context-aware translation (batched for cost efficiency)
- `generator.go` — Sentence/example generation for vocabulary items
- `config.go` — Provider selection, API key management, model params

### 5.3 Content Manager (`internal/content/`)

Manages course content, reading materials, and vocabulary lists.

**Key components:**
- `course.go` — Course/module/skill/challenge data structures and loading
- `reader.go` — Article/reading content management (text extraction, metadata)
- `importer.go` — Import Anki decks, JSON course packs, OPML podcast lists
- `yaml_loader.go` — Parse YAML/JSON course definitions into DB records

### 5.4 Stats & Tracking (`internal/stats/`)

Tracks learning activity, streaks, and progress.

**Key components:**
- `tracker.go` — Session tracking (start/end, cards reviewed, time spent)
- `streaks.go` — Daily streak calculation, freeze logic
- `reports.go` — Review history, retention rates, forecast

### 5.5 Config Service (`internal/config/`)

**Already implemented** — persists theme, daily goal, words per session. Will be extended with AI provider settings, audio preferences, and course selection.

---

## 6. Spaced Repetition System (SRS) Design

### 6.1 Algorithm: Modified SM-2

We implement a variant of the SuperMemo SM-2 algorithm that powers Anki. All user-modifiable parameters are stored in the settings table.

#### Card Data Model

```
Card {
  id             TEXT PRIMARY KEY    -- UUID
  deck_id        TEXT NOT NULL       -- FK -> decks.id
  front          TEXT NOT NULL       -- Prompt (Chinese word/phrase)
  back           TEXT NOT NULL       -- Answer (English word + phonetic + example)
  example_sentence TEXT              -- Context sentence
  hints          TEXT                -- Additional hints (image URL, mnemonic)
  card_type      TEXT NOT NULL       -- 'vocabulary' | 'sentence' | 'cloze'
  created_at     TEXT NOT NULL       -- ISO 8601
  updated_at     TEXT NOT NULL
  
  -- SRS state (non-null for graduated cards)
  ease           REAL DEFAULT 2.5    -- Ease factor (minimum 1.3)
  interval       INTEGER DEFAULT 0   -- Current interval in days
  repetitions    INTEGER DEFAULT 0   -- Consecutive correct answers
  due_at         TEXT                -- Next review datetime (ISO 8601)
  lapsed         INTEGER DEFAULT 0   -- Leech counter
  suspended      INTEGER DEFAULT 0   -- Manual suspend flag
}
```

#### Review Rating Scale

| Rating | Name | Behavior |
|--------|------|----------|
| 1 | Again | Reset to learning phase; interval = 0 |
| 2 | Hard | Interval * 1.2; ease -= 0.15 |
| 3 | Good | Interval * ease; ease unchanged |
| 4 | Easy | Interval * ease * 1.3; ease += 0.15 |

#### Learning Phases

```
New Card
   │
   ▼
┌─────────────────────────────────────┐
│ Learning Phase (intra-day)          │
│                                     │
│ Step 1: Again → back to step 1      │
│         Good → Step 2 (10 min)      │
│         Easy → Graduated (1 day)    │
│                                     │
│ Step 2: Again → back to step 1      │
│         Good → Graduated (1 day)    │
│         Easy → Graduated (3 days)   │
└─────────────────────────────────────┘
   │ (passed step 2)
   ▼
┌─────────────────────────────────────┐
│ Graduated (Reviewing Phase)         │
│                                     │
│ Again → back to Learning Phase      │
│         (interval reset to 0)       │
│ Hard  → interval * HARD_MULT       │
│ Good  → interval * ease             │
│ Easy  → interval * ease * EASY_BONUS│
└─────────────────────────────────────┘
```

#### Interval Calculation (Go pseudo-code)

```go
func NextInterval(card *Card, rating Rating) time.Duration {
    switch {
    case card.Interval == 0:
        // Learning phase
        switch rating {
        case Again: return 0           // repeat step
        case Good:  return 10 * minute // step 2: 10 min
        case Easy:  return 1 * day     // graduate to reviewing
        default:    return 10 * minute
        }
    case card.Interval < 1*day:
        // Second learning step
        switch rating {
        case Again: return 0              // back to step 1
        case Good:  return 1 * day        // graduate
        case Easy:  return 3 * day        // graduate, bonus
        default:    return 1 * day
        }
    default:
        // Reviewing phase
        switch rating {
        case Again:
            card.Ease = max(1.3, card.Ease - 0.20)
            return 0  // back to learning
        case Hard:
            card.Ease = max(1.3, card.Ease - 0.15)
            return time.Duration(float64(card.Interval) * 1.2) * day
        case Good:
            return time.Duration(float64(card.Interval) * card.Ease) * day
        case Easy:
            card.Ease = min(10.0, card.Ease + 0.15)
            return time.Duration(float64(card.Interval) * card.Ease * 1.3) * day
        }
    }
}
```

#### Leech Detection

A card is flagged as a potential **leech** when:
- It has been lapsed (rated Again after graduation) ≥ `leech_threshold` times (default: 8)
- Or its average interval over the last 5 reviews is < 21 days with ≥ 3 lapses

Leeched cards are:
1. Moved to a separate "Leech" deck/filter
2. Optionally suspended (user configurable)
3. Logged for review: the user should add mnemonics or break the card down

#### Daily Workload & New Card Limit

- **New cards per day**: configurable (default: 20, matches `WordsPerSession`)
- **Review cap**: configurable (default: 100 reviews/day)
- **Fuzz**: Intervals are jittered ±25% to avoid card clumping (Anki-compatible)
- **Graduating interval**: default `[1, 4]` days (fuzzed)

### 6.2 Deck Organization

```
All Decks
├── Vocabulary (default)
│   ├── CET-4
│   ├── CET-6
│   ├── IELTS
│   ├── TOEFL
│   └── GRE
├── Courses
│   ├── {Course Name}
│   │   └── {Module} → {Skill} → cards
├── Reading Clippings  (auto-saved from Reading module)
├── Errors             (auto-saved incorrect answers)
├── Leech              (leeched cards)
└── Custom             (user-created decks)
```

### 6.3 Review Session Flow

```
1. User opens "Recite Words" tab
2. SRS scheduler queries due cards:
   - Cards where due_at <= now
   - Limit by WordsPerSession config
   - Prioritize: lapsed > learning phase > reviewing (oldest due first)
3. Frontend displays card front
4. User self-assesses → tap to reveal back
5. User rates: Again / Hard / Good / Easy
6. Backend updates card SRS state + records review log
7. Advance to next card
8. Session ends when:
   - All due cards reviewed
   - User closes tab
   - Daily goal reached
```

### 6.4 Review Logging

Every review is recorded for analytics and spaced-repetition fine-tuning:

```sql
CREATE TABLE review_log (
    id          TEXT PRIMARY KEY,
    card_id     TEXT NOT NULL REFERENCES cards(id),
    rating      INTEGER NOT NULL,        -- 1-4
    elapsed_ms  INTEGER NOT NULL,        -- time spent on this card
    review_at   TEXT NOT NULL,            -- ISO 8601
    day_seq     INTEGER NOT NULL,         -- session sequence number
    ease_before REAL,                     -- ease before this review
    interval_before INTEGER               -- interval before this review
);
```

---

## 7. Data Model

### 7.1 Complete SQLite Schema (Additions to Existing Schema)

```sql
-- Existing: settings table

-- === SRS Engine ===
CREATE TABLE decks (
    id          TEXT PRIMARY KEY,
    name        TEXT NOT NULL,
    description TEXT DEFAULT '',
    parent_id   TEXT REFERENCES decks(id),
    sort_order  INTEGER DEFAULT 0,
    card_count  INTEGER DEFAULT 0,
    created_at  TEXT NOT NULL DEFAULT (datetime('now')),
    updated_at  TEXT NOT NULL DEFAULT (datetime('now'))
);

CREATE TABLE cards (
    id              TEXT PRIMARY KEY,
    deck_id         TEXT NOT NULL REFERENCES decks(id),
    front           TEXT NOT NULL,
    back            TEXT NOT NULL,
    example_sentence TEXT DEFAULT '',
    hints           TEXT DEFAULT '',
    card_type       TEXT NOT NULL DEFAULT 'vocabulary'
                    CHECK (card_type IN ('vocabulary', 'sentence', 'cloze', 'listening')),
    -- SRS state
    ease            REAL DEFAULT 2.5,
    interval        INTEGER DEFAULT 0,   -- days
    repetitions     INTEGER DEFAULT 0,
    due_at          TEXT,                 -- next review datetime
    lapsed          INTEGER DEFAULT 0,
    suspended       INTEGER DEFAULT 0,
    -- Metadata
    source          TEXT DEFAULT '',      -- e.g. "course:basics", "reading:article123", "import:anki"
    created_at      TEXT NOT NULL DEFAULT (datetime('now')),
    updated_at      TEXT NOT NULL DEFAULT (datetime('now'))
);

CREATE TABLE review_log (
    id              TEXT PRIMARY KEY,
    card_id         TEXT NOT NULL REFERENCES cards(id),
    rating          INTEGER NOT NULL CHECK (rating BETWEEN 1 AND 4),
    elapsed_ms      INTEGER NOT NULL,
    review_at       TEXT NOT NULL DEFAULT (datetime('now')),
    day_seq         INTEGER NOT NULL,
    ease_before     REAL,
    interval_before INTEGER
);

CREATE INDEX idx_cards_due ON cards(deck_id, due_at) WHERE suspended = 0;
CREATE INDEX idx_review_log_card ON review_log(card_id);
CREATE INDEX idx_review_log_date ON review_log(review_at);

-- === Course Content ===
CREATE TABLE courses (
    id          TEXT PRIMARY KEY,
    title       TEXT NOT NULL,
    description TEXT DEFAULT '',
    language    TEXT NOT NULL DEFAULT 'en',
    source_lang TEXT NOT NULL DEFAULT 'zh',
    level       TEXT DEFAULT 'beginner'
                CHECK (level IN ('beginner', 'elementary', 'intermediate', 'advanced')),
    version     TEXT DEFAULT '1.0',
    author      TEXT DEFAULT '',
    created_at  TEXT NOT NULL DEFAULT (datetime('now'))
);

CREATE TABLE modules (
    id          TEXT PRIMARY KEY,
    course_id   TEXT NOT NULL REFERENCES courses(id),
    title       TEXT NOT NULL,
    description TEXT DEFAULT '',
    sort_order  INTEGER DEFAULT 0
);

CREATE TABLE skills (
    id          TEXT PRIMARY KEY,
    module_id   TEXT NOT NULL REFERENCES modules(id),
    title       TEXT NOT NULL,
    description TEXT DEFAULT '',
    sort_order  INTEGER DEFAULT 0
);

CREATE TABLE challenges (
    id          TEXT PRIMARY KEY,
    skill_id    TEXT NOT NULL REFERENCES skills(id),
    challenge_type TEXT NOT NULL
                  CHECK (challenge_type IN ('cards', 'options', 'short_input', 'listening', 'chips')),
    prompt      TEXT NOT NULL,
    answer      TEXT NOT NULL,
    alternatives TEXT DEFAULT '',       -- JSON array of acceptable answers
    audio_ref   TEXT DEFAULT '',        -- path to audio file
    image_ref   TEXT DEFAULT '',        -- path to image asset
    sort_order  INTEGER DEFAULT 0
);

-- === Vocabulary ===
CREATE TABLE vocabulary (
    id          TEXT PRIMARY KEY,
    word        TEXT NOT NULL,
    phonetic    TEXT DEFAULT '',
    pos         TEXT DEFAULT '',        -- part of speech
    definition  TEXT NOT NULL,          -- Chinese definition
    definition_en TEXT DEFAULT '',       -- English definition
    examples    TEXT DEFAULT '',        -- JSON array of example sentences
    level       TEXT DEFAULT ''
                CHECK (level IN ('cet4', 'cet6', 'ielts', 'toefl', 'gre', 'custom')),
    created_at  TEXT NOT NULL DEFAULT (datetime('now'))
);

-- Many-to-many: cards can reference vocabulary entries
CREATE TABLE card_vocab (
    card_id     TEXT NOT NULL REFERENCES cards(id),
    vocab_id    TEXT NOT NULL REFERENCES vocabulary(id),
    PRIMARY KEY (card_id, vocab_id)
);

-- === Reading / Articles ===
CREATE TABLE articles (
    id          TEXT PRIMARY KEY,
    title       TEXT NOT NULL,
    url         TEXT DEFAULT '',
    content     TEXT NOT NULL,           -- original text
    translation TEXT DEFAULT '',         -- AI-translated (cached)
    level       TEXT DEFAULT ''
                CHECK (level IN ('beginner', 'elementary', 'intermediate', 'advanced')),
    source      TEXT DEFAULT 'manual',   -- 'manual', 'clipboard', 'web', 'youtube'
    word_count  INTEGER DEFAULT 0,
    ai_provider TEXT DEFAULT '',
    created_at  TEXT NOT NULL DEFAULT (datetime('now')),
    updated_at  TEXT NOT NULL DEFAULT (datetime('now'))
);

-- Clippings: words/phrases saved during reading
CREATE TABLE clippings (
    id          TEXT PRIMARY KEY,
    article_id  TEXT REFERENCES articles(id),
    text        TEXT NOT NULL,           -- selected text (English)
    translation TEXT DEFAULT '',         -- AI contextual translation
    explanation TEXT DEFAULT '',         -- AI detailed explanation
    note        TEXT DEFAULT '',         -- user note
    saved_to_card INTEGER DEFAULT 0,     -- 1 = added to SRS queue
    created_at  TEXT NOT NULL DEFAULT (datetime('now'))
);

-- === Learning Stats ===
CREATE TABLE daily_stats (
    date        TEXT PRIMARY KEY,        -- 'YYYY-MM-DD'
    study_seconds INTEGER DEFAULT 0,
    cards_reviewed INTEGER DEFAULT 0,
    new_cards_learned INTEGER DEFAULT 0,
    correct_count   INTEGER DEFAULT 0,
    wrong_count     INTEGER DEFAULT 0,
    streak_day  INTEGER DEFAULT 1,
    goal_reached INTEGER DEFAULT 0
);

CREATE TABLE sessions (
    id              TEXT PRIMARY KEY,
    started_at      TEXT NOT NULL,
    ended_at        TEXT,
    cards_reviewed  INTEGER DEFAULT 0,
    new_cards       INTEGER DEFAULT 0,
    correct_count   INTEGER DEFAULT 0,
    wrong_count     INTEGER DEFAULT 0
);
```

### 7.2 Relationship Diagram

```
courses 1──N modules 1──N skills 1──N challenges
                                                  (cards reference skills via source field)
decks 1──N cards 1──N review_log
cards N──M vocabulary
articles 1──N clippings
```

---

## 8. UI/UX Design

### 8.1 Navigation Structure (Existing + Proposed)

The navbar already has two sections: **Tools** and **English Study**. English Study gets the following sub-pages:

| Tab | Page | Priority | Description |
|-----|------|----------|-------------|
| 📖 Reading | `ReadingPage` | P0 | Immersive reading with AI contextual lookup |
| 🧠 Recite Words | `RecitePage` | P0 | SRS flashcard review |
| 📚 Vocabulary | `VocabPage` | P1 | Browsable word bank with search/filter |
| 🎧 Listening | `ListeningPage` | P1 | Audio content + SRS for listening cards |
| ✏️ Grammar | `GrammarPage` | P2 | Grammar challenges from course content |

### 8.2 Page Designs

#### 8.2.1 Recite Words (P0)

```
┌─────────────────────────────────────────────┐
│  ◀ Back to Dashboard    Recite Words  🔄 20 │
│  ─────────────────────────────────────────── │
│                                              │
│  ┌──────────────────────────────────────┐    │
│  │                                      │    │
│  │           (Card Area)                │    │
│  │                                      │    │
│  │      ▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓       │    │
│  │      ▓     "出类拔萃"        ▓       │    │
│  │      ▓                        ▓       │    │
│  │      ▓   (Tap to reveal)      ▓       │    │
│  │      ▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓       │    │
│  │                                      │    │
│  │  ─── After reveal ───               │    │
│  │                                      │    │
│  │      "stand out from the crowd"      │    │
│  │      /stænd aʊt frɒm ðə kraʊd/       │    │
│  │      "His talent made him stand      │    │
│  │       out from the crowd."           │    │
│  │                                      │    │
│  └──────────────────────────────────────┘    │
│                                              │
│  ┌──────┐ ┌──────┐ ┌──────┐ ┌──────┐        │
│  │ Again │ │ Hard │ │ Good │ │ Easy │        │
│  │   ⏳  │ │  😅  │ │  ✅  │ │  🚀  │        │
│  └──────┘ └──────┘ └──────┘ └──────┘        │
│                                              │
│  Session progress: ████████░░░░ 8/20        │
│  Today: 45 cards reviewed                    │
└─────────────────────────────────────────────┘
```

**States:**
- **Empty state**: "No cards due! Add words from Reading or import a deck. 🎉"
- **New card intro**: First reveal shows phonetic + definition; second tap shows example sentence
- **Card with audio**: 🔊 icon, auto-plays TTS on reveal
- **Session complete**: Confetti animation + stats summary + "Nice work! Come back later for more."

#### 8.2.2 Reading (P0)

```
┌──────────────────────────────────────────────┐
│  ◀ Back          Reading                     │
│  ──────────────────────────────────────────── │
│  [📄 New Article] [📋 From Clipboard] [🔗 URL]│
│  ──────────────────────────────────────────── │
│                                               │
│  ┌───────────────────────────────────────┐    │
│  │ Title: How Go Handles Memory          │    │
│  │ Level: Intermediate  Words: 843       │    │
│  │ ───────────────────────────────────── │    │
│  │                                       │    │
│  │ Go's memory model is designed for     │    │
│  │ ═══════════════════════════════       │    │
│  │ concurrent access. It uses...         │    │
│  │                                       │    │
│  │ (click any word → shows popover)      │    │
│  │ ┌─────────────────────┐               │    │
│  │ │ concurrent          │               │    │
│  │ │ /kənˈkʌrənt/  adj.  │               │    │
│  │ │ 并发的, 同时发生的   │               │    │
│  │ │ "concurrent access" │               │    │
│  │ │ → 并发访问          │               │    │
│  │ │ [💾 Save to SRS]  🔊│               │    │
│  │ └─────────────────────┘               │    │
│  │                                       │    │
│  │ (tap sentence → bottom panel)         │    │
│  │ ──────────────────────────────────    │    │
│  │ 📝 "Go's memory model is designed     │    │
│  │      for concurrent access."          │    │
│  │ 🔄 并发访问 — explains memory access  │    │
│  │    pattern where multiple goroutines  │    │
│  │    read/write shared data.            │    │
│  │ [🗣 Read Aloud] [📌 Save Sentence]    │    │
│  └───────────────────────────────────────┘    │
└──────────────────────────────────────────────┘
```

**States:**
- **Empty state**: "No articles yet. Paste a URL, import from clipboard, or write something."
- **Loading**: Skeleton placeholder for article content
- **No AI key configured**: Show translation without AI detail; prompt to add API key in Config
- **Word lookup**: Popover with phonetic, definition, example, AI explanation, save-to-SRS button

#### 8.2.3 Vocabulary (P1)

```
┌──────────────────────────────────────────────┐
│  ◀ Back              Vocabulary Bank         │
│  ──────────────────────────────────────────── │
│  [🔍 Search...] [Filter: All ▼] [Sort: A-Z ▼]│
│  ──────────────────────────────────────────── │
│                                              │
│  Level chips: [CET-4] [CET-6] [IELTS] [All] │
│                                              │
│  ┌───────────────────────────────────────┐   │
│  │ stand out         /stænd aʊt/  📌     │   │
│  │ 出类拔萃, 突出 — v. phrase           │   │
│  │ [🔊] [📝 Add to SRS] [✏️ Edit]       │   │
│  ├───────────────────────────────────────┤   │
│  │ concurrent        /kənˈkʌrənt/  📌   │   │
│  │ 并发的 — adj.                        │   │
│  │ [🔊] [📝 Add to SRS] [✏️ Edit]       │   │
│  └───────────────────────────────────────┘   │
│                                              │
│  Showing 25 of 342 words                      │
└──────────────────────────────────────────────┘
```

**States:**
- **Empty state**: "Your vocabulary bank is empty. Words are added when you look them up in Reading or import word lists."
- **Filtered empty**: "No words match this filter. Try a different level or search term."
- **Loading**: Skeleton list items

#### 8.2.4 Listening (P1)

```
┌──────────────────────────────────────────────┐
│  ◀ Back              Listening               │
│  ──────────────────────────────────────────── │
│  [🎵 YouTube URL] [📁 Import Audio]          │
│  ──────────────────────────────────────────── │
│                                              │
│  ┌───────────────────────────────────────┐   │
│  │ Now Playing: Go Conference 2025 Keynote│   │
│  │ Level: Intermediate  Duration: 15:32  │   │
│  │ ───────────────────────────────────── │   │
│  │ ████████████░░░░░░░░░░ 00:34 / 15:32 │   │
│  │ [⏪] [▶/⏸] [⏩] [🔁 Loop] [⏬ 1.0x]  │   │
│  │ ───────────────────────────────────── │   │
│  │ 📝 Transcript (bilingual):            │   │
│  │ ┌────────────────────────────────┐    │   │
│  │ │ So we need to think about     │    │   │
│  │ │ [concurrent access] patterns  │    │   │
│  │ │ ───────────────────────────── │    │   │
│  │ │ 所以我们需要考虑并发访问的模式│    │   │
│  │ └────────────────────────────────┘    │   │
│  │ [📌 concurrent access → SRS]          │   │
│  └───────────────────────────────────────┘   │
└──────────────────────────────────────────────┘
```

- Audio player with speed control, looping, transcript
- YouTube subtitle integration (Read Frog pattern)
- Clickable transcript words → lookup popover (same as Reading)

#### 8.2.5 Grammar (P2)

Course-based grammar challenges using the LibreLingo challenge types: cards, options, short input, chips (drag-to-order syntax).

### 8.3 Theme & Visual Design

Already using Mantine with auto/light/dark theme support. Additional considerations:

- **Card design**: Clean, minimal, large typography. Card should feel like a physical flashcard.
- **Progress indicators**: Circular progress for daily goal, bar progress for session
- **Animations**: Card flip animation (CSS transform), subtle confetti on session complete
- **Responsive**: While primarily desktop, window should resize gracefully

---

## 9. Implementation Phases

### Phase 1: SRS Foundation (Weeks 1–2) — P0

**Goal**: A working flashcard review loop with SM-2 scheduling.

| Task | Files | Description |
|------|-------|-------------|
| 1.1 SRS data model | `internal/database/schema.sql`, `db/queries/srs.sql` | decks, cards, review_log tables and queries |
| 1.2 SRS engine (Go) | `internal/srs/` | Card, Scheduler, Calculator, Leech detector |
| 1.3 SRS service | `appservice.go` | Expose SRS operations to frontend via Wails bindings |
| 1.4 Recite page (frontend) | `frontend/src/pages/RecitePage.tsx` | Card display, tap-to-reveal, rating buttons |
| 1.5 Seed data & import | `internal/content/` | Built-in CET-4/CET-6 word lists; Anki APKG import |
| 1.6 Basic stats tracking | `internal/stats/` | Session logging, daily goal tracking |

**Deliverable**: User can open the app, see due cards, review them with Again/Hard/Good/Easy, see next interval.

### Phase 2: Reading with AI (Weeks 3–4) — P0

**Goal**: Immersive reading with AI contextual lookup and save-to-SRS.

| Task | Files | Description |
|------|-------|-------------|
| 2.1 AI service (Go) | `internal/ai/` | Multi-provider HTTP client, batch translation |
| 2.2 Article management | `internal/database/schema.sql`, `db/queries/articles.sql` | articles, clippings tables |
| 2.3 Reading page (frontend) | `frontend/src/pages/ReadingPage.tsx` | Article display, word click popover |
| 2.4 AI explanation popover | `frontend/src/components/WordPopover.tsx` | Popover with AI explanation |
| 2.5 Save-to-SRS flow | `appservice.go` | Clip word → create card → add to review queue |
| 2.6 Clipboard/URL import | `internal/content/reader.go` | Import from clipboard or URL |

**Deliverable**: User can paste an article, click words for AI explanations, save them to SRS deck.

### Phase 3: Vocabulary & Course System (Weeks 5–6) — P1

**Goal**: Browsable vocabulary bank, structured course content, multi-challenge types.

| Task | Files | Description |
|------|-------|-------------|
| 3.1 Vocabulary page | `frontend/src/pages/VocabPage.tsx` | Search, filter, sort, edit |
| 3.2 Course loader | `internal/content/course.go` | YAML/JSON → DB pipeline |
| 3.3 Challenge types | `frontend/src/components/challenges/` | Cards, Options, ShortInput, Chips components |
| 3.4 Listening page | `frontend/src/pages/ListeningPage.tsx` | YouTube integration, transcript display |
| 3.5 Course browser | `frontend/src/pages/CoursesPage.tsx` | Browse/select courses and modules |

**Deliverable**: Course-based learning with multiple challenge types, vocabulary browsing.

### Phase 4: Analytics & Gamification (Weeks 7–8) — P1

**Goal**: Streaks, statistics, retention graphs, study planning.

| Task | Files | Description |
|------|-------|-------------|
| 4.1 Daily stats engine | `internal/stats/tracker.go` | Streak calculation, goal tracking |
| 4.2 Stats dashboard | `frontend/src/pages/StatsPage.tsx` | Charts (review count, retention %, time) |
| 4.3 Retention forecast | `internal/stats/forecast.go` | Predict future review load |
| 4.4 Gamification | `frontend/src/components/` | Streak display, session complete celebration |
| 4.5 Study scheduler | `internal/srs/scheduler.go` | Smart daily scheduling based on workload forecast |

**Deliverable**: Full analytics dashboard, streak tracking, study planning.

### Phase 5: Polish & Extensions (Weeks 9–10) — P2

**Goal**: Grammar module, import/export, performance optimization.

| Task | Files | Description |
|------|-------|-------------|
| 5.1 Grammar page | `frontend/src/pages/GrammarPage.tsx` | Grammar challenges, rule cards |
| 5.2 Anki import/export | `internal/content/importer.go` | APKG format support |
| 5.3 Backup & sync | `internal/database/` | DB export, cloud sync (optional, future) |
| 5.4 Performance | All | Card review optimization, batch queries |
| 5.5 Onboarding | `frontend/src/pages/Onboarding.tsx` | First-run wizard |

**Deliverable**: Feature-complete v1.0.

---

## 10. Technical Decisions

### 10.1 Why SQLite for SRS State

- **Offline-first**: SRS is useless without local data access. SQLite works offline by definition.
- **Single-user**: This is a desktop app for one person. No concurrency, no server.
- **Already in project**: The project already uses modernc.org/sqlite + sqlc. Adding tables is zero infrastructural cost.
- **Portable**: The entire database is one file. Backup = copy file.

### 10.2 Why Not PouchDB / CouchDB (LibreLingo approach)

LibreLingo uses PouchDB for cross-device sync. Talus Echo starts as desktop-only. If sync is desired later, SQLite → some sync layer (e.g., SQLite replication, or export/import) is more appropriate for a Wails app than embedding PouchDB.

### 10.3 AI Provider Abstraction

The AI service uses an **OpenAI-compatible HTTP client** as the common interface. This means any provider with an OpenAI-compatible API (OpenAI, Claude via Anthropic API, DeepSeek, Ollama, Groq, etc.) works with the same client — only the base URL and API key change.

```go
type AIClient struct {
    BaseURL   string
    APIKey    string
    Model     string
    HTTPClient *http.Client
}

func (c *AIClient) Explain(ctx context.Context, text string, level string) (*Explanation, error)
func (c *AIClient) Translate(ctx context.Context, text string, context string) (*Translation, error)
func (c *AIClient) GenerateSentences(ctx context.Context, word string) ([]string, error)
```

### 10.4 Batch Requests for Cost Efficiency

Following Read Frog's pattern, the AI service batches multiple translation/explanation requests into one API call:

```go
// Instead of N API calls:
for _, word := range words {
    resp, _ := client.Translate(ctx, word)
}

// Batch into one call:
resp, _ := client.BatchTranslate(ctx, words)  // {"input": ["word1", "word2", ...]}
```

### 10.5 Content as YAML/JSON (LibreLingo pattern)

Course content is defined in YAML files (contributor-friendly) and compiled into the SQLite database on first load:

```yaml
course:
  title: "English for Programmers"
  level: intermediate
  modules:
    - title: "Git & Version Control"
      skills:
        - title: "Basic Git Commands"
          challenges:
            - type: cards
              prompt: "提交"
              answer: "commit"
            - type: short_input
              prompt: "How do you say '分支' in English?"
              answer: "branch"
              alternatives: ["branches"]
```

### 10.6 Internationalization (i18n)

While the target user is Chinese-speaking, the UI should support i18n from the start for future expansion. Use Mantine's built-in localization or react-i18next. For v1, Chinese + English UI strings.

---

## 11. Open Questions

These need further discussion and prototyping before locking in:

| Question | Options | Decision Needed By |
|----------|---------|-------------------|
| **Which AI provider to default?** | Local Ollama (free but limited), OpenAI (paid, reliable), Claude (good at explanations) | Phase 2 start |
| **Audio TTS: local or cloud?** | Edge TTS (free, 150+ voices, used by Read Frog), local eSpeak (offline, robotic), OpenAI TTS (paid, natural) | Phase 1 end |
| **Course content: built-in vs community?** | Ship with CET-4/CET-6 + IELTS lists built-in; accept community YAML PRs | Phase 3 start |
| **Anki import: full APKG or just CSV?** | APKG is SQLite-based under the hood, doable. CSV is simpler. Both? | Phase 5 |
| **Cloud sync: needed for v1?** | "No" for v1. Add optional sync layer post-v1 using SQLite WAL replication or manual file export. | Post-v1 |
| **Mobile companion app?** | Separate Flutter/Swift app or PWA? Not in scope for v1. | Post-v1 |
| **Pronunciation: IPA or Pinyin-based?** | IPA for accuracy (programmers can learn IPA). Show both? | Phase 1 |

---

## Appendix A: Glossary

| Term | Definition |
|------|------------|
| SRS | Spaced Repetition System — algorithm for scheduling reviews at optimal intervals |
| SM-2 | SuperMemo 2 algorithm — the basis for Anki's scheduling |
| Ease factor | Multiplier applied to interval when a card is correctly recalled (default 2.5) |
| Leech | A card that is consistently forgotten, flagged for special handling |
| Graduated | A card that has passed the initial learning phase and entered long-term review |
| i+1 | Input that is just slightly above the learner's current level (Krashen's hypothesis) |
| Cloze deletion | A fill-in-the-blank style card where part of the sentence is hidden |
| Chips | Drag-and-drop sentence construction challenge (LibreLingo term) |

## Appendix B: Reference Repository Links

| Repository | URL | Key Learnings |
|------------|-----|---------------|
| Echo-Loop | [github.com/echo-loop/Echo-Loop](https://github.com/echo-loop/Echo-Loop) | SM-2 SRS engine, card pipeline |
| LibreLingo | [github.com/kantord/LibreLingo](https://github.com/kantord/LibreLingo) | Course hierarchy, multi-challenge types, stale skill SRS |
| A Programmer's Guide to English | [github.com/yujiangshui/A-Programmers-Guide-to-English](https://github.com/yujiangshui/A-Programmers-Guide-to-English) | Learning methodology, three-pillar training, corpus-first |
| Read Frog | [github.com/mengxi-ream/read-frog](https://github.com/mengxi-ream/read-frog) | Immersive reading, AI contextual explanations, level-adaptive |
