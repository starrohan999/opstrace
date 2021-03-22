UPDATE public.exporter e
  SET tenant = t.name
  FROM public.tenant t
  WHERE e.tenant_id = t.id;

UPDATE public.credential c
  SET tenant = t.name
  FROM public.tenant t
  WHERE c.tenant_id = t.id;
