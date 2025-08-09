#### go_gin_api.news 
新闻表

| 序号 | 名称 | 描述 | 类型 | 键 | 为空 | 额外 | 默认值 |
| :--: | :--: | :--: | :--: | :--: | :--: | :--: | :--: |
| 1 | id | 自增主键 | int unsigned | PRI | NO | auto_increment |  |
| 2 | title | 新闻标题 | varchar(255) |  | NO |  |  |
| 3 | sub_title | 新闻副标题 | varchar(255) |  | NO |  |  |
| 4 | web_link | 新闻网页链接 | varchar(1024) |  | NO |  |  |
| 5 | binance_id | Binance平台新闻ID | varchar(50) |  | NO |  |  |
| 6 | publish_time | 新闻发布时间 | datetime |  | NO |  | 0000-00-00 00:00:00 |
| 7 | created_at | 记录创建时间 | datetime |  | NO | DEFAULT_GENERATED | CURRENT_TIMESTAMP |
| 8 | type | 类型 1-币安新闻 | tinyint unsigned |  | NO |  | 1 |
