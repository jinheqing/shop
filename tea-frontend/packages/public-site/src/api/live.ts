import { api } from './client'

export interface LiveRoom {
  id: number
  room_name: string
  room_type: string
  obs_rtmp_url?: string
  obs_rtmp_key?: string
  camera_rtmp_url?: string
  status: string
  peak_viewers?: number
}

export const getSlowPresets = () => api.get<any, any[]>('/slow-presets')
export const listLiveRooms = () => api.get<any, { items: LiveRoom[] }>('/live-rooms')
export const getLivekitToken = (roomName: string, identity: string) =>
  api.post<any, any>('/livekit/token', { room_name: roomName, identity })
