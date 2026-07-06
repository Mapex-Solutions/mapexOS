/**
 * LoRaWAN frequency plans accepted by The Things Stack — the exact set the
 * mapexLNS embeds (src/shared/frequencyplans/data/frequency-plans.yml). Kept on
 * the front as a static reference (no backend endpoint): a gateway's
 * frequencyPlanId and a device's region are both chosen from these ids. The
 * `gateways` flag marks plans valid for a gateway radio config; device-side
 * plan variants have it false. Regenerate from the LNS index if the embedded
 * plans change.
 */
export interface FrequencyPlanOption {
  /** Plan id sent to the API (LorawanGatewayConfig.frequencyPlanId / device region). */
  id: string;
  /** Human-readable label for the dropdown. */
  name: string;
  /** Whether the plan is valid for a gateway radio config. */
  gateways: boolean;
}

export const LORAWAN_FREQUENCY_PLANS: readonly FrequencyPlanOption[] = [
  { id: 'EU_863_870', name: 'Europe 863-870 MHz (SF12 for RX2)', gateways: true },
  { id: 'EU_863_870_TTN', name: 'Europe 863-870 MHz (SF9 for RX2 - recommended)', gateways: false },
  { id: 'EU_863_870_ROAMING_DRAFT', name: 'Europe 863-870 MHz, 6 channels for roaming (Draft)', gateways: false },
  { id: 'EU_868_1', name: 'Europe 868.1 MHz', gateways: true },
  { id: 'EU_433', name: 'Europe 433 MHz (ITU region 1)', gateways: true },
  { id: 'US_902_928_FSB_1', name: 'United States 902-928 MHz, FSB 1', gateways: true },
  { id: 'US_902_928_FSB_2', name: 'United States 902-928 MHz, FSB 2 (used by TTN)', gateways: true },
  { id: 'US_902_928_FSB_3', name: 'United States 902-928 MHz, FSB 3', gateways: true },
  { id: 'US_902_928_FSB_4', name: 'United States 902-928 MHz, FSB 4', gateways: true },
  { id: 'US_902_928_FSB_5', name: 'United States 902-928 MHz, FSB 5', gateways: true },
  { id: 'US_902_928_FSB_6', name: 'United States 902-928 MHz, FSB 6', gateways: true },
  { id: 'US_902_928_FSB_7', name: 'United States 902-928 MHz, FSB 7', gateways: true },
  { id: 'US_902_928_FSB_8', name: 'United States 902-928 MHz, FSB 8', gateways: true },
  { id: 'US_903_0', name: 'United States 903.0 MHz', gateways: true },
  { id: 'US_904_6', name: 'United States 904.6 MHz', gateways: true },
  { id: 'US_906_2', name: 'United States 906.2 MHz', gateways: true },
  { id: 'US_907_8', name: 'United States 907.8 MHz', gateways: true },
  { id: 'US_909_4', name: 'United States 909.4 MHz', gateways: true },
  { id: 'US_911_0', name: 'United States 911.0 MHz', gateways: true },
  { id: 'US_912_6', name: 'United States 912.6 MHz', gateways: true },
  { id: 'US_914_2', name: 'United States 914.2 MHz', gateways: true },
  { id: 'AU_915_928_FSB_1', name: 'Australia 915-928 MHz, FSB 1', gateways: true },
  { id: 'AU_915_928_FSB_2', name: 'Australia 915-928 MHz, FSB 2 (used by TTN)', gateways: true },
  { id: 'AU_915_928_FSB_3', name: 'Australia 915-928 MHz, FSB 3', gateways: true },
  { id: 'AU_915_928_FSB_4', name: 'Australia 915-928 MHz, FSB 4', gateways: true },
  { id: 'AU_915_928_FSB_5', name: 'Australia 915-928 MHz, FSB 5', gateways: true },
  { id: 'AU_915_928_FSB_6', name: 'Australia 915-928 MHz, FSB 6', gateways: true },
  { id: 'AU_915_928_FSB_7', name: 'Australia 915-928 MHz, FSB 7', gateways: true },
  { id: 'AU_915_928_FSB_8', name: 'Australia 915-928 MHz, FSB 8', gateways: true },
  { id: 'AU_915_928_FSB_1_BR', name: 'Brazil 915-928 MHz, FSB 1 (dwell time disabled)', gateways: true },
  { id: 'AU_915_928_FSB_2_BR', name: 'Brazil 915-928 MHz, FSB 2 (dwell time disabled)', gateways: true },
  { id: 'AU_915_928_FSB_3_BR', name: 'Brazil 915-928 MHz, FSB 3 (dwell time disabled)', gateways: true },
  { id: 'AU_915_928_FSB_4_BR', name: 'Brazil 915-928 MHz, FSB 4 (dwell time disabled)', gateways: true },
  { id: 'AU_915_928_FSB_5_BR', name: 'Brazil 915-928 MHz, FSB 5 (dwell time disabled)', gateways: true },
  { id: 'AU_915_928_FSB_6_BR', name: 'Brazil 915-928 MHz, FSB 6 (dwell time disabled)', gateways: true },
  { id: 'AU_915_928_FSB_7_BR', name: 'Brazil 915-928 MHz, FSB 7 (dwell time disabled)', gateways: true },
  { id: 'AU_915_928_FSB_8_BR', name: 'Brazil 915-928 MHz, FSB 8 (dwell time disabled)', gateways: true },
  { id: 'CN_470_510_FSB_1', name: 'China 470-510 MHz, FSB 1', gateways: true },
  { id: 'CN_470_510_FSB_11', name: 'China 470-510 MHz, FSB 11 (used by TTN)', gateways: true },
  { id: 'AS_920_923', name: 'Asia 920-923 MHz', gateways: true },
  { id: 'AS_920_923_LBT', name: 'Asia 920-923 MHz with LBT', gateways: true },
  { id: 'AS_920_923_TTN_JP_1', name: 'Japan 920-923 MHz with LBT (channels 31-38)', gateways: true },
  { id: 'AS_920_923_TTN_JP_1_LAND_MOBILE', name: 'Japan 920-923 MHz with LBT (channels 31-38), Max EIRP 27 dBm', gateways: true },
  { id: 'AS_920_923_TTN_JP_2', name: 'Japan 920-923 MHz with LBT (channels 24-27 and 35-38)', gateways: true },
  { id: 'AS_920_923_TTN_JP_3', name: 'Japan 920-923 MHz with LBT (channels 24-31)', gateways: true },
  { id: 'AS_920_923_TTN_JP_3_LAND_MOBILE', name: 'Japan 920-923 MHz with LBT (channels 24-31), Max EIRP 27 dBm', gateways: true },
  { id: 'AS_923', name: 'Asia 915-928 MHz (AS923 Group 1) with only default channels', gateways: true },
  { id: 'AS_923_2', name: 'Asia 920-923 MHz (AS923 Group 2) with only default channels', gateways: true },
  { id: 'AS_923_3', name: 'Asia 915-921 MHz (AS923 Group 3) with only default channels', gateways: true },
  { id: 'AS_923_4', name: 'Asia 917-920 MHz (AS923 Group 4) with only default channels', gateways: true },
  { id: 'AS_923_NDT', name: 'Asia 915-928 MHz (AS923 Group 1) with only default channels and dwell time disabled', gateways: true },
  { id: 'AS_923_DT', name: 'Asia 915-928 MHz (AS923 Group 1) with only default channels and dwell time enabled', gateways: true },
  { id: 'AS_923_2_NDT', name: 'Asia 920-923 MHz (AS923 Group 2) with only default channels and dwell time disabled', gateways: true },
  { id: 'AS_923_2_DT', name: 'Asia 920-923 MHz (AS923 Group 2) with only default channels and dwell time enabled', gateways: true },
  { id: 'AS_923_3_NDT', name: 'Asia 920-923 MHz (AS923 Group 3) with only default channels and dwell time disabled', gateways: true },
  { id: 'AS_923_3_DT', name: 'Asia 920-923 MHz (AS923 Group 3) with only default channels and dwell time enabled', gateways: true },
  { id: 'AS_923_4_NDT', name: 'Asia 920-923 MHz (AS923 Group 4) with only default channels and dwell time disabled', gateways: true },
  { id: 'AS_923_4_DT', name: 'Asia 920-923 MHz (AS923 Group 4) with only default channels and dwell time enabled', gateways: true },
  { id: 'AS_923_925', name: 'Asia 923-925 MHz', gateways: true },
  { id: 'AS_923_925_LBT', name: 'Asia 923-925 MHz with LBT', gateways: true },
  { id: 'AS_920_923_TTN_AU', name: 'Asia 920-923 MHz (used by TTN Australia)', gateways: true },
  { id: 'AS_923_925_TTN_AU', name: 'Asia 923-925 MHz (used by TTN Australia - secondary channels)', gateways: false },
  { id: 'KR_920_923_TTN', name: 'South Korea 920-923 MHz', gateways: true },
  { id: 'MA_869_870_DRAFT', name: 'Morocco 869-870 MHz', gateways: true },
  { id: 'IN_865_867', name: 'India 865-867 MHz', gateways: true },
  { id: 'RU_864_870_TTN', name: 'Russia 864-870 MHz', gateways: true },
  { id: 'ISM_2400_3CH_DRAFT2', name: 'LoRa 2.4 GHz with 3 channels (Draft 2)', gateways: true },
  { id: 'IL_917_920_TTN', name: 'Israel 917-920 MHz (channels 1-4 and 11-14)', gateways: true },
  { id: 'IL_917_920_TTN_2', name: 'Israel 917-920 MHz (channels 1-7 and 9)', gateways: true },
  { id: 'SG_920_923_TTN', name: 'Singapore 920-923 MHz', gateways: true },
  { id: 'UZ_923_1', name: 'Uzbekistan 915-928 MHz with only default channels', gateways: true },
  { id: 'US_902_928_CUSTOM_ADC_500KHZ', name: 'United States 902-928 MHz, Custom ADC 500kHz', gateways: true },
  { id: 'US_902_928_CUSTOM_ADC_500KHZ_904_6', name: 'United States 904.6 MHz, Custom ADC 500kHz', gateways: true },
  { id: 'US_902_928_CUSTOM_VIASAT_FSB_1', name: 'United States 902-928 MHz, Custom Viasat FSB 1', gateways: true },
  { id: 'US_902_928_CUSTOM_VIASAT_FSB_2', name: 'United States 902-928 MHz, Custom Viasat FSB 2', gateways: true },
];
