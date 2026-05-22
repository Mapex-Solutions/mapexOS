export interface ChannelTelegramProps {
  botName: string;
  chatNames: string[];
  botToken: string;
  parseMode: string;
  disableNotification: boolean;
  messageTemplate: string;
}