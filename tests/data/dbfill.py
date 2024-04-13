import json
import random

import psycopg2

import ranjg


conn = psycopg2.connect(
    dbname="banner",
    user="test_user",
    password="crakme",
    host="localhost",
    port="4444"
)

schema = {
    'type': 'object',
    'properties': {
        'name': {
            'type': 'string',
            'minLength': 1,
            'maxLength': 32,
        },
        'age': {
            'type': 'integer',
            'minimum': 0,
            'maximum': 50,
        },
        'comment': {
            'type': 'string',
            'minLength': 1
        },
        "additional": {
            'type': 'object',
            'required': ['key1'],
            'properties': {
                'key1': {
                    'type': 'array',
                    'items': {
                        'type': 'integer',
                    },
                    'minLength': 10,
                    'maxLength': 32,
                }
            }
        }
    },
    'required': ['name', 'age']
}

with conn.cursor() as curs:
    for i in range(1, 1001):
        gen = ranjg.gen(schema)

        is_active = True if i <= 20 else bool(random.randint(0, 1))

        try:
            row = json.dumps(gen), is_active
            curs.execute("INSERT INTO banners (content, is_active) VALUES (%s, %s) ON CONFLICT DO NOTHING;", row)
            curs.execute("INSERT INTO tag_feature_banners (tag_id, feature_id, banner_id) VALUES (%s, %s, %s);", (i, i, i))
            curs.execute("INSERT INTO tag_feature_banners (tag_id, feature_id, banner_id) VALUES (%s, %s, %s);", (i+1, i, i))
            curs.execute("INSERT INTO tag_feature_banners (tag_id, feature_id, banner_id) VALUES (%s, %s, %s);", (i+2, i, i))

        except psycopg2.DatabaseError as exc:
            print(row, "не вставлено: ", exc)
            raise exc


with conn.cursor() as curs:
    cnt = 0
    while cnt < 500:
        try:
            row = random.randint(1, 100), random.randint(1, 10), random.randint(1, 1000)
            curs.execute(
                "INSERT INTO tag_feature_banners (tag_id, feature_id, banner_id) VALUES (%s, %s, %s) ON CONFLICT DO NOTHING;",
                row)
            cnt += 1
        except psycopg2.DatabaseError as exc:
            print(row, "не вставлено: ", exc)
            raise exc

conn.commit()
conn.close()
