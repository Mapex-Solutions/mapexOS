export interface ChannelSlackProps {
  workspace: string;
  channelsName: string[];
  webhookUrl: string;
  messageTemplate: string;
  botName: string;
}