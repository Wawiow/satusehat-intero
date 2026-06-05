import { ReloadOutlined } from "@ant-design/icons";
import { Button, Space, Typography } from "antd";
import type { ReactNode } from "react";

type DataToolbarProps = {
  title: string;
  description?: string;
  loading?: boolean;
  onReload?: () => void;
  extra?: ReactNode;
};

export function DataToolbar({ title, description, loading, onReload, extra }: DataToolbarProps) {
  return (
    <div className="data-toolbar">
      <div>
        <Typography.Title level={3}>{title}</Typography.Title>
        {description ? <Typography.Text type="secondary">{description}</Typography.Text> : null}
      </div>
      <Space wrap>
        {extra}
        {onReload ? (
          <Button icon={<ReloadOutlined />} loading={loading} onClick={onReload}>
            Muat Ulang
          </Button>
        ) : null}
      </Space>
    </div>
  );
}
