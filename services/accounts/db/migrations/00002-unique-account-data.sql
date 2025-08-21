ALTER TABLE accounts
    ADD CONSTRAINT accounts_login_unique UNIQUE (login);

ALTER TABLE accounts
    ADD CONSTRAINT accounts_email_unique UNIQUE (email);

ALTER TABLE accounts
    ADD CONSTRAINT accounts_phone_unique UNIQUE (phone_number);

ALTER TABLE profiles
    ADD CONSTRAINT profile_id_unique UNIQUE (profile_id);
