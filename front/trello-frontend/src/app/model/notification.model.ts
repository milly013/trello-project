// Notification model
export interface Notification {
    id?: string;
    user_id?: string;
    type?: string;
    message?: string;
    created_at?: Date;
    is_read?: boolean;
  }