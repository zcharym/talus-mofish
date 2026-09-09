import { IconCloud, IconRefresh } from '@tabler/icons-react';
import {
  Alert,
  Badge,
  Box,
  Burger,
  Button,
  Group,
  Loader,
  Paper,
  ScrollArea,
  Stack,
  Table,
  Text,
  Title,
} from '@mantine/core';
import { useCloudflareDashboard } from '../../../hooks/useCloudflareDashboard';
import type { CloudflareWorker } from '../../../utils/api';
import classes from './CloudflareDashboard.module.css';

interface CloudflareDashboardProps {
  configured: boolean;
  onOpenSidebar?: () => void;
  onOpenManagement: () => void;
}

function formatCount(value: number | undefined): string {
  return new Intl.NumberFormat().format(Math.round(value || 0));
}

function formatRate(value: number | undefined): string {
  const rate = (value || 0) * 100;
  if (rate === 0) {
    return '0%';
  }
  if (rate < 0.1) {
    return `${rate.toFixed(2)}%`;
  }
  return `${rate.toFixed(1)}%`;
}

function formatCpu(microseconds: number | undefined): string {
  if (!microseconds) {
    return '—';
  }
  return `${(microseconds / 1000).toFixed(1)} ms`;
}

function formatTimestamp(value: string | undefined): string {
  if (!value) {
    return '—';
  }
  const date = new Date(value);
  if (Number.isNaN(date.getTime())) {
    return value;
  }
  return date.toLocaleString();
}

