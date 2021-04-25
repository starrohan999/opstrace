alter table "public"."exporter" drop constraint "exporter_credential_id_fkey",
          add constraint "exporter_tenant_credential_fkey"
          foreign key ("credential", "tenant")
          references "public"."credential"
          ("name", "tenant")
          on update cascade
          on delete restrict;
