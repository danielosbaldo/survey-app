SELECT id,
       prompt,
       type,
       order_num,
       created_at,
       updated_at
FROM public.questions
LIMIT 1000;