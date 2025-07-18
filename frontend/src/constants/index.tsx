export const RESOURCES = {
  USERS: 'users',
  ROLES: 'roles',
  PERMISSIONS: 'permissions',
  DASHBOARD: 'dashboard',
  PROFILE: 'profile'
} as const;

export const ACTIONS = {
  CREATE: 'create',
  READ: 'read',
  UPDATE: 'update',
  DELETE: 'delete',
  LIST: 'list',
  MANAGE: 'manage',
  VIEW: 'view',
  EDIT: 'edit'
} as const;

export const USER_STATUS = {
  ACTIVE: 'active',
  INACTIVE: 'inactive',
  SUSPENDED: 'suspended',
  PENDING: 'pending'
} as const;

export const USER_STATUS_COLORS = {
  [USER_STATUS.ACTIVE]: 'bg-green-100 text-green-800',
  [USER_STATUS.INACTIVE]: 'bg-gray-100 text-gray-800',
  [USER_STATUS.SUSPENDED]: 'bg-red-100 text-red-800',
  [USER_STATUS.PENDING]: 'bg-yellow-100 text-yellow-800'
} as const;

export const LANGUAGES = {
  EN: 'en',
  VI: 'vi'
} as const;

export const DEFAULT_VALUES = {
  LANGUAGE: LANGUAGES.EN,
  TIMEZONE: 'UTC',
  PAGE_SIZE: 10,
  PAGE_SIZE_OPTIONS: [1, 5, 10, 25, 50] as number[]
} as const;

export const PERMISSION_WILDCARDS = {
  ALL: '*'
} as const;

export type ResourceType = typeof RESOURCES[keyof typeof RESOURCES];
export type ActionType = typeof ACTIONS[keyof typeof ACTIONS];
export type UserStatusType = typeof USER_STATUS[keyof typeof USER_STATUS];
export type LanguageType = typeof LANGUAGES[keyof typeof LANGUAGES];