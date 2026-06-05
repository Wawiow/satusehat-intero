import { Button, Col, Row, Space, Statistic, Typography } from "antd";
import { useCallback, useEffect, useState } from "react";
import { DataToolbar } from "@/components/DataToolbar";
import { JsonPanel } from "@/components/JsonPanel";
import {
  formatApiError,
  getLocalEncounters,
  getLocalLocations,
  getLocalPatients,
  getLocalPractitioners,
  getToken,
} from "@/lib/api";

type Counts = {
  patients: number;
  practitioners: number;
  locations: number;
  encounters: number;
};

const emptyCounts: Counts = { patients: 0, practitioners: 0, locations: 0, encounters: 0 };

export function OverviewScreen() {
  const [counts, setCounts] = useState<Counts>(emptyCounts);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string>();
  const [tokenResult, setTokenResult] = useState<unknown>();
  const [tokenError, setTokenError] = useState<string>();

  const loadCounts = useCallback(async () => {
    setLoading(true);
    setError(undefined);
    try {
      const [patients, practitioners, locations, encounters] = await Promise.all([
        getLocalPatients(),
        getLocalPractitioners(),
        getLocalLocations(),
        getLocalEncounters(),
      ]);
      setCounts({
        patients: patients.length,
        practitioners: practitioners.length,
        locations: locations.length,
        encounters: encounters.length,
      });
    } catch (caught) {
      setError(formatApiError(caught));
    } finally {
      setLoading(false);
    }
  }, []);

  useEffect(() => {
    void loadCounts();
  }, [loadCounts]);

  async function handleToken() {
    setTokenError(undefined);
    try {
      setTokenResult(await getToken());
    } catch (caught) {
      setTokenResult(undefined);
      setTokenError(formatApiError(caught));
    }
  }

  return (
    <div className="resource-grid">
      <DataToolbar
        title="Ringkasan Operasional"
        description="Pantau kesiapan koneksi, data master lokal, dan akses token SatuSehat untuk layanan rumah sakit."
        loading={loading}
        onReload={loadCounts}
        extra={<Button onClick={handleToken}>Cek Token SatuSehat</Button>}
      />
      {error ? <JsonPanel title="Gagal Memuat Ringkasan" data={null} error={error} /> : null}
      <div className="panel">
        <Row gutter={[16, 16]}>
          <Col xs={12} lg={6}>
            <Statistic title="Pasien Terdaftar" value={counts.patients} loading={loading} />
          </Col>
          <Col xs={12} lg={6}>
            <Statistic title="Tenaga Medis" value={counts.practitioners} loading={loading} />
          </Col>
          <Col xs={12} lg={6}>
            <Statistic title="Unit Layanan" value={counts.locations} loading={loading} />
          </Col>
          <Col xs={12} lg={6}>
            <Statistic title="Kunjungan" value={counts.encounters} loading={loading} />
          </Col>
        </Row>
      </div>
      <Space direction="vertical" size={6}>
        <Typography.Text strong>Catatan operasional</Typography.Text>
        <Typography.Text type="secondary">
          Endpoint SatuSehat dapat gagal bila kredensial, organisasi, atau data upstream belum tersedia. Gunakan ringkasan ini sebagai pemeriksaan awal sebelum layanan digunakan operator.
        </Typography.Text>
      </Space>
      <JsonPanel title="Respons Token SatuSehat" data={tokenResult} error={tokenError} />
    </div>
  );
}
