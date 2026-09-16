---
id: product
updated: 2026-09-16
---

# Product

店名：高木桑收藏品店（TKG Collection Store）。地址叙事：广州市白云区三元里。

卖纸本收藏与邮政服务，不是通用电商、不是投资平台。标价是零售整理价，不承诺升值，不提供鉴定证书。

## Catalog (shared)

唯一商品源：`services/rag/data/catalog.json`。

| id | 角色 |
| --- | --- |
| `postcard-forward` | 代寄服务，18 元起 |
| `mail-route-test` | 邮路测试服务 |
| 其他 id | 可加购的实物 / 定制（明信片、邮票、礼盒、信件） |

服务 id 集合在 `$lib/server/catalog.ts` 的 `SERVICE_IDS`。

## Fulfillment

- 实寄只走中国邮政：贴票、销戳、交寄
- 不提供商业快递、不承诺快递式追踪号
- `/shop/sent` 提交的是申请（`inquiries`），不是已支付订单
- 未登录可提交；登录后用 JWT 挂到 `user_id`，在 `/userspace` 可见

## Voice

中文、具体、克制。避免「全网最低 / 必升值 / 包邮到门」。知识库回答要标明来自切块，价格以商品条目为准。
