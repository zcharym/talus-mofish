import { NumberInput, Stack, Text } from "@mantine/core";
import type { AppConfigForm } from "../../hooks/useAppConfig";

interface EnglishLearningTabProps {
  dailyGoalMinutes: number;
  wordsPerSession: number;
  onChange: <K extends keyof AppConfigForm>(key: K, value: AppConfigForm[K]) => void;
}

export function EnglishLearningTab({
  dailyGoalMinutes,
  wordsPerSession,
  onChange,
}: EnglishLearningTabProps) {
  return (
    <Stack gap="md">
      <Text size="sm" c="dimmed">
        Settings for the English Learning domain in Management and Agent quick actions.
      </Text>

      <NumberInput
        label="Daily study goal"
        description="Target minutes per day for English Learning sessions"
        value={dailyGoalMinutes}
        onChange={(value) => onChange("dailyGoalMinutes", Number(value) || 30)}
        min={1}
        max={480}
        suffix=" min"
      />

      <NumberInput
        label="Words per session"
        description="Number of words to practice in each recite session"
        value={wordsPerSession}
        onChange={(value) => onChange("wordsPerSession", Number(value) || 20)}
        min={1}
        max={200}
      />
    </Stack>
  );
}
