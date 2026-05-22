export interface UserProfileData {
  firstName: string;
  lastName: string;
  email: string;
  phone: string;
  avatar: string | null;
}

export interface PasswordData {
  current: string;
  new: string;
  confirm: string;
}

export interface ApiKey {
  id: number;
  name: string;
  key: string;
  created: string;
}

export interface NotificationType {
  name: string;
  enabled: boolean;
}

export interface NotificationChannel {
  id: number;
  name: string;
  icon: string;
  enabled: boolean;
  types: NotificationType[];
}

export interface NavigationItem {
  label: string;
  value: string;
  icon: string;
  description: string;
}
