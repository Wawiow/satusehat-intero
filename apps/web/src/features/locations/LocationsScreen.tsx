import { Button, Form, Input, InputNumber, Table, Tabs } from "antd";
import type { ColumnsType } from "antd/es/table";
import { useCallback, useEffect, useState } from "react";
import { DataToolbar } from "@/components/DataToolbar";
import { JsonPanel } from "@/components/JsonPanel";
import {
  createLocation,
  formatApiError,
  getLocalLocations,
  searchLocations,
  type CreateLocationRequest,
  type LocationResponse,
} from "@/lib/api";

const locationColumns: ColumnsType<LocationResponse> = [
  { title: "ID", dataIndex: "id", key: "id", width: 260 },
  { title: "Kode Unit", dataIndex: "identifier_value", key: "identifier_value", width: 160 },
  { title: "Nama Unit Layanan", dataIndex: "name", key: "name", width: 220 },
  { title: "Deskripsi", dataIndex: "description", key: "description", width: 420 },
  { title: "Telepon", dataIndex: "phone", key: "phone", width: 120 },
];

type LocationSearch = {
  id?: string;
  identifier?: string;
  page?: number;
  limit?: number;
};

export function LocationsScreen() {
  const [locations, setLocations] = useState<LocationResponse[]>([]);
  const [loading, setLoading] = useState(false);
  const [result, setResult] = useState<unknown>();
  const [error, setError] = useState<string>();
  const [searchForm] = Form.useForm<LocationSearch>();
  const [createForm] = Form.useForm<CreateLocationRequest>();

  const loadLocations = useCallback(async () => {
    setLoading(true);
    setError(undefined);
    try {
      setLocations(await getLocalLocations());
    } catch (caught) {
      setError(formatApiError(caught));
    } finally {
      setLoading(false);
    }
  }, []);

  useEffect(() => {
    void loadLocations();
  }, [loadLocations]);

  async function handleSearch(values: LocationSearch) {
    setError(undefined);
    try {
      setResult(await searchLocations(values));
    } catch (caught) {
      setResult(undefined);
      setError(formatApiError(caught));
    }
  }

  async function handleCreate(values: CreateLocationRequest) {
    setError(undefined);
    try {
      const data = await createLocation(values);
      setResult(data);
      createForm.resetFields();
      await loadLocations();
    } catch (caught) {
      setResult(undefined);
      setError(formatApiError(caught));
    }
  }

  return (
    <div className="resource-grid">
      <DataToolbar title="Unit Layanan" description="Kelola poli, ruang, atau instalasi rumah sakit yang akan dipetakan sebagai Location di SatuSehat." loading={loading} onReload={loadLocations} />
      <div className="table-wrap">
        <Table rowKey={(row) => row.id || row.identifier_value} columns={locationColumns} dataSource={locations} loading={loading} scroll={{ x: 1180 }} />
      </div>
      <Tabs
        items={[
          {
            key: "search",
            label: "Cari Unit",
            children: (
              <div className="panel">
                <Form form={searchForm} layout="vertical" initialValues={{ page: 1, limit: 10 }} onFinish={handleSearch}>
                  <div className="form-grid">
                    <Form.Item name="id" label="ID"><Input /></Form.Item>
                    <Form.Item name="identifier" label="Kode Unit"><Input /></Form.Item>
                    <Form.Item name="page" label="Halaman"><InputNumber min={1} className="number-input" /></Form.Item>
                    <Form.Item name="limit" label="Limit"><InputNumber min={-1} className="number-input" /></Form.Item>
                  </div>
                  <Button type="primary" htmlType="submit">Cari Unit Layanan</Button>
                </Form>
              </div>
            ),
          },
          {
            key: "create",
            label: "Tambah Unit",
            children: (
              <div className="panel">
                <Form form={createForm} layout="vertical" onFinish={handleCreate}>
                  <Form.Item name="identifier_value" label="Kode Unit" rules={[{ required: true, message: "Kode unit wajib diisi" }]}><Input /></Form.Item>
                  <Form.Item name="name" label="Nama Unit Layanan" rules={[{ required: true, message: "Nama unit wajib diisi" }]}><Input /></Form.Item>
                  <Form.Item name="description" label="Deskripsi Layanan" rules={[{ required: true, message: "Deskripsi wajib diisi" }]}><Input.TextArea rows={3} /></Form.Item>
                  <Form.Item name="phone" label="Telepon Unit"><Input /></Form.Item>
                  <Button type="primary" htmlType="submit">Simpan Unit Layanan</Button>
                </Form>
              </div>
            ),
          },
        ]}
      />
      <JsonPanel title="Respons Unit Layanan" data={result} error={error} />
    </div>
  );
}
