import { KeyboardEvent, useCallback, useEffect, useState } from 'react';
import { IconCheck, IconEraser } from '@tabler/icons-react';
import { Box, Button, Group, Loader, Select, Text, Title } from '@mantine/core';
import { SudokuService } from '../../../../bindings/github.com/songwei.ma/talus-mofish/backend/services';
import { notify } from '../../../services/notifications';
import classes from './SudokuBoard.module.css';

export interface SudokuGameState {
  id: string;
  session_id: string;
  difficulty: string;
  puzzle: string;
  board: string;
  status: string;
  created_at: string;
  updated_at: string;
}

interface SudokuBoardProps {
  sessionId: string;
  sessionTitle: string | null;
  onSessionUpdated: () => Promise<void>;
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

export function SudokuBoard({ sessionId, sessionTitle, onSessionUpdated }: SudokuBoardProps) {
  const [game, setGame] = useState<SudokuGameState | null>(null);
  const [difficulty, setDifficulty] = useState('easy');
  const [selected, setSelected] = useState<number | null>(null);
  const [conflicts, setConflicts] = useState<number[]>([]);
  const [loading, setLoading] = useState(true);
  const [busy, setBusy] = useState(false);
  const [statusText, setStatusText] = useState('');

  const loadGame = useCallback(async () => {
    setLoading(true);
    setConflicts([]);
    setStatusText('');
    try {
      const next = (await SudokuService.GetSudokuGame(sessionId)) as SudokuGameState;
      setGame(next);
      setDifficulty(next.difficulty || 'easy');
    } catch (err) {
      notify.failed('Failed to load Sudoku', String(err));
    } finally {
      setLoading(false);
    }
  }, [sessionId]);

  useEffect(() => {
    void loadGame();
  }, [loadGame]);

  const applyGame = (next: SudokuGameState) => {
    setGame(next);
    setDifficulty(next.difficulty || difficulty);
    if (next.status === 'solved') {
      setConflicts([]);
      setStatusText('Solved');
    }
  };

  const setCell = async (index: number, value: number) => {
    if (!game || game.status === 'solved' || isGiven(game.puzzle, index) || busy) {
      return;
    }
    setBusy(true);
    try {
      const next = (await SudokuService.SetSudokuCell(sessionId, index, value)) as SudokuGameState;
      applyGame(next);
      setConflicts((current) => current.filter((item) => item !== index));
      if (next.status === 'solved') {
        await onSessionUpdated();
      }
    } catch (err) {
      notify.failed('Could not update cell', String(err));
    } finally {
      setBusy(false);
    }
  };

  const handleCheck = async () => {
    if (!game || busy) {
      return;
    }
    setBusy(true);
    try {
      const result = await SudokuService.CheckSudokuGame(sessionId);
      const nextGame = result.game as SudokuGameState;
      applyGame(nextGame);
      const nextConflicts = (result.conflicts ?? []) as number[];
      setConflicts(nextConflicts);
      if (result.solved) {
        setStatusText('Solved');
        await onSessionUpdated();
      } else if (nextConflicts.length === 0) {
        setStatusText('No mistakes yet — keep going');
      } else {
        setStatusText(`${nextConflicts.length} incorrect ${nextConflicts.length === 1 ? 'cell' : 'cells'}`);
      }
    } catch (err) {
      notify.failed('Check failed', String(err));
    } finally {
      setBusy(false);
    }
  };

  const handleNewPuzzle = async () => {
    if (busy) {
      return;
    }
    setBusy(true);
    setConflicts([]);
    setStatusText('');
    try {
      const next = (await SudokuService.NewSudokuPuzzle(sessionId, difficulty)) as SudokuGameState;
      applyGame(next);
      setSelected(null);
      await onSessionUpdated();
    } catch (err) {
      notify.failed('Could not fetch a new puzzle', String(err));
    } finally {
      setBusy(false);
    }
  };

  const handleKeyDown = (event: KeyboardEvent<HTMLDivElement>) => {
    if (selected === null || !game) {
      return;
    }
    if (event.key >= '1' && event.key <= '9') {
      event.preventDefault();
      void setCell(selected, Number(event.key));
      return;
    }
    if (event.key === 'Backspace' || event.key === 'Delete' || event.key === '0') {
      event.preventDefault();
      void setCell(selected, 0);
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
        <Title order={4}>{sessionTitle || 'Sudoku'}</Title>
        <Group gap="sm" wrap="nowrap">
          <Select
            size="xs"
            w={120}
            data={DIFFICULTIES}
            value={difficulty}
            onChange={(value) => setDifficulty(value || 'easy')}
            allowDeselect={false}
            disabled={busy}
          />
          <Button size="xs" variant="light" loading={busy} onClick={() => void handleNewPuzzle()}>
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
                  disabled={solved || selected === null || busy}
                  onClick={() => selected !== null && void setCell(selected, digit + 1)}
                >
                  {digit + 1}
                </Button>
              ))}
              <Button
                variant="default"
                size="sm"
                leftSection={<IconEraser size={14} />}
                disabled={solved || selected === null || busy}
                onClick={() => selected !== null && void setCell(selected, 0)}
              >
                Clear
              </Button>
              <Button
                size="sm"
                leftSection={<IconCheck size={14} />}
                disabled={solved || busy}
                onClick={() => void handleCheck()}
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
