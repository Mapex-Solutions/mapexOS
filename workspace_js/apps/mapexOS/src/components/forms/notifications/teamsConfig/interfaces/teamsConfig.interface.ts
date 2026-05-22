export interface ChannelTeamsProps {
  teamName: string;
  channelsName: string[];
  webhookUrl: string;
  messageTemplate: string;
  adaptiveCard: boolean;
}