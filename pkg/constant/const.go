package constant

import "errors"

// const API_GATEWAY = "http://api_gateway_service:9000"

// const OAUTH_SERVICE = "http://oauth_service:8080/auth"
// const PRODUCT_SERVICE = "http://product_service:8081/api/products"
// const USER_SERVICE = "http://user_service:8082/api/users"
// const NEWS_SERVICE = "http://news_service:8083/api/news"
// const CART_ITEM_SERVICE = "http://cart_items_service:8084/api/cartItems"
// const CART_SERVICE = "http://cart_service:8085/api/carts"
// const CATEGORY_SERVICE = "http://category_service:8086/api/categories"
// const COURIER_SERVICE = "http://courier_service:8087/api/couriers"
// const DISCOUNT_SERVICE = "http://discount_service:8088/api/discounts"
// const FREIGHT_RATE_SERVICE = "http://freight_rate_service:8089/api/freightRates"
// const ORDER_SERVICE = "http://order_service:8090/api/orders"
// const ORDER_DETAILS_SERVICE = "http://order_detail_service:8091/api/orderDetails"
// const PRODUCT_DISCOUNT_SERVICE = "http://product_discount_service:8092/api/productDiscounts"
// const REVIEW_SERVICE = "http://review_service:8093/api/reviews"
// const PAYMENT_SERVICE = "http://payment_service:8094/api/payments"
// const VOUCHER_SERVICE = "http://voucher_service:8095/api/vouchers"
// const MAIL_SERVICE = "http://mail_service:8096/api/mail"
// const MOMO_SERVICE = "http://momo_service:8097/api/momo"
// const VNPAY_SERVICE = "http://vnpay_service:8098/api/vnpay"
// const AUTH_SERVICE = "http://auth_service:8099/api/authentication"

const API_GATEWAY = "http://localhost:9000"

const OAUTH_SERVICE = "http://localhost:8080/auth"
const PRODUCT_SERVICE = "http://localhost:8081/api/products"
const USER_SERVICE = "http://localhost:8082/api/users"
const NEWS_SERVICE = "http://localhost:8083/api/news"
const CART_ITEM_SERVICE = "http://localhost:8084/api/cartItems"
const CART_SERVICE = "http://localhost:8085/api/carts"
const CATEGORY_SERVICE = "http://localhost:8086/api/categories"
const COURIER_SERVICE = "http://localhost:8087/api/couriers"
const DISCOUNT_SERVICE = "http://localhost:8088/api/discounts"
const FREIGHT_RATE_SERVICE = "http://localhost:8089/api/freightRates"
const ORDER_SERVICE = "http://localhost:8090/api/orders"
const ORDER_DETAILS_SERVICE = "http://localhost:8091/api/orderDetails"
const PRODUCT_DISCOUNT_SERVICE = "http://localhost:8092/api/productDiscounts"
const REVIEW_SERVICE = "http://localhost:8093/api/reviews"
const PAYMENT_SERVICE = "http://localhost:8094/api/payments"
const VOUCHER_SERVICE = "http://localhost:8095/api/vouchers"
const MAIL_SERVICE = "http://localhost:8096/api/mail"
const MOMO_SERVICE = "http://localhost:8097/api/momo"
const VNPAY_SERVICE = "http://localhost:8098/api/vnpay"
const AUTH_SERVICE = "http://localhost:8099/api/authentication"

const PAYMENT_RESPONSE_REJECT_URL = "http://localhost:3000/reject"
const PAYMENT_RESPONSE_CONFIRM_URL = "http://localhost:3000/confirm"

const PAYMENT_STATUS_PENDING = "Pending"
const PAYMENT_STATUS_COMPLETED = "Complete"
const PAYMENT_STATUS_FAILED = "Failed"

const ORDER_STATUS_PENDING = "Pending"
const ORDER_STATUS_COMPLETED = "Complete"
const ORDER_STATUS_FAILED = "Failed"
const ORDER_STATUS_CANCELED = "Canceled"

const PAYMENT_METHOD_MOMO = "MoMo"
const PAYMENT_METHOD_VNPAY = "VnPay"

const VOUCHER_DISCOUNT_TYPE_PERCENTAGE = "Percentage"
const VOUCHER_DISCOUNT_TYPE_FIXED = "Fixed"

const DEFAULT_USER_IMAGE = "https://firebasestorage.googleapis.com/v0/b/storage-8b808.appspot.com/o/OIP.jpeg?alt=media&token=60195a0a-2fd6-4c66-9e3a-0f7f80eb8473"

var ErrNoProductDiscountsFound = errors.New("no product discounts found")
