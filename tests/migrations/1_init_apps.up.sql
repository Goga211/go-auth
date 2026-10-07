INSERT INTO apps (id, name, secret)
VALUES (1, 'test', 'secret_test')
ON CONFLICT DO NOTHING;