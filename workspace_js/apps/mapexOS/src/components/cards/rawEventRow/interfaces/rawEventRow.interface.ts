export interface RawEventProps {
  id: string;
  asset: {
    name: string
    description: string
    icon: string
    type: string
  };
  type: string;
  status: 'high' | 'medium' | 'low';
  protocol: string;
  created: string;
  values: Record<string, any>;
}