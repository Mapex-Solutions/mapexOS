export interface LakeHouseFrequency {
  type: 'minute' | 'hour' | 'day' | 'week' | 'month' | 'year';
  interval?: number | undefined; // Para intervalos customizados (ex: a cada 2 horas)
  weekdays?: string[] | undefined; // Para frequência semanal
  dayOfMonth?: number | undefined; // Para frequência mensal
  time?: string | undefined; // Horário específico (HH:mm)
}