export interface User {
  id: string;
  email?: string;
  firstName: string;
  lastName: string;
  phone: string;
  role: UserRole;
  isActive: boolean;
  createdAt: Date;
}

export type UserRole =
  | 'admin'
  | 'customer'
  | 'agent';