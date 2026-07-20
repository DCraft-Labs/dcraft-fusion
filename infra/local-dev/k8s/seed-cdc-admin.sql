-- Local CDC bootstrap: default admin for Docker Desktop smoke tests.
-- Password: Admin@123
DO $$
DECLARE
  v_role_id uuid;
  v_user_id uuid;
BEGIN
  INSERT INTO roles (role_name, display_name, role_level, description, is_active, is_system_role, created_at, updated_at)
  VALUES ('superadmin', 'Super Admin', 'superadmin', 'Full system access', true, true, NOW(), NOW())
  ON CONFLICT (role_name) DO NOTHING;
  SELECT role_id INTO v_role_id FROM roles WHERE role_name = 'superadmin';

  INSERT INTO users (
    username, email, password_hash,
    first_name, last_name, full_name,
    is_active, is_superuser, is_email_verified,
    failed_login_attempts, preferences,
    created_at, updated_at
  )
  VALUES (
    'admin', 'admin@dcraftfusion.io',
    '$2b$12$DnILea2s7A9Jq7rPYNtdieWK3vWLQ4Aljlq2/SkJqQIuw17KKtdHO',
    'Fusion', 'Admin', 'Fusion Admin',
    true, true, true,
    '0', '{}'::jsonb,
    NOW(), NOW()
  )
  ON CONFLICT (username) DO UPDATE SET
    password_hash = EXCLUDED.password_hash,
    email = EXCLUDED.email,
    is_superuser = true,
    is_active = true,
    updated_at = NOW()
  RETURNING user_id INTO v_user_id;

  IF v_user_id IS NULL THEN
    SELECT user_id INTO v_user_id FROM users WHERE username = 'admin';
  END IF;

  INSERT INTO user_roles (user_id, role_id)
  VALUES (v_user_id, v_role_id)
  ON CONFLICT DO NOTHING;
END $$;
