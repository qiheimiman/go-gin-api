#### go_gin_api.bubbles1000 
加密泡泡排名

| 序号 | 名称 | 描述 | 类型 | 键 | 为空 | 额外 | 默认值 |
| :--: | :--: | :--: | :--: | :--: | :--: | :--: | :--: |
| 1 | id |  | int unsigned | PRI | NO | auto_increment |  |
| 2 | name |  | varchar(255) |  | NO |  |  |
| 3 | slug |  | varchar(255) | UNI | NO |  |  |
| 4 | symbol |  | varchar(50) | UNI | NO |  |  |
| 5 | dominance |  | decimal(10,6) |  | YES |  | 0.000000 |
| 6 | image |  | varchar(255) |  | YES |  |  |
| 7 | rank |  | int |  | NO |  | 0 |
| 8 | price |  | decimal(20,10) |  | YES |  | 0.0000000000 |
| 9 | marketcap |  | bigint |  | YES |  | 0 |
| 10 | volume |  | bigint |  | YES |  | 0 |
| 11 | cg_id |  | varchar(255) |  | YES |  |  |
| 12 | symbols |  | text |  | YES |  |  |
| 13 | performance |  | text |  | YES |  |  |
| 14 | rank_diffs |  | text |  | YES |  |  |
| 15 | exchange_prices |  | text |  | YES |  |  |
