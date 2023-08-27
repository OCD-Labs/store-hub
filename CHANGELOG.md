## StoreHub Changelog

### **Sun 27 Aug 2023**

#### Updated
---
1. Endpoint **`GET /stores`** Response Structure:

  **Before**:
  ```json
  {
    "status": "string",
    "data": {
      "message": "string",
      "result": {
        "stores": [
          {
            "id": 0,
            "name": "string",
            "description": "string",
            "profile_image_url": "string",
            "is_verified": true,
            "category": "string",
            "is_frozen": true,
            "created_at": "2023-08-26T09:21:58.150Z",
            "user_access_levels": [
              0
            ]
          }
        ],
        "metadata": {
          "current_page": 0,
          "page_size": 0,
          "first_page": 0,
          "last_page": 0,
          "total_records": 0
        }
      }
    }
  }
  ```

  **Now**:
  ```json
  {
    "status": "string",
    "data": {
      "message": "string",
      "result": {
        "stores": [
          {
            "store": {
              "id": 0,
              "name": "string",
              "description": "string",
              "profile_image_url": "string",
              "is_verified": true,
              "category": "string",
              "is_frozen": true,
              "created_at": "2023-08-26T09:42:00.733Z",
              "user_access_levels": [
                0
              ]
            },
            "store_owners": [
              {
                "account_id": "string",
                "profile_img_url": "string"
              }
            ]
          }
        ],
        "metadata": {
          "current_page": 0,
          "page_size": 0,
          "first_page": 0,
          "last_page": 0,
          "total_records": 0
        }
      }
    }
  }
  ```

#### Reasons for Change

- To provide a clearer distinction between store details and the owners associated with each store.
- To enhance the data structure for better scalability and clarity in understanding the relationship between stores and their owners.
---
2. Endpoint **`POST /users/{id}/stores`** Response Structure:

  **Before:**

  ```json
  {
    "status": "string",
    "data": {
      "message": "string",
      "result": {
        "store": {
          "id": 0,
          "name": "string",
          "description": "string",
          "profile_image_url": "string",
          "is_verified": true,
          "category": "string",
          "is_frozen": true,
          "created_at": "2023-08-27T08:43:52.196Z",
          "user_access_levels": [
            0
          ]
        },
        "store_owners": [
          {
            "user_id": 0,
            "store_id": 0,
            "access_level": 0,
            "added_at": "2023-08-27T08:43:52.196Z"
          }
        ]
      }
    }
  }
  ```
  **Now:**

  ```json
  {
    "status": "string",
    "data": {
      "message": "string",
      "result": {
        "store": {
          "id": 0,
          "name": "string",
          "description": "string",
          "profile_image_url": "string",
          "is_verified": true,
          "category": "string",
          "is_frozen": true,
          "created_at": "2023-08-27T08:33:37.795Z",
          "user_access_levels": [
            0
          ]
        },
        "store_owners": [
          {
            "account_id": "string",
            "profile_img_url": "string",
            "access_levels": [
              0
            ],
            "added_at": "2023-08-27T08:33:37.795Z"
          }
        ]
      }
    }
  }
  ```

3. Changed endpoint **`POST /users/{id}/stores`** to **`POST /inventory/stores`** eliminating the need for `id` in the path variables.
4. Changed endpoint **`GET /users/{id}/stores`** to **`GET /inventory/stores`** eliminating the need for `id` in the path variables.
5. Changed endpoint **`POST /users/user_id}/stores/{store_id}/items`** to **`POST /inventory/stores/{store_id}/items`** eliminating the need for `user_id` in the path variables.
6. Changed endpoint **`GET /users/user_id}/stores/{store_id}/items`** to **`GET /inventory/stores/{store_id}/items`** eliminating the need for `user_id` in the path variables.
7. Changed endpoint **`PATCH /users/user_id}/stores/{store_id}/items/{item_id}`** to **`PATCH /inventory/stores/{store_id}/items/{item_id}`** eliminating the need for `user_id` in the path variables.
8. Changed endpoint **`DELETE /users/user_id}/stores/{store_id}/items/{item_id}`** to **`DELETE /inventory/stores/{store_id}/items/{item_id}`** eliminating the need for `user_id` in the path variables.
8. Changed endpoint **`PATCH /users/user_id}/stores/{store_id}`** to **`PATCH /inventory/stores/{store_id}`** eliminating the need for `user_id` in the path variables.
9. Changed endpoint **`POST /stores/{store_id}/orders`** to **`POST /inventory/stores/{store_id}/orders`**.
10. Changed endpoint **`GET /stores/{store_id}/orders`** to **`GET /inventory/stores/{store_id}/orders`**.
11. Changed endpoint **`GET /stores/{store_id}/orders/{store_id}`** to **`GET /inventory/stores/{store_id}/orders/{order_id}`**.
12. Changed endpoint **`PATCH /stores/{store_id}/orders/{store_id}`** to **`PATCH /inventory/stores/{store_id}/orders/{order_id}`**.
13. Changed endpoint **`GET /users/user_id}/stores/{store_id}/sales`** to **`GET /inventory/stores/{store_id}/sales`** eliminating the need for `user_id` in the path variables.
14. Changed endpoint **`GET /users/user_id}/stores/{store_id}/sales/{sale_id}`** to **`GET /inventory/stores/{store_id}/sales/{sale_id}`** eliminating the need for `user_id` in the path variables.
15. Changed endpoint **`GET /users/user_id}/stores/{store_id}/sales-overview`** to **`GET /inventory/stores/{store_id}/sales-overview`** eliminating the need for `user_id` in the path variables.

> Make sure to check and update the path variables for each changed endpoint stated above.
