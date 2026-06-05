import { Alert, Typography } from "antd";

type JsonPanelProps = {
  title: string;
  data: unknown;
  error?: string;
};

export function JsonPanel({ title, data, error }: JsonPanelProps) {
  return (
    <section className="json-panel">
      <Typography.Title level={5}>{title}</Typography.Title>
      {error ? <Alert type="error" showIcon message={error} /> : null}
      <pre>{JSON.stringify(data ?? null, null, 2)}</pre>
    </section>
  );
}
