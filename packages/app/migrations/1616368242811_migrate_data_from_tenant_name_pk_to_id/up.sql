UPDATE public.credential c
  SET tenant_id = t.id
  FROM public.tenant t
  WHERE c.tenant = t.name;

UPDATE public.exporter e
  SET tenant_id = t.id
  FROM public.tenant t
  WHERE e.tenant = t.name;
