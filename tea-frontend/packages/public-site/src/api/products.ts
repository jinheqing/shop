import { api } from './client'

export interface CustomProduct {
  id: number
  product_token?: string
  title: string
  tea_type: string
  tea_shape: string
  mountain_location?: string
  master_name?: string
  raw_tea_source?: string
  unit_price: number
  quantity: number
  shipping_cost: number
  lead_time?: string
  harvest_date?: string
  roasting_date?: string
  status: string
  version: number
  inner_packaging?: string
  outer_packaging?: string
}

export const getByToken = (token: string) => api.get<any, CustomProduct>(`/custom-products/by-token/${token}`)
export const getPublished = () => api.get<any, CustomProduct[]>('/custom-products/published')
