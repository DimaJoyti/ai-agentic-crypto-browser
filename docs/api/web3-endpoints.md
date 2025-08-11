# Web3 API Endpoints

This document describes wallet and transaction endpoints exposed by the Web3 service.
All endpoints are protected and require authentication via JWT. Use the Authorization header:

Authorization: Bearer <token>

Base path: /web3

## GET /web3/wallets

List user's connected wallets with optional filters and pagination.

Query parameters:
- chain_id (int, optional): filter by chain ID
- primary (bool, optional): filter by primary wallet (true/false)
- page (int, optional, default 1)
- page_size (int, optional, default 20, max 100)

Responses:
- 200 OK
{
  "wallets": [
    {
      "id": "uuid",
      "user_id": "uuid",
      "address": "0x...",
      "chain_id": 1,
      "wallet_type": "metamask",
      "is_primary": true,
      "created_at": "RFC3339",
      "updated_at": "RFC3339"
    }
  ],
  "pagination": {"page": 1, "page_size": 20, "total_items": 1, "total_pages": 1}
}
- 400 Bad Request: invalid parameters
- 401 Unauthorized: missing/invalid token
- 500 Internal Server Error

Example:
GET /web3/wallets?chain_id=1&primary=true&page=1&page_size=10

## GET /web3/transactions

List user's transactions with filters and pagination.

Query parameters:
- chain_id (int, optional)
- status (string, optional): pending|confirmed|failed
- page (int, optional, default 1)
- page_size (int, optional, default 20, max 100)
- from_time (RFC3339, optional)
- to_time (RFC3339, optional)

Responses:
- 200 OK
{
  "transactions": [
    {
      "id": "uuid",
      "user_id": "uuid",
      "wallet_id": "uuid",
      "tx_hash": "0x...",
      "chain_id": 1,
      "from_address": "0x...",
      "to_address": "0x...",
      "value": "string (wei)",
      "status": "pending|confirmed|failed",
      "block_number": 123,
      "transaction_type": "transfer",
      "created_at": "RFC3339",
      "updated_at": "RFC3339"
    }
  ],
  "pagination": {"page": 1, "page_size": 20, "total_items": 1, "total_pages": 1}
}
- 400 Bad Request
- 401 Unauthorized
- 500 Internal Server Error

Example:
GET /web3/transactions?chain_id=1&status=confirmed&page=1&page_size=10

## Notes
- Authentication: All endpoints require JWT. The service uses middleware.JWT in the server.
- Rate limiting: Global rate limiting is applied via middleware. Be mindful of 429 responses.
- Errors: JSON error messages are returned via http.Error with appropriate status.

If time permits, an OpenAPI specification can be added at docs/api/openapi.yaml for tooling and schema validation.

