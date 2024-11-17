# Order Management System

This is an Order Management System built with Go, Gin, GORM, Postgres. It provides APIs to create, get, and cancel orders with a test user authentication system.


### Project Structure
```
order-management/
├── config/
│   ├── config.go          # Project configuration
│   ├── database.go        # Database configuration
├── database/
│   ├── migrations.go      # Database migrations
├── handlers/
│   ├── orders.go          # Order handlers
│   ├── auth.go            # Authentication handlers
├── middleware/
│   ├── auth.go            # Authentication middleware
├── models/
│   ├── order.go           # Order model
├── utils/
│   ├── validation.go      # Validation utilities
├── .env                   # Environment variables
├── .env.example           # Example environment variables
├── go.mod                 # Go module file
├── go.sum                 # Go dependencies file
├── main.go                # Entry point of the application
├── README.md              # Project documentation
```
### Prerequisites

- Go 1.19 or later
- PostgreSQL

### Installation

## 1. **Clone the repository:**
```bash
git clone https://github.com/shadhin-int/order-management.git
cd order-management
  ```

### 2. Set Up Environment Variables

Copy the `.env.example` file to `.env` and adjust the environment variables according to your setup.

    cp .env.example .env

### 3. Install Dependencies

Make sure all Go dependencies are installed by running:

```
go mod tidy
```

## Run the Project

#### Start the server

```bash
go run main.go
```
The server will start on `localhost:8080`.

> **_NOTE:_**  When run the project first time, it will automatically create the required tables in the database.

### Testing the API

You can use Postman or `curl` to test the API endpoints.

1. ***Login***

    - **URL:** `{BASE_URL}/api/v1/login`
    - **Method:** `POST`
    - **Sample Request Body:**
        ```json
        {
           "username": "01901901901@mailinator.com",
           "password": "321dsaf"
        }
      ```
    - **Sample Response:**
       ```json
      {
         "token_type": "Bearer",
         "expires_in": 1,
         "access_token": "eyJhb***********",
         "refresh_token": "refresh_eyJhbGc************"
      }
      ```
2. ***Logout***
    - **URL:** `{BASE_URL}/api/v1/logout`
    - **Method:** `POST`
    - **Headers:**
        ```
        Authorization: Bearer {ACCESS_TOKEN}
        ```
   - **Response:**
     ```json
     {
        "message": "Logout successful",
        "type": "success",
        "code": 200
     }
     ```

3. ***Create an Order***
    - **URL:** `{BASE_URL}/api/v1/orders`
    - **Method:** `POST`
    - **Headers:**
        ```
        Authorization: Bearer {ACCESS_TOKEN}
        ```
    - **Sample Request Body:**
        ```json
        {
           "store_id": 131172,
           "merchant_order_id": "12345",
           "recipient_name": "John Doe",
           "recipient_phone": "01901901901",
           "recipient_address": "banani, gulshan 2, dhaka, bangladesh",
           "recipient_city": 1,
           "recipient_zone": 1,
           "recipient_area": 1,
           "delivery_type": 48,
           "item_type": 2,
           "special_instruction": "Leave at the front door",
           "item_quantity": 1,
           "item_weight": 0.5,
           "amount_to_collect": 1000,
           "item_description": "Books"
        }
        ```
    - **Sample Success Response:**
        ```json
      {
          "message": "Order created successfully",
          "type": "success",
          "code": 200,
          "data": {
              "consignment_id": "DA241117wEEEEE",
              "delivery_fee": 60,
              "merchant_order_id": "12345",
              "order_status": "Pending"
          }
      }
      ```

   - **Sample Validation Error Response:**
       ```json
       {
           "message": "Please fix the given errors",
           "type": "error",
           "code": 422,
           "errors": {
               "store_id": [
                   "Store ID is required", "Wrong Store ID"
               ],
               "recipient_name": [
                   "Recipient Name is required"
               ],
               "recipient_phone": [
                   "Recipient Phone is required"
               ],
               "recipient_address": [
                   "Recipient Address is required"
               ],
               "recipient_city": [
                   "Recipient City is required"
               ],
               "recipient_zone": [
                   "Recipient Zone is required"
               ],
               "recipient_area": [
                   "Recipient Area is required"
               ],
               "delivery_type": [
                   "Delivery Type is required"
               ],
               "item_type": [
                   "Item Type is required"
               ],
               "item_quantity": [
                   "Item Quantity is required"
               ],
               "item_weight": [
                   "Item Weight is required"
               ],
               "amount_to_collect": [
                   "Amount to Collect is required"
               ],
               "item_description": [
                   "Item Description is required"
               ]
           }
       }
       ```
   - **Sample Phone Validation Error Response:**
       ```json
       {
           "message": "Please fix the given errors",
           "type": "error",
           "code": 422,
           "errors": {
               "recipient_phone": [
                   "Recipient Phone must be a valid phone number"
               ]
           }
       }
       ```
   - **Unauthorized Response:**
       ```json
       {
           "message": "Unauthorized",
           "type": "error",
           "code": 401
       }
       ```
