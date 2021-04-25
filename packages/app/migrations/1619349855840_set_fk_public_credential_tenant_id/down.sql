alter table "public"."credential" drop constraint "credential_tenant_id_fkey",
          add constraint "credential_tenant_fkey"
          foreign key ("tenant")
          references "public"."tenant"
          ("name")
          on update cascade
          on delete cascade;
