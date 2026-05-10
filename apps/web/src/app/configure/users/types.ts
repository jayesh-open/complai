export interface Role {
  id: string;
  tenant_id: string;
  name: string;
  display_name: string;
  description?: string;
  is_system: boolean;
  created_at: string;
  updated_at: string;
}

export interface UserRole {
  id: string;
  name: string;
  display_name: string;
}

export interface User {
  id: string;
  email: string;
  first_name: string;
  last_name: string;
  status: "active" | "inactive";
  role: UserRole | null;
  created_at: string;
  updated_at: string;
}
