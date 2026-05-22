export interface ChannelEmailProps {
  from: string;
  subject: string;
  to: string[];
  cc: string[];
  bcc: string[];
  template: string;
  smtp: {
    host: string;
    port: number;
    secure: boolean;
    auth: {
      user: string;
      password: string;
    };
  };
  attachments: boolean;
}