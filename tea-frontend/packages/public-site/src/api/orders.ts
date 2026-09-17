import { api } from './client'

export interface CreateOrderReq {
  custom_product_id: number
  unit_price: number
  quantity: number
  shipping_cost: number
  billing_address: any
  delivery_address: any
}

export const createOrder = (body: CreateOrderReq) => api.post<any, any>('/orders', body)
export const getOrder = (id: number) => api.get<any, any>(`/orders/${id}`)
export const initPayment = (id: number, gateway: string) => api.post<any, any>(`/orders/${id}/payment/init`, { gateway })
