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
  json_patch(
    s.data,
    json_object(
      'project_name',
      p.data ->> 'name',
      'project_order',
      p.data ->> 'child_order'
    )
  ) AS data
FROM
  sections s
  LEFT JOIN projects p ON s.data ->> 'project_id' = p.id;

CREATE VIEW IF NOT EXISTS tasks_view AS
SELECT
  t.id,
  json_patch(
    t.data,
    json_object(
      -- project
      'project',
      json(p.data),
      -- section
      'section',
      json(s.data),
      -- labels
      'labels',
      coalesce(
        (
          SELECT
            json_group_array(json(lv.data))
          FROM
            json_each(t.data -> 'labels') AS je
            LEFT JOIN labels_view lv ON je.value = lv.id
        ),
        json_array()
      ),
      -- task
      'sub_task_status',
      json_object(
        'total',
        count(child.id),
        'completed',
        sum(
          CASE
            WHEN child.data -> 'checked' = 'true' THEN 1
            ELSE 0
          END
        )
      )
    )
  ) AS data
FROM
  tasks t
  LEFT JOIN projects p ON t.data ->> 'project_id' = p.id
  LEFT JOIN sections s ON t.data ->> 'section_id' = s.id
  LEFT JOIN tasks child ON t.id = child.data ->> 'parent_id'
GROUP BY
  t.id;
