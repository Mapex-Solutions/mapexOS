import { computed } from 'vue';
import { useTS } from '@utils/translation';

/**
 * Custom composable for CreateEditTrigger page translations.
 * Provides type-safe reactive access to all translated strings with formatting utilities.
 *
 * Structure mirrors:
 * - File: src/pages/automations/triggers/createEditTriggerPage/CreateEditTriggerPage.vue
 * - JSON: src/i18n/{locale}/pages/automations/createEditTrigger.json
 * - Composable: src/composables/i18n/pages/automations/triggers/createEditTrigger/useCreateEditTriggerTranslations.ts
 */
export function useCreateEditTriggerTranslations() {
  const ts = useTS({ capitalize: true });
  const tsTitle = useTS({ titleCase: true });
  const tsRaw = useTS({ capitalize: false });

  return {
    /**
     * Page header translations
     */
    page: {
      title: computed(() => ts('pages.automations.createEditTrigger.page.title')),
      titleEdit: computed(() => ts('pages.automations.createEditTrigger.page.titleEdit')),
      description: computed(() => tsRaw('pages.automations.createEditTrigger.page.description')),
      button: computed(() => ts('pages.automations.createEditTrigger.page.button')),
    },

    /**
     * Stepper header translations
     */
    stepper: {
      title: computed(() => tsTitle('pages.automations.createEditTrigger.stepper.title')),
      subtitle: computed(() => tsRaw('pages.automations.createEditTrigger.stepper.subtitle')),
      requiredInfo: computed(() => tsRaw('pages.automations.createEditTrigger.stepper.requiredInfo')),
      currentStep: computed(() => ts('pages.automations.createEditTrigger.stepper.currentStep')),
    },

    /**
     * Step-specific translations
     */
    steps: {
      /**
       * Step 1: Select Category
       */
      step1: {
        label: computed(() => ts('pages.automations.createEditTrigger.steps.step1.label')),
        description: computed(() => tsRaw('pages.automations.createEditTrigger.steps.step1.description')),
        intro: computed(() => tsRaw('pages.automations.createEditTrigger.steps.step1.intro')),
        tip: {
          prefix: computed(() => tsRaw('pages.automations.createEditTrigger.steps.step1.tip.prefix')),
          text: computed(() => tsRaw('pages.automations.createEditTrigger.steps.step1.tip.text')),
        },
      },

      /**
       * Step 2: Select Type
       */
      step2: {
        label: computed(() => ts('pages.automations.createEditTrigger.steps.step2.label')),
        description: computed(() => tsRaw('pages.automations.createEditTrigger.steps.step2.description')),
        introTechnical: computed(() => tsRaw('pages.automations.createEditTrigger.steps.step2.introTechnical')),
        introCommunication: computed(() => tsRaw('pages.automations.createEditTrigger.steps.step2.introCommunication')),
        note: {
          prefix: computed(() => tsRaw('pages.automations.createEditTrigger.steps.step2.note.prefix')),
          text: computed(() => tsRaw('pages.automations.createEditTrigger.steps.step2.note.text')),
        },
      },

      /**
       * Step 3: Basic Info
       */
      step3: {
        label: computed(() => ts('pages.automations.createEditTrigger.steps.step3.label')),
        description: computed(() => tsRaw('pages.automations.createEditTrigger.steps.step3.description')),
        /**
         * Intro paragraph displayed above the form
         */
        intro: computed(() => tsRaw('pages.automations.createEditTrigger.steps.step3.intro')),
        /**
         * Field-level labels, placeholders, hints, and validation messages
         * Mirrors: Step3BasicInfo.vue
         */
        fields: {
          nameLabel: computed(() => ts('pages.automations.createEditTrigger.steps.step3.fields.nameLabel')),
          namePlaceholder: computed(() => tsRaw('pages.automations.createEditTrigger.steps.step3.fields.namePlaceholder')),
          nameHint: computed(() => tsRaw('pages.automations.createEditTrigger.steps.step3.fields.nameHint')),
          nameRequired: computed(() => tsRaw('pages.automations.createEditTrigger.steps.step3.fields.nameRequired')),
          nameMinLength: computed(() => tsRaw('pages.automations.createEditTrigger.steps.step3.fields.nameMinLength')),
          statusLabel: computed(() => ts('pages.automations.createEditTrigger.steps.step3.fields.statusLabel')),
          descriptionLabel: computed(() => ts('pages.automations.createEditTrigger.steps.step3.fields.descriptionLabel')),
          descriptionPlaceholder: computed(() => tsRaw('pages.automations.createEditTrigger.steps.step3.fields.descriptionPlaceholder')),
          descriptionHint: computed(() => tsRaw('pages.automations.createEditTrigger.steps.step3.fields.descriptionHint')),
          sharedTemplateLabel: computed(() => ts('pages.automations.createEditTrigger.steps.step3.fields.sharedTemplateLabel')),
          sharedTemplateHint: computed(() => tsRaw('pages.automations.createEditTrigger.steps.step3.fields.sharedTemplateHint')),
        },
        /**
         * Tip banner displayed below the form
         */
        tip: {
          prefix: computed(() => ts('pages.automations.createEditTrigger.steps.step3.tip.prefix')),
          text: computed(() => tsRaw('pages.automations.createEditTrigger.steps.step3.tip.text')),
        },
      },

      /**
       * Step 4: Configuration
       */
      step4: {
        label: computed(() => ts('pages.automations.createEditTrigger.steps.step4.label')),
        description: computed(() => tsRaw('pages.automations.createEditTrigger.steps.step4.description')),
        comingSoon: computed(() => tsRaw('pages.automations.createEditTrigger.steps.step4.comingSoon')),
        comingSoonDetail: computed(() => tsRaw('pages.automations.createEditTrigger.steps.step4.comingSoonDetail')),
      },

      /**
       * Step 5: Review
       */
      step5: {
        label: computed(() => ts('pages.automations.createEditTrigger.steps.step5.label')),
        description: computed(() => tsRaw('pages.automations.createEditTrigger.steps.step5.description')),
        subtitle: computed(() => tsRaw('pages.automations.createEditTrigger.steps.step5.subtitle')),
        successMessage: computed(() => tsRaw('pages.automations.createEditTrigger.steps.step5.successMessage')),
        sections: {
          categoryType: computed(() => ts('pages.automations.createEditTrigger.steps.step5.sections.categoryType')),
          basicInfo: computed(() => ts('pages.automations.createEditTrigger.steps.step5.sections.basicInfo')),
          configuration: computed(() => ts('pages.automations.createEditTrigger.steps.step5.sections.configuration')),
        },
        fields: {
          category: computed(() => ts('pages.automations.createEditTrigger.steps.step5.fields.category')),
          triggerType: computed(() => ts('pages.automations.createEditTrigger.steps.step5.fields.triggerType')),
          name: computed(() => ts('pages.automations.createEditTrigger.steps.step5.fields.name')),
          description: computed(() => ts('pages.automations.createEditTrigger.steps.step5.fields.description')),
          status: computed(() => ts('pages.automations.createEditTrigger.steps.step5.fields.status')),
          configuration: computed(() => ts('pages.automations.createEditTrigger.steps.step5.fields.configuration')),
        },
      },
    },

    /**
     * Step 4 protocol-specific configuration translations
     * Mirrors: Step4Configuration/configs/{Http,Mqtt,Nats,Rabbitmq,Slack,Email,Teams,Websocket}Config.vue
     */
    step4Configs: {
      /**
       * HTTP trigger configuration
       */
      http: {
        endpointLabel: computed(() => ts('pages.automations.createEditTrigger.step4Configs.http.endpointLabel')),
        endpointPlaceholder: computed(() => tsRaw('pages.automations.createEditTrigger.step4Configs.http.endpointPlaceholder')),
        endpointHint: computed(() => tsRaw('pages.automations.createEditTrigger.step4Configs.http.endpointHint')),
        endpointRequired: computed(() => tsRaw('pages.automations.createEditTrigger.step4Configs.http.endpointRequired')),
        methodLabel: computed(() => ts('pages.automations.createEditTrigger.step4Configs.http.methodLabel')),
        timeoutLabel: computed(() => ts('pages.automations.createEditTrigger.step4Configs.http.timeoutLabel')),
        timeoutPlaceholder: computed(() => tsRaw('pages.automations.createEditTrigger.step4Configs.http.timeoutPlaceholder')),
        timeoutHint: computed(() => tsRaw('pages.automations.createEditTrigger.step4Configs.http.timeoutHint')),
        headersTitle: computed(() => ts('pages.automations.createEditTrigger.step4Configs.http.headersTitle')),
        addHeaderButton: computed(() => ts('pages.automations.createEditTrigger.step4Configs.http.addHeaderButton')),
        headerNamePlaceholder: computed(() => tsRaw('pages.automations.createEditTrigger.step4Configs.http.headerNamePlaceholder')),
        headerValuePlaceholder: computed(() => tsRaw('pages.automations.createEditTrigger.step4Configs.http.headerValuePlaceholder')),
        noHeadersMessage: computed(() => tsRaw('pages.automations.createEditTrigger.step4Configs.http.noHeadersMessage')),
        bodyTitle: computed(() => ts('pages.automations.createEditTrigger.step4Configs.http.bodyTitle')),
        formatJsonButton: computed(() => ts('pages.automations.createEditTrigger.step4Configs.http.formatJsonButton')),
        bodyPlaceholder: computed(() => tsRaw('pages.automations.createEditTrigger.step4Configs.http.bodyPlaceholder')),
        bodyHint: computed(() => tsRaw('pages.automations.createEditTrigger.step4Configs.http.bodyHint')),
        invalidJson: computed(() => tsRaw('pages.automations.createEditTrigger.step4Configs.http.invalidJson')),
      },

      /**
       * MQTT trigger configuration
       */
      mqtt: {
        brokerLabel: computed(() => ts('pages.automations.createEditTrigger.step4Configs.mqtt.brokerLabel')),
        brokerPlaceholder: computed(() => tsRaw('pages.automations.createEditTrigger.step4Configs.mqtt.brokerPlaceholder')),
        brokerHint: computed(() => tsRaw('pages.automations.createEditTrigger.step4Configs.mqtt.brokerHint')),
        brokerRequired: computed(() => tsRaw('pages.automations.createEditTrigger.step4Configs.mqtt.brokerRequired')),
        portLabel: computed(() => ts('pages.automations.createEditTrigger.step4Configs.mqtt.portLabel')),
        portPlaceholder: computed(() => tsRaw('pages.automations.createEditTrigger.step4Configs.mqtt.portPlaceholder')),
        portHint: computed(() => tsRaw('pages.automations.createEditTrigger.step4Configs.mqtt.portHint')),
        portRequired: computed(() => tsRaw('pages.automations.createEditTrigger.step4Configs.mqtt.portRequired')),
        topicLabel: computed(() => ts('pages.automations.createEditTrigger.step4Configs.mqtt.topicLabel')),
        topicPlaceholder: computed(() => tsRaw('pages.automations.createEditTrigger.step4Configs.mqtt.topicPlaceholder')),
        topicHint: computed(() => tsRaw('pages.automations.createEditTrigger.step4Configs.mqtt.topicHint')),
        topicRequired: computed(() => tsRaw('pages.automations.createEditTrigger.step4Configs.mqtt.topicRequired')),
        qosLabel: computed(() => ts('pages.automations.createEditTrigger.step4Configs.mqtt.qosLabel')),
        clientIdLabel: computed(() => ts('pages.automations.createEditTrigger.step4Configs.mqtt.clientIdLabel')),
        clientIdPlaceholder: computed(() => tsRaw('pages.automations.createEditTrigger.step4Configs.mqtt.clientIdPlaceholder')),
        clientIdHint: computed(() => tsRaw('pages.automations.createEditTrigger.step4Configs.mqtt.clientIdHint')),
        useTlsLabel: computed(() => ts('pages.automations.createEditTrigger.step4Configs.mqtt.useTlsLabel')),
        usernameLabel: computed(() => ts('pages.automations.createEditTrigger.step4Configs.mqtt.usernameLabel')),
        usernamePlaceholder: computed(() => tsRaw('pages.automations.createEditTrigger.step4Configs.mqtt.usernamePlaceholder')),
        usernameHint: computed(() => tsRaw('pages.automations.createEditTrigger.step4Configs.mqtt.usernameHint')),
        passwordLabel: computed(() => ts('pages.automations.createEditTrigger.step4Configs.mqtt.passwordLabel')),
        passwordPlaceholder: computed(() => tsRaw('pages.automations.createEditTrigger.step4Configs.mqtt.passwordPlaceholder')),
        passwordHint: computed(() => tsRaw('pages.automations.createEditTrigger.step4Configs.mqtt.passwordHint')),
        messageTitle: computed(() => ts('pages.automations.createEditTrigger.step4Configs.mqtt.messageTitle')),
        formatJsonButton: computed(() => ts('pages.automations.createEditTrigger.step4Configs.mqtt.formatJsonButton')),
        messagePlaceholder: computed(() => tsRaw('pages.automations.createEditTrigger.step4Configs.mqtt.messagePlaceholder')),
        messageHint: computed(() => tsRaw('pages.automations.createEditTrigger.step4Configs.mqtt.messageHint')),
      },

      /**
       * NATS trigger configuration
       */
      nats: {
        serverLabel: computed(() => ts('pages.automations.createEditTrigger.step4Configs.nats.serverLabel')),
        serverPlaceholder: computed(() => tsRaw('pages.automations.createEditTrigger.step4Configs.nats.serverPlaceholder')),
        serverHint: computed(() => tsRaw('pages.automations.createEditTrigger.step4Configs.nats.serverHint')),
        serverRequired: computed(() => tsRaw('pages.automations.createEditTrigger.step4Configs.nats.serverRequired')),
        subjectLabel: computed(() => ts('pages.automations.createEditTrigger.step4Configs.nats.subjectLabel')),
        subjectPlaceholder: computed(() => tsRaw('pages.automations.createEditTrigger.step4Configs.nats.subjectPlaceholder')),
        subjectHint: computed(() => tsRaw('pages.automations.createEditTrigger.step4Configs.nats.subjectHint')),
        subjectRequired: computed(() => tsRaw('pages.automations.createEditTrigger.step4Configs.nats.subjectRequired')),
        useTlsLabel: computed(() => ts('pages.automations.createEditTrigger.step4Configs.nats.useTlsLabel')),
        usernameLabel: computed(() => ts('pages.automations.createEditTrigger.step4Configs.nats.usernameLabel')),
        usernamePlaceholder: computed(() => tsRaw('pages.automations.createEditTrigger.step4Configs.nats.usernamePlaceholder')),
        usernameHint: computed(() => tsRaw('pages.automations.createEditTrigger.step4Configs.nats.usernameHint')),
        passwordLabel: computed(() => ts('pages.automations.createEditTrigger.step4Configs.nats.passwordLabel')),
        passwordPlaceholder: computed(() => tsRaw('pages.automations.createEditTrigger.step4Configs.nats.passwordPlaceholder')),
        passwordHint: computed(() => tsRaw('pages.automations.createEditTrigger.step4Configs.nats.passwordHint')),
        tokenLabel: computed(() => ts('pages.automations.createEditTrigger.step4Configs.nats.tokenLabel')),
        tokenPlaceholder: computed(() => tsRaw('pages.automations.createEditTrigger.step4Configs.nats.tokenPlaceholder')),
        tokenHint: computed(() => tsRaw('pages.automations.createEditTrigger.step4Configs.nats.tokenHint')),
        messageTitle: computed(() => ts('pages.automations.createEditTrigger.step4Configs.nats.messageTitle')),
        formatJsonButton: computed(() => ts('pages.automations.createEditTrigger.step4Configs.nats.formatJsonButton')),
        messagePlaceholder: computed(() => tsRaw('pages.automations.createEditTrigger.step4Configs.nats.messagePlaceholder')),
        messageHint: computed(() => tsRaw('pages.automations.createEditTrigger.step4Configs.nats.messageHint')),
      },

      /**
       * RabbitMQ trigger configuration
       */
      rabbitmq: {
        hostLabel: computed(() => ts('pages.automations.createEditTrigger.step4Configs.rabbitmq.hostLabel')),
        hostPlaceholder: computed(() => tsRaw('pages.automations.createEditTrigger.step4Configs.rabbitmq.hostPlaceholder')),
        hostHint: computed(() => tsRaw('pages.automations.createEditTrigger.step4Configs.rabbitmq.hostHint')),
        hostRequired: computed(() => tsRaw('pages.automations.createEditTrigger.step4Configs.rabbitmq.hostRequired')),
        portLabel: computed(() => ts('pages.automations.createEditTrigger.step4Configs.rabbitmq.portLabel')),
        portPlaceholder: computed(() => tsRaw('pages.automations.createEditTrigger.step4Configs.rabbitmq.portPlaceholder')),
        portHint: computed(() => tsRaw('pages.automations.createEditTrigger.step4Configs.rabbitmq.portHint')),
        portRequired: computed(() => tsRaw('pages.automations.createEditTrigger.step4Configs.rabbitmq.portRequired')),
        vhostLabel: computed(() => ts('pages.automations.createEditTrigger.step4Configs.rabbitmq.vhostLabel')),
        vhostPlaceholder: computed(() => tsRaw('pages.automations.createEditTrigger.step4Configs.rabbitmq.vhostPlaceholder')),
        vhostHint: computed(() => tsRaw('pages.automations.createEditTrigger.step4Configs.rabbitmq.vhostHint')),
        useTlsLabel: computed(() => ts('pages.automations.createEditTrigger.step4Configs.rabbitmq.useTlsLabel')),
        usernameLabel: computed(() => ts('pages.automations.createEditTrigger.step4Configs.rabbitmq.usernameLabel')),
        usernamePlaceholder: computed(() => tsRaw('pages.automations.createEditTrigger.step4Configs.rabbitmq.usernamePlaceholder')),
        usernameHint: computed(() => tsRaw('pages.automations.createEditTrigger.step4Configs.rabbitmq.usernameHint')),
        usernameRequired: computed(() => tsRaw('pages.automations.createEditTrigger.step4Configs.rabbitmq.usernameRequired')),
        passwordLabel: computed(() => ts('pages.automations.createEditTrigger.step4Configs.rabbitmq.passwordLabel')),
        passwordPlaceholder: computed(() => tsRaw('pages.automations.createEditTrigger.step4Configs.rabbitmq.passwordPlaceholder')),
        passwordHint: computed(() => tsRaw('pages.automations.createEditTrigger.step4Configs.rabbitmq.passwordHint')),
        passwordRequired: computed(() => tsRaw('pages.automations.createEditTrigger.step4Configs.rabbitmq.passwordRequired')),
        publishingTitle: computed(() => ts('pages.automations.createEditTrigger.step4Configs.rabbitmq.publishingTitle')),
        publishModeLabel: computed(() => ts('pages.automations.createEditTrigger.step4Configs.rabbitmq.publishModeLabel')),
        publishModeExchange: computed(() => ts('pages.automations.createEditTrigger.step4Configs.rabbitmq.publishModeExchange')),
        publishModeExchangeDescription: computed(() => tsRaw('pages.automations.createEditTrigger.step4Configs.rabbitmq.publishModeExchangeDescription')),
        publishModeQueue: computed(() => ts('pages.automations.createEditTrigger.step4Configs.rabbitmq.publishModeQueue')),
        publishModeQueueDescription: computed(() => tsRaw('pages.automations.createEditTrigger.step4Configs.rabbitmq.publishModeQueueDescription')),
        exchangeLabel: computed(() => ts('pages.automations.createEditTrigger.step4Configs.rabbitmq.exchangeLabel')),
        exchangePlaceholder: computed(() => tsRaw('pages.automations.createEditTrigger.step4Configs.rabbitmq.exchangePlaceholder')),
        exchangeHint: computed(() => tsRaw('pages.automations.createEditTrigger.step4Configs.rabbitmq.exchangeHint')),
        exchangeRequired: computed(() => tsRaw('pages.automations.createEditTrigger.step4Configs.rabbitmq.exchangeRequired')),
        exchangeTypeLabel: computed(() => ts('pages.automations.createEditTrigger.step4Configs.rabbitmq.exchangeTypeLabel')),
        exchangeTypeDirect: computed(() => ts('pages.automations.createEditTrigger.step4Configs.rabbitmq.exchangeTypeDirect')),
        exchangeTypeFanout: computed(() => ts('pages.automations.createEditTrigger.step4Configs.rabbitmq.exchangeTypeFanout')),
        exchangeTypeTopic: computed(() => ts('pages.automations.createEditTrigger.step4Configs.rabbitmq.exchangeTypeTopic')),
        exchangeTypeHeaders: computed(() => ts('pages.automations.createEditTrigger.step4Configs.rabbitmq.exchangeTypeHeaders')),
        routingKeyLabel: computed(() => ts('pages.automations.createEditTrigger.step4Configs.rabbitmq.routingKeyLabel')),
        routingKeyPlaceholder: computed(() => tsRaw('pages.automations.createEditTrigger.step4Configs.rabbitmq.routingKeyPlaceholder')),
        routingKeyHint: computed(() => tsRaw('pages.automations.createEditTrigger.step4Configs.rabbitmq.routingKeyHint')),
        queueLabel: computed(() => ts('pages.automations.createEditTrigger.step4Configs.rabbitmq.queueLabel')),
        queuePlaceholder: computed(() => tsRaw('pages.automations.createEditTrigger.step4Configs.rabbitmq.queuePlaceholder')),
        queueHint: computed(() => tsRaw('pages.automations.createEditTrigger.step4Configs.rabbitmq.queueHint')),
        queueRequired: computed(() => tsRaw('pages.automations.createEditTrigger.step4Configs.rabbitmq.queueRequired')),
        directQueueBannerPrefix: computed(() => ts('pages.automations.createEditTrigger.step4Configs.rabbitmq.directQueueBannerPrefix')),
        directQueueBannerText: computed(() => tsRaw('pages.automations.createEditTrigger.step4Configs.rabbitmq.directQueueBannerText')),
        messageTitle: computed(() => ts('pages.automations.createEditTrigger.step4Configs.rabbitmq.messageTitle')),
        formatJsonButton: computed(() => ts('pages.automations.createEditTrigger.step4Configs.rabbitmq.formatJsonButton')),
        messagePlaceholder: computed(() => tsRaw('pages.automations.createEditTrigger.step4Configs.rabbitmq.messagePlaceholder')),
        messageHint: computed(() => tsRaw('pages.automations.createEditTrigger.step4Configs.rabbitmq.messageHint')),
      },

      /**
       * Slack trigger configuration
       */
      slack: {
        webhookUrlLabel: computed(() => ts('pages.automations.createEditTrigger.step4Configs.slack.webhookUrlLabel')),
        webhookUrlPlaceholder: computed(() => tsRaw('pages.automations.createEditTrigger.step4Configs.slack.webhookUrlPlaceholder')),
        webhookUrlHint: computed(() => tsRaw('pages.automations.createEditTrigger.step4Configs.slack.webhookUrlHint')),
        webhookUrlRequired: computed(() => tsRaw('pages.automations.createEditTrigger.step4Configs.slack.webhookUrlRequired')),
        channelLabel: computed(() => ts('pages.automations.createEditTrigger.step4Configs.slack.channelLabel')),
        channelPlaceholder: computed(() => tsRaw('pages.automations.createEditTrigger.step4Configs.slack.channelPlaceholder')),
        channelHint: computed(() => tsRaw('pages.automations.createEditTrigger.step4Configs.slack.channelHint')),
        usernameLabel: computed(() => ts('pages.automations.createEditTrigger.step4Configs.slack.usernameLabel')),
        usernamePlaceholder: computed(() => tsRaw('pages.automations.createEditTrigger.step4Configs.slack.usernamePlaceholder')),
        usernameHint: computed(() => tsRaw('pages.automations.createEditTrigger.step4Configs.slack.usernameHint')),
        iconEmojiLabel: computed(() => ts('pages.automations.createEditTrigger.step4Configs.slack.iconEmojiLabel')),
        iconEmojiPlaceholder: computed(() => tsRaw('pages.automations.createEditTrigger.step4Configs.slack.iconEmojiPlaceholder')),
        iconEmojiHint: computed(() => tsRaw('pages.automations.createEditTrigger.step4Configs.slack.iconEmojiHint')),
        messageLabel: computed(() => ts('pages.automations.createEditTrigger.step4Configs.slack.messageLabel')),
        messagePlaceholder: computed(() => tsRaw('pages.automations.createEditTrigger.step4Configs.slack.messagePlaceholder')),
        messageHint: computed(() => tsRaw('pages.automations.createEditTrigger.step4Configs.slack.messageHint')),
        messageRequired: computed(() => tsRaw('pages.automations.createEditTrigger.step4Configs.slack.messageRequired')),
      },

      /**
       * Email trigger configuration
       */
      email: {
        smtpSectionTitle: computed(() => ts('pages.automations.createEditTrigger.step4Configs.email.smtpSectionTitle')),
        smtpHostLabel: computed(() => ts('pages.automations.createEditTrigger.step4Configs.email.smtpHostLabel')),
        smtpHostPlaceholder: computed(() => tsRaw('pages.automations.createEditTrigger.step4Configs.email.smtpHostPlaceholder')),
        smtpHostHint: computed(() => tsRaw('pages.automations.createEditTrigger.step4Configs.email.smtpHostHint')),
        smtpHostRequired: computed(() => tsRaw('pages.automations.createEditTrigger.step4Configs.email.smtpHostRequired')),
        smtpPortLabel: computed(() => ts('pages.automations.createEditTrigger.step4Configs.email.smtpPortLabel')),
        smtpPortPlaceholder: computed(() => tsRaw('pages.automations.createEditTrigger.step4Configs.email.smtpPortPlaceholder')),
        smtpPortHint: computed(() => tsRaw('pages.automations.createEditTrigger.step4Configs.email.smtpPortHint')),
        smtpPortRequired: computed(() => tsRaw('pages.automations.createEditTrigger.step4Configs.email.smtpPortRequired')),
        smtpPortInvalid: computed(() => tsRaw('pages.automations.createEditTrigger.step4Configs.email.smtpPortInvalid')),
        fromAddrLabel: computed(() => ts('pages.automations.createEditTrigger.step4Configs.email.fromAddrLabel')),
        fromAddrPlaceholder: computed(() => tsRaw('pages.automations.createEditTrigger.step4Configs.email.fromAddrPlaceholder')),
        fromAddrHint: computed(() => tsRaw('pages.automations.createEditTrigger.step4Configs.email.fromAddrHint')),
        fromAddrRequired: computed(() => tsRaw('pages.automations.createEditTrigger.step4Configs.email.fromAddrRequired')),
        usernameLabel: computed(() => ts('pages.automations.createEditTrigger.step4Configs.email.usernameLabel')),
        usernamePlaceholder: computed(() => tsRaw('pages.automations.createEditTrigger.step4Configs.email.usernamePlaceholder')),
        usernameHint: computed(() => tsRaw('pages.automations.createEditTrigger.step4Configs.email.usernameHint')),
        passwordLabel: computed(() => ts('pages.automations.createEditTrigger.step4Configs.email.passwordLabel')),
        passwordPlaceholder: computed(() => tsRaw('pages.automations.createEditTrigger.step4Configs.email.passwordPlaceholder')),
        passwordHint: computed(() => tsRaw('pages.automations.createEditTrigger.step4Configs.email.passwordHint')),
        contentSectionTitle: computed(() => ts('pages.automations.createEditTrigger.step4Configs.email.contentSectionTitle')),
        toLabel: computed(() => ts('pages.automations.createEditTrigger.step4Configs.email.toLabel')),
        toPlaceholder: computed(() => tsRaw('pages.automations.createEditTrigger.step4Configs.email.toPlaceholder')),
        toHint: computed(() => tsRaw('pages.automations.createEditTrigger.step4Configs.email.toHint')),
        toRequired: computed(() => tsRaw('pages.automations.createEditTrigger.step4Configs.email.toRequired')),
        ccLabel: computed(() => ts('pages.automations.createEditTrigger.step4Configs.email.ccLabel')),
        ccPlaceholder: computed(() => tsRaw('pages.automations.createEditTrigger.step4Configs.email.ccPlaceholder')),
        bccLabel: computed(() => ts('pages.automations.createEditTrigger.step4Configs.email.bccLabel')),
        bccPlaceholder: computed(() => tsRaw('pages.automations.createEditTrigger.step4Configs.email.bccPlaceholder')),
        subjectLabel: computed(() => ts('pages.automations.createEditTrigger.step4Configs.email.subjectLabel')),
        subjectPlaceholder: computed(() => tsRaw('pages.automations.createEditTrigger.step4Configs.email.subjectPlaceholder')),
        subjectRequired: computed(() => tsRaw('pages.automations.createEditTrigger.step4Configs.email.subjectRequired')),
        bodyLabel: computed(() => ts('pages.automations.createEditTrigger.step4Configs.email.bodyLabel')),
        bodyPlaceholder: computed(() => tsRaw('pages.automations.createEditTrigger.step4Configs.email.bodyPlaceholder')),
        htmlBodyLabel: computed(() => ts('pages.automations.createEditTrigger.step4Configs.email.htmlBodyLabel')),
        htmlBodyPlaceholder: computed(() => tsRaw('pages.automations.createEditTrigger.step4Configs.email.htmlBodyPlaceholder')),
        htmlBodyHint: computed(() => tsRaw('pages.automations.createEditTrigger.step4Configs.email.htmlBodyHint')),
        bodyRequired: computed(() => tsRaw('pages.automations.createEditTrigger.step4Configs.email.bodyRequired')),
      },

      /**
       * Microsoft Teams trigger configuration
       */
      teams: {
        webhookUrlLabel: computed(() => ts('pages.automations.createEditTrigger.step4Configs.teams.webhookUrlLabel')),
        webhookUrlPlaceholder: computed(() => tsRaw('pages.automations.createEditTrigger.step4Configs.teams.webhookUrlPlaceholder')),
        webhookUrlHint: computed(() => tsRaw('pages.automations.createEditTrigger.step4Configs.teams.webhookUrlHint')),
        webhookUrlRequired: computed(() => tsRaw('pages.automations.createEditTrigger.step4Configs.teams.webhookUrlRequired')),
        titleLabel: computed(() => ts('pages.automations.createEditTrigger.step4Configs.teams.titleLabel')),
        titlePlaceholder: computed(() => tsRaw('pages.automations.createEditTrigger.step4Configs.teams.titlePlaceholder')),
        titleRequired: computed(() => tsRaw('pages.automations.createEditTrigger.step4Configs.teams.titleRequired')),
        textLabel: computed(() => ts('pages.automations.createEditTrigger.step4Configs.teams.textLabel')),
        textPlaceholder: computed(() => tsRaw('pages.automations.createEditTrigger.step4Configs.teams.textPlaceholder')),
        textHint: computed(() => tsRaw('pages.automations.createEditTrigger.step4Configs.teams.textHint')),
        textRequired: computed(() => tsRaw('pages.automations.createEditTrigger.step4Configs.teams.textRequired')),
        themeColorLabel: computed(() => ts('pages.automations.createEditTrigger.step4Configs.teams.themeColorLabel')),
        themeColorPlaceholder: computed(() => tsRaw('pages.automations.createEditTrigger.step4Configs.teams.themeColorPlaceholder')),
        themeColorHint: computed(() => tsRaw('pages.automations.createEditTrigger.step4Configs.teams.themeColorHint')),
      },

      /**
       * WebSocket trigger configuration
       */
      websocket: {
        urlLabel: computed(() => ts('pages.automations.createEditTrigger.step4Configs.websocket.urlLabel')),
        urlPlaceholder: computed(() => tsRaw('pages.automations.createEditTrigger.step4Configs.websocket.urlPlaceholder')),
        urlHint: computed(() => tsRaw('pages.automations.createEditTrigger.step4Configs.websocket.urlHint')),
        urlRequired: computed(() => tsRaw('pages.automations.createEditTrigger.step4Configs.websocket.urlRequired')),
        headersTitle: computed(() => ts('pages.automations.createEditTrigger.step4Configs.websocket.headersTitle')),
        addHeaderButton: computed(() => ts('pages.automations.createEditTrigger.step4Configs.websocket.addHeaderButton')),
        headerNamePlaceholder: computed(() => tsRaw('pages.automations.createEditTrigger.step4Configs.websocket.headerNamePlaceholder')),
        headerValuePlaceholder: computed(() => tsRaw('pages.automations.createEditTrigger.step4Configs.websocket.headerValuePlaceholder')),
        messageTitle: computed(() => ts('pages.automations.createEditTrigger.step4Configs.websocket.messageTitle')),
        formatJsonButton: computed(() => ts('pages.automations.createEditTrigger.step4Configs.websocket.formatJsonButton')),
        messagePlaceholder: computed(() => tsRaw('pages.automations.createEditTrigger.step4Configs.websocket.messagePlaceholder')),
        messageHint: computed(() => tsRaw('pages.automations.createEditTrigger.step4Configs.websocket.messageHint')),
      },
    },

    /**
     * Category translations
     */
    categories: {
      technical: computed(() => ts('pages.automations.createEditTrigger.categories.technical')),
      communication: computed(() => ts('pages.automations.createEditTrigger.categories.communication')),
    },

    /**
     * Status options translations
     */
    statusOptions: {
      active: {
        label: computed(() => ts('pages.automations.createEditTrigger.statusOptions.active.label')),
        value: true,
      },
      inactive: {
        label: computed(() => ts('pages.automations.createEditTrigger.statusOptions.inactive.label')),
        value: false,
      },
    },

    /**
     * Navigation buttons translations
     */
    navigation: {
      previous: computed(() => ts('pages.automations.createEditTrigger.navigation.previous')),
      next: computed(() => ts('pages.automations.createEditTrigger.navigation.next')),
      save: computed(() => ts('pages.automations.createEditTrigger.navigation.save')),
      update: computed(() => ts('pages.automations.createEditTrigger.navigation.update')),
    },

    /**
     * Notification messages translations
     */
    notifications: {
      loading: computed(() => tsRaw('pages.automations.createEditTrigger.notifications.loading')),
      loadFailed: computed(() => tsRaw('pages.automations.createEditTrigger.notifications.loadFailed')),
      created: computed(() => tsRaw('pages.automations.createEditTrigger.notifications.created')),
      creationFailed: computed(() => tsRaw('pages.automations.createEditTrigger.notifications.creationFailed')),
      updated: computed(() => tsRaw('pages.automations.createEditTrigger.notifications.updated')),
      updateFailed: computed(() => tsRaw('pages.automations.createEditTrigger.notifications.updateFailed')),
      validationFailed: computed(() => tsRaw('pages.automations.createEditTrigger.notifications.validationFailed')),
      networkError: computed(() => tsRaw('pages.automations.createEditTrigger.notifications.networkError')),
      alreadyExists: computed(() => tsRaw('pages.automations.createEditTrigger.notifications.alreadyExists')),
      fillRequiredFields: computed(() => tsRaw('pages.automations.createEditTrigger.notifications.fillRequiredFields')),
    },

    /**
     * Validation translations
     */
    validation: {
      selectCategory: computed(() => tsRaw('pages.automations.createEditTrigger.validation.selectCategory')),
      selectType: computed(() => tsRaw('pages.automations.createEditTrigger.validation.selectType')),
      fillRequired: computed(() => tsRaw('pages.automations.createEditTrigger.validation.fillRequired')),
      completeConfig: computed(() => tsRaw('pages.automations.createEditTrigger.validation.completeConfig')),
    },
  };
}
