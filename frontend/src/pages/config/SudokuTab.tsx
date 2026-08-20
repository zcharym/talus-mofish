import { PasswordInput, Stack, Text } from "@mantine/core";
import type { AppConfigForm } from "../../hooks/useAppConfig";

interface SudokuTabProps {
  sudokuAPIKey: string;
  onChange: <K extends keyof AppConfigForm>(key: K, value: AppConfigForm[K]) => void;
}

export function SudokuTab({ sudokuAPIKey, onChange }: SudokuTabProps) {
  return (
    <Stack gap="md">
      <Text size="sm" c="dimmed">
        Optional API key for the YouDoSudoku puzzle service used by Agent window games.
      </Text>
      <PasswordInput
        label="YouDoSudoku API key"
        description="Generate a key at youdosudoku.com. Leave blank if the API accepts unauthenticated requests."
        value={sudokuAPIKey}
        onChange={(event) => onChange("sudokuAPIKey", event.currentTarget.value)}
      />
    </Stack>
  );
}
