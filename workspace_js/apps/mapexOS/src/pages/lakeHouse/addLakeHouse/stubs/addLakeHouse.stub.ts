export const LAKE_HOUSE_LIST_STUB = [
  {
    type: 'slack' as const,
    name: 'Slack',
    description: 'Envia alertas e atualizações em tempo real para o canal Slack da equipe.',
    icon: 'mdi-slack',
    status: 'Active',
    data: {
      workspace: 'Acme Corp',
      channelsName: ['#alerts'],
      webhookUrl: 'https://hooks.slack.com/services/T00000000/B00000000/XXXXXXXXXXXXXXXXXXXXXXXX',
      messageTemplate: 'Alert: {{title}} - {{message}}',
      botName: 'Notification Bot'
    },
    created: new Date().toDateString()
  },
  {
    type: 'teams' as const,
    name: 'Microsoft Teams',
    description: 'Publica notificações diretamente no canal do Microsoft Teams.',
    icon: 'mdi-microsoft-teams',
    status: 'Active',
    data: {
      teamName: 'Acme Corp Team',
      channelsName: ['Alerts Channel', 'Alerts Channel 2', 'Alerts Channel 3'],
      webhookUrl: 'https://outlook.office.com/webhook/00000000-0000-0000-0000-000000000000@00000000-0000-0000-0000-000000000000/IncomingWebhook/00000000000000000000000000000000/00000000-0000-0000-0000-000000000000',
      messageTemplate: '{"title":"{{title}}","text":"{{message}}"}',
      adaptiveCard: false
    },
    created: new Date().toDateString()
  },
  {
    type: 'email' as const,
    name: 'Email',
    description: 'Dispara mensagens por e-mail para a lista de destinatários configurada.',
    icon: 'mdi-email',
    status: 'Active',
    data: {
      from: 'bot@acme.com',
      subject: 'Alerta de Sistema',
      to: ['ops@acme.com', 'dev@acme.com', 'ceo@acme.com', 'ops@acme.com', 'dev@acme.com', 'ceo@acme.com'],
      cc: [],
      bcc: [],
      template: '<h1>{{title}}</h1><p>{{message}}</p>',
      smtp: {
        host: 'smtp.acme.com',
        port: 587,
        secure: true,
        auth: {
          user: 'notifications@acme.com',
          password: '********'
        }
      },
      attachments: true
    },
    created: new Date().toDateString()
  },
  {
    type: 'push' as const,
    name: 'Push Notifications',
    description: 'Envia notificações push para dispositivos móveis cadastrados.',
    icon: 'mdi-bell-ring',
    status: 'Active',
    data: {
      appName: 'Acme Mobile App',
      deviceCount: 1245,
      apiKey: 'xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx',
      serviceProvider: 'firebase',
      priority: 'high',
      ttl: 3600,
      badge: true,
      sound: 'default',
      clickAction: 'OPEN_APP'
    },
    created: new Date().toDateString(),
    updated: new Date().toDateString()
  },
  {
    type: 'telegram' as const,
    name: 'Telegram',
    description: 'Envia mensagens via bot do Telegram para chats configurados.',
    icon: 'mdi-send',
    status: 'Active',
    data: {
      botName: '@AlertsBot',
      chatNames: ['Equipe Desenvolvimento', 'Alertas Críticos', 'Suporte Técnico', 'Gerência', 'QA Team', 'DevOps'],
      botToken: '0000000000:AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA',
      parseMode: 'Markdown',
      disableNotification: false,
      messageTemplate: '*{{title}}*\n{{message}}'
    },
    created: new Date().toDateString()
  },
  {
    type: 'webhook' as const,
    name: 'Webhook',
    description: 'Envia requisições HTTP para endpoints externos quando eventos ocorrem.',
    icon: 'mdi-webhook',
    status: 'Active',
    data: {
      name: 'Sistema de Alertas',
      method: 'POST',
      url: 'https://api.empresa.com/webhooks/alerts',
      headers: {
        'Content-Type': 'application/json',
        'Authorization': 'Bearer xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx'
      },
      payload: '{"title":"{{title}}","message":"{{message}}","timestamp":"{{timestamp}}"}',
      timeout: 30000,
      retryCount: 3
    },
    created: new Date().toDateString()
  },
  {
    type: 'sms' as const,
    name: 'SMS',
    description: 'Envia mensagens SMS para números de telefone cadastrados.',
    icon: 'mdi-message-text',
    status: 'Active',
    data: {
      provider: 'twilio',
      accountSid: 'ACxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx',
      authToken: 'xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx',
      from: '+1234567890',
      to: ['+1234567890', '+0987654321'],
      messageTemplate: '{{title}}: {{message}}',
      maxLength: 160
    },
    created: new Date().toDateString()
  },
];
