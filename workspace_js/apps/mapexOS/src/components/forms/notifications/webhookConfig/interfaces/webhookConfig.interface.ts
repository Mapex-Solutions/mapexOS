export interface ChannelWebhookProps {
  name: string;
  method: string;
  url: string;
  headers: Record<string, string>;
  payload: string;
  timeout: number;
  retryCount: number;
}