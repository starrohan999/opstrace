ALTER TABLE "public"."exporter" ADD COLUMN "credential" text;
ALTER TABLE "public"."exporter" ALTER COLUMN "credential" DROP NOT NULL;
ALTER TABLE "public"."exporter" ADD CONSTRAINT exporter_tenant_credential_fkey
              FOREIGN KEY (credential, tenant)
              REFERENCES "public"."credential"
              (name, tenant) ON DELETE restrict ON UPDATE cascade;

