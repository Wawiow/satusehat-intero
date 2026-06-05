import { Button, Form, Input, Select, Table, Tabs } from "antd";
import type { ColumnsType } from "antd/es/table";
import { useCallback, useEffect, useState } from "react";
import { DataToolbar } from "@/components/DataToolbar";
import { JsonPanel } from "@/components/JsonPanel";
import {
  createEncounter,
  formatApiError,
  getEncounterById,
  getLocalEncounters,
  updateEncounterStatus,
  type CreateEncounterRequest,
  type EncounterResponse,
} from "@/lib/api";

const encounterColumns: ColumnsType<EncounterResponse> = [
  { title: "ID", dataIndex: "id", key: "id", width: 260 },
  { title: "Nomor Registrasi", dataIndex: "identifier_value", key: "identifier_value", width: 180 },
  { title: "Status", dataIndex: "status", key: "status", width: 130 },
  { title: "Pasien", dataIndex: "subject_id", key: "subject_id", width: 220 },
  { title: "Unit Layanan", dataIndex: "location_id", key: "location_id", width: 260 },
  { title: "Waktu Mulai", dataIndex: "start_time", key: "start_time", width: 180 },
];

type GetEncounterValues = { id: string };
type UpdateEncounterValues = { id: string; status: string };

export function EncountersScreen() {
  const [encounters, setEncounters] = useState<EncounterResponse[]>([]);
  const [loading, setLoading] = useState(false);
  const [result, setResult] = useState<unknown>();
  const [error, setError] = useState<string>();
  const [createForm] = Form.useForm<CreateEncounterRequest>();

  const loadEncounters = useCallback(async () => {
    setLoading(true);
    setError(undefined);
    try {
      setEncounters(await getLocalEncounters());
    } catch (caught) {
      setError(formatApiError(caught));
    } finally {
      setLoading(false);
    }
  }, []);

  useEffect(() => {
    void loadEncounters();
  }, [loadEncounters]);

  async function handleGet(values: GetEncounterValues) {
    setError(undefined);
    try {
      setResult(await getEncounterById(values.id));
    } catch (caught) {
      setResult(undefined);
      setError(formatApiError(caught));
    }
  }

  async function handleCreate(values: CreateEncounterRequest) {
    setError(undefined);
    try {
      const data = await createEncounter(values);
      setResult(data);
      createForm.resetFields();
      await loadEncounters();
    } catch (caught) {
      setResult(undefined);
      setError(formatApiError(caught));
    }
  }

  async function handleUpdate(values: UpdateEncounterValues) {
    setError(undefined);
    try {
      const data = await updateEncounterStatus(values.id, { status: values.status });
      setResult(data);
      await loadEncounters();
    } catch (caught) {
      setResult(undefined);
      setError(formatApiError(caught));
    }
  }

  return (
    <div className="resource-grid">
      <DataToolbar title="Kunjungan Pasien" description="Buat registrasi kunjungan, pantau Encounter, dan perbarui status pelayanan pasien." loading={loading} onReload={loadEncounters} />
      <div className="table-wrap">
        <Table rowKey={(row) => row.id || row.identifier_value} columns={encounterColumns} dataSource={encounters} loading={loading} scroll={{ x: 1230 }} />
      </div>
      <Tabs
        items={[
          {
            key: "get",
            label: "Cari Kunjungan",
            children: (
              <div className="panel">
                <Form layout="vertical" onFinish={handleGet}>
                  <Form.Item name="id" label="Encounter ID" rules={[{ required: true, message: "ID wajib diisi" }]}><Input /></Form.Item>
                  <Button type="primary" htmlType="submit">Ambil Detail Kunjungan</Button>
                </Form>
              </div>
            ),
          },
          {
            key: "create",
            label: "Registrasi Kunjungan",
            children: (
              <div className="panel">
                <Form form={createForm} layout="vertical" onFinish={handleCreate}>
                  <div className="form-grid">
                    <Form.Item name="identifier_value" label="Nomor Registrasi" rules={[{ required: true, message: "Nomor registrasi wajib diisi" }]}><Input /></Form.Item>
                    <Form.Item name="subject_id" label="ID Pasien" rules={[{ required: true, message: "ID pasien wajib diisi" }]}><Input placeholder="Patient/123" /></Form.Item>
                    <Form.Item name="location_id" label="ID Unit Layanan" rules={[{ required: true, message: "ID unit layanan wajib diisi" }]}><Input placeholder="Location/abc" /></Form.Item>
                    <Form.Item name="practitioner_id" label="ID Tenaga Medis" rules={[{ required: true, message: "ID tenaga medis wajib diisi" }]}><Input placeholder="Practitioner/abc" /></Form.Item>
                    <Form.Item name="start_time" label="Waktu Mulai" rules={[{ required: true, message: "Waktu mulai wajib diisi" }]}><Input placeholder="2026-06-05T10:30:00Z" /></Form.Item>
                  </div>
                  <Button type="primary" htmlType="submit">Simpan Kunjungan</Button>
                </Form>
              </div>
            ),
          },
          {
            key: "update",
            label: "Ubah Status",
            children: (
              <div className="panel">
                <Form layout="vertical" onFinish={handleUpdate}>
                  <Form.Item name="id" label="Encounter ID" rules={[{ required: true, message: "ID wajib diisi" }]}><Input /></Form.Item>
                  <Form.Item name="status" label="Status" rules={[{ required: true, message: "Status wajib diisi" }]}>
                    <Select options={["arrived", "in-progress", "finished", "cancelled"].map((value) => ({ value, label: value }))} />
                  </Form.Item>
                  <Button type="primary" htmlType="submit">Perbarui Status</Button>
                </Form>
              </div>
            ),
          },
        ]}
      />
      <JsonPanel title="Respons Kunjungan" data={result} error={error} />
    </div>
  );
}
