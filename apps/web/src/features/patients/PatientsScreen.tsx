import { Button, Form, Input, Select, Space, Table, Tabs } from "antd";
import type { ColumnsType } from "antd/es/table";
import { useCallback, useEffect, useState } from "react";
import { DataToolbar } from "@/components/DataToolbar";
import { JsonPanel } from "@/components/JsonPanel";
import {
  createPatient,
  formatApiError,
  getLocalPatients,
  searchPatient,
  type CreatePatientRequest,
  type PersonResponse,
} from "@/lib/api";

const personColumns: ColumnsType<PersonResponse> = [
  { title: "NIK", dataIndex: "nik", key: "nik", width: 170 },
  { title: "Nomor IHS", dataIndex: "ihs_number", key: "ihs_number", width: 150 },
  { title: "ID", dataIndex: "id", key: "id", width: 180 },
  { title: "Nama Pasien", dataIndex: "name", key: "name", width: 220 },
  { title: "Jenis Kelamin", dataIndex: "gender", key: "gender", width: 130 },
  { title: "Tanggal Lahir", dataIndex: "birth_date", key: "birth_date", width: 140 },
  { title: "Telepon", dataIndex: "phone", key: "phone", width: 150 },
  { title: "Alamat", dataIndex: "address", key: "address", width: 280 },
];

export function PatientsScreen() {
  const [patients, setPatients] = useState<PersonResponse[]>([]);
  const [loading, setLoading] = useState(false);
  const [result, setResult] = useState<unknown>();
  const [error, setError] = useState<string>();
  const [searchForm] = Form.useForm();
  const [createForm] = Form.useForm<CreatePatientRequest>();

  const loadPatients = useCallback(async () => {
    setLoading(true);
    setError(undefined);
    try {
      setPatients(await getLocalPatients());
    } catch (caught) {
      setError(formatApiError(caught));
    } finally {
      setLoading(false);
    }
  }, []);

  useEffect(() => {
    void loadPatients();
  }, [loadPatients]);

  async function handleSearch(values: { nik: string; name?: string }) {
    setError(undefined);
    try {
      const data = await searchPatient(values);
      setResult(data);
    } catch (caught) {
      setResult(undefined);
      setError(formatApiError(caught));
    }
  }

  async function handleCreate(values: CreatePatientRequest) {
    setError(undefined);
    try {
      const data = await createPatient(values);
      setResult(data);
      createForm.resetFields();
      await loadPatients();
    } catch (caught) {
      setResult(undefined);
      setError(formatApiError(caught));
    }
  }

  return (
    <div className="resource-grid">
      <DataToolbar
        title="Data Pasien"
        description="Kelola data pasien lokal, lakukan pencarian NIK ke SatuSehat, dan daftarkan pasien baru untuk proses layanan rumah sakit."
        loading={loading}
        onReload={loadPatients}
      />
      <div className="table-wrap">
        <Table rowKey={(row) => row.nik || row.id} columns={personColumns} dataSource={patients} loading={loading} scroll={{ x: 1280 }} />
      </div>
      <Tabs
        items={[
          {
            key: "search",
            label: "Cari Pasien",
            children: (
              <div className="panel">
                <Form form={searchForm} layout="vertical" onFinish={handleSearch}>
                  <Form.Item name="nik" label="NIK" rules={[{ required: true, message: "NIK wajib diisi" }]}>
                    <Input placeholder="9104025209000006" />
                  </Form.Item>
                  <Form.Item name="name" label="Nama Pasien">
                    <Input placeholder="Nama pasien untuk pencarian SatuSehat" />
                  </Form.Item>
                  <Button type="primary" htmlType="submit">Cari Pasien</Button>
                </Form>
              </div>
            ),
          },
          {
            key: "create",
            label: "Daftarkan Pasien",
            children: (
              <div className="panel">
                <Form form={createForm} layout="vertical" onFinish={handleCreate}>
                  <Space direction="vertical" size={0} className="form-stack">
                    <Form.Item name="nik" label="NIK" rules={[{ required: true, message: "NIK wajib diisi" }]}><Input /></Form.Item>
                    <Form.Item name="name" label="Nama Pasien" rules={[{ required: true, message: "Nama pasien wajib diisi" }]}><Input /></Form.Item>
                    <Form.Item name="gender" label="Jenis Kelamin" rules={[{ required: true, message: "Jenis kelamin wajib diisi" }]}>
                      <Select options={[{ value: "male", label: "Laki-laki" }, { value: "female", label: "Perempuan" }]} />
                    </Form.Item>
                    <Form.Item name="birth_date" label="Tanggal Lahir" rules={[{ required: true, message: "Tanggal lahir wajib diisi" }]}><Input placeholder="1945-08-19" /></Form.Item>
                    <Form.Item name="phone" label="Telepon"><Input /></Form.Item>
                    <Form.Item name="address" label="Alamat"><Input.TextArea rows={2} /></Form.Item>
                    <Form.Item name="city" label="Kota/Kabupaten"><Input /></Form.Item>
                    <Form.Item name="province_code" label="Kode Provinsi"><Input /></Form.Item>
                    <Form.Item name="city_code" label="Kode Kota/Kabupaten"><Input /></Form.Item>
                    <Form.Item name="district_code" label="Kode Kecamatan"><Input /></Form.Item>
                    <Form.Item name="village_code" label="Kode Kelurahan/Desa"><Input /></Form.Item>
                    <Form.Item name="rt" label="RT"><Input /></Form.Item>
                    <Form.Item name="rw" label="RW"><Input /></Form.Item>
                    <Form.Item name="postal_code" label="Kode Pos"><Input /></Form.Item>
                  </Space>
                  <Button type="primary" htmlType="submit">Simpan Pasien</Button>
                </Form>
              </div>
            ),
          },
        ]}
      />
      <JsonPanel title="Respons Data Pasien" data={result} error={error} />
    </div>
  );
}
