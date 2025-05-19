CREATE VIEW IF NOT EXISTS labels_view AS
WITH
  label_count AS (
    SELECT
      je.value AS id,
      count(tasks.id) AS count
    FROM
      tasks,
      json_each(tasks.data -> 'labels') AS je
    GROUP BY
      je.value
  )
SELECT
  l.id,
  json_patch(
    l.data,
    json_object('count', coalesce(lc.count, 0))
  ) AS data
FROM
  labels l
  LEFT JOIN label_count lc ON l.id = lc.id
UNION ALL
SELECT
  lc.id AS id,
  json_object(
    'name',
    lc.id,
    'is_shared',
    json('true'),
    'count',
    lc.count
  ) AS data
FROM
  label_count lc
WHERE
  NOT EXISTS (
    SELECT
      1
    FROM
      labels l
    WHERE
      l.id = lc.id
  );

CREATE VIEW IF NOT EXISTS sections_view AS
SELECT
  s.id,
  s.data AS section,
  p.data AS project
FROM
  sections s
  LEFT JOIN projects p ON s.data ->> 'project_id' = p.id;

CREATE VIEW IF NOT EXISTS tasks_view AS
WITH
  sub_task_status AS (
    SELECT
      child.data ->> 'parent_id' AS parent_id,
      count(child.id) AS total,
      sum(
        CASE
          WHEN child.data -> 'checked' = 'true' THEN 1
          ELSE 0
        END
      ) AS completed
    FROM
      tasks child
    WHERE
      child.data ->> 'parent_id' IS NOT NULL
    GROUP BY
      child.data ->> 'parent_id'
  ),
  task_labels AS (
    SELECT
      t.id AS task_id,
      (
        SELECT
          json_group_array(json(lv.data))
        FROM
          json_each(t.data -> 'labels') AS je
          LEFT JOIN labels_view lv ON je.value = lv.id
      ) AS labels
    FROM
      tasks t
  )
SELECT
  t.id,
  json_patch(
    t.data,
    json_object(
      'sub_task_status',
      json_object(
        'total',
        coalesce(sts.total, 0),
        'completed',
        coalesce(sts.completed, 0)
      ),
      'labels',
      json(tl.labels)
    )
  ) AS task,
  p.data AS project,
  s.data AS section
FROM
  tasks t
  LEFT JOIN sub_task_status sts ON t.id = sts.parent_id
  LEFT JOIN task_labels tl ON t.id = tl.task_id
  LEFT JOIN projects p ON t.data ->> 'project_id' = p.id
  LEFT JOIN sections s ON t.data ->> 'section_id' = s.id
