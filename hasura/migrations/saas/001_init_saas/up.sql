-- Shared functions
CREATE OR REPLACE FUNCTION trigger_set_timestamp()
RETURNS TRIGGER AS $$
BEGIN
  NEW.updated_at = NOW();
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;


-- Saas Account
CREATE TABLE public.saas_account(
  id uuid NOT NULL DEFAULT gen_random_uuid(),
  name text NOT NULL UNIQUE,
  id_address_invoice uuid,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  PRIMARY KEY (id)
);

CREATE TRIGGER saas_account_set_timestamp BEFORE UPDATE ON saas_account FOR EACH ROW EXECUTE PROCEDURE trigger_set_timestamp();


-- Saas Role
CREATE TABLE public.saas_role(
  id text NOT NULL,
  description text NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  PRIMARY KEY (id)
);

CREATE TRIGGER saas_role_set_timestamp BEFORE UPDATE ON saas_role FOR EACH ROW EXECUTE PROCEDURE trigger_set_timestamp();

INSERT INTO public.saas_role (id, description) VALUES
  ('admin',         'Role with all priviledges'),
  ('account_owner', 'Owner of an account, it has  all priviledges within an account'),
  ('account_admin', 'Administrator of the account'),
  ('account_user',  'Simple use within an account'),
  ('logged_in',     'Logged in user that has no account permission'),
  ('anonymous',     'User that is not logged in');


-- Saas Address
CREATE TABLE public.saas_address(
  id uuid NOT NULL DEFAULT gen_random_uuid(),
  id_account uuid,
  id_user text NOT NULL,
  address text NOT NULL,
  city text NOT NULL,
  country text NOT NULL,
  postal_code text NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  PRIMARY KEY (id),
  CONSTRAINT id_account_saas_account_id_fkey FOREIGN KEY(id_account) REFERENCES saas_account(id)
);

CREATE TRIGGER saas_address_set_timestamp BEFORE UPDATE ON saas_address FOR EACH ROW EXECUTE PROCEDURE trigger_set_timestamp();

ALTER TABLE saas_account ADD CONSTRAINT id_address_invoice_saas_address_id_fkey FOREIGN KEY(id_address_invoice) REFERENCES saas_address(id);

-- Saas Membership
CREATE TABLE public.saas_membership(
  id uuid NOT NULL DEFAULT gen_random_uuid(),
  id_account uuid NOT NULL,
  id_user text NOT NULL,
  id_role text NOT NULL,
  selected_at TIMESTAMPTZ,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  PRIMARY KEY (id),
  UNIQUE(id_account, id_user),
  CONSTRAINT id_account_saas_account_id_fkey FOREIGN KEY(id_account) REFERENCES saas_account(id),
  CONSTRAINT id_role_saas_role_id_fkey FOREIGN KEY(id_role) REFERENCES saas_role(id) ON UPDATE NO ACTION ON DELETE CASCADE DEFERRABLE INITIALLY DEFERRED
);

CREATE INDEX idx_saas_membership_id_account ON saas_membership(id_account);
CREATE INDEX idx_saas_membership_id_account_selected_at ON saas_membership(id_account,selected_at);

CREATE TRIGGER saas_membership_set_timestamp BEFORE UPDATE ON saas_membership FOR EACH ROW EXECUTE PROCEDURE trigger_set_timestamp();

-- Saas User Account
CREATE OR REPLACE VIEW public.saas_user_account AS
SELECT
  saas_account.id,
  saas_account.name,
  saas_account.created_at,
  saas_account.updated_at,
  saas_membership.id_user
FROM
  saas_membership, saas_account
WHERE
  (saas_membership.id_account = saas_account.id);


-- Subscription Plan
CREATE TABLE public.subscription_plan(
  id text NOT NULL,
  description text NOT NULL,
  is_active boolean NOT NULL DEFAULT FALSE,
  stripe_code text,
  trial_days integer,
  price numeric,
  currency text,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  PRIMARY KEY (id)
);

CREATE TRIGGER subscription_plan_set_timestamp BEFORE UPDATE ON subscription_plan FOR EACH ROW EXECUTE PROCEDURE trigger_set_timestamp();

INSERT INTO subscription_plan (id, description) VALUES
  ('free',      'free plan'),
  ('basic',     'basic plan'),
  ('enterprise','enterprise plan'),
  ('premium',   'premium plan');


-- Subscription Active Plan
CREATE
OR REPLACE VIEW public.subscription_active_plan AS
SELECT
  subscription_plan.id,
  subscription_plan.description,
  subscription_plan.price,
  subscription_plan.trial_days
FROM
  subscription_plan
WHERE
  (subscription_plan.is_active = true);


-- Subscription Customer
CREATE TABLE public.subscription_customer(
  id uuid NOT NULL DEFAULT gen_random_uuid(),
  id_account uuid NOT NULL UNIQUE,
  stripe_customer text NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  PRIMARY KEY (id),
  CONSTRAINT id_account_saas_account_id_fkey FOREIGN KEY(id_account) REFERENCES saas_account(id) ON UPDATE NO ACTION ON DELETE CASCADE DEFERRABLE INITIALLY DEFERRED
);

CREATE TRIGGER subscription_customer_set_timestamp BEFORE UPDATE ON subscription_customer FOR EACH ROW EXECUTE PROCEDURE trigger_set_timestamp();


-- Subscription Status
CREATE TABLE public.subscription_status(
  id uuid NOT NULL DEFAULT gen_random_uuid(),
  id_account uuid NOT NULL UNIQUE,
  status text NOT NULL,
  is_active boolean NOT NULL DEFAULT FALSE,
  id_plan text NOT NULL,
  stripe_subscription_id text,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  PRIMARY KEY (id),
  CONSTRAINT id_plan_subscription_plan_id_fkey FOREIGN KEY(id_plan) REFERENCES subscription_plan(id) ON UPDATE NO ACTION ON DELETE CASCADE DEFERRABLE INITIALLY DEFERRED,
  CONSTRAINT id_account_saas_account_id_fkey FOREIGN KEY(id_account) REFERENCES saas_account(id) ON UPDATE NO ACTION ON DELETE CASCADE DEFERRABLE INITIALLY DEFERRED
);

CREATE INDEX idx_subscription_status_stripe_subscription_id ON subscription_status(stripe_subscription_id);

CREATE TRIGGER subscription_status_set_timestamp BEFORE UPDATE ON subscription_status FOR EACH ROW EXECUTE PROCEDURE trigger_set_timestamp();


-- Subscription Event
CREATE TABLE public.subscription_event(
  id uuid NOT NULL DEFAULT gen_random_uuid(),
  data jsonb NOT NULL,
  type text NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  PRIMARY KEY (id)
);

CREATE TRIGGER subscription_event_set_timestamp BEFORE UPDATE ON subscription_event FOR EACH ROW EXECUTE PROCEDURE trigger_set_timestamp();
