package constant

import "errors"

const API_GATEWAY = "http://localhost:8000"

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

var ErrNoProductDiscountsFound = errors.New("no product discounts found")
