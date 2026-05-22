export interface ProtocolConfig {
  type: 'http' | 'mqtt' | 'lorawan';
  http?: any;
  mqtt?: {
    clientId: string;
    username: string;
    password: string;
  };
}

export interface AssetRowProps {
  id: string;
  name: string;
  enabled: boolean;
  description?: string;
  assetUUID: string;
  assetTemplateId?: string;
  category?: string;
  assetType?: string;
  protocol?: ProtocolConfig;
  routeGroupIds?: string[];
  latitude?: number;
  longitude?: number;
  created?: string;
  updated?: string;

  // For display purposes
  templateName?: string;
  categoryName?: string;
  assetTypeName?: string;
  icon?: string;
  iconColor?: string;
}
