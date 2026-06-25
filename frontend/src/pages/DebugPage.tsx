import { useState } from "react";
import { Stack, Tabs, Text } from "@mantine/core";
import { debugDemos } from "./debug/registry";

export function DebugPage() {
  const [activeDemo, setActiveDemo] = useState(debugDemos[0]?.id ?? "");

  const selected = debugDemos.find((demo) => demo.id === activeDemo) ?? debugDemos[0];
  const DemoComponent = selected?.Component;

  return (
    <Stack gap="md" maw={720}>
      <Text c="dimmed" size="sm">
        Component playground — each tab loads an isolated demo module. Enable debug mode in
        Configuration to show this page in the sidebar.
      </Text>

      <Tabs value={activeDemo} onChange={(value) => setActiveDemo(value ?? activeDemo)}>
        <Tabs.List>
          {debugDemos.map((demo) => (
            <Tabs.Tab key={demo.id} value={demo.id}>
              {demo.title}
            </Tabs.Tab>
          ))}
        </Tabs.List>

        {debugDemos.map((demo) => (
          <Tabs.Panel key={demo.id} value={demo.id} pt="md">
            <Text size="sm" c="dimmed" mb="md">
              {demo.description}
            </Text>
            {demo.id === activeDemo && DemoComponent ? <DemoComponent /> : null}
          </Tabs.Panel>
        ))}
      </Tabs>
    </Stack>
  );
}
