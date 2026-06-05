import {
  ApiOutlined,
  EnvironmentOutlined,
  HomeOutlined,
  MedicineBoxOutlined,
  TeamOutlined,
  UserOutlined,
} from "@ant-design/icons";
import { Layout, Menu, Typography } from "antd";
import type { MenuProps } from "antd";
import { useMemo, useState } from "react";
import { OverviewScreen } from "@/features/overview/OverviewScreen";
import { PatientsScreen } from "@/features/patients/PatientsScreen";
import { PractitionersScreen } from "@/features/practitioners/PractitionersScreen";
import { LocationsScreen } from "@/features/locations/LocationsScreen";
import { EncountersScreen } from "@/features/encounters/EncountersScreen";
import { ApiToolsScreen } from "@/features/tools/ApiToolsScreen";

const { Header, Content, Sider } = Layout;

type ScreenKey = "overview" | "patients" | "practitioners" | "locations" | "encounters" | "tools";

const menuItems: MenuProps["items"] = [
  { key: "overview", icon: <HomeOutlined />, label: "Ringkasan" },
  { key: "patients", icon: <UserOutlined />, label: "Pasien" },
  { key: "practitioners", icon: <TeamOutlined />, label: "Tenaga Medis" },
  { key: "locations", icon: <EnvironmentOutlined />, label: "Unit Layanan" },
  { key: "encounters", icon: <MedicineBoxOutlined />, label: "Kunjungan" },
  { key: "tools", icon: <ApiOutlined />, label: "Alat Integrasi" },
];

function EmptyScreen({ name }: { name: string }) {
  return <Typography.Title level={2}>{name}</Typography.Title>;
}

export function AdminShell() {
  const [activeScreen, setActiveScreen] = useState<ScreenKey>("overview");

  const title = useMemo(() => {
    const active = menuItems?.find((item) => item && "key" in item && item.key === activeScreen);
    return active && "label" in active ? String(active.label) : "Ringkasan";
  }, [activeScreen]);

  const screen = useMemo(() => {
    switch (activeScreen) {
      case "overview":
        return <OverviewScreen />;
      case "patients":
        return <PatientsScreen />;
      case "practitioners":
        return <PractitionersScreen />;
      case "locations":
        return <LocationsScreen />;
      case "encounters":
        return <EncountersScreen />;
      case "tools":
        return <ApiToolsScreen />;
      default:
        return <EmptyScreen name="Ringkasan" />;
    }
  }, [activeScreen]);

  return (
    <Layout className="admin-layout">
      <Sider className="admin-sider" width={248} breakpoint="lg" collapsedWidth={0}>
        <div className="brand-block">
          <Typography.Title level={4}>SatuSehat Intero</Typography.Title>
          <Typography.Text>Dashboard Integrasi Rumah Sakit</Typography.Text>
        </div>
        <Menu
          mode="inline"
          selectedKeys={[activeScreen]}
          items={menuItems}
          onClick={({ key }) => setActiveScreen(key as ScreenKey)}
        />
      </Sider>
      <Layout>
        <Header className="admin-header">
          <div>
            <Typography.Title level={2}>{title}</Typography.Title>
            <Typography.Text type="secondary">
              Konsol operasional untuk sinkronisasi layanan rumah sakit dengan SatuSehat. API: {process.env.NEXT_PUBLIC_API_BASE_URL ?? "http://localhost:8083/api"}
            </Typography.Text>
          </div>
        </Header>
        <Content className="admin-content">
          <main className="screen-stack">{screen}</main>
        </Content>
      </Layout>
    </Layout>
  );
}
