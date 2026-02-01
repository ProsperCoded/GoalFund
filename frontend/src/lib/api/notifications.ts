import apiClient from "./config"
import type { Notification, NotificationListResponse } from "./types"

export const notificationsApi = {
  // Get user's notifications
  getNotifications: async (params?: {
    page?: number
    limit?: number
    unread_only?: boolean
  }): Promise<NotificationListResponse> => {
    const response = await apiClient.get<NotificationListResponse>("/notifications", { params })
    return response.data
  },

  // Mark notification as read
  markAsRead: async (id: string): Promise<{ message: string }> => {
    const response = await apiClient.put<{ message: string }>(`/notifications/${id}/read`)
    return response.data
  },

  // Mark all as read
  markAllAsRead: async (): Promise<{ message: string }> => {
    const response = await apiClient.put<{ message: string }>("/notifications/read-all")
    return response.data
  },

  // Get unread count
  getUnreadCount: async (): Promise<{ count: number }> => {
    const response = await apiClient.get<{ count: number }>("/notifications/unread-count")
    return response.data
  },

  // Delete notification
  delete: async (id: string): Promise<{ message: string }> => {
    const response = await apiClient.delete<{ message: string }>(`/notifications/${id}`)
    return response.data
  },
}

// WebSocket connection for real-time notifications
export class NotificationWebSocket {
  private ws: WebSocket | null = null
  private reconnectAttempts = 0
  private maxReconnectAttempts = 5
  private reconnectDelay = 1000
  private onMessageCallback: ((notification: Notification) => void) | null = null
  private onConnectionChangeCallback: ((connected: boolean) => void) | null = null

  connect(token: string) {
    const wsUrl = import.meta.env.VITE_WS_URL || "ws://localhost/api/v1/notifications/ws"
    
    try {
      this.ws = new WebSocket(`${wsUrl}?token=${token}`)
      
      this.ws.onopen = () => {
        console.log("WebSocket connected")
        this.reconnectAttempts = 0
        this.onConnectionChangeCallback?.(true)
      }

      this.ws.onmessage = (event) => {
        try {
          const notification = JSON.parse(event.data) as Notification
          this.onMessageCallback?.(notification)
        } catch (error) {
          console.error("Failed to parse notification:", error)
        }
      }

      this.ws.onclose = () => {
        console.log("WebSocket disconnected")
        this.onConnectionChangeCallback?.(false)
        this.attemptReconnect(token)
      }

      this.ws.onerror = (error) => {
        console.error("WebSocket error:", error)
      }
    } catch (error) {
      console.error("Failed to connect WebSocket:", error)
    }
  }

  private attemptReconnect(token: string) {
    if (this.reconnectAttempts < this.maxReconnectAttempts) {
      this.reconnectAttempts++
      const delay = this.reconnectDelay * Math.pow(2, this.reconnectAttempts - 1)
      console.log(`Reconnecting in ${delay}ms (attempt ${this.reconnectAttempts})`)
      setTimeout(() => this.connect(token), delay)
    }
  }

  onMessage(callback: (notification: Notification) => void) {
    this.onMessageCallback = callback
  }

  onConnectionChange(callback: (connected: boolean) => void) {
    this.onConnectionChangeCallback = callback
  }

  disconnect() {
    if (this.ws) {
      this.ws.close()
      this.ws = null
    }
  }

  isConnected(): boolean {
    return this.ws?.readyState === WebSocket.OPEN
  }
}

export const notificationWs = new NotificationWebSocket()
