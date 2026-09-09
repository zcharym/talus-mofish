import { KeyboardEvent } from 'react';
import { IconCheck, IconEraser } from '@tabler/icons-react';
import { Box, Burger, Button, Group, Loader, Select, Text, Title } from '@mantine/core';
import { useSudokuGame } from '../../../hooks/useSudokuGame';
import classes from './SudokuBoard.module.css';

interface SudokuBoardProps {
  sessionId: string;
  sessionTitle: string | null;
  onSessionUpdated: () => Promise<void>;
  onOpenSidebar?: () => void;
}

const DIFFICULTIES = [
  { value: 'easy', label: 'Easy' },
  { value: 'medium', label: 'Medium' },
  { value: 'hard', label: 'Hard' },
];

function isGiven(puzzle: string, index: number): boolean {
  const cell = puzzle[index];
  return cell !== undefined && cell !== '0';
}

export function SudokuBoard({ sessionId, sessionTitle, onSessionUpdated, onOpenSidebar }: SudokuBoardProps) {
  const {
    game,
    difficulty,
    setDifficulty,
    selected,
    setSelected,
    conflicts,
    loading,
    busy,
    statusText,
    setCell,
    checkGame,
    newPuzzle,
  } = useSudokuGame(sessionId, onSessionUpdated);

  const handleKeyDown = (event: KeyboardEvent<HTMLDivElement>) => {
    if (selected === null || !game) {
      return;
    }
    if (event.key >= '1' && event.key <= '9') {
      event.preventDefault();
      if (!isGiven(game.puzzle, selected) && game.status !== 'solved' && !busy) {
        void setCell(selected, Number(event.key));
      }
      return;
    }
    if (event.key === 'Backspace' || event.key === 'Delete' || event.key === '0') {
      event.preventDefault();
      if (!isGiven(game.puzzle, selected) && game.status !== 'solved' && !busy) {
        void setCell(selected, 0);
      }
      return;
    }
    const row = Math.floor(selected / 9);
    const col = selected % 9;
    let next = selected;
    if (event.key === 'ArrowLeft') {
      next = row * 9 + Math.max(0, col - 1);
    } else if (event.key === 'ArrowRight') {
      next = row * 9 + Math.min(8, col + 1);
    } else if (event.key === 'ArrowUp') {
      next = Math.max(0, row - 1) * 9 + col;
    } else if (event.key === 'ArrowDown') {
      next = Math.min(8, row + 1) * 9 + col;
    } else {
      return;
    }
    event.preventDefault();
    setSelected(next);
  };

  const solved = game?.status === 'solved';

  return (
    <Box className={classes.page}>
      <Box className={classes.header}>
        <Group gap="sm" wrap="nowrap" className={classes.headerTitle}>
          {onOpenSidebar ? (
            <Burger opened={false} onClick={onOpenSidebar} size="sm" aria-label="Open chats" />
          ) : null}
          <Title order={4}>{sessionTitle || 'Sudoku'}</Title>
        </Group>
        <Group gap="sm" wrap="wrap">
          <Select
            size="xs"
            w={120}
            data={DIFFICULTIES}
            value={difficulty}
            onChange={(value) => setDifficulty(value || 'easy')}
            allowDeselect={false}
            disabled={busy}
          />
          <Button size="xs" variant="light" loading={busy} onClick={() => void newPuzzle()}>
            New puzzle
          </Button>
        </Group>
      </Box>

      <Box className={classes.body}>
        {loading || !game ? (
          <Loader size="sm" />
        ) : (
          <>
            <Box
              className={classes.board}
              tabIndex={0}
              role="grid"
              aria-label="Sudoku board"
              onKeyDown={handleKeyDown}
            >
              {Array.from({ length: 81 }, (_, index) => {
                const given = isGiven(game.puzzle, index);
                const value = game.board[index] === '0' ? '' : game.board[index];
                const col = index % 9;
                const row = Math.floor(index / 9);
                return (
                  <button
                    key={index}
                    type="button"
                    className={classes.cell}
                    data-given={given || undefined}
                    data-selected={selected === index || undefined}
                    data-conflict={conflicts.includes(index) || undefined}
                    data-box-right={col === 2 || col === 5 || undefined}
                    data-box-bottom={row === 2 || row === 5 || undefined}
                    disabled={solved}
                    onClick={() => setSelected(index)}
                  >
                    {value}
                  </button>
                );
              })}
            </Box>

            <Group className={classes.pad} gap="xs">
              {Array.from({ length: 9 }, (_, digit) => (
                <Button
                  key={digit + 1}
                  variant="default"
                  size="sm"
                  w={36}
                  px={0}
                  disabled={solved || selected === null || busy || (selected !== null && isGiven(game.puzzle, selected))}
                  onClick={() => selected !== null && void setCell(selected, digit + 1)}
                >
                  {digit + 1}
                </Button>
              ))}
              <Button
                variant="default"
                size="sm"
                leftSection={<IconEraser size={14} />}
                disabled={solved || selected === null || busy || (selected !== null && isGiven(game.puzzle, selected))}
                onClick={() => selected !== null && void setCell(selected, 0)}
              >
                Clear
              </Button>
              <Button
                size="sm"
                leftSection={<IconCheck size={14} />}
                disabled={solved || busy}
                onClick={() => void checkGame()}
              >
                Check
              </Button>
            </Group>

            <Text className={classes.status} size="sm" c={solved ? 'teal' : 'dimmed'}>
              {statusText}
            </Text>
          </>
        )}
      </Box>
    </Box>
  );
}