export function CloudflareDashboard({
  configured,
  onOpenSidebar,
  onOpenManagement,
}: CloudflareDashboardProps) {
  const { loading, refreshing, snapshot, loadError, loadDashboard } = useCloudflareDashboard(configured);

  const header = (
    <Box className={classes.header}>
      <Group gap="sm" wrap="nowrap" className={classes.headerTitle}>
        {onOpenSidebar ? (
          <Burger opened={false} onClick={onOpenSidebar} size="sm" aria-label="Open chats" />
        ) : null}
        <IconCloud size={20} />
        <Title order={4}>Cloudflare</Title>
      </Group>
      {configured ? (
        <Group gap="sm" wrap="wrap">
          {snapshot?.fetchedAt ? (
            <Text size="xs" c="dimmed">
              Updated {formatTimestamp(snapshot.fetchedAt)}
            </Text>
          ) : null}
          <Button
            size="xs"
            variant="light"
            leftSection={<IconRefresh size={14} />}
            loading={refreshing}
            onClick={() => void loadDashboard(true)}
          >
            Refresh
          </Button>
        </Group>
      ) : null}
    </Box>
  );

  if (!configured) {
    return (
      <Box className={classes.page}>
        {header}
        <Box className={classes.empty}>
          <Title order={4}>Connect a Cloudflare account</Title>
          <Text size="sm" c="dimmed" maw={420}>
            Add an account ID and API token in Management → Configuration → Cloudflare to monitor
            Workers, 24h traffic, D1, and KV.
          </Text>
          <Button onClick={onOpenManagement}>Open Management</Button>
        </Box>
      </Box>
    );
  }

  if (loading) {
    return (
      <Box className={classes.page}>
        {header}
        <Box className={classes.empty}>
          <Loader size="sm" />
        </Box>
      </Box>
    );
  }

  if (loadError && !snapshot) {
    return (
      <Box className={classes.page}>
        {header}
        <Box className={classes.empty}>
          <Alert color="red" title="Could not load dashboard">
            {loadError}
          </Alert>
          <Button variant="light" onClick={() => void loadDashboard()}>
            Try again
          </Button>
        </Box>
      </Box>
    );
  }

  const errors = snapshot?.errors;
  const workers = snapshot?.workers ?? [];
  const known = snapshot?.known ?? [];
  const d1 = snapshot?.d1 ?? [];
  const kv = snapshot?.kv ?? [];
  const kpis = snapshot?.kpis;
  const accountName = snapshot?.account?.name || snapshot?.account?.id || 'Cloudflare account';

  return (
    <Box className={classes.page}>
      {header}
      <ScrollArea className={classes.body} type="auto">
        <Stack gap="md" className={classes.scrollInner}>
          {errors?.account ? (
            <Alert color="red" title="Account">
              {errors.account}
            </Alert>
          ) : (
            <Group gap="sm" wrap="wrap">
              <Badge variant="light">{accountName}</Badge>
              {snapshot?.account?.id ? (
                <Text size="xs" c="dimmed">
                  {snapshot.account.id}
                </Text>
              ) : null}
            </Group>
          )}

          <Box className={classes.kpiGrid}>
            <KpiCard label="Workers" value={formatCount(kpis?.workerCount)} />
            <KpiCard label="Requests (24h)" value={formatCount(kpis?.requests24h)} />
            <KpiCard label="Errors (24h)" value={formatCount(kpis?.errors24h)} />
            <KpiCard label="Error rate" value={formatRate(kpis?.errorRate)} />
          </Box>

          {errors?.analytics ? (
            <Alert color="yellow" title="Analytics">
              {errors.analytics}
            </Alert>
          ) : null}

          <section>
            <Title order={5} mb="xs">
              Known services
            </Title>
            <Box className={classes.knownGrid}>
              {known.map((service) => (
                <Paper key={service.id} withBorder p="sm">
                  <Group justify="space-between" wrap="nowrap" mb={4}>
                    <Text fw={600} size="sm">
                      {service.title}
                    </Text>
                    <Badge color={service.present ? 'teal' : 'gray'} variant="light">
                      {service.present ? 'Found' : 'Missing'}
                    </Badge>
                  </Group>
                  <Text size="xs" c="dimmed">
                    {service.id}
                    {service.bindingsNote ? ` · ${service.bindingsNote}` : ''}
                  </Text>
                  {service.worker ? (
                    <Text size="xs" mt={6}>
                      {formatCount(service.worker.requests24h)} req /{' '}
                      {formatCount(service.worker.errors24h)} err (24h)
                    </Text>
                  ) : (
                    <Text size="xs" c="dimmed" mt={6}>
                      Not found on this account.
                    </Text>
                  )}
                </Paper>
              ))}
            </Box>
          </section>

          <section>
            <Title order={5} mb="xs">
              Workers
            </Title>
            {errors?.workers ? (
              <Alert color="yellow" title="Workers">
                {errors.workers}
              </Alert>
            ) : workers.length === 0 ? (
              <Text size="sm" c="dimmed">
                No Workers scripts on this account.
              </Text>
            ) : (
              <Box className={classes.tableWrap}>
                <Table striped highlightOnHover withTableBorder>
                  <Table.Thead>
                    <Table.Tr>
                      <Table.Th>Name</Table.Th>
                      <Table.Th>Modified</Table.Th>
                      <Table.Th>Requests</Table.Th>
                      <Table.Th>Errors</Table.Th>
                      <Table.Th>CPU p99</Table.Th>
                      <Table.Th>Cron</Table.Th>
                    </Table.Tr>
                  </Table.Thead>
                  <Table.Tbody>
                    {workers.map((worker) => (
                      <WorkerRow key={worker.name} worker={worker} />
                    ))}
                  </Table.Tbody>
                </Table>
              </Box>
            )}
          </section>

          <section>
            <Title order={5} mb="xs">
              D1
            </Title>
            {errors?.d1 ? (
              <Alert color="yellow" title="D1">
                {errors.d1}
              </Alert>
            ) : d1.length === 0 ? (
              <Text size="sm" c="dimmed">
                No D1 databases visible with this token.
              </Text>
            ) : (
              <Stack gap={4}>
                {d1.map((item) => (
                  <Text key={item.id || item.name} size="sm">
                    {item.name || item.id}
                    {item.id && item.name ? (
                      <Text span size="xs" c="dimmed">
                        {' '}
                        · {item.id}
                      </Text>
                    ) : null}
                  </Text>
                ))}
              </Stack>
            )}
          </section>

          <section>
            <Title order={5} mb="xs">
              KV
            </Title>
            {errors?.kv ? (
              <Alert color="yellow" title="KV">
                {errors.kv}
              </Alert>
            ) : kv.length === 0 ? (
              <Text size="sm" c="dimmed">
                No KV namespaces visible with this token.
              </Text>
            ) : (
              <Stack gap={4}>
                {kv.map((item) => (
                  <Text key={item.id || item.title} size="sm">
                    {item.title || item.id}
                    {item.id && item.title ? (
                      <Text span size="xs" c="dimmed">
                        {' '}
                        · {item.id}
                      </Text>
                    ) : null}
                  </Text>
                ))}
              </Stack>
            )}
          </section>
        </Stack>
      </ScrollArea>
    </Box>
  );
}

function KpiCard({ label, value }: { label: string; value: string }) {
  return (
    <Paper withBorder p="sm">
      <Text size="xs" c="dimmed">
        {label}
      </Text>
      <Text fw={700} size="lg">
        {value}
      </Text>
    </Paper>
  );
}

function WorkerRow({ worker }: { worker: CloudflareWorker }) {
  const cron = worker.cron?.filter(Boolean) ?? [];
  return (
    <Table.Tr>
      <Table.Td>{worker.name}</Table.Td>
      <Table.Td>{formatTimestamp(worker.modifiedOn)}</Table.Td>
      <Table.Td>{formatCount(worker.requests24h)}</Table.Td>
      <Table.Td>{formatCount(worker.errors24h)}</Table.Td>
      <Table.Td>{formatCpu(worker.cpuTimeP99)}</Table.Td>
      <Table.Td>{cron.length > 0 ? cron.join(', ') : '—'}</Table.Td>
    </Table.Tr>
  );
}
