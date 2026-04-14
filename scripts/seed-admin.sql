-- Default admin user (password: admin123)
-- bcrypt hash of "admin123" with cost 10
INSERT INTO admin_users (username, password_hash, role)
VALUES ('admin', '$2a$10$nUTthNgz8t59pRYyncjPwe0knC5c9xmrMuCnu489lrKaJ0F1wV0/q', 'admin')
ON CONFLICT DO NOTHING;

-- Default mirrors
INSERT INTO mirrors (domain, is_active, is_primary, region)
VALUES
  ('localhost:8080', true, true, 'global')
ON CONFLICT DO NOTHING;
