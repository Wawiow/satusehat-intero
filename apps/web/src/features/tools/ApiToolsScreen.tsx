import { Button, Form, Input, Select, Space, Typography } from "antd";
import { useState } from "react";
import { JsonPanel } from "@/components/JsonPanel";
import { API_BASE_URL } from "@/lib/api";

type Method = "GET" | "POST" | "PUT";

type ToolForm = {
  method: Method;
  endpoint: string;
  pathId?: string;
  query?: string;
  body?: string;
};

const endpointOptions = [
  "/token",
  "/local/patients",
  "/patients",
  "/local/practitioners",
  "/practitioners",
  "/local/locations",
  "/locations",
  "/local/encounters",
  "/encounters",
  "/encounters/{id}",
].map((value) => ({ value, label: value }));

export function ApiToolsScreen() {
  const [result, setResult] = useState<unknown>();
  const [error, setError] = useState<string>();
  const [loading, setLoading] = useState(false);
  const [method, setMethod] = useState<Method>("GET");
  const [endpoint, setEndpoint] = useState("/local/patients");

  async function handleExecute(values: ToolForm) {
    setLoading(true);
    setError(undefined);
    setResult(undefined);

    try {
      const path = values.endpoint.replace("{id}", encodeURIComponent(values.pathId ?? ""));
      const query = values.method === "GET" && values.query ? `?${values.query.replace(/^\?/, "")}` : "";
      const response = await fetch(`${API_BASE_URL}${path}${query}`, {
        method: values.method,
        headers: { "Content-Type": "application/json" },
        body: values.method === "GET" || !values.body ? undefined : values.body,
      });
      const contentType = response.headers.get("content-type") ?? "";
      const payload = contentType.includes("application/json") ? await response.json() : await response.text();
      setResult({ status: response.status, ok: response.ok, body: payload });
    } catch (caught) {
      setError(caught instanceof Error ? caught.message : "Unknown request error");
    } finally {
      setLoading(false);
    }
  }

  return (
    <div className="resource-grid">
      <div className="panel">
        <Space direction="vertical" size={4}>
          <Typography.Title level={3}>Alat Integrasi</Typography.Title>
          <Typography.Text type="secondary">Runner endpoint untuk tim TI rumah sakit saat validasi koneksi backend. Base URL: {API_BASE_URL}</Typography.Text>
        </Space>
      </div>
      <div className="panel">
        <Form<ToolForm>
          layout="vertical"
          initialValues={{ method: "GET", endpoint: "/local/patients" }}
          onValuesChange={(_, values) => {
            setMethod(values.method ?? "GET");
            setEndpoint(values.endpoint ?? "/local/patients");
          }}
          onFinish={handleExecute}
        >
          <div className="form-grid">
            <Form.Item name="method" label="Method" rules={[{ required: true, message: "Method wajib diisi" }]}>
              <Select options={["GET", "POST", "PUT"].map((value) => ({ value, label: value }))} />
            </Form.Item>
            <Form.Item name="endpoint" label="Endpoint" rules={[{ required: true, message: "Endpoint wajib diisi" }]}> 
              <Select options={endpointOptions} />
            </Form.Item>
            {endpoint === "/encounters/{id}" ? (
              <Form.Item name="pathId" label="Path ID" rules={[{ required: true, message: "Path ID wajib diisi" }]}>
                <Input />
              </Form.Item>
            ) : null}
            {method === "GET" ? <Form.Item name="query" label="Query String"><Input placeholder="nik=123&name=Nama" /></Form.Item> : null}
          </div>
          {method !== "GET" ? (
            <Form.Item name="body" label="JSON Body">
              <Input.TextArea rows={10} placeholder={'{\n  "status": "finished"\n}'} />
            </Form.Item>
          ) : null}
          <Button type="primary" htmlType="submit" loading={loading}>Kirim Request</Button>
        </Form>
      </div>
      <JsonPanel title="Respons API" data={result} error={error} />
    </div>
  );
}
