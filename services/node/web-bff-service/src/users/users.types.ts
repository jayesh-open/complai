export interface User {
  id: string;
  tenant_id: string;
  email: string;
  email_verified: boolean;
  first_name: string;
  last_name: string;
  status: string;
  created_at: string;
  updated_at: string;
}

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

export interface Permission {
  id: string;
  tenant_id: string;
  resource: string;
  action: string;
  description?: string;
  created_at: string;
}

export interface UserWithRole {
  id: string;
  email: string;
  first_name: string;
  last_name: string;
  status: string;
  created_at: string;
  updated_at: string;
  role: { id: string; name: string; display_name: string } | null;
}

export interface RoleWithPermissions {
  role: Role;
  permissions: Permission[];
}

export interface CreateUserPayload {
  email: string;
  first_name: string;
  last_name: string;
}

export interface CreateRolePayload {
  name: string;
  display_name: string;
  description?: string;
}

export interface PermissionPair {
  resource: string;
  action: string;
}
