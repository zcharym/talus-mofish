import { useCallback, useEffect, useState } from 'react';
import { SudokuService, toApiError, type SudokuGame } from '../utils/api';
import { notify } from '../services/notifications';

export function useSudokuGame(sessionId: string, onSessionUpdated: () => Promise<void>) {
  const [game, setGame] = useState<SudokuGame | null>(null);
  const [difficulty, setDifficulty] = useState('easy');
  const [selected, setSelected] = useState<number | null>(null);
  const [conflicts, setConflicts] = useState<number[]>([]);
  const [loading, setLoading] = useState(true);
  const [busy, setBusy] = useState(false);
  const [statusText, setStatusText] = useState('');

  const applyGame = useCallback((next: SudokuGame) => {
    setGame(next);
    setDifficulty(next.difficulty || 'easy');
    if (next.status === 'solved') {
      setConflicts([]);
      setStatusText('Solved');
    }
  }, []);

  const loadGame = useCallback(async () => {
    setLoading(true);
    setConflicts([]);
    setStatusText('');
    try {
      const next = await SudokuService.GetSudokuGame(sessionId);
      applyGame(next);
    } catch (err) {
      notify.failed('Failed to load Sudoku', toApiError(err));
    } finally {
      setLoading(false);
    }
  }, [applyGame, sessionId]);

  useEffect(() => {
    void loadGame();
  }, [loadGame]);

  const setCell = useCallback(
    async (index: number, value: number) => {
      if (!game || game.status === 'solved' || busy) {
        return;
      }
      const given = game.puzzle[index] !== undefined && game.puzzle[index] !== '0';
      if (given) {
        return;
      }
      setBusy(true);
      try {
        const next = await SudokuService.SetSudokuCell(sessionId, index, value);
        applyGame(next);
        setConflicts((current) => current.filter((item) => item !== index));
        if (next.status === 'solved') {
          await onSessionUpdated();
        }
      } catch (err) {
        notify.failed('Could not update cell', toApiError(err));
      } finally {
        setBusy(false);
      }
    },
    [applyGame, busy, game, onSessionUpdated, sessionId],
  );

  const checkGame = useCallback(async () => {
    if (!game || busy) {
      return;
    }
    setBusy(true);
    try {
      const result = await SudokuService.CheckSudokuGame(sessionId);
      applyGame(result.game);
      const nextConflicts = result.conflicts ?? [];
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
      notify.failed('Check failed', toApiError(err));
    } finally {
      setBusy(false);
    }
  }, [applyGame, busy, game, onSessionUpdated, sessionId]);

  const newPuzzle = useCallback(async () => {
    if (busy) {
      return;
    }
    setBusy(true);
    setConflicts([]);
    setStatusText('');
    try {
      const next = await SudokuService.NewSudokuPuzzle(sessionId, difficulty);
      applyGame(next);
      setSelected(null);
      await onSessionUpdated();
    } catch (err) {
      notify.failed('Could not fetch a new puzzle', toApiError(err));
    } finally {
      setBusy(false);
    }
  }, [applyGame, busy, difficulty, onSessionUpdated, sessionId]);

  return {
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
  };
}