4. ***Get Orders***
    - **URL:** `{BASE_URL}/api/v1/orders/all`
    - **Method:** `GET`
    - **Headers:**
        ```
        Authorization: Bearer {ACCESS_TOKEN}
        ```
    - **Query Parameters:**
        ```
        order_status: Pending/Cancelled/Completed
        limit: 1
        page: 10
        ```
    - **Sample Success Response:**
        ```json
      {
          "message": "Orders fetched successfully",
          "type": "success",
          "code": 200,
          "data": {
              "data":[
                  {
                      "cod_fee": 12,
                      "delivery_fee": 60,
                      "discount": 0,
                      "instruction": "Leave at the front door",
                      "item_type": "Parcel",
                      "merchant_order_id": "12345",
                      "order_amount": 1000,
                      "order_consignment_id": "DA241117wEEEEE",
                      "order_created_at": "2024-11-17 17:52:02",
                      "order_description": "Books",
                      "order_status": "Cancelled",
                      "order_type": "Regular",
                      "order_type_id": 1,
                      "promo_discount": 0,
                      "recipient_address": "banani, gulshan 2, dhaka, bangladesh",
                      "recipient_name": "John Doe",
                      "recipient_phone": "01901901901",
                      "total_fee": 72
                  }
              ],
              "total": 1,
              "current_page": 1,
              "per_page": 10,
              "total_in_page": 1,
              "total_pages": 1
          }
      }
      ```
    - **Unauthorized Response:**
        ```json
        {
            "message": "Unauthorized",
            "type": "error",
            "code": 401
        }
        ```
5. ***Cancel an Order***
    - **URL:** `{BASE_URL}/api/v1/orders/{consignment_id}/cancel`
    - **Method:** `PUT`
    - **Headers:**
        ```
        Authorization: Bearer {ACCESS_TOKEN}
        ```
    - **Sample Success Response:**
        ```json
      {
          "message": "Order cancelled successfully",
          "type": "success",
          "code": 200
      }
      ```
    - **Sample Error Response:**
        ```json
        {
            "message": "Order not found",
            "type": "error",
            "code": 404
        }
        ```
    - ***Already Cancelled Order:***
      ```json
      {
          "message": "Order already cancelled",
          "type": "error",
          "code": 400
      }
      ```
    - **Unauthorized Response:**
        ```json
        {
            "message": "Unauthorized",
            "type": "error",
            "code": 401
        }
        ```
    - **Except Pending Order To Cancel Response:**
        ```json
        {
            "message": "Please contact cx to cancel the order",
            "type": "error",
            "code": 400
        }
        ```
    
 ### Further Improvements
- Though for user authentication used test data, it should be implemented with a proper authentication system.
- Currently only one user is available in the system. It should be implemented with multiple users.
- For order maintain by a single table, it should be implemented with multiple tables with proper relationships.
- For the order status, it should be implemented with a separate table.
- For the city, zone, and area, it should be implemented with a separate table.
- Need to implement blacklisting the token.
- Need to implement the refresh token.
- For logging, it should be implemented with a proper logging system.