import { Button, Form, Input, InputNumber, Select, Table } from "antd";
import type { ColumnsType } from "antd/es/table";
import { useCallback, useEffect, useState } from "react";
import { DataToolbar } from "@/components/DataToolbar";
import { JsonPanel } from "@/components/JsonPanel";
import { formatApiError, getLocalPractitioners, searchPractitioners, type PersonResponse } from "@/lib/api";

const practitionerColumns: ColumnsType<PersonResponse> = [
  { title: "NIK", dataIndex: "nik", key: "nik", width: 170 },
  { title: "Nomor IHS", dataIndex: "ihs_number", key: "ihs_number", width: 150 },
  { title: "ID", dataIndex: "id", key: "id", width: 180 },
  { title: "Nama Tenaga Medis", dataIndex: "name", key: "name", width: 220 },
  { title: "Jenis Kelamin", dataIndex: "gender", key: "gender", width: 130 },
  { title: "Tanggal Lahir", dataIndex: "birth_date", key: "birth_date", width: 140 },
  { title: "Telepon", dataIndex: "phone", key: "phone", width: 150 },
  { title: "Alamat", dataIndex: "address", key: "address", width: 280 },
];

type PractitionerSearch = {
  id?: string;
  nik?: string;
  name?: string;
  gender?: string;
  birthdate?: string;
  page?: number;
  limit?: number;
};

export function PractitionersScreen() {
  const [practitioners, setPractitioners] = useState<PersonResponse[]>([]);
  const [loading, setLoading] = useState(false);
  const [result, setResult] = useState<unknown>();
  const [error, setError] = useState<string>();
  const [form] = Form.useForm<PractitionerSearch>();

  const loadPractitioners = useCallback(async () => {
    setLoading(true);
    setError(undefined);
    try {
      setPractitioners(await getLocalPractitioners());
    } catch (caught) {
      setError(formatApiError(caught));
    } finally {
      setLoading(false);
    }
  }, []);

  useEffect(() => {
    void loadPractitioners();
  }, [loadPractitioners]);

  async function handleSearch(values: PractitionerSearch) {
    setError(undefined);
    try {
      setResult(await searchPractitioners(values));
    } catch (caught) {
      setResult(undefined);
      setError(formatApiError(caught));
    }
  }

  return (
    <div className="resource-grid">
      <DataToolbar
        title="Tenaga Medis"
        description="Pantau cache tenaga medis dan cari profil praktisi di SatuSehat berdasarkan ID, NIK, atau data identitas."
        loading={loading}
        onReload={loadPractitioners}
      />
      <div className="table-wrap">
        <Table rowKey={(row) => row.nik || row.id} columns={practitionerColumns} dataSource={practitioners} loading={loading} scroll={{ x: 1280 }} />
      </div>
      <div className="panel">
        <Form form={form} layout="vertical" initialValues={{ page: 1, limit: 10 }} onFinish={handleSearch}>
          <div className="form-grid">
            <Form.Item name="id" label="ID / Nomor IHS"><Input /></Form.Item>
            <Form.Item name="nik" label="NIK"><Input /></Form.Item>
            <Form.Item name="name" label="Nama Tenaga Medis"><Input /></Form.Item>
            <Form.Item name="gender" label="Jenis Kelamin">
              <Select allowClear options={[{ value: "male", label: "Laki-laki" }, { value: "female", label: "Perempuan" }]} />
            </Form.Item>
            <Form.Item name="birthdate" label="Tanggal Lahir"><Input placeholder="1995-02-02" /></Form.Item>
            <Form.Item name="page" label="Halaman"><InputNumber min={1} className="number-input" /></Form.Item>
            <Form.Item name="limit" label="Limit"><InputNumber min={-1} className="number-input" /></Form.Item>
          </div>
          <Button type="primary" htmlType="submit">Cari Tenaga Medis</Button>
        </Form>
      </div>
      <JsonPanel title="Respons Tenaga Medis" data={result} error={error} />
    </div>
  );
}
